package auth

import (
	"errors"
	sa "intehub/internal/app/service/auth"
	"intehub/internal/utils/http"
	"log/slog"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	HandleAuthAPI(r *gin.RouterGroup)
}
type handler struct {
	saService sa.Service
}

func New(
	// authService service.AuthService,
	saService sa.Service,
) Handler {
	return &handler{
		// authService: authService,
		saService: saService,
	}
}

func (h *handler) HandleAuthAPI(r *gin.RouterGroup) {
	j := http.NewJSONHandler(r)
	j.POST("/login", h.Login)
}

func (h *handler) Login(ctx *gin.Context) (interface{}, error) {
	var req sa.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, errors.New("param error")
	}
	resp, err := h.saService.Login(ctx.Request.Context(), req.Username, req.Password)
	if err != nil {
		slog.Error(err.Error())
	}
	return resp, nil

}

// // Login 用户登录
// func (h *Handler) Login(c *gin.Context) {
// 	var req LoginRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	token, user, err := h.authService.Login(req.Username, req.Password)
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, LoginResponse{
// 		Token: token,
// 		User:  user,
// 	})
// }

// // Logout 用户登出
// func (h *Handler) Logout(c *gin.Context) {
// 	c.JSON(http.StatusOK, gin.H{"message": "登出成功"})
// }

// // GetProfile 获取用户信息
// func (h *Handler) GetProfile(c *gin.Context) {
// 	userID, exists := c.Get("user_id")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
// 		return
// 	}

// 	user, err := h.authService.GetProfile(userID.(uint))
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, user)
// }
