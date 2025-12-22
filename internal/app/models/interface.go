package models

import (
	"intehub/internal/app/models/app"
	"intehub/internal/app/models/field"
	"intehub/internal/app/models/push"
	"intehub/internal/app/models/schedule"
	"intehub/internal/app/models/system"
	"intehub/internal/app/models/user"
)

type Model interface {
	UserModel() user.Model
	AppModel() app.Model
	FieldModel() field.Model
	PushModel() push.Model
	ScheduleModel() schedule.Model
	SystemModel() system.Model
}
