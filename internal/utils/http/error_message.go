package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 错误码
const (
	ParamsError         = 400
	Unauthorized        = 401
	Forbidden           = 403
	NotFound            = 404
	MethodNotAllowed    = 405
	InternalServerError = 500
)

var ErrMessage = map[int]string{
	ParamsError:         "error.params_error",
	Unauthorized:        "error.unauthorized",
	Forbidden:           "error.forbidden",
	NotFound:            "error.not_found",
	MethodNotAllowed:    "error.method_not_allowed",
	InternalServerError: "error.internal_server_error",
}

func HandlerNotFound(ctx *gin.Context) {
	response := NormalResponse
	response.Code = NotFound
	response.Message = ErrMessage[response.Code]
	ctx.JSON(http.StatusNotFound, response)
}

func HandlerMethodNotAllowed(ctx *gin.Context) {
	response := NormalResponse
	response.Code = MethodNotAllowed
	response.Message = ErrMessage[response.Code]
	ctx.JSON(http.StatusMethodNotAllowed, response)
}
