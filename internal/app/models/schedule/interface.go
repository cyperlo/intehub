package schedule

type Model interface {
	GetTasks(userID uint) ([]*ScheduleTask, error)
	GetTask(id uint) (*ScheduleTask, error)
	CreateTask(task *ScheduleTask) error
	UpdateTask(task *ScheduleTask) error
	DeleteTask(id uint) error
	ToggleTask(id uint) error
	GetLogs(taskID *uint, page, pageSize int) ([]*ScheduleLog, int64, error)
	CreateLog(log *ScheduleLog) error
}
