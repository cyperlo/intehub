package app

import (
	appModel "intehub/internal/app/models/app"
	"intehub/internal/app/models/appstore"
)

type ConfigInput struct {
	Key       string
	Value     string
	Type      string
	Encrypted bool
}

type Service interface {
	Create(app *appModel.App) error
	CreateWithConfigs(app *appModel.App, configs []ConfigInput) error
	GetByID(id uint) (*appModel.App, error)
	GetByIDWithConfigs(id uint) (*appModel.App, []*appstore.AppConfig, error)
	List(userID uint) ([]*appModel.App, error)
	Update(app *appModel.App) error
	UpdateWithConfigs(app *appModel.App, configs []ConfigInput) error
	Delete(id uint) error
	Run(id uint) (*appModel.AppLog, error)
	RunWithInput(id uint, input map[string]interface{}) (*appModel.AppLog, error)
	GetLogs(appID *uint, page, pageSize int) ([]*appModel.AppLog, int64, error)
	PublishToStore(appID uint, displayName, description, icon, category, version, author, tags, configSchema string) (*appstore.AppTemplate, error)
}
