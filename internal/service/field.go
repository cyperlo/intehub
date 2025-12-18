package service

import (
	"intehub/internal/model"

	"gorm.io/gorm"
)

type FieldService interface {
	GetAll(userID uint) ([]*model.FieldSchema, error)
	GetByID(id uint) (*model.FieldSchema, error)
	Create(field *model.FieldSchema) error
	Update(field *model.FieldSchema) error
	Delete(id uint) error
}

type fieldService struct {
	db *gorm.DB
}

func NewFieldService(db *gorm.DB) FieldService {
	return &fieldService{db: db}
}

func (s *fieldService) GetAll(userID uint) ([]*model.FieldSchema, error) {
	var fields []*model.FieldSchema
	query := s.db.Order("created_at DESC")
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	err := query.Find(&fields).Error
	return fields, err
}

func (s *fieldService) GetByID(id uint) (*model.FieldSchema, error) {
	var field model.FieldSchema
	err := s.db.First(&field, id).Error
	return &field, err
}

func (s *fieldService) Create(field *model.FieldSchema) error {
	return s.db.Create(field).Error
}

func (s *fieldService) Update(field *model.FieldSchema) error {
	return s.db.Save(field).Error
}

func (s *fieldService) Delete(id uint) error {
	return s.db.Delete(&model.FieldSchema{}, id).Error
}
