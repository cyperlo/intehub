package schedule

import scheduleModel "intehub/internal/app/models/schedule"

type service struct {
	model scheduleModel.Model
}

func New(model scheduleModel.Model) Service {
	return &service{model: model}
}

func (s *service) GetTasks(userID uint) ([]*scheduleModel.ScheduleTask, error) {
	return s.model.GetTasks(userID)
}

func (s *service) GetTask(id uint) (*scheduleModel.ScheduleTask, error) {
	return s.model.GetTask(id)
}

func (s *service) CreateTask(task *scheduleModel.ScheduleTask) error {
	return s.model.CreateTask(task)
}

func (s *service) UpdateTask(task *scheduleModel.ScheduleTask) error {
	return s.model.UpdateTask(task)
}

func (s *service) DeleteTask(id uint) error {
	return s.model.DeleteTask(id)
}

func (s *service) ToggleTask(id uint) error {
	return s.model.ToggleTask(id)
}

func (s *service) GetLogs(taskID *uint, page, pageSize int) ([]*scheduleModel.ScheduleLog, int64, error) {
	return s.model.GetLogs(taskID, page, pageSize)
}
