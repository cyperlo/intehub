package router

import "github.com/gin-gonic/gin"

func handleAPI(r *HttpRouter, rg *gin.RouterGroup, lrMw gin.HandlerFunc,
	expireMw gin.HandlerFunc, authKeyMw gin.HandlerFunc, sourceTypeMw gin.HandlerFunc,
) {
	{
		// 健康检查
		healthGroup := rg.Group("/health")
		r.healthHandler.HandlerHealthAPI(healthGroup)
	}

	{
		// 认证模块（无需鉴权）
		authGroup := rg.Group("/auth")
		r.authHandler.HandleAuthAPI(authGroup)
	}

	// 需要鉴权的路由
	// TODO: 添加认证中间件
	// authorized := rg.Group("")
	// authorized.Use(middleware.AuthMiddleware())
	{
		// 应用管理
		appGroup := rg.Group("/apps")
		r.appHandler.HandleAppAPI(appGroup)

		// 字段管理
		fieldGroup := rg.Group("/fields")
		r.fieldHandler.HandleFieldAPI(fieldGroup)

		// 推送管理
		pushGroup := rg.Group("/push")
		r.pushHandler.HandlePushAPI(pushGroup)

		// 定时任务
		scheduleGroup := rg.Group("/schedule")
		r.scheduleHandler.HandleScheduleAPI(scheduleGroup)

		// 系统管理
		systemGroup := rg.Group("/system")
		r.systemHandler.HandleSystemAPI(systemGroup)
	}
}
