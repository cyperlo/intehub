package health

import (
	"intehub/internal/utils/http"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	HandlerHealthAPI(r *gin.RouterGroup)
}

type handler struct {
}

func (h *handler) HandlerHealthAPI(r *gin.RouterGroup) {
	j := http.NewJSONHandler(r)
	j.GET("", h.HealthCheck)
}

func (h *handler) HealthCheck(ctx *gin.Context) (interface{}, error) {
	return "health", nil
}

func New() Handler {
	return &handler{}
}
