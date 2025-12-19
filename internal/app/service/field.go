package service

import (
	"intehub/internal/app/models"

	"gorm.io/gorm"
)

type FieldService interface {
	GetAll(userID uint) ([]*models.FieldSchema, error)
	GetByID(id uint) (*models.FieldSchema, error)
	Create(field *models.FieldSchema) error
	Update(field *models.FieldSchema) error
	Delete(id uint) error
}

type fieldService struct {
	db *gorm.DB
}

func NewFieldService(db *gorm.DB) FieldService {
	return &fieldService{db: db}
}

func (s *fieldService) GetAll(userID uint) ([]*models.FieldSchema, error) {
	var fields []*models.FieldSchema
	query := s.db.Order("created_at DESC")
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	err := query.Find(&fields).Error
	return fields, err
}

func (s *fieldService) GetByID(id uint) (*models.FieldSchema, error) {
	var field models.FieldSchema
	err := s.db.First(&field, id).Error
	return &field, err
}

func (s *fieldService) Create(field *models.FieldSchema) error {
	return s.db.Create(field).Error
}

func (s *fieldService) Update(field *models.FieldSchema) error {
	return s.db.Save(field).Error
}

func (s *fieldService) Delete(id uint) error {
	return s.db.Delete(&models.FieldSchema{}, id).Error
}
