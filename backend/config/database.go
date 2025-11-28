package config

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() {
	var err error

	// 从环境变量读取数据库路径，默认使用 intehub.db
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "intehub.db"
	}

	// 使用SQLite数据库，生产环境可以切换到MySQL或PostgreSQL
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	log.Printf("数据库连接成功: %s", dbPath)
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
