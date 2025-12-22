package system

import (
	"gorm.io/gorm"
)

type model struct {
	db *gorm.DB
}

func New(db *gorm.DB) Model {
	return &model{db: db}
}

// Menu methods
func (m *model) GetMenus() ([]*Menu, error) {
	var menus []*Menu
	err := m.db.Order("sort ASC").Find(&menus).Error
	return menus, err
}

func (m *model) CreateMenu(menu *Menu) error {
	return m.db.Create(menu).Error
}

func (m *model) UpdateMenu(menu *Menu) error {
	return m.db.Save(menu).Error
}

func (m *model) DeleteMenu(id uint) error {
	return m.db.Delete(&Menu{}, id).Error
}

// Log methods
func (m *model) GetLogs(page, pageSize int) ([]*SystemLog, int64, error) {
	var logs []*SystemLog
	var total int64

	m.db.Model(&SystemLog{}).Count(&total)

	offset := (page - 1) * pageSize
	err := m.db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error

	return logs, total, err
}

func (m *model) CreateLog(log *SystemLog) error {
	return m.db.Create(log).Error
}
