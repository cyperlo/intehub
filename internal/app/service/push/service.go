package push

import (
	"bytes"
	"encoding/json"
	fieldModel "intehub/internal/app/models/field"
	pushModel "intehub/internal/app/models/push"
	"io"
	"net/http"
	"strings"
	"time"
)

type service struct {
	model pushModel.Model
}

func New(model pushModel.Model) Service {
	return &service{model: model}
}

func (s *service) GetConfigs(userID uint) ([]*pushModel.PushConfig, error) {
	return s.model.GetConfigs(userID)
}

func (s *service) GetConfig(id uint) (*pushModel.PushConfig, error) {
	return s.model.GetConfig(id)
}

func (s *service) CreateConfig(config *pushModel.PushConfig) error {
	return s.model.CreateConfig(config)
}

func (s *service) UpdateConfig(config *pushModel.PushConfig) error {
	return s.model.UpdateConfig(config)
}

func (s *service) DeleteConfig(id uint) error {
	return s.model.DeleteConfig(id)
}

func (s *service) Send(configID uint, data map[string]interface{}) error {
	config, err := s.model.GetConfig(configID)
	if err != nil {
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

	history := &pushModel.PushHistory{
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

	s.model.CreateHistory(history)
	return err
}

func (s *service) GetHistory(configID *uint, page, pageSize int) ([]*pushModel.PushHistory, int64, error) {
	return s.model.GetHistory(configID, page, pageSize)
}

func (s *service) GetConfigFields(configID uint) ([]*fieldModel.FieldSchema, error) {
	return s.model.GetConfigFields(configID)
}

func (s *service) UpdateConfigFields(configID uint, fieldIDs []uint) error {
	return s.model.UpdateConfigFields(configID, fieldIDs)
}
