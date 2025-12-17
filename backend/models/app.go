package models

import (
	"time"

	"gorm.io/gorm"
)

// App 应用
type App struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"not null" json:"name"`           // 应用名称
	Description string         `json:"description"`                    // 应用描述
	Code        string         `gorm:"type:text;not null" json:"code"` // 应用代码
	Language    string         `gorm:"default:go" json:"language"`     // 编程语言
	Enabled     bool           `gorm:"default:true" json:"enabled"`    // 是否启用
	UserID      uint           `gorm:"not null" json:"user_id"`        // 创建者ID
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// AppLog 应用执行日志
type AppLog struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	AppID      uint      `gorm:"not null;index" json:"app_id"` // 应用ID
	AppName    string    `json:"app_name"`                     // 应用名称
	Status     string    `gorm:"not null" json:"status"`       // 状态：success, failed
	Output     string    `gorm:"type:text" json:"output"`      // 输出内容
	Error      string    `gorm:"type:text" json:"error"`       // 错误信息
	Duration   int64     `json:"duration"`                     // 执行时长（毫秒）
	StartedAt  time.Time `json:"started_at"`                   // 开始时间
	FinishedAt time.Time `json:"finished_at"`                  // 结束时间
	CreatedAt  time.Time `json:"created_at"`
}
