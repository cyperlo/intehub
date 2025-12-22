package field

type Model interface {
	Create(field *FieldSchema) error
	GetByID(id uint) (*FieldSchema, error)
	GetAll(userID uint) ([]*FieldSchema, error)
	Update(field *FieldSchema) error
	Delete(id uint) error
}
