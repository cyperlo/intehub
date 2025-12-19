package http

import (
	"github.com/gin-gonic/gin"
)

var NormalResponse = Response{
	Code:    0,
	Message: "Succeed",
}

type Response struct {
	Code    int         `json:"code" example:"0"`          // 响应状态码
	Message string      `json:"message" example:"succeed"` // 响应信息
	Detail  interface{} `json:"detail,omitempty"`          // 错误详情
	Data    interface{} `json:"data,omitempty"`            // 响应数据
}

func NewInternalErrorResponse(ctx *gin.Context, detail string) *Response {
	return &Response{
		Code:    InternalServerError,
		Message: ErrMessage[InternalServerError],
		Detail:  detail,
	}
}
