package system

import (
	"errors"
	systemModel "intehub/internal/app/models/system"
	userModel "intehub/internal/app/models/user"
	systemService "intehub/internal/app/service/system"
	httputil "intehub/internal/utils/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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

	// User routes
	j.GET("/users", h.ListUsers)
	j.POST("/users", h.CreateUser)
	j.PUT("/users/:id", h.UpdateUser)
	j.DELETE("/users/:id", h.DeleteUser)
}

// User handlers
func (h *Handler) ListUsers(c *gin.Context) (interface{}, error) {
	users, err := h.systemService.GetUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (h *Handler) CreateUser(c *gin.Context) (interface{}, error) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Nickname string `json:"nickname"`
		Password string `json:"password" binding:"required"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, errors.New("参数错误")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	user := &userModel.DataObject{
		Username: req.Username,
		Nickname: req.Nickname,
		Password: string(hashedPassword),
		Role:     req.Role,
	}

	if user.Role == "" {
		user.Role = "user"
	}

	if err := h.systemService.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (h *Handler) UpdateUser(c *gin.Context) (interface{}, error) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req struct {
		Nickname string `json:"nickname"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, errors.New("参数错误")
	}

	user, err := h.systemService.GetUserByID(uint(id))
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 构建更新字段
	updates := make(map[string]interface{})

	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
		user.Nickname = req.Nickname
	}
	if req.Role != "" {
		updates["role"] = req.Role
		user.Role = req.Role
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New("密码加密失败")
		}
		updates["password"] = string(hashedPassword)
		user.Password = string(hashedPassword)
	}

	if len(updates) == 0 {
		return user, nil
	}

	if err := h.systemService.UpdateUserFields(uint(id), updates); err != nil {
		return nil, err
	}

	return user, nil
}

func (h *Handler) DeleteUser(c *gin.Context) (interface{}, error) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.systemService.DeleteUser(uint(id)); err != nil {
		return nil, err
	}
	return gin.H{"message": "删除成功"}, nil
}
