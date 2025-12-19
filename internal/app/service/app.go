package service

import (
	"intehub/internal/app/models"

	"gorm.io/gorm"
)

type AppService interface {
	Create(app *models.App) error
	GetByID(id uint) (*models.App, error)
	List(userID uint) ([]*models.App, error)
	Update(app *models.App) error
	Delete(id uint) error
	Run(id uint) (*models.AppLog, error)
	GetLogs(appID *uint, page, pageSize int) ([]*models.AppLog, int64, error)
}

type appService struct {
	db *gorm.DB
}

func NewAppService(db *gorm.DB) AppService {
	return &appService{
		db: db,
	}
}

func (s *appService) Create(app *models.App) error {
	return s.db.Create(app).Error
}

func (s *appService) GetByID(id uint) (*models.App, error) {
	var app models.App
	err := s.db.First(&app, id).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (s *appService) List(userID uint) ([]*models.App, error) {
	var apps []*models.App
	query := s.db.Order("created_at DESC")
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	err := query.Find(&apps).Error
	return apps, err
}

func (s *appService) Update(app *models.App) error {
	return s.db.Save(app).Error
}

func (s *appService) Delete(id uint) error {
	return s.db.Delete(&models.App{}, id).Error
}

func (s *appService) Run(id uint) (*models.AppLog, error) {
	return nil, nil
}

func (s *appService) GetLogs(appID *uint, page, pageSize int) ([]*models.AppLog, int64, error) {
	var logs []*models.AppLog
	var total int64

	query := s.db.Model(&models.AppLog{})
	if appID != nil && *appID > 0 {
		query = query.Where("app_id = ?", *appID)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error

	return logs, total, err
}
