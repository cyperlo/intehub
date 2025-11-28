package models

import (
	"time"

	"gorm.io/gorm"
)

// PushConfig 推送配置
type PushConfig struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"not null" json:"name"`                         // 配置名称
	Description string         `json:"description"`                                  // 描述
	URL         string         `gorm:"not null" json:"url"`                          // 推送URL
	Method      string         `gorm:"not null;default:POST" json:"method"`          // 推送方法: GET, POST, PUT
	Headers     string         `gorm:"type:text" json:"headers"`                     // 请求头 (JSON格式)
	ContentType string         `gorm:"default:application/json" json:"content_type"` // 内容类型
	Template    string         `gorm:"type:text" json:"template"`                    // 推送内容模板
	Enabled     bool           `gorm:"default:true" json:"enabled"`                  // 是否启用
	UserID      uint           `gorm:"not null" json:"user_id"`                      // 创建者ID
}

// PushHistory 推送历史
type PushHistory struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	ConfigID   uint      `gorm:"not null" json:"config_id"` // 配置ID
	ConfigName string    `json:"config_name"`               // 配置名称
	URL        string    `gorm:"not null" json:"url"`       // 推送URL
	Method     string    `gorm:"not null" json:"method"`    // 推送方法
	Content    string    `gorm:"type:text" json:"content"`  // 推送内容
	StatusCode int       `json:"status_code"`               // 响应状态码
	Response   string    `gorm:"type:text" json:"response"` // 响应内容
	Success    bool      `json:"success"`                   // 是否成功
	Error      string    `gorm:"type:text" json:"error"`    // 错误信息
	Duration   int64     `json:"duration"`                  // 耗时（毫秒）
	UserID     uint      `gorm:"not null" json:"user_id"`   // 执行者ID
}
