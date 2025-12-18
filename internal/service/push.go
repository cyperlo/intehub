package service

import (
	"bytes"
	"encoding/json"
	"intehub/internal/model"
	"io"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"
)

type PushService interface {
	GetConfigs(userID uint) ([]*model.PushConfig, error)
	GetConfig(id uint) (*model.PushConfig, error)
	CreateConfig(config *model.PushConfig) error
	UpdateConfig(config *model.PushConfig) error
	DeleteConfig(id uint) error
	Send(configID uint, data map[string]interface{}) error
	GetHistory(configID *uint, page, pageSize int) ([]*model.PushHistory, int64, error)
	GetConfigFields(configID uint) ([]*model.FieldSchema, error)
	UpdateConfigFields(configID uint, fieldIDs []uint) error
}

type pushService struct {
	db *gorm.DB
}

func NewPushService(db *gorm.DB) PushService {
	return &pushService{db: db}
}

func (s *pushService) GetConfigs(userID uint) ([]*model.PushConfig, error) {
	var configs []*model.PushConfig
	query := s.db.Order("created_at DESC")
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	err := query.Find(&configs).Error
	return configs, err
}

func (s *pushService) GetConfig(id uint) (*model.PushConfig, error) {
	var config model.PushConfig
	err := s.db.First(&config, id).Error
	return &config, err
}

func (s *pushService) CreateConfig(config *model.PushConfig) error {
	return s.db.Create(config).Error
}

func (s *pushService) UpdateConfig(config *model.PushConfig) error {
	return s.db.Save(config).Error
}

func (s *pushService) DeleteConfig(id uint) error {
	return s.db.Delete(&model.PushConfig{}, id).Error
}

func (s *pushService) Send(configID uint, data map[string]interface{}) error {
	var config model.PushConfig
	if err := s.db.First(&config, configID).Error; err != nil {
		return err
	}

	// 替换模板
	content := config.Template
	for k, v := range data {
		content = strings.ReplaceAll(content, "{{"+k+"}}", v.(string))
	}

	startTime := time.Now()
	req, _ := http.NewRequest(config.Method, config.URL, bytes.NewBufferString(content))
	req.Header.Set("Content-Type", "application/json")

	// 添加自定义 headers
	if config.Headers != "" {
		var headers map[string]string
		json.Unmarshal([]byte(config.Headers), &headers)
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	duration := time.Since(startTime).Milliseconds()

	history := &model.PushHistory{
		ConfigID:   config.ID,
		ConfigName: config.Name,
		URL:        config.URL,
		Method:     config.Method,
		Content:    content,
		Duration:   duration,
	}

	if err != nil {
		history.Success = false
		history.Response = err.Error()
	} else {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		history.Success = resp.StatusCode >= 200 && resp.StatusCode < 300
		history.StatusCode = resp.StatusCode
		history.Response = string(body)
	}

	s.db.Create(history)
	return err
}

func (s *pushService) GetHistory(configID *uint, page, pageSize int) ([]*model.PushHistory, int64, error) {
	var history []*model.PushHistory
	var total int64

	query := s.db.Model(&model.PushHistory{})
	if configID != nil && *configID > 0 {
		query = query.Where("config_id = ?", *configID)
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&history).Error

	return history, total, err
}

func (s *pushService) GetConfigFields(configID uint) ([]*model.FieldSchema, error) {
	var fields []*model.FieldSchema
	err := s.db.Table("field_schemas").
		Joins("JOIN config_field_relations ON field_schemas.id = config_field_relations.field_id").
		Where("config_field_relations.config_id = ?", configID).
		Find(&fields).Error
	return fields, err
}

func (s *pushService) UpdateConfigFields(configID uint, fieldIDs []uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 删除旧关联
		if err := tx.Where("config_id = ?", configID).Delete(&model.ConfigFieldRelation{}).Error; err != nil {
			return err
		}
		// 创建新关联
		for _, fieldID := range fieldIDs {
			if err := tx.Create(&model.ConfigFieldRelation{
				ConfigID: configID,
				FieldID:  fieldID,
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
