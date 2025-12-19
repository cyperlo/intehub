package service

import (
	"intehub/internal/app/models"

	"gorm.io/gorm"
)

type ScheduleService interface {
	GetTasks(userID uint) ([]*models.ScheduleTask, error)
	GetTask(id uint) (*models.ScheduleTask, error)
	CreateTask(task *models.ScheduleTask) error
	UpdateTask(task *models.ScheduleTask) error
	DeleteTask(id uint) error
	ToggleTask(id uint) error
	GetLogs(taskID *uint, page, pageSize int) ([]*models.ScheduleLog, int64, error)
}

type scheduleService struct {
	db *gorm.DB
}

func NewScheduleService(db *gorm.DB) ScheduleService {
	return &scheduleService{db: db}
}

func (s *scheduleService) GetTasks(userID uint) ([]*models.ScheduleTask, error) {
	var tasks []*models.ScheduleTask
	query := s.db.Order("created_at DESC")
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	err := query.Find(&tasks).Error
	return tasks, err
}

func (s *scheduleService) GetTask(id uint) (*models.ScheduleTask, error) {
	var task models.ScheduleTask
	err := s.db.First(&task, id).Error
	return &task, err
}

func (s *scheduleService) CreateTask(task *models.ScheduleTask) error {
	return s.db.Create(task).Error
}

func (s *scheduleService) UpdateTask(task *models.ScheduleTask) error {
	return s.db.Save(task).Error
}

func (s *scheduleService) DeleteTask(id uint) error {
	return s.db.Delete(&models.ScheduleTask{}, id).Error
}

func (s *scheduleService) ToggleTask(id uint) error {
	var task models.ScheduleTask
	if err := s.db.First(&task, id).Error; err != nil {
		return err
	}
	task.Enabled = !task.Enabled
	return s.db.Save(&task).Error
}

func (s *scheduleService) GetLogs(taskID *uint, page, pageSize int) ([]*models.ScheduleLog, int64, error) {
	var logs []*models.ScheduleLog
	var total int64

	query := s.db.Model(&models.ScheduleLog{})
	if taskID != nil && *taskID > 0 {
		query = query.Where("task_id = ?", *taskID)
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error

	return logs, total, err
}
