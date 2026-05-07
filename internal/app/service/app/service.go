package app

import (
	"context"
	"encoding/json"
	"fmt"
	cryptoutil "intehub/internal/utils/crypto"
	"intehub/appl"
	appModel "intehub/internal/app/models/app"
	"intehub/internal/app/models/appstore"
	"intehub/pkg/params"
	"strconv"
	"time"
)

type service struct {
	model         appModel.Model
	appstoreModel appstore.Model
	cryptoKey     string
}

func New(model appModel.Model, appstoreModel appstore.Model, cryptoKey string) Service {
	return &service{
		model:         model,
		appstoreModel: appstoreModel,
		cryptoKey:     cryptoKey,
	}
}

func (s *service) Create(app *appModel.App) error {
	return s.model.Create(app)
}

func (s *service) CreateWithConfigs(app *appModel.App, configs []ConfigInput) error {
	// 创建应用
	if err := s.model.Create(app); err != nil {
		return err
	}

	// 创建配置
	for _, cfg := range configs {
		if cfg.Key == "" {
			continue
		}
		config := &appstore.AppConfig{
			AppID:     app.ID,
			Key:       cfg.Key,
			Value:     cfg.Value,
			Type:      cfg.Type,
			Encrypted: cfg.Encrypted,
		}
		if config.Type == "" {
			config.Type = "string"
		}
		// 如果启用加密，则加密存储
		if cfg.Encrypted && cfg.Value != "" && s.cryptoKey != "" {
			encryptedValue, err := cryptoutil.AESEncrypt([]byte(cfg.Value), s.cryptoKey)
			if err != nil {
				return fmt.Errorf("加密配置值失败: %w", err)
			}
			config.Value = encryptedValue
		}
		if err := s.appstoreModel.CreateConfig(config); err != nil {
			return fmt.Errorf("创建配置失败: %w", err)
		}
	}

	return nil
}

func (s *service) GetByID(id uint) (*appModel.App, error) {
	return s.model.GetByID(id)
}

func (s *service) GetByIDWithConfigs(id uint) (*appModel.App, []*appstore.AppConfig, error) {
	app, err := s.model.GetByID(id)
	if err != nil {
		return nil, nil, err
	}

	configs, err := s.appstoreModel.GetConfigsByAppID(id)
	if err != nil {
		return app, []*appstore.AppConfig{}, nil
	}

	return app, configs, nil
}

// GetByIDWithConfigsMasked 获取应用配置，加密字段返回遮蔽值
func (s *service) GetByIDWithConfigsMasked(id uint) (*appModel.App, []*appstore.AppConfig, error) {
	app, configs, err := s.GetByIDWithConfigs(id)
	if err != nil {
		return nil, nil, err
	}

	// 对加密字段进行遮蔽
	for _, config := range configs {
		if config.Encrypted {
			config.Value = "******"
		}
	}

	return app, configs, nil
}

func (s *service) List(userID uint) ([]*appModel.App, error) {
	return s.model.List(userID)
}

func (s *service) Update(app *appModel.App) error {
	return s.model.Update(app)
}

func (s *service) UpdateWithConfigs(app *appModel.App, configs []ConfigInput) error {
	// 更新应用
	if err := s.model.Update(app); err != nil {
		return err
	}

	// 删除旧配置
	if err := s.appstoreModel.DeleteConfigsByAppID(app.ID); err != nil {
		return fmt.Errorf("删除旧配置失败: %w", err)
	}

	// 创建新配置
	for _, cfg := range configs {
		if cfg.Key == "" {
			continue
		}
		config := &appstore.AppConfig{
			AppID:     app.ID,
			Key:       cfg.Key,
			Value:     cfg.Value,
			Type:      cfg.Type,
			Encrypted: cfg.Encrypted,
		}
		if config.Type == "" {
			config.Type = "string"
		}
		// 如果启用加密，则加密存储
		if cfg.Encrypted && cfg.Value != "" && s.cryptoKey != "" {
			encryptedValue, err := cryptoutil.AESEncrypt([]byte(cfg.Value), s.cryptoKey)
			if err != nil {
				return fmt.Errorf("加密配置值失败: %w", err)
			}
			config.Value = encryptedValue
		}
		if err := s.appstoreModel.CreateConfig(config); err != nil {
			return fmt.Errorf("创建配置失败: %w", err)
		}
	}

	return nil
}

func (s *service) Delete(id uint) error {
	return s.model.Delete(id)
}

func (s *service) Run(id uint) (*appModel.AppLog, error) {
	return s.RunWithInput(id, make(map[string]interface{}))
}

func (s *service) RunWithInput(id uint, input map[string]interface{}) (*appModel.AppLog, error) {
	// 获取应用信息
	app, err := s.model.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("应用不存在: %w", err)
	}

	if !app.Enabled {
		return nil, fmt.Errorf("应用已禁用")
	}

	// 创建应用日志
	log := &appModel.AppLog{
		AppID:     app.ID,
		AppName:   app.Name,
		Status:    "running",
		StartedAt: time.Now(),
	}

	startTime := time.Now()

	// 根据语言类型创建应用实例
	var appIns *appl.AppIns
	switch app.Language {
	case "go", "":
		goApp, err := appl.NewGoAppIns(app.Code)
		if err != nil {
			log.Status = "error"
			log.Error = fmt.Sprintf("创建应用实例失败: %v", err)
			log.FinishedAt = time.Now()
			log.Duration = time.Since(startTime).Milliseconds()
			s.model.CreateLog(log)
			return log, fmt.Errorf("创建应用实例失败: %w", err)
		}
		appIns = &appl.AppIns{Runnable: goApp}
	case "javascript":
		jsApp, err := appl.NewJavaScriptAppIns(app.Code)
		if err != nil {
			log.Status = "error"
			log.Error = fmt.Sprintf("创建 JavaScript 应用实例失败: %v", err)
			log.FinishedAt = time.Now()
			log.Duration = time.Since(startTime).Milliseconds()
			s.model.CreateLog(log)
			return log, fmt.Errorf("创建 JavaScript 应用实例失败: %w", err)
		}
		appIns = &appl.AppIns{Runnable: jsApp}
	default:
		log.Status = "error"
		log.Error = fmt.Sprintf("不支持的语言类型: %s", app.Language)
		log.FinishedAt = time.Now()
		log.Duration = time.Since(startTime).Milliseconds()
		s.model.CreateLog(log)
		return log, fmt.Errorf("不支持的语言类型: %s", app.Language)
	}

	// 准备输入参数：合并应用配置和请求参数
	inputParams := make(params.Params)

	// 1. 先加载应用配置作为默认参数
	configs, err := s.appstoreModel.GetConfigsByAppID(app.ID)
	if err == nil && len(configs) > 0 {
		for _, config := range configs {
			configValue := config.Value

			// 如果是加密配置，则解密
			if config.Encrypted && configValue != "" && s.cryptoKey != "" {
				decrypted, err := cryptoutil.AESDecrypt(configValue, s.cryptoKey)
				if err == nil {
					configValue = string(decrypted)
				} else {
					log.Error = fmt.Sprintf("解密配置失败: %s", err.Error())
					s.model.CreateLog(log)
					return log, fmt.Errorf("解密配置 [%s] 失败: %w", config.Key, err)
				}
			}

			// 根据配置类型转换值
			var value interface{}
			switch config.Type {
			case "number":
				// 尝试转换为数字
				if num, err := strconv.ParseFloat(configValue, 64); err == nil {
					value = num
				} else {
					value = configValue
				}
			case "boolean":
				// 转换为布尔值
				value = configValue == "true" || configValue == "1"
			case "json":
				// 尝试解析 JSON
				var jsonValue interface{}
				if err := json.Unmarshal([]byte(configValue), &jsonValue); err == nil {
					value = jsonValue
				} else {
					value = configValue
				}
			default:
				// string 或其他类型，直接使用字符串
				value = configValue
			}
			inputParams[config.Key] = value
		}
	}

	// 2. 请求参数会覆盖配置参数
	for k, v := range input {
		inputParams[k] = v
	}

	// 执行应用
	ctx := context.Background()
	output, execErr := appIns.Run(ctx, inputParams)

	// 记录执行结果
	log.FinishedAt = time.Now()
	log.Duration = time.Since(startTime).Milliseconds()

	if execErr != nil {
		log.Status = "error"
		log.Error = execErr.Error()
	} else {
		log.Status = "success"
		// 将输出转换为 JSON 字符串保存到日志
		if outputBytes, err := json.Marshal(output); err == nil {
			log.Output = string(outputBytes)
		} else {
			log.Output = fmt.Sprintf("%v", output)
		}
	}

	// 保存日志
	if err := s.model.CreateLog(log); err != nil {
		return log, fmt.Errorf("保存日志失败: %w", err)
	}

	return log, nil
}

func (s *service) GetLogs(appID *uint, page, pageSize int) ([]*appModel.AppLog, int64, error) {
	return s.model.GetLogs(appID, page, pageSize)
}

func (s *service) PublishToStore(appID uint, displayName, description, icon, category, version, author, tags, configSchema string) (*appstore.AppTemplate, error) {
	// 获取应用信息
	app, err := s.model.GetByID(appID)
	if err != nil {
		return nil, fmt.Errorf("应用不存在: %w", err)
	}

	// 检查应用是否已启用
	if !app.Enabled {
		return nil, fmt.Errorf("应用已禁用，无法发布")
	}

	// 验证代码是否可以编译
	if app.Language == "go" || app.Language == "" {
		_, err := appl.NewGoAppIns(app.Code)
		if err != nil {
			return nil, fmt.Errorf("应用代码编译失败: %w", err)
		}
	}

	// 检查是否已存在同名模板
	existingTemplate, _ := s.appstoreModel.GetTemplateByName(app.Name)
	if existingTemplate != nil {
		return nil, fmt.Errorf("应用商店中已存在名为 %s 的模板", app.Name)
	}

	// 创建应用商店模板
	template := &appstore.AppTemplate{
		Name:         app.Name,
		DisplayName:  displayName,
		Description:  description,
		Icon:         icon,
		Code:         app.Code,
		Language:     app.Language,
		Category:     category,
		Version:      version,
		Author:       author,
		Tags:         tags,
		ConfigSchema: configSchema,
		Enabled:      true,
	}

	// 如果语言为空，默认为 go
	if template.Language == "" {
		template.Language = "go"
	}

	// 保存到应用商店
	if err := s.appstoreModel.CreateTemplate(template); err != nil {
		return nil, fmt.Errorf("发布到应用商店失败: %w", err)
	}

	// 注意：应用配置不会自动复制到模板
	// 用户从模板安装应用时，需要重新配置参数

	return template, nil
}
