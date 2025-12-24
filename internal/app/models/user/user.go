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

func (m *model) GetUserByID(id uint) (*DataObject, error) {
	var user DataObject
	err := m.db.First(&user, id).Error
	if err != nil {
		return nil, errors.New("user is not found")
	}
	return &user, nil
}

func (m *model) GetUsers() ([]*DataObject, error) {
	var users []*DataObject
	err := m.db.Find(&users).Error
	return users, err
}

func (m *model) CreateUser(user *DataObject) error {
	return m.db.Create(user).Error
}

func (m *model) UpdateUser(user *DataObject) error {
	return m.db.Save(user).Error
}

func (m *model) UpdateUserFields(id uint, fields map[string]interface{}) error {
	return m.db.Model(&DataObject{}).Where("id = ?", id).Updates(fields).Error
}

func (m *model) DeleteUser(id uint) error {
	return m.db.Delete(&DataObject{}, id).Error
}
