package models

import (
	"time"

	"gorm.io/gorm"
)

// SystemLog 系统日志
type SystemLog struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Level     string         `gorm:"not null;index" json:"level"`  // 日志级别：info, warning, error
	Module    string         `gorm:"not null;index" json:"module"` // 模块名称
	Action    string         `gorm:"not null" json:"action"`       // 操作
	Content   string         `gorm:"type:text" json:"content"`     // 日志内容
	UserID    uint           `json:"user_id"`                      // 操作用户ID
	IP        string         `json:"ip"`                           // IP地址
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Menu 菜单配置
type Menu struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"not null" json:"name"`             // 菜单名称
	Path      string         `gorm:"not null" json:"path"`             // 路由路径
	Icon      string         `json:"icon"`                             // 图标
	ParentID  uint           `gorm:"default:0;index" json:"parent_id"` // 父菜单ID
	Sort      int            `gorm:"default:0" json:"sort"`            // 排序
	Visible   bool           `gorm:"default:true" json:"visible"`      // 是否可见
	Roles     string         `gorm:"type:text" json:"roles"`           // 允许访问的角色（JSON数组）
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
