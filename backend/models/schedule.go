package models

import (
	"time"

	"gorm.io/gorm"
)

// ScheduleTask 定时任务
type ScheduleTask struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"not null" json:"name"`            // 任务名称
	Description string         `json:"description"`                     // 任务描述
	CronExpr    string         `gorm:"not null" json:"cron_expr"`       // Cron表达式
	TaskType    string         `gorm:"not null" json:"task_type"`       // 任务类型：push（推送）
	ConfigID    uint           `gorm:"not null;index" json:"config_id"` // 关联的配置ID
	Enabled     bool           `gorm:"default:false" json:"enabled"`    // 是否启用
	LastRunAt   *time.Time     `json:"last_run_at"`                     // 上次运行时间
	NextRunAt   *time.Time     `json:"next_run_at"`                     // 下次运行时间
	UserID      uint           `gorm:"not null" json:"user_id"`         // 创建者ID
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// ScheduleLog 定时任务执行日志
type ScheduleLog struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	TaskID     uint      `gorm:"not null;index" json:"task_id"` // 任务ID
	TaskName   string    `json:"task_name"`                     // 任务名称
	Status     string    `gorm:"not null" json:"status"`        // 状态：success, failed
	Message    string    `gorm:"type:text" json:"message"`      // 执行信息
	Duration   int64     `json:"duration"`                      // 执行时长（毫秒）
	StartedAt  time.Time `json:"started_at"`                    // 开始时间
	FinishedAt time.Time `json:"finished_at"`                   // 结束时间
	CreatedAt  time.Time `json:"created_at"`
}
