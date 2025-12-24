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

	app, configs, err := h.appService.GetByIDWithConfigs(uint(id))
	if err != nil {
		return nil, errors.New("应用不存在")
	}

	return gin.H{
		"app":     app,
		"configs": configs,
	}, nil
}

// Create 创建应用
func (h *Handler) Create(c *gin.Context) (interface{}, error) {
	var req struct {
		appModel.App
		Configs []struct {
			Key       string `json:"key"`
			Value     string `json:"value"`
			Type      string `json:"type"`
			Encrypted bool   `json:"encrypted"`
		} `json:"configs"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, errors.New("参数错误")
	}

	userID, exists := c.Get("user_id")
	if exists && userID != nil {
		req.App.UserID = userID.(uint)
	}

	// 转换配置
	configs := make([]appService.ConfigInput, len(req.Configs))
	for i, cfg := range req.Configs {
		configs[i] = appService.ConfigInput{
			Key:       cfg.Key,
			Value:     cfg.Value,
			Type:      cfg.Type,
			Encrypted: cfg.Encrypted,
		}
	}

	if err := h.appService.CreateWithConfigs(&req.App, configs); err != nil {
		return nil, err
	}

	return req.App, nil
}

// Update 更新应用
func (h *Handler) Update(c *gin.Context) (interface{}, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return nil, errors.New("无效的ID")
	}

	var req struct {
		appModel.App
		Configs []struct {
			Key       string `json:"key"`
			Value     string `json:"value"`
			Type      string `json:"type"`
			Encrypted bool   `json:"encrypted"`
		} `json:"configs"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, errors.New("参数错误")
	}

	req.App.ID = uint(id)

	// 转换配置
	configs := make([]appService.ConfigInput, len(req.Configs))
	for i, cfg := range req.Configs {
		configs[i] = appService.ConfigInput{
			Key:       cfg.Key,
			Value:     cfg.Value,
			Type:      cfg.Type,
			Encrypted: cfg.Encrypted,
		}
	}

	if err := h.appService.UpdateWithConfigs(&req.App, configs); err != nil {
		return nil, err
	}

	return req.App, nil
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

	// 获取输入参数
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		// 如果没有提供参数，使用空 map
		input = make(map[string]interface{})
	}

	log, err := h.appService.RunWithInput(uint(id), input)
	if err != nil {
		return nil, err
	}

	// 返回详细的执行结果
	// 包含日志信息和原始输出，方便前端处理
	return gin.H{
		"log":      log,          // 完整的日志对象
		"status":   log.Status,   // 执行状态
		"output":   log.Output,   // 输出结果（字符串格式）
		"error":    log.Error,    // 错误信息
		"duration": log.Duration, // 执行时长（毫秒）
	}, nil
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

// Publish 发布应用到应用商店
func (h *Handler) Publish(c *gin.Context) (interface{}, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return nil, errors.New("无效的ID")
	}

	var req struct {
		DisplayName  string `json:"display_name" binding:"required"`
		Description  string `json:"description"`
		Icon         string `json:"icon"`
		Category     string `json:"category"`
		Version      string `json:"version" binding:"required"`
		Author       string `json:"author"`
		Tags         string `json:"tags"`
		ConfigSchema string `json:"config_schema"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, errors.New("参数错误")
	}

	template, err := h.appService.PublishToStore(uint(id), req.DisplayName, req.Description, req.Icon,
		req.Category, req.Version, req.Author, req.Tags, req.ConfigSchema)
	if err != nil {
		return nil, err
	}

	return template, nil
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
	j.POST("/:id/publish", h.Publish)
	j.GET("/logs", h.GetLogs)
}
