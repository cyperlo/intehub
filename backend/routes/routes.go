package routes

import (
	"intehub/controllers"
	"intehub/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// 认证路由（无需token）
		auth := api.Group("/auth")
		{
			auth.POST("/login", controllers.Login)
			auth.POST("/logout", controllers.Logout)
		}

		// 需要认证的路由
		authorized := api.Group("")
		authorized.Use(middleware.AuthMiddleware())
		{
			// 用户信息
			authorized.GET("/user/current", controllers.GetCurrentUser)

			// 字段定义管理
			fields := authorized.Group("/fields")
			{
				fields.GET("", controllers.GetFieldSchemas)
				fields.GET("/:id", controllers.GetFieldSchema)
				fields.POST("", controllers.CreateFieldSchema)
				fields.PUT("/:id", controllers.UpdateFieldSchema)
				fields.DELETE("/:id", controllers.DeleteFieldSchema)
			}

			// 推送配置管理
			push := authorized.Group("/push")
			{
				push.GET("/configs", controllers.GetPushConfigs)
				push.GET("/configs/:id", controllers.GetPushConfig)
				push.POST("/configs", controllers.CreatePushConfig)
				push.PUT("/configs/:id", controllers.UpdatePushConfig)
				push.DELETE("/configs/:id", controllers.DeletePushConfig)

				// 配置关联的字段
				push.GET("/configs/:id/fields", controllers.GetConfigFields)
				push.PUT("/configs/:id/fields", controllers.UpdateConfigFields)

				// 推送操作
				push.POST("/send", controllers.SendPush)

				// 推送历史
				push.GET("/history", controllers.GetPushHistory)
			}
		}
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
