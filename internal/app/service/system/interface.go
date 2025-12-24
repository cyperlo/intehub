package system

import (
	systemModel "intehub/internal/app/models/system"
	userModel "intehub/internal/app/models/user"
)

type Service interface {
	// Menu methods
	GetMenus() ([]*systemModel.Menu, error)
	CreateMenu(menu *systemModel.Menu) error
	UpdateMenu(menu *systemModel.Menu) error
	DeleteMenu(id uint) error

	// Log methods
	GetLogs(page, pageSize int) ([]*systemModel.SystemLog, int64, error)

	// User methods
	GetUsers() ([]*userModel.DataObject, error)
	GetUserByID(id uint) (*userModel.DataObject, error)
	CreateUser(user *userModel.DataObject) error
	UpdateUser(user *userModel.DataObject) error
	UpdateUserFields(id uint, fields map[string]interface{}) error
	DeleteUser(id uint) error
}
