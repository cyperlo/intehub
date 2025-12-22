package schedule

import "gorm.io/gorm"

type model struct {
	db *gorm.DB
}

func New(db *gorm.DB) Model {
	return &model{db: db}
}

func (m *model) GetTasks(userID uint) ([]*ScheduleTask, error) {
	var tasks []*ScheduleTask
	query := m.db.Order("created_at DESC")
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	err := query.Find(&tasks).Error
	return tasks, err
}

func (m *model) GetTask(id uint) (*ScheduleTask, error) {
	var task ScheduleTask
	err := m.db.First(&task, id).Error
	return &task, err
}

func (m *model) CreateTask(task *ScheduleTask) error {
	return m.db.Create(task).Error
}

func (m *model) UpdateTask(task *ScheduleTask) error {
	return m.db.Save(task).Error
}

func (m *model) DeleteTask(id uint) error {
	return m.db.Delete(&ScheduleTask{}, id).Error
}

func (m *model) ToggleTask(id uint) error {
	var task ScheduleTask
	if err := m.db.First(&task, id).Error; err != nil {
		return err
	}
	task.Enabled = !task.Enabled
	return m.db.Save(&task).Error
}

func (m *model) GetLogs(taskID *uint, page, pageSize int) ([]*ScheduleLog, int64, error) {
	var logs []*ScheduleLog
	var total int64

	query := m.db.Model(&ScheduleLog{})
	if taskID != nil && *taskID > 0 {
		query = query.Where("task_id = ?", *taskID)
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error

	return logs, total, err
}

func (m *model) CreateLog(log *ScheduleLog) error {
	return m.db.Create(log).Error
}
