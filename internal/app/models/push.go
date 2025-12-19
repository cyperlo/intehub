package models

import (
	"time"

	"gorm.io/gorm"
)

type PushConfig struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description"`
	URL         string         `gorm:"not null" json:"url"`
	Method      string         `gorm:"default:POST" json:"method"`
	Headers     string         `gorm:"type:text" json:"headers"`
	Template    string         `gorm:"type:text" json:"template"`
	Enabled     bool           `gorm:"default:true" json:"enabled"`
	UserID      uint           `json:"user_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type PushHistory struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	ConfigID   uint      `json:"config_id"`
	ConfigName string    `json:"config_name"`
	URL        string    `json:"url"`
	Method     string    `json:"method"`
	Content    string    `gorm:"type:text" json:"content"`
	Response   string    `gorm:"type:text" json:"response"`
	Success    bool      `json:"success"`
	StatusCode int       `json:"status_code"`
	Duration   int64     `json:"duration"`
	UserID     uint      `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type ConfigFieldRelation struct {
	ID        uint `gorm:"primarykey" json:"id"`
	ConfigID  uint `json:"config_id"`
	FieldID   uint `json:"field_id"`
	CreatedAt time.Time
}
