package auth

import (
	"errors"
	sa "intehub/internal/app/service/auth"
	"intehub/internal/utils/http"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	HandleAuthAPI(r *gin.RouterGroup)
}

type handler struct {
	saService sa.Service
}

func New(saService sa.Service) Handler {
	return &handler{
		saService: saService,
	}
}

func (h *handler) HandleAuthAPI(r *gin.RouterGroup) {
	j := http.NewJSONHandler(r)
	j.POST("/login", h.Login)
	j.POST("/logout", h.Logout)
}

func (h *handler) Login(ctx *gin.Context) (interface{}, error) {
	var req sa.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, errors.New("param error")
	}
	resp, err := h.saService.Login(ctx.Request.Context(), req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (h *handler) Logout(ctx *gin.Context) (interface{}, error) {
	return gin.H{"message": "登出成功"}, nil
}
