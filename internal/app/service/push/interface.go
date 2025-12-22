package push

import (
	fieldModel "intehub/internal/app/models/field"
	pushModel "intehub/internal/app/models/push"
)

type Service interface {
	GetConfigs(userID uint) ([]*pushModel.PushConfig, error)
	GetConfig(id uint) (*pushModel.PushConfig, error)
	CreateConfig(config *pushModel.PushConfig) error
	UpdateConfig(config *pushModel.PushConfig) error
	DeleteConfig(id uint) error
	Send(configID uint, data map[string]interface{}) error
	GetHistory(configID *uint, page, pageSize int) ([]*pushModel.PushHistory, int64, error)
	GetConfigFields(configID uint) ([]*fieldModel.FieldSchema, error)
	UpdateConfigFields(configID uint, fieldIDs []uint) error
}
