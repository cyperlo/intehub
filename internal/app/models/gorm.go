package models

import (
	"intehub/internal/app/models/app"
	"intehub/internal/app/models/appstore"
	"intehub/internal/app/models/field"
	"intehub/internal/app/models/push"
	"intehub/internal/app/models/schedule"
	"intehub/internal/app/models/system"
	"intehub/internal/app/models/user"
	"intehub/internal/app/models/workflow"

	"gorm.io/gorm"
)

type GormDB struct {
	db *gorm.DB
}

func New(db *gorm.DB) Model {
	return &GormDB{
		db: db,
	}
}

func (g *GormDB) UserModel() user.Model {
	return user.New(g.db)
}

func (g *GormDB) AppModel() app.Model {
	return app.New(g.db)
}

func (g *GormDB) AppStoreModel() appstore.Model {
	return appstore.New(g.db)
}

func (g *GormDB) FieldModel() field.Model {
	return field.New(g.db)
}

func (g *GormDB) PushModel() push.Model {
	return push.New(g.db)
}

func (g *GormDB) ScheduleModel() schedule.Model {
	return schedule.New(g.db)
}

func (g *GormDB) SystemModel() system.Model {
	return system.New(g.db)
}

func (g *GormDB) WorkflowModel() workflow.Model {
	return workflow.New(g.db)
}
