package push

import (
	fieldModel "intehub/internal/app/models/field"

	"gorm.io/gorm"
)

type model struct {
	db *gorm.DB
}

func New(db *gorm.DB) Model {
	return &model{db: db}
}

func (m *model) GetConfigs(userID uint) ([]*PushConfig, error) {
	var configs []*PushConfig
	query := m.db.Order("created_at DESC")
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	err := query.Find(&configs).Error
	return configs, err
}

func (m *model) GetConfig(id uint) (*PushConfig, error) {
	var config PushConfig
	err := m.db.First(&config, id).Error
	return &config, err
}

func (m *model) CreateConfig(config *PushConfig) error {
	return m.db.Create(config).Error
}

func (m *model) UpdateConfig(config *PushConfig) error {
	return m.db.Save(config).Error
}

func (m *model) DeleteConfig(id uint) error {
	return m.db.Delete(&PushConfig{}, id).Error
}

func (m *model) GetHistory(configID *uint, page, pageSize int) ([]*PushHistory, int64, error) {
	var history []*PushHistory
	var total int64

	query := m.db.Model(&PushHistory{})
	if configID != nil && *configID > 0 {
		query = query.Where("config_id = ?", *configID)
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&history).Error

	return history, total, err
}

func (m *model) CreateHistory(history *PushHistory) error {
	return m.db.Create(history).Error
}

func (m *model) GetConfigFields(configID uint) ([]*fieldModel.FieldSchema, error) {
	var fields []*fieldModel.FieldSchema
	err := m.db.Table("field_schemas").
		Joins("JOIN config_field_relations ON field_schemas.id = config_field_relations.field_id").
		Where("config_field_relations.config_id = ?", configID).
		Find(&fields).Error
	return fields, err
}

func (m *model) UpdateConfigFields(configID uint, fieldIDs []uint) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		// 删除旧关联
		if err := tx.Where("config_id = ?", configID).Delete(&ConfigFieldRelation{}).Error; err != nil {
			return err
		}
		// 创建新关联
		for _, fieldID := range fieldIDs {
			if err := tx.Create(&ConfigFieldRelation{
				ConfigID: configID,
				FieldID:  fieldID,
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
