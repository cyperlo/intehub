package system

type Model interface {
	// Menu methods
	GetMenus() ([]*Menu, error)
	CreateMenu(menu *Menu) error
	UpdateMenu(menu *Menu) error
	DeleteMenu(id uint) error

	// Log methods
	GetLogs(page, pageSize int) ([]*SystemLog, int64, error)
	CreateLog(log *SystemLog) error
}
