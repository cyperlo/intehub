package app

import (
	"errors"
	appModel "intehub/internal/app/models/app"
	appService "intehub/internal/app/service/app"
	httputil "intehub/internal/utils/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	appService appService.Service
}

func NewHandler(appService appService.Service) *Handler {
	return &Handler{
		appService: appService,
	}
}

// List 获取应用列表
func (h *Handler) List(c *gin.Context) (interface{}, error) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var uid uint
	if role != "admin" && userID != nil {
		uid = userID.(uint)
	}

	apps, err := h.appService.List(uid)
	if err != nil {
		return nil, err
	}

	return apps, nil
}

// Get 获取应用详情
func (h *Handler) Get(c *gin.Context) (interface{}, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return nil, errors.New("无效的ID")
	}

	app, err := h.appService.GetByID(uint(id))
	if err != nil {
		return nil, errors.New("应用不存在")
	}

	return app, nil
}

// Create 创建应用
func (h *Handler) Create(c *gin.Context) (interface{}, error) {
	var app appModel.App
	if err := c.ShouldBindJSON(&app); err != nil {
		return nil, errors.New("参数错误")
	}

	userID, exists := c.Get("user_id")
	if exists && userID != nil {
		app.UserID = userID.(uint)
	}

	if err := h.appService.Create(&app); err != nil {
		return nil, err
	}

	return app, nil
}

// Update 更新应用
func (h *Handler) Update(c *gin.Context) (interface{}, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return nil, errors.New("无效的ID")
	}

	var app appModel.App
	if err := c.ShouldBindJSON(&app); err != nil {
		return nil, errors.New("参数错误")
	}

	app.ID = uint(id)
	if err := h.appService.Update(&app); err != nil {
		return nil, err
	}

	return app, nil
}

// Delete 删除应用
func (h *Handler) Delete(c *gin.Context) (interface{}, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return nil, errors.New("无效的ID")
	}

	if err := h.appService.Delete(uint(id)); err != nil {
		return nil, err
	}

	return gin.H{"message": "删除成功"}, nil
}

// Run 运行应用
func (h *Handler) Run(c *gin.Context) (interface{}, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return nil, errors.New("无效的ID")
	}

	log, err := h.appService.Run(uint(id))
	if err != nil {
		return nil, err
	}

	return log, nil
}

// GetLogs 获取应用日志
func (h *Handler) GetLogs(c *gin.Context) (interface{}, error) {
	var appID *uint
	if id := c.Query("app_id"); id != "" {
		if parsed, err := strconv.ParseUint(id, 10, 32); err == nil {
			uid := uint(parsed)
			appID = &uid
		}
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	logs, total, err := h.appService.GetLogs(appID, page, pageSize)
	if err != nil {
		return nil, err
	}

	return gin.H{
		"list":  logs,
		"total": total,
		"page":  page,
	}, nil
}

// HandleAppAPI 注册路由
func (h *Handler) HandleAppAPI(r *gin.RouterGroup) {
	j := httputil.NewJSONHandler(r)
	j.GET("", h.List)
	j.GET("/:id", h.Get)
	j.POST("", h.Create)
	j.PUT("/:id", h.Update)
	j.DELETE("/:id", h.Delete)
	j.POST("/:id/run", h.Run)
	j.GET("/logs", h.GetLogs)
}
