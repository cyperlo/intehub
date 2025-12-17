package config

import (
	"log"
	"os"
	"time"

	// "gorm.io/driver/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() {
	var err error

	// 从环境变量读取数据库路径，默认使用 intehub.db
	dsn := os.Getenv("dsn")
	if dsn == "" {
		dsn = "intehub.db"
	}

	// 添加时区参数
	if dsn != "intehub.db" && !contains(dsn, "TimeZone") {
		if contains(dsn, "?") {
			dsn += "&TimeZone=Asia/Shanghai"
		} else {
			dsn += "?TimeZone=Asia/Shanghai"
		}
	}

	// 使用PostgreSQL数据库
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			// 使用本地时区
			loc, _ := time.LoadLocation("Asia/Shanghai")
			return time.Now().In(loc)
		},
	})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	log.Printf("数据库连接成功: %s", dsn)
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
