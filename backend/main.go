package main

import (
	"intehub/config"
	"intehub/models"
	"intehub/routes"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	config.InitDB()

	// 自动迁移数据库
	if err := config.DB.AutoMigrate(
		&models.User{},
		&models.PushConfig{},
		&models.PushHistory{},
		&models.FieldSchema{},
		&models.ConfigFieldRelation{},
	); err != nil {
		log.Fatal("数据库迁移失败:", err)
	}

	// 创建默认管理员账户
	models.CreateDefaultAdmin()

	// 创建Gin引擎
	r := gin.Default()

	// 配置CORS
	// 跨域配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 注册路由
	routes.SetupRoutes(r)

	// 启动服务器
	log.Println("服务器启动在 http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}
