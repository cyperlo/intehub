package model

import (
	"time"

	"gorm.io/gorm"
)

type FieldSchema struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	Name         string         `gorm:"not null" json:"name"`
	Key          string         `gorm:"not null;uniqueIndex" json:"key"`
	Type         string         `gorm:"not null" json:"type"`
	Description  string         `json:"description"`
	Required     bool           `gorm:"default:false" json:"required"`
	DefaultValue string         `json:"default_value"`
	Options      string         `gorm:"type:text" json:"options"`
	Placeholder  string         `json:"placeholder"`
	Validation   string         `json:"validation"`
	UserID       uint           `json:"user_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
