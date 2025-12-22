package push

import fieldModel "intehub/internal/app/models/field"

type Model interface {
	GetConfigs(userID uint) ([]*PushConfig, error)
	GetConfig(id uint) (*PushConfig, error)
	CreateConfig(config *PushConfig) error
	UpdateConfig(config *PushConfig) error
	DeleteConfig(id uint) error
	GetHistory(configID *uint, page, pageSize int) ([]*PushHistory, int64, error)
	CreateHistory(history *PushHistory) error
	GetConfigFields(configID uint) ([]*fieldModel.FieldSchema, error)
	UpdateConfigFields(configID uint, fieldIDs []uint) error
}
