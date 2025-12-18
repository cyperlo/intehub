package service

import (
	"intehub/internal/app/model"

	"gorm.io/gorm"
)

type SystemService interface {
	// User
	GetUsers() ([]*model.User, error)
	GetUser(id uint) (*model.User, error)
	CreateUser(user *model.User) error
	UpdateUser(user *model.User) error
	DeleteUser(id uint) error

	// Menu
	GetMenus() ([]*model.Menu, error)
	CreateMenu(menu *model.Menu) error
	UpdateMenu(menu *model.Menu) error
	DeleteMenu(id uint) error

	// Log
	GetLogs(page, pageSize int) ([]*model.SystemLog, int64, error)
	CreateLog(log *model.SystemLog) error
}

type systemService struct {
	db *gorm.DB
}

func NewSystemService(db *gorm.DB) SystemService {
	return &systemService{db: db}
}

func (s *systemService) GetUsers() ([]*model.User, error) {
	var users []*model.User
	err := s.db.Find(&users).Error
	return users, err
}

func (s *systemService) GetUser(id uint) (*model.User, error) {
	var user model.User
	err := s.db.First(&user, id).Error
	return &user, err
}

func (s *systemService) CreateUser(user *model.User) error {
	return s.db.Create(user).Error
}

func (s *systemService) UpdateUser(user *model.User) error {
	return s.db.Save(user).Error
}

func (s *systemService) DeleteUser(id uint) error {
	return s.db.Delete(&model.User{}, id).Error
}

func (s *systemService) GetMenus() ([]*model.Menu, error) {
	var menus []*model.Menu
	err := s.db.Order("sort ASC").Find(&menus).Error
	return menus, err
}

func (s *systemService) CreateMenu(menu *model.Menu) error {
	return s.db.Create(menu).Error
}

func (s *systemService) UpdateMenu(menu *model.Menu) error {
	return s.db.Save(menu).Error
}

func (s *systemService) DeleteMenu(id uint) error {
	return s.db.Delete(&model.Menu{}, id).Error
}

func (s *systemService) GetLogs(page, pageSize int) ([]*model.SystemLog, int64, error) {
	var logs []*model.SystemLog
	var total int64

	s.db.Model(&model.SystemLog{}).Count(&total)
	offset := (page - 1) * pageSize
	err := s.db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error

	return logs, total, err
}

func (s *systemService) CreateLog(log *model.SystemLog) error {
	return s.db.Create(log).Error
}
