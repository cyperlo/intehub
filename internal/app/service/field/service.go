package field

import fieldModel "intehub/internal/app/models/field"

type service struct {
	model fieldModel.Model
}

func New(model fieldModel.Model) Service {
	return &service{model: model}
}

func (s *service) Create(field *fieldModel.FieldSchema) error {
	return s.model.Create(field)
}

func (s *service) GetByID(id uint) (*fieldModel.FieldSchema, error) {
	return s.model.GetByID(id)
}

func (s *service) GetAll(userID uint) ([]*fieldModel.FieldSchema, error) {
	return s.model.GetAll(userID)
}

func (s *service) Update(field *fieldModel.FieldSchema) error {
	return s.model.Update(field)
}

func (s *service) Delete(id uint) error {
	return s.model.Delete(id)
}
