package schedule

import (
	"errors"
	scheduleModel "intehub/internal/app/models/schedule"
	scheduleService "intehub/internal/app/service/schedule"
	httputil "intehub/internal/utils/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	scheduleService scheduleService.Service
}

func NewHandler(scheduleService scheduleService.Service) *Handler {
	return &Handler{scheduleService: scheduleService}
}

func (h *Handler) ListTasks(c *gin.Context) (interface{}, error) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var uid uint
	if role != "admin" && userID != nil {
		uid = userID.(uint)
	}

	tasks, err := h.scheduleService.GetTasks(uid)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (h *Handler) GetTask(c *gin.Context) (interface{}, error) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	task, err := h.scheduleService.GetTask(uint(id))
	if err != nil {
		return nil, errors.New("任务不存在")
	}
	return task, nil
}

func (h *Handler) CreateTask(c *gin.Context) (interface{}, error) {
	var task scheduleModel.ScheduleTask
	if err := c.ShouldBindJSON(&task); err != nil {
		return nil, errors.New("参数错误")
	}

	userID, exists := c.Get("user_id")
	if exists && userID != nil {
		task.UserID = userID.(uint)
	}

	if err := h.scheduleService.CreateTask(&task); err != nil {
		return nil, err
	}
	return task, nil
}

func (h *Handler) UpdateTask(c *gin.Context) (interface{}, error) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var task scheduleModel.ScheduleTask
	if err := c.ShouldBindJSON(&task); err != nil {
		return nil, errors.New("参数错误")
	}

	task.ID = uint(id)
	if err := h.scheduleService.UpdateTask(&task); err != nil {
		return nil, err
	}
	return task, nil
}

func (h *Handler) DeleteTask(c *gin.Context) (interface{}, error) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.scheduleService.DeleteTask(uint(id)); err != nil {
		return nil, err
	}
	return gin.H{"message": "删除成功"}, nil
}

func (h *Handler) ToggleTask(c *gin.Context) (interface{}, error) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.scheduleService.ToggleTask(uint(id)); err != nil {
		return nil, err
	}
	return gin.H{"message": "操作成功"}, nil
}

func (h *Handler) GetLogs(c *gin.Context) (interface{}, error) {
	var taskID *uint
	if id := c.Query("task_id"); id != "" {
		if parsed, err := strconv.ParseUint(id, 10, 32); err == nil {
			uid := uint(parsed)
			taskID = &uid
		}
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	logs, total, err := h.scheduleService.GetLogs(taskID, page, pageSize)
	if err != nil {
		return nil, err
	}

	return gin.H{
		"list":  logs,
		"total": total,
		"page":  page,
	}, nil
}

// HandleScheduleAPI 注册路由
func (h *Handler) HandleScheduleAPI(r *gin.RouterGroup) {
	j := httputil.NewJSONHandler(r)
	j.GET("/tasks", h.ListTasks)
	j.GET("/tasks/:id", h.GetTask)
	j.POST("/tasks", h.CreateTask)
	j.PUT("/tasks/:id", h.UpdateTask)
	j.DELETE("/tasks/:id", h.DeleteTask)
	j.POST("/tasks/:id/toggle", h.ToggleTask)
	j.GET("/logs", h.GetLogs)
}
