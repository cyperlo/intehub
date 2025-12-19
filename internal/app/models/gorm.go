package models

import (
	"intehub/internal/app/models/user"

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
