package app

import (
	appModel "intehub/internal/app/models/app"
)

type service struct {
	model appModel.Model
}

func New(model appModel.Model) Service {
	return &service{model: model}
}

func (s *service) Create(app *appModel.App) error {
	return s.model.Create(app)
}

func (s *service) GetByID(id uint) (*appModel.App, error) {
	return s.model.GetByID(id)
}

func (s *service) List(userID uint) ([]*appModel.App, error) {
	return s.model.List(userID)
}

func (s *service) Update(app *appModel.App) error {
	return s.model.Update(app)
}

func (s *service) Delete(id uint) error {
	return s.model.Delete(id)
}

func (s *service) Run(id uint) (*appModel.AppLog, error) {
	// TODO: 实现应用运行逻辑
	return nil, nil
}

func (s *service) GetLogs(appID *uint, page, pageSize int) ([]*appModel.AppLog, int64, error) {
	return s.model.GetLogs(appID, page, pageSize)
}
