package field

import "gorm.io/gorm"

type model struct {
	db *gorm.DB
}

func New(db *gorm.DB) Model {
	return &model{db: db}
}

func (m *model) Create(field *FieldSchema) error {
	return m.db.Create(field).Error
}

func (m *model) GetByID(id uint) (*FieldSchema, error) {
	var field FieldSchema
	err := m.db.First(&field, id).Error
	if err != nil {
		return nil, err
	}
	return &field, nil
}

func (m *model) GetAll(userID uint) ([]*FieldSchema, error) {
	var fields []*FieldSchema
	query := m.db.Order("created_at DESC")
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	err := query.Find(&fields).Error
	return fields, err
}

func (m *model) Update(field *FieldSchema) error {
	return m.db.Save(field).Error
}

func (m *model) Delete(id uint) error {
	return m.db.Delete(&FieldSchema{}, id).Error
}
