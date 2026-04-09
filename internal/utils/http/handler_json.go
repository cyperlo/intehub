package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type HandlerJSONFunc func(ctx *gin.Context) (interface{}, error)

type JSONHandler gin.RouterGroup

// BadRequest 创建业务错误（返回 400 状态码）
type BadRequest struct {
	Message string
}

func (e *BadRequest) Error() string {
	return e.Message
}

// NotFoundError 创建未找到错误（返回 404 状态码）
type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

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

func JSONFunc(f HandlerJSONFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		jsonResp, errWithStack := f(ctx)
		switch err := errors.Cause(errWithStack).(type) {
		case nil:
			resp := NormalResponse
			resp.Data = jsonResp
			ctx.JSON(http.StatusOK, resp)
			return
		case *BadRequest:
			ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Message})
			return
		case *NotFoundError:
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": err.Message})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
			return
		}
	}
}
