package system

import systemModel "intehub/internal/app/models/system"

type Service interface {
	// Menu methods
	GetMenus() ([]*systemModel.Menu, error)
	CreateMenu(menu *systemModel.Menu) error
	UpdateMenu(menu *systemModel.Menu) error
	DeleteMenu(id uint) error

	// Log methods
	GetLogs(page, pageSize int) ([]*systemModel.SystemLog, int64, error)
}
