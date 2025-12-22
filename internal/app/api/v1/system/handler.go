package system

import (
	"errors"
	systemModel "intehub/internal/app/models/system"
	systemService "intehub/internal/app/service/system"
	httputil "intehub/internal/utils/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	systemService systemService.Service
}

func NewHandler(systemService systemService.Service) *Handler {
	return &Handler{systemService: systemService}
}

// Menu handlers
func (h *Handler) ListMenus(c *gin.Context) (interface{}, error) {
	menus, err := h.systemService.GetMenus()
	if err != nil {
		return nil, err
	}
	return menus, nil
}

func (h *Handler) CreateMenu(c *gin.Context) (interface{}, error) {
	var menu systemModel.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		return nil, errors.New("参数错误")
	}

	if err := h.systemService.CreateMenu(&menu); err != nil {
		return nil, err
	}
	return menu, nil
}

func (h *Handler) UpdateMenu(c *gin.Context) (interface{}, error) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var menu systemModel.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		return nil, errors.New("参数错误")
	}

	menu.ID = uint(id)
	if err := h.systemService.UpdateMenu(&menu); err != nil {
		return nil, err
	}
	return menu, nil
}

func (h *Handler) DeleteMenu(c *gin.Context) (interface{}, error) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.systemService.DeleteMenu(uint(id)); err != nil {
		return nil, err
	}
	return gin.H{"message": "删除成功"}, nil
}

// Log handlers
func (h *Handler) GetLogs(c *gin.Context) (interface{}, error) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	logs, total, err := h.systemService.GetLogs(page, pageSize)
	if err != nil {
		return nil, err
	}

	return gin.H{
		"list":  logs,
		"total": total,
		"page":  page,
	}, nil
}

// HandleSystemAPI 注册路由
func (h *Handler) HandleSystemAPI(r *gin.RouterGroup) {
	j := httputil.NewJSONHandler(r)

	// Menu routes
	j.GET("/menus", h.ListMenus)
	j.POST("/menus", h.CreateMenu)
	j.PUT("/menus/:id", h.UpdateMenu)
	j.DELETE("/menus/:id", h.DeleteMenu)

	// Log routes
	j.GET("/logs", h.GetLogs)
}
