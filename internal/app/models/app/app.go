package app

import "gorm.io/gorm"

type model struct {
	db *gorm.DB
}

func New(db *gorm.DB) Model {
	return &model{db: db}
}

func (m *model) Create(app *App) error {
	return m.db.Create(app).Error
}

func (m *model) GetByID(id uint) (*App, error) {
	var app App
	err := m.db.First(&app, id).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (m *model) List(userID uint) ([]*App, error) {
	var apps []*App
	query := m.db.Where("enabled = ?", true).Order("created_at DESC")
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	err := query.Find(&apps).Error
	return apps, err
}

func (m *model) Update(app *App) error {
	return m.db.Save(app).Error
}

func (m *model) Delete(id uint) error {
	return m.db.Delete(&App{}, id).Error
}

func (m *model) GetLogs(appID *uint, page, pageSize int) ([]*AppLog, int64, error) {
	var logs []*AppLog
	var total int64

	query := m.db.Model(&AppLog{})
	if appID != nil && *appID > 0 {
		query = query.Where("app_id = ?", *appID)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error

	return logs, total, err
}

func (m *model) CreateLog(log *AppLog) error {
	return m.db.Create(log).Error
}
