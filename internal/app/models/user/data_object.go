package user

import (
	"intehub/internal/utils/database"

	"golang.org/x/crypto/bcrypt"
)

type DataObject struct {
	database.BaseModel
	Username string `gorm:"uniqueIndex;not null" json:"username"`
	Nickname string `json:"nickname"`
	Password string `json:"-"`
	Role     string `gorm:"default:user" json:"role"` // admin, user
}

func (u *DataObject) TableName() string {
	return "user"
}

func (u *DataObject) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
