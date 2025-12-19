package service

import (
	"intehub/internal/app/models"

	"gorm.io/gorm"
)

type SystemService interface {
	// User
	GetUsers() ([]*models.User, error)
	GetUser(id uint) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	DeleteUser(id uint) error

	// Menu
	GetMenus() ([]*models.Menu, error)
	CreateMenu(menu *models.Menu) error
	UpdateMenu(menu *models.Menu) error
	DeleteMenu(id uint) error

	// Log
	GetLogs(page, pageSize int) ([]*models.SystemLog, int64, error)
	CreateLog(log *models.SystemLog) error
}

type systemService struct {
	db *gorm.DB
}

func NewSystemService(db *gorm.DB) SystemService {
	return &systemService{db: db}
}

func (s *systemService) GetUsers() ([]*models.User, error) {
	var users []*models.User
	err := s.db.Find(&users).Error
	return users, err
}

func (s *systemService) GetUser(id uint) (*models.User, error) {
	var user models.User
	err := s.db.First(&user, id).Error
	return &user, err
}

func (s *systemService) CreateUser(user *models.User) error {
	return s.db.Create(user).Error
}

func (s *systemService) UpdateUser(user *models.User) error {
	return s.db.Save(user).Error
}

func (s *systemService) DeleteUser(id uint) error {
	return s.db.Delete(&models.User{}, id).Error
}

func (s *systemService) GetMenus() ([]*models.Menu, error) {
	var menus []*models.Menu
	err := s.db.Order("sort ASC").Find(&menus).Error
	return menus, err
}

func (s *systemService) CreateMenu(menu *models.Menu) error {
	return s.db.Create(menu).Error
}

func (s *systemService) UpdateMenu(menu *models.Menu) error {
	return s.db.Save(menu).Error
}

func (s *systemService) DeleteMenu(id uint) error {
	return s.db.Delete(&models.Menu{}, id).Error
}

func (s *systemService) GetLogs(page, pageSize int) ([]*models.SystemLog, int64, error) {
	var logs []*models.SystemLog
	var total int64

	s.db.Model(&models.SystemLog{}).Count(&total)
	offset := (page - 1) * pageSize
	err := s.db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error

	return logs, total, err
}

func (s *systemService) CreateLog(log *models.SystemLog) error {
	return s.db.Create(log).Error
}
