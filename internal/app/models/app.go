package models

import (
	"time"

	"gorm.io/gorm"
)

type App struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description"`
	Code        string         `gorm:"type:text" json:"code"`
	Language    string         `gorm:"default:go" json:"language"`
	Enabled     bool           `gorm:"default:true" json:"enabled"`
	UserID      uint           `json:"user_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type AppLog struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	AppID      uint      `json:"app_id"`
	AppName    string    `json:"app_name"`
	Status     string    `json:"status"` // success, error
	Output     string    `gorm:"type:text" json:"output"`
	Error      string    `gorm:"type:text" json:"error"`
	Duration   int64     `json:"duration"` // 执行时长(毫秒)
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`
	CreatedAt  time.Time `json:"created_at"`
}
