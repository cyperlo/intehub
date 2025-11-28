package models

import (
	"intehub/config"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Username  string         `gorm:"uniqueIndex;not null" json:"username"`
	Password  string         `gorm:"not null" json:"-"`
	Nickname  string         `json:"nickname"`
	Role      string         `gorm:"default:user" json:"role"` // admin, user
}

// HashPassword 加密密码
func (u *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

// CheckPassword 验证密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// CreateDefaultAdmin 创建默认管理员账户
func CreateDefaultAdmin() {
	db := config.GetDB()

	var count int64
	db.Model(&User{}).Count(&count)

	if count == 0 {
		admin := User{
			Username: "admin",
			Nickname: "管理员",
			Role:     "admin",
		}
		if err := admin.HashPassword("admin123"); err != nil {
			log.Printf("创建管理员密码失败: %v", err)
			return
		}

		if err := db.Create(&admin).Error; err != nil {
			log.Printf("创建默认管理员失败: %v", err)
		} else {
			log.Println("默认管理员账户创建成功 (admin/admin123)")
		}
	}
}
