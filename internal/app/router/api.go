package router

import "github.com/gin-gonic/gin"

func handleAPI(r *HttpRouter, rg *gin.RouterGroup) {
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

	{
		// 应用管理
		appGroup := rg.Group("/apps")
		r.appHandler.HandleAppAPI(appGroup)

		// 应用商店
		appstoreGroup := rg.Group("/appstore")
		r.appstoreHandler.HandleAppStoreAPI(appstoreGroup)

		// 工作流
		workflowGroup := rg.Group("/workflows")
		r.workflowHandler.HandleWorkflowAPI(workflowGroup)

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
