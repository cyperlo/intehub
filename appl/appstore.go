package appl

import (
	"fmt"
	"intehub/internal/app/models/appstore"
	"log/slog"
)

// AppStore 应用商店管理器
type AppStore struct {
	model appstore.Model
}

func NewAppStore(model appstore.Model) *AppStore {
	return &AppStore{
		model: model,
	}
}

// InstallFromTemplate 从模板安装应用
func (s *AppStore) InstallFromTemplate(templateID uint, appName string, configs map[string]string) (*GoAppIns, error) {
	template, err := s.model.GetTemplateByID(templateID)
	if err != nil {
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	if !template.Enabled {
		return nil, fmt.Errorf("template is disabled")
	}

	// 创建应用实例
	runFunc, err := code2Runnable(template.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to create runnable: %w", err)
	}

	app := &GoAppIns{
		code:    template.Code,
		runFunc: runFunc,
	}

	slog.Info("installed app from template", "template", template.Name, "app", appName)
	return app, nil
}

// GetTemplate 获取模板
func (s *AppStore) GetTemplate(id uint) (*appstore.AppTemplate, error) {
	return s.model.GetTemplateByID(id)
}

// ListTemplates 列出所有模板
func (s *AppStore) ListTemplates(category string) ([]*appstore.AppTemplate, error) {
	return s.model.ListTemplates(category)
}
