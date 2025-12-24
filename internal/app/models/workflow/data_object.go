package workflow

import (
	"time"

	"gorm.io/gorm"
)

// Workflow 工作流定义
type Workflow struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description"`
	Nodes       string         `gorm:"type:text" json:"nodes"` // JSON 格式的节点配置
	Enabled     bool           `gorm:"default:true" json:"enabled"`
	UserID      uint           `json:"user_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// WorkflowNode 工作流节点配置（用于 JSON 序列化）
type WorkflowNode struct {
	ID     string                 `json:"id"`     // 节点唯一ID
	Type   string                 `json:"type"`   // 节点类型: app, condition, transform
	AppID  uint                   `json:"app_id"` // 关联的应用ID
	Config map[string]interface{} `json:"config"` // 节点配置
	Next   []string               `json:"next"`   // 下一个节点ID列表
}

// WorkflowLog 工作流执行日志
type WorkflowLog struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	WorkflowID uint      `json:"workflow_id"`
	Name       string    `json:"name"`
	Status     string    `json:"status"` // running, success, error
	Input      string    `gorm:"type:text" json:"input"`
	Output     string    `gorm:"type:text" json:"output"`
	Error      string    `gorm:"type:text" json:"error"`
	NodeLogs   string    `gorm:"type:text" json:"node_logs"` // JSON 格式的节点执行日志
	Duration   int64     `json:"duration"`                   // 执行时长(毫秒)
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`
	CreatedAt  time.Time `json:"created_at"`
}

// NodeLog 节点执行日志（用于 JSON 序列化）
type NodeLog struct {
	NodeID     string                 `json:"node_id"`
	AppID      uint                   `json:"app_id"`
	AppName    string                 `json:"app_name"`
	Status     string                 `json:"status"` // success, error, skipped
	Input      map[string]interface{} `json:"input"`
	Output     map[string]interface{} `json:"output"`
	Error      string                 `json:"error"`
	Duration   int64                  `json:"duration"`
	StartedAt  time.Time              `json:"started_at"`
	FinishedAt time.Time              `json:"finished_at"`
}
