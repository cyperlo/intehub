package models

import (
	"time"

	"gorm.io/gorm"
)

type ScheduleTask struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description"`
	CronExpr    string         `gorm:"column:cron_expr;not null" json:"cron_expr"`
	TaskType    string         `gorm:"column:task_type;not null" json:"task_type"` // push, app
	ConfigID    *uint          `gorm:"column:config_id" json:"config_id"`
	AppID       *uint          `gorm:"column:app_id" json:"app_id"`
	FieldData   string         `gorm:"column:field_data;type:text" json:"field_data"`
	Enabled     bool           `gorm:"default:false" json:"enabled"`
	LastRunAt   *time.Time     `gorm:"column:last_run_at" json:"last_run_at"`
	NextRunAt   *time.Time     `gorm:"column:next_run_at" json:"next_run_at"`
	UserID      uint           `gorm:"column:user_id" json:"user_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type ScheduleLog struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	TaskID     uint      `gorm:"column:task_id" json:"task_id"`
	TaskName   string    `gorm:"column:task_name" json:"task_name"`
	Status     string    `json:"status"` // success, error
	Message    string    `gorm:"type:text" json:"message"`
	Duration   int64     `json:"duration"` // 执行时长(毫秒)
	StartedAt  time.Time `gorm:"column:started_at" json:"started_at"`
	FinishedAt time.Time `gorm:"column:finished_at" json:"finished_at"`
	CreatedAt  time.Time `json:"created_at"`
}
