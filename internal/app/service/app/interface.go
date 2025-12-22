package app

import (
	appModel "intehub/internal/app/models/app"
)

type Service interface {
	Create(app *appModel.App) error
	GetByID(id uint) (*appModel.App, error)
	List(userID uint) ([]*appModel.App, error)
	Update(app *appModel.App) error
	Delete(id uint) error
	Run(id uint) (*appModel.AppLog, error)
	GetLogs(appID *uint, page, pageSize int) ([]*appModel.AppLog, int64, error)
}
