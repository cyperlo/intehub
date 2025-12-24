package appstore

import (
	"fmt"
	"intehub/appl"
	"intehub/internal/app/models"
	"intehub/internal/app/models/app"
	"intehub/internal/app/models/appstore"
)

type CreateTemplateRequest struct {
	Name         string `json:"name" binding:"required"`
	DisplayName  string `json:"display_name" binding:"required"`
	Description  string `json:"description"`
	Icon         string `json:"icon"`
	Code         string `json:"code" binding:"required"`
	Language     string `json:"language"`
	Category     string `json:"category"`
	Version      string `json:"version"`
	Author       string `json:"author"`
	Tags         string `json:"tags"`
	ConfigSchema string `json:"config_schema"`
	Enabled      *bool  `json:"enabled"`
}

type UpdateTemplateRequest struct {
	DisplayName  *string `json:"display_name"`
	Description  *string `json:"description"`
	Icon         *string `json:"icon"`
	Code         *string `json:"code"`
	Language     *string `json:"language"`
	Category     *string `json:"category"`
	Version      *string `json:"version"`
	Author       *string `json:"author"`
	Tags         *string `json:"tags"`
	ConfigSchema *string `json:"config_schema"`
	Enabled      *bool   `json:"enabled"`
}

type CreateConfigRequest struct {
	AppID     uint   `json:"app_id" binding:"required"`
	Key       string `json:"key" binding:"required"`
	Value     string `json:"value" binding:"required"`
	Type      string `json:"type"`
	Encrypted bool   `json:"encrypted"`
}

type UpdateConfigRequest struct {
	Value     *string `json:"value"`
	Type      *string `json:"type"`
	Encrypted *bool   `json:"encrypted"`
}

type InstallRequest struct {
	TemplateID  uint              `json:"template_id" binding:"required"`
	Name        string            `json:"name" binding:"required"`
	Description string            `json:"description"`
	UserID      uint              `json:"user_id" binding:"required"`
	Configs     map[string]string `json:"configs"`
}

type Service struct {
	model    models.Model
	appStore *appl.AppStore
}

func New(model models.Model) *Service {
	return &Service{
		model:    model,
		appStore: appl.NewAppStore(model.AppStoreModel()),
	}
}

// Template methods
func (s *Service) ListTemplates(category string) ([]*appstore.AppTemplate, error) {
	return s.model.AppStoreModel().ListTemplates(category)
}

func (s *Service) GetTemplate(id uint) (*appstore.AppTemplate, error) {
	return s.model.AppStoreModel().GetTemplateByID(id)
}

func (s *Service) CreateTemplate(req *CreateTemplateRequest) (*appstore.AppTemplate, error) {
	template := &appstore.AppTemplate{
		Name:         req.Name,
		DisplayName:  req.DisplayName,
		Description:  req.Description,
		Icon:         req.Icon,
		Code:         req.Code,
		Language:     req.Language,
		Category:     req.Category,
		Version:      req.Version,
		Author:       req.Author,
		Tags:         req.Tags,
		ConfigSchema: req.ConfigSchema,
		Enabled:      true,
	}

	if req.Enabled != nil {
		template.Enabled = *req.Enabled
	}

	if template.Language == "" {
		template.Language = "go"
	}

	// 验证代码是否可以编译
	if template.Language == "go" {
		_, err := appl.NewGoAppIns(template.Code)
		if err != nil {
			return nil, fmt.Errorf("invalid code: %w", err)
		}
	}

	if err := s.model.AppStoreModel().CreateTemplate(template); err != nil {
		return nil, err
	}

	return template, nil
}

func (s *Service) UpdateTemplate(id uint, req *UpdateTemplateRequest) (*appstore.AppTemplate, error) {
	template, err := s.model.AppStoreModel().GetTemplateByID(id)
	if err != nil {
		return nil, err
	}

	if req.DisplayName != nil {
		template.DisplayName = *req.DisplayName
	}
	if req.Description != nil {
		template.Description = *req.Description
	}
	if req.Icon != nil {
		template.Icon = *req.Icon
	}
	if req.Code != nil {
		// 验证代码
		if template.Language == "go" {
			_, err := appl.NewGoAppIns(*req.Code)
			if err != nil {
				return nil, fmt.Errorf("invalid code: %w", err)
			}
		}
		template.Code = *req.Code
	}
	if req.Language != nil {
		template.Language = *req.Language
	}
	if req.Category != nil {
		template.Category = *req.Category
	}
	if req.Version != nil {
		template.Version = *req.Version
	}
	if req.Author != nil {
		template.Author = *req.Author
	}
	if req.Tags != nil {
		template.Tags = *req.Tags
	}
	if req.ConfigSchema != nil {
		template.ConfigSchema = *req.ConfigSchema
	}
	if req.Enabled != nil {
		template.Enabled = *req.Enabled
	}

	if err := s.model.AppStoreModel().UpdateTemplate(template); err != nil {
		return nil, err
	}

	return template, nil
}

func (s *Service) DeleteTemplate(id uint) error {
	return s.model.AppStoreModel().DeleteTemplate(id)
}

// Config methods
func (s *Service) GetAppConfigs(appID uint) ([]*appstore.AppConfig, error) {
	return s.model.AppStoreModel().GetConfigsByAppID(appID)
}

func (s *Service) CreateConfig(req *CreateConfigRequest) (*appstore.AppConfig, error) {
	config := &appstore.AppConfig{
		AppID:     req.AppID,
		Key:       req.Key,
		Value:     req.Value,
		Type:      req.Type,
		Encrypted: req.Encrypted,
	}

	if config.Type == "" {
		config.Type = "string"
	}

	if err := s.model.AppStoreModel().CreateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func (s *Service) UpdateConfig(id uint, req *UpdateConfigRequest) (*appstore.AppConfig, error) {
	config, err := s.model.AppStoreModel().GetConfig(0, "")
	if err != nil {
		return nil, err
	}

	if req.Value != nil {
		config.Value = *req.Value
	}
	if req.Type != nil {
		config.Type = *req.Type
	}
	if req.Encrypted != nil {
		config.Encrypted = *req.Encrypted
	}

	if err := s.model.AppStoreModel().UpdateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func (s *Service) DeleteConfig(id uint) error {
	return s.model.AppStoreModel().DeleteConfig(id)
}

// InstallFromTemplate 从模板安装应用
func (s *Service) InstallFromTemplate(req *InstallRequest) (*app.App, error) {
	// 获取模板
	template, err := s.model.AppStoreModel().GetTemplateByID(req.TemplateID)
	if err != nil {
		return nil, fmt.Errorf("template not found: %w", err)
	}

	// 创建应用
	newApp := &app.App{
		Name:        req.Name,
		Description: req.Description,
		Code:        template.Code,
		Language:    template.Language,
		Enabled:     true,
		UserID:      req.UserID,
	}

	if err := s.model.AppModel().Create(newApp); err != nil {
		return nil, fmt.Errorf("failed to create app: %w", err)
	}

	// 保存配置
	for key, value := range req.Configs {
		config := &appstore.AppConfig{
			AppID: newApp.ID,
			Key:   key,
			Value: value,
			Type:  "string",
		}
		if err := s.model.AppStoreModel().CreateConfig(config); err != nil {
			// 如果配置保存失败，记录日志但不影响应用创建
			fmt.Printf("failed to save config %s: %v\n", key, err)
		}
	}

	return newApp, nil
}
