package system

import (
	systemModel "intehub/internal/app/models/system"
	userModel "intehub/internal/app/models/user"
)

type service struct {
	model     systemModel.Model
	userModel userModel.Model
}

func New(model systemModel.Model, userModel userModel.Model) Service {
	return &service{
		model:     model,
		userModel: userModel,
	}
}

// Menu methods
func (s *service) GetMenus() ([]*systemModel.Menu, error) {
	return s.model.GetMenus()
}

func (s *service) CreateMenu(menu *systemModel.Menu) error {
	return s.model.CreateMenu(menu)
}

func (s *service) UpdateMenu(menu *systemModel.Menu) error {
	return s.model.UpdateMenu(menu)
}

func (s *service) DeleteMenu(id uint) error {
	return s.model.DeleteMenu(id)
}

// Log methods
func (s *service) GetLogs(page, pageSize int) ([]*systemModel.SystemLog, int64, error) {
	return s.model.GetLogs(page, pageSize)
}

// User methods
func (s *service) GetUsers() ([]*userModel.DataObject, error) {
	return s.userModel.GetUsers()
}

func (s *service) GetUserByID(id uint) (*userModel.DataObject, error) {
	return s.userModel.GetUserByID(id)
}

func (s *service) CreateUser(user *userModel.DataObject) error {
	return s.userModel.CreateUser(user)
}

func (s *service) UpdateUser(user *userModel.DataObject) error {
	return s.userModel.UpdateUser(user)
}

func (s *service) UpdateUserFields(id uint, fields map[string]interface{}) error {
	return s.userModel.UpdateUserFields(id, fields)
}

func (s *service) DeleteUser(id uint) error {
	return s.userModel.DeleteUser(id)
}
