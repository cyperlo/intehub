package appstore

import (
	"time"

	"gorm.io/gorm"
)

// AppTemplate 应用商店中的应用模板
type AppTemplate struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	Name         string         `gorm:"not null;uniqueIndex" json:"name"`
	DisplayName  string         `gorm:"not null" json:"display_name"`
	Description  string         `json:"description"`
	Icon         string         `json:"icon"`
	Code         string         `gorm:"type:text;not null" json:"code"`
	Language     string         `gorm:"default:go" json:"language"`
	Category     string         `json:"category"` // 分类：工具、数据处理、通知等
	Version      string         `json:"version"`
	Author       string         `json:"author"`
	Tags         string         `json:"tags"`                           // 标签，逗号分隔
	ConfigSchema string         `gorm:"type:text" json:"config_schema"` // JSON Schema 配置模板
	Enabled      bool           `gorm:"default:true" json:"enabled"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// AppConfig 应用配置
type AppConfig struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	AppID     uint           `gorm:"not null;index" json:"app_id"`
	Key       string         `gorm:"not null" json:"key"`
	Value     string         `gorm:"type:text" json:"value"`
	Type      string         `gorm:"default:string" json:"type"` // string, number, boolean, json
	Encrypted bool           `gorm:"default:false" json:"encrypted"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
