package app

import (
	"fmt"
	"intehub/internal/app/context"
	r "intehub/internal/app/router"
	"log/slog"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	engine          *gin.Engine
	appCtx          *context.AppContext
	httpRouter      *r.HttpRouter
	scheduleService interface{ StartScheduler() error }
}

func NewServerWire(appContext *context.AppContext, httpRouter *r.HttpRouter, scheduleService interface{ StartScheduler() error }) *Server {
	if !appContext.Config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// CORS 配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	s := &Server{
		engine:          r,
		appCtx:          appContext,
		httpRouter:      httpRouter,
		scheduleService: scheduleService,
	}

	// s.setupRoutes()
	return s
}

// func (s *Server) setupRoutes() {
// 	api := s.engine.Group("/api")
// 	authMiddleware := middleware.AuthMiddleware(s.appCtx.Config.JWT.Secret)
// 	adminMiddleware := middleware.AdminMiddleware()

// 	// 认证路由
// 	authGroup := api.Group("/auth")
// 	{
// 		authGroup.POST("/login", s.appCtx.AuthHandler.Login)
// 		authGroup.POST("/logout", s.appCtx.AuthHandler.Logout)
// 		authGroup.GET("/profile", authMiddleware, s.appCtx.AuthHandler.GetProfile)
// 	}

// 	// 应用路由
// 	appGroup := api.Group("/apps", authMiddleware)
// 	{
// 		appGroup.GET("", s.appCtx.AppHandler.List)
// 		appGroup.GET("/:id", s.appCtx.AppHandler.Get)
// 		appGroup.POST("", s.appCtx.AppHandler.Create)
// 		appGroup.PUT("/:id", s.appCtx.AppHandler.Update)
// 		appGroup.DELETE("/:id", s.appCtx.AppHandler.Delete)
// 		appGroup.POST("/:id/run", s.appCtx.AppHandler.Run)
// 		appGroup.GET("/logs", s.appCtx.AppHandler.GetLogs)
// 	}

// 	// 集成配置路由
// 	pushGroup := api.Group("/push", authMiddleware)
// 	{
// 		pushGroup.GET("/configs", s.appCtx.PushHandler.ListConfigs)
// 		pushGroup.GET("/configs/:id", s.appCtx.PushHandler.GetConfig)
// 		pushGroup.POST("/configs", s.appCtx.PushHandler.CreateConfig)
// 		pushGroup.PUT("/configs/:id", s.appCtx.PushHandler.UpdateConfig)
// 		pushGroup.DELETE("/configs/:id", s.appCtx.PushHandler.DeleteConfig)
// 		pushGroup.POST("/send", s.appCtx.PushHandler.Send)
// 		pushGroup.GET("/history", s.appCtx.PushHandler.GetHistory)
// 		pushGroup.GET("/configs/:id/fields", s.appCtx.PushHandler.GetConfigFields)
// 		pushGroup.PUT("/configs/:id/fields", s.appCtx.PushHandler.UpdateConfigFields)
// 	}

// 	// 字段定义路由
// 	fieldGroup := api.Group("/fields", authMiddleware)
// 	{
// 		fieldGroup.GET("", s.appCtx.FieldHandler.List)
// 		fieldGroup.GET("/:id", s.appCtx.FieldHandler.Get)
// 		fieldGroup.POST("", s.appCtx.FieldHandler.Create)
// 		fieldGroup.PUT("/:id", s.appCtx.FieldHandler.Update)
// 		fieldGroup.DELETE("/:id", s.appCtx.FieldHandler.Delete)
// 	}

// 	// 定时任务路由
// 	scheduleGroup := api.Group("/schedule", authMiddleware)
// 	{
// 		scheduleGroup.GET("/tasks", s.appCtx.ScheduleHandler.ListTasks)
// 		scheduleGroup.GET("/tasks/:id", s.appCtx.ScheduleHandler.GetTask)
// 		scheduleGroup.POST("/tasks", s.appCtx.ScheduleHandler.CreateTask)
// 		scheduleGroup.PUT("/tasks/:id", s.appCtx.ScheduleHandler.UpdateTask)
// 		scheduleGroup.DELETE("/tasks/:id", s.appCtx.ScheduleHandler.DeleteTask)
// 		scheduleGroup.POST("/tasks/:id/toggle", s.appCtx.ScheduleHandler.ToggleTask)
// 		scheduleGroup.GET("/logs", s.appCtx.ScheduleHandler.GetLogs)
// 	}

// 	// 系统管理路由（需要管理员权限）
// 	systemGroup := api.Group("/system", authMiddleware, adminMiddleware)
// 	{
// 		systemGroup.GET("/users", s.appCtx.SystemHandler.ListUsers)
// 		systemGroup.POST("/users", s.appCtx.SystemHandler.CreateUser)
// 		systemGroup.PUT("/users/:id", s.appCtx.SystemHandler.UpdateUser)
// 		systemGroup.DELETE("/users/:id", s.appCtx.SystemHandler.DeleteUser)

// 		systemGroup.GET("/menus", s.appCtx.SystemHandler.ListMenus)
// 		systemGroup.POST("/menus", s.appCtx.SystemHandler.CreateMenu)
// 		systemGroup.PUT("/menus/:id", s.appCtx.SystemHandler.UpdateMenu)
// 		systemGroup.DELETE("/menus/:id", s.appCtx.SystemHandler.DeleteMenu)

// 		systemGroup.GET("/logs", s.appCtx.SystemHandler.GetLogs)
// 	}

// 	// 健康检查
// 	s.engine.GET("/health", func(c *gin.Context) {
// 		c.JSON(200, gin.H{"status": "ok"})
// 	})
// }

func (s *Server) MustRun() {
	port := s.appCtx.Config.Server.Port
	addr := fmt.Sprintf(":%d", port)

	// 启动定时任务调度器
	if err := s.scheduleService.StartScheduler(); err != nil {
		slog.Error("failed to start scheduler", "error", err)
	} else {
		slog.Info("scheduler started successfully")
	}

	slog.Info("server starting", "port", port)

	// if err := s.engine.Run(addr); err != nil {
	// 	slog.Error("server run error", "error", err)
	// 	panic(err)
	// }

	if err := s.httpRouter.Run(addr); err != nil {
		panic(err)
	}
}
