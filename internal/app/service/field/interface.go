package field

import fieldModel "intehub/internal/app/models/field"

type Service interface {
	Create(field *fieldModel.FieldSchema) error
	GetByID(id uint) (*fieldModel.FieldSchema, error)
	GetAll(userID uint) ([]*fieldModel.FieldSchema, error)
	Update(field *fieldModel.FieldSchema) error
	Delete(id uint) error
}
