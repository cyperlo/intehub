package appstore

import "gorm.io/gorm"

type model struct {
	db *gorm.DB
}

func New(db *gorm.DB) Model {
	return &model{db: db}
}

// AppTemplate methods
func (m *model) CreateTemplate(template *AppTemplate) error {
	return m.db.Create(template).Error
}

func (m *model) GetTemplateByID(id uint) (*AppTemplate, error) {
	var template AppTemplate
	err := m.db.First(&template, id).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (m *model) GetTemplateByName(name string) (*AppTemplate, error) {
	var template AppTemplate
	err := m.db.Where("name = ?", name).First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (m *model) ListTemplates(category string) ([]*AppTemplate, error) {
	var templates []*AppTemplate
	query := m.db.Where("enabled = ?", true).Order("created_at DESC")
	if category != "" {
		query = query.Where("category = ?", category)
	}
	err := query.Find(&templates).Error
	return templates, err
}

func (m *model) UpdateTemplate(template *AppTemplate) error {
	return m.db.Save(template).Error
}

func (m *model) DeleteTemplate(id uint) error {
	return m.db.Delete(&AppTemplate{}, id).Error
}

// AppConfig methods
func (m *model) CreateConfig(config *AppConfig) error {
	return m.db.Create(config).Error
}

func (m *model) GetConfig(appID uint, key string) (*AppConfig, error) {
	var config AppConfig
	err := m.db.Where("app_id = ? AND key = ?", appID, key).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (m *model) GetConfigsByAppID(appID uint) ([]*AppConfig, error) {
	var configs []*AppConfig
	err := m.db.Where("app_id = ?", appID).Find(&configs).Error
	return configs, err
}

func (m *model) UpdateConfig(config *AppConfig) error {
	return m.db.Save(config).Error
}

func (m *model) DeleteConfig(id uint) error {
	return m.db.Delete(&AppConfig{}, id).Error
}

func (m *model) DeleteConfigsByAppID(appID uint) error {
	return m.db.Where("app_id = ?", appID).Delete(&AppConfig{}).Error
}
