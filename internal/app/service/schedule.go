package service

import (
	"intehub/internal/app/model"

	"gorm.io/gorm"
)

type ScheduleService interface {
	GetTasks(userID uint) ([]*model.ScheduleTask, error)
	GetTask(id uint) (*model.ScheduleTask, error)
	CreateTask(task *model.ScheduleTask) error
	UpdateTask(task *model.ScheduleTask) error
	DeleteTask(id uint) error
	ToggleTask(id uint) error
	GetLogs(taskID *uint, page, pageSize int) ([]*model.ScheduleLog, int64, error)
}

type scheduleService struct {
	db *gorm.DB
}

func NewScheduleService(db *gorm.DB) ScheduleService {
	return &scheduleService{db: db}
}

func (s *scheduleService) GetTasks(userID uint) ([]*model.ScheduleTask, error) {
	var tasks []*model.ScheduleTask
	query := s.db.Order("created_at DESC")
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	err := query.Find(&tasks).Error
	return tasks, err
}

func (s *scheduleService) GetTask(id uint) (*model.ScheduleTask, error) {
	var task model.ScheduleTask
	err := s.db.First(&task, id).Error
	return &task, err
}

func (s *scheduleService) CreateTask(task *model.ScheduleTask) error {
	return s.db.Create(task).Error
}

func (s *scheduleService) UpdateTask(task *model.ScheduleTask) error {
	return s.db.Save(task).Error
}

func (s *scheduleService) DeleteTask(id uint) error {
	return s.db.Delete(&model.ScheduleTask{}, id).Error
}

func (s *scheduleService) ToggleTask(id uint) error {
	var task model.ScheduleTask
	if err := s.db.First(&task, id).Error; err != nil {
		return err
	}
	task.Enabled = !task.Enabled
	return s.db.Save(&task).Error
}

func (s *scheduleService) GetLogs(taskID *uint, page, pageSize int) ([]*model.ScheduleLog, int64, error) {
	var logs []*model.ScheduleLog
	var total int64

	query := s.db.Model(&model.ScheduleLog{})
	if taskID != nil && *taskID > 0 {
		query = query.Where("task_id = ?", *taskID)
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error

	return logs, total, err
}
