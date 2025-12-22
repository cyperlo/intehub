package push

import (
	"errors"
	pushModel "intehub/internal/app/models/push"
	pushService "intehub/internal/app/service/push"
	httputil "intehub/internal/utils/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	pushService pushService.Service
}

func NewHandler(pushService pushService.Service) *Handler {
	return &Handler{pushService: pushService}
}

func (h *Handler) ListConfigs(c *gin.Context) (interface{}, error) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var uid uint
	if role != "admin" && userID != nil {
		uid = userID.(uint)
	}

	configs, err := h.pushService.GetConfigs(uid)
	if err != nil {
		return nil, err
	}
	return configs, nil
}

func (h *Handler) GetConfig(c *gin.Context) (interface{}, error) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	config, err := h.pushService.GetConfig(uint(id))
	if err != nil {
		return nil, errors.New("配置不存在")
	}
	return config, nil
}

func (h *Handler) CreateConfig(c *gin.Context) (interface{}, error) {
	var config pushModel.PushConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		return nil, errors.New("参数错误")
	}

	userID, exists := c.Get("user_id")
	if exists && userID != nil {
		config.UserID = userID.(uint)
	}

	if err := h.pushService.CreateConfig(&config); err != nil {
		return nil, err
	}
	return config, nil
}

func (h *Handler) UpdateConfig(c *gin.Context) (interface{}, error) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var config pushModel.PushConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		return nil, errors.New("参数错误")
	}

	config.ID = uint(id)
	if err := h.pushService.UpdateConfig(&config); err != nil {
		return nil, err
	}
	return config, nil
}

func (h *Handler) DeleteConfig(c *gin.Context) (interface{}, error) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.pushService.DeleteConfig(uint(id)); err != nil {
		return nil, err
	}
	return gin.H{"message": "删除成功"}, nil
}

func (h *Handler) Send(c *gin.Context) (interface{}, error) {
	var req struct {
		ConfigID uint                   `json:"config_id"`
		Data     map[string]interface{} `json:"data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, errors.New("参数错误")
	}

	if err := h.pushService.Send(req.ConfigID, req.Data); err != nil {
		return nil, err
	}
	return gin.H{"message": "发送成功"}, nil
}

func (h *Handler) GetHistory(c *gin.Context) (interface{}, error) {
	var configID *uint
	if id := c.Query("config_id"); id != "" {
		if parsed, err := strconv.ParseUint(id, 10, 32); err == nil {
			uid := uint(parsed)
			configID = &uid
		}
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	history, total, err := h.pushService.GetHistory(configID, page, pageSize)
	if err != nil {
		return nil, err
	}

	return gin.H{
		"data":      history,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}, nil
}

func (h *Handler) GetConfigFields(c *gin.Context) (interface{}, error) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	fields, err := h.pushService.GetConfigFields(uint(id))
	if err != nil {
		return nil, err
	}
	return fields, nil
}

func (h *Handler) UpdateConfigFields(c *gin.Context) (interface{}, error) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var req struct {
		FieldIDs []uint `json:"field_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, errors.New("参数错误")
	}

	if err := h.pushService.UpdateConfigFields(uint(id), req.FieldIDs); err != nil {
		return nil, err
	}
	return gin.H{"message": "更新成功"}, nil
}

// HandlePushAPI 注册路由
func (h *Handler) HandlePushAPI(r *gin.RouterGroup) {
	j := httputil.NewJSONHandler(r)
	j.GET("/configs", h.ListConfigs)
	j.GET("/configs/:id", h.GetConfig)
	j.POST("/configs", h.CreateConfig)
	j.PUT("/configs/:id", h.UpdateConfig)
	j.DELETE("/configs/:id", h.DeleteConfig)
	j.POST("/send", h.Send)
	j.GET("/history", h.GetHistory)
	j.GET("/configs/:id/fields", h.GetConfigFields)
	j.PUT("/configs/:id/fields", h.UpdateConfigFields)
}
