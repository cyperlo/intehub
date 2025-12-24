package workflow

import (
	"errors"
	workflowModel "intehub/internal/app/models/workflow"
	workflowService "intehub/internal/app/service/workflow"
	httputil "intehub/internal/utils/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	workflowService workflowService.Service
}

func NewHandler(workflowService workflowService.Service) *Handler {
	return &Handler{
		workflowService: workflowService,
	}
}

// List 获取工作流列表
func (h *Handler) List(c *gin.Context) (interface{}, error) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var uid uint
	if role != "admin" && userID != nil {
		uid = userID.(uint)
	}

	workflows, err := h.workflowService.List(uid)
	if err != nil {
		return nil, err
	}

	return workflows, nil
}

// Get 获取工作流详情
func (h *Handler) Get(c *gin.Context) (interface{}, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return nil, errors.New("无效的ID")
	}

	workflow, err := h.workflowService.GetByID(uint(id))
	if err != nil {
		return nil, errors.New("工作流不存在")
	}

	return workflow, nil
}

// Create 创建工作流
func (h *Handler) Create(c *gin.Context) (interface{}, error) {
	var workflow workflowModel.Workflow

	if err := c.ShouldBindJSON(&workflow); err != nil {
		return nil, errors.New("参数错误")
	}

	userID, exists := c.Get("user_id")
	if exists && userID != nil {
		workflow.UserID = userID.(uint)
	}

	if err := h.workflowService.Create(&workflow); err != nil {
		return nil, err
	}

	return workflow, nil
}

// Update 更新工作流
func (h *Handler) Update(c *gin.Context) (interface{}, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return nil, errors.New("无效的ID")
	}

	var workflow workflowModel.Workflow
	if err := c.ShouldBindJSON(&workflow); err != nil {
		return nil, errors.New("参数错误")
	}

	workflow.ID = uint(id)

	if err := h.workflowService.Update(&workflow); err != nil {
		return nil, err
	}

	return workflow, nil
}

// Delete 删除工作流
func (h *Handler) Delete(c *gin.Context) (interface{}, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return nil, errors.New("无效的ID")
	}

	if err := h.workflowService.Delete(uint(id)); err != nil {
		return nil, err
	}

	return gin.H{"message": "删除成功"}, nil
}

// Run 运行工作流
func (h *Handler) Run(c *gin.Context) (interface{}, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return nil, errors.New("无效的ID")
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		input = make(map[string]interface{})
	}

	log, err := h.workflowService.Run(uint(id), input)
	if err != nil {
		return nil, err
	}

	return gin.H{
		"log":      log,
		"status":   log.Status,
		"output":   log.Output,
		"error":    log.Error,
		"duration": log.Duration,
	}, nil
}

// GetLogs 获取工作流日志
func (h *Handler) GetLogs(c *gin.Context) (interface{}, error) {
	var workflowID *uint
	if id := c.Query("workflow_id"); id != "" {
		if parsed, err := strconv.ParseUint(id, 10, 32); err == nil {
			uid := uint(parsed)
			workflowID = &uid
		}
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	logs, total, err := h.workflowService.GetLogs(workflowID, page, pageSize)
	if err != nil {
		return nil, err
	}

	return gin.H{
		"list":  logs,
		"total": total,
		"page":  page,
	}, nil
}

// HandleWorkflowAPI 注册路由
func (h *Handler) HandleWorkflowAPI(r *gin.RouterGroup) {
	j := httputil.NewJSONHandler(r)
	j.GET("", h.List)
	j.GET("/:id", h.Get)
	j.POST("", h.Create)
	j.PUT("/:id", h.Update)
	j.DELETE("/:id", h.Delete)
	j.POST("/:id/run", h.Run)
	j.GET("/logs", h.GetLogs)
}
