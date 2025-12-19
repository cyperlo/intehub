package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type HandlerJSONFunc func(ctx *gin.Context) (interface{}, error)

type JSONHandler gin.RouterGroup

func (j *JSONHandler) GET(relativePath string, handler HandlerJSONFunc) {
	g := (*gin.RouterGroup)(j)
	g.GET(relativePath, JSONFunc(handler))
}

func (j *JSONHandler) POST(relativePath string, handler HandlerJSONFunc) {
	g := (*gin.RouterGroup)(j)
	g.POST(relativePath, JSONFunc(handler))
}

func (j *JSONHandler) PUT(relativePath string, handler HandlerJSONFunc) {
	g := (*gin.RouterGroup)(j)
	g.PUT(relativePath, JSONFunc(handler))
}

func (j *JSONHandler) PATCH(relativePath string, handler HandlerJSONFunc) {
	g := (*gin.RouterGroup)(j)
	g.PATCH(relativePath, JSONFunc(handler))
}

func (j *JSONHandler) DELETE(relativePath string, handler HandlerJSONFunc) {
	g := (*gin.RouterGroup)(j)
	g.DELETE(relativePath, JSONFunc(handler))
}

func NewJSONHandler(routerGroup *gin.RouterGroup) *JSONHandler {
	return (*JSONHandler)(routerGroup)
}

// TODO: 增加不同错误类型
func JSONFunc(f HandlerJSONFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		jsonResp, errWithStack := f(ctx)
		switch err := errors.Cause(errWithStack).(type) {
		case nil:
			resp := NormalResponse
			resp.Data = jsonResp
			ctx.JSON(http.StatusOK, resp)
			return
		default:
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}
	}
}
