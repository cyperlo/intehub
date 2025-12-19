package app

import (
	model "intehub/internal/app/models"
	"intehub/internal/app/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	appService service.AppService
}

func NewHandler(appService service.AppService) *Handler {
	return &Handler{
		appService: appService,
	}
}

// List 获取应用列表
func (h *Handler) List(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var uid uint
	if role != "admin" {
		uid = userID.(uint)
	}

	apps, err := h.appService.List(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, apps)
}

// Get 获取应用详情
func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	app, err := h.appService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "应用不存在"})
		return
	}

	c.JSON(http.StatusOK, app)
}

// Create 创建应用
func (h *Handler) Create(c *gin.Context) {
	var app model.App
	if err := c.ShouldBindJSON(&app); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	app.UserID = userID.(uint)

	if err := h.appService.Create(&app); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, app)
}

// Update 更新应用
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var app model.App
	if err := c.ShouldBindJSON(&app); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	app.ID = uint(id)
	if err := h.appService.Update(&app); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, app)
}

// Delete 删除应用
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	if err := h.appService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// Run 运行应用
func (h *Handler) Run(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	log, err := h.appService.Run(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, log)
}

// GetLogs 获取应用日志
func (h *Handler) GetLogs(c *gin.Context) {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"list":  logs,
		"total": total,
		"page":  page,
	})
}
