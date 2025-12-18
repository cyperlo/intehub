package model

import "time"

type SystemLog struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Module    string    `json:"module"`
	Action    string    `json:"action"`
	Content   string    `gorm:"type:text" json:"content"`
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	IP        string    `json:"ip"`
	CreatedAt time.Time `json:"created_at"`
}

type Menu struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Path      string    `json:"path"`
	Icon      string    `json:"icon"`
	ParentID  *uint     `json:"parent_id"`
	Sort      int       `gorm:"default:0" json:"sort"`
	Visible   bool      `gorm:"default:true" json:"visible"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
