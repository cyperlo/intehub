package service

import (
	"intehub/internal/app/model"

	"gorm.io/gorm"
)

type AppService interface {
	Create(app *model.App) error
	GetByID(id uint) (*model.App, error)
	List(userID uint) ([]*model.App, error)
	Update(app *model.App) error
	Delete(id uint) error
	Run(id uint) (*model.AppLog, error)
	GetLogs(appID *uint, page, pageSize int) ([]*model.AppLog, int64, error)
}

type appService struct {
	db *gorm.DB
}

func NewAppService(db *gorm.DB) AppService {
	return &appService{
		db: db,
	}
}

func (s *appService) Create(app *model.App) error {
	return s.db.Create(app).Error
}

func (s *appService) GetByID(id uint) (*model.App, error) {
	var app model.App
	err := s.db.First(&app, id).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (s *appService) List(userID uint) ([]*model.App, error) {
	var apps []*model.App
	query := s.db.Order("created_at DESC")
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	err := query.Find(&apps).Error
	return apps, err
}

func (s *appService) Update(app *model.App) error {
	return s.db.Save(app).Error
}

func (s *appService) Delete(id uint) error {
	return s.db.Delete(&model.App{}, id).Error
}

func (s *appService) Run(id uint) (*model.AppLog, error) {
	return nil, nil
}

func (s *appService) GetLogs(appID *uint, page, pageSize int) ([]*model.AppLog, int64, error) {
	var logs []*model.AppLog
	var total int64

	query := s.db.Model(&model.AppLog{})
	if appID != nil && *appID > 0 {
		query = query.Where("app_id = ?", *appID)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error

	return logs, total, err
}
