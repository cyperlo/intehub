package schedule

import scheduleModel "intehub/internal/app/models/schedule"

type Service interface {
	GetTasks(userID uint) ([]*scheduleModel.ScheduleTask, error)
	GetTask(id uint) (*scheduleModel.ScheduleTask, error)
	CreateTask(task *scheduleModel.ScheduleTask) error
	UpdateTask(task *scheduleModel.ScheduleTask) error
	DeleteTask(id uint) error
	ToggleTask(id uint) error
	GetLogs(taskID *uint, page, pageSize int) ([]*scheduleModel.ScheduleLog, int64, error)
	ExecuteTask(task *scheduleModel.ScheduleTask) error
	StartScheduler() error
	StopScheduler()
}
