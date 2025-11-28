package models

import (
	"time"

	"gorm.io/gorm"
)

// FieldType 字段类型
type FieldType string

const (
	FieldTypeText     FieldType = "text"     // 单行文本
	FieldTypeTextarea FieldType = "textarea" // 多行文本
	FieldTypeNumber   FieldType = "number"   // 数字
	FieldTypeDate     FieldType = "date"     // 日期
	FieldTypeSelect   FieldType = "select"   // 下拉选择
	FieldTypeURL      FieldType = "url"      // URL
	FieldTypeEmail    FieldType = "email"    // 邮箱
)

// FieldSchema 字段定义
type FieldSchema struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Name         string         `gorm:"not null" json:"name"`              // 字段名称
	Key          string         `gorm:"uniqueIndex;not null" json:"key"`   // 字段key（用于模板变量）
	Type         FieldType      `gorm:"not null;default:text" json:"type"` // 字段类型
	Description  string         `json:"description"`                       // 字段描述
	Required     bool           `gorm:"default:false" json:"required"`     // 是否必填
	DefaultValue string         `json:"default_value"`                     // 默认值
	Options      string         `gorm:"type:text" json:"options"`          // 选项（JSON格式，用于select类型）
	Placeholder  string         `json:"placeholder"`                       // 占位符
	Validation   string         `gorm:"type:text" json:"validation"`       // 验证规则（JSON格式）
	UserID       uint           `gorm:"not null" json:"user_id"`           // 创建者ID
}

// ConfigFieldRelation 配置-字段关联表
type ConfigFieldRelation struct {
	ID       uint      `gorm:"primarykey" json:"id"`
	ConfigID uint      `gorm:"not null;index" json:"config_id"` // 配置ID
	FieldID  uint      `gorm:"not null;index" json:"field_id"`  // 字段ID
	Order    int       `gorm:"default:0" json:"order"`          // 显示顺序
	CreateAt time.Time `json:"created_at"`
}
