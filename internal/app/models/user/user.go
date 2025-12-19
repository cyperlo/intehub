package user

import (
	"github.com/pkg/errors"

	"gorm.io/gorm"
)

type model struct {
	db *gorm.DB
}

func New(db *gorm.DB) Model {
	return &model{
		db: db,
	}
}

func (m *model) GetUserByUsername(username string) (*DataObject, error) {
	var user DataObject
	err := m.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, errors.New("user is not found")
	}

	return &user, nil
}
