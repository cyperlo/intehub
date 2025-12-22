package app

type Model interface {
	Create(app *App) error
	GetByID(id uint) (*App, error)
	List(userID uint) ([]*App, error)
	Update(app *App) error
	Delete(id uint) error
	GetLogs(appID *uint, page, pageSize int) ([]*AppLog, int64, error)
	CreateLog(log *AppLog) error
}
