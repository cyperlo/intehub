package push

import (
	model "intehub/internal/app/models"
	"intehub/internal/app/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	pushService service.PushService
}

func NewHandler(pushService service.PushService) *Handler {
	return &Handler{pushService: pushService}
}

func (h *Handler) ListConfigs(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var uid uint
	if role != "admin" {
		uid = userID.(uint)
	}

	configs, err := h.pushService.GetConfigs(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, configs)
}

func (h *Handler) GetConfig(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	config, err := h.pushService.GetConfig(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "配置不存在"})
		return
	}
	c.JSON(http.StatusOK, config)
}

func (h *Handler) CreateConfig(c *gin.Context) {
	var config model.PushConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	config.UserID = userID.(uint)

	if err := h.pushService.CreateConfig(&config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, config)
}

func (h *Handler) UpdateConfig(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var config model.PushConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.ID = uint(id)
	if err := h.pushService.UpdateConfig(&config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, config)
}

func (h *Handler) DeleteConfig(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.pushService.DeleteConfig(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func (h *Handler) Send(c *gin.Context) {
	var req struct {
		ConfigID uint                   `json:"config_id"`
		Data     map[string]interface{} `json:"data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.pushService.Send(req.ConfigID, req.Data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "发送成功"})
}

func (h *Handler) GetHistory(c *gin.Context) {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      history,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *Handler) GetConfigFields(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	fields, err := h.pushService.GetConfigFields(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, fields)
}

func (h *Handler) UpdateConfigFields(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var req struct {
		FieldIDs []uint `json:"field_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.pushService.UpdateConfigFields(uint(id), req.FieldIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}
