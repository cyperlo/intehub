package server

import (
	"intehub/internal/app/api/v1/app"
	"intehub/internal/app/api/v1/auth"
	"intehub/internal/app/api/v1/field"
	"intehub/internal/app/api/v1/push"
	"intehub/internal/app/api/v1/schedule"
	"intehub/internal/app/api/v1/system"
	"intehub/internal/app/config"

	"gorm.io/gorm"
)

type AppContext struct {
	Config *config.Config
	DB     *gorm.DB

	// Handlers
	AuthHandler     *auth.Handler
	AppHandler      *app.Handler
	PushHandler     *push.Handler
	FieldHandler    *field.Handler
	SystemHandler   *system.Handler
	ScheduleHandler *schedule.Handler
}
