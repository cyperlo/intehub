package middleware

import (
	"intehub/internal/app/context"

	"github.com/gin-gonic/gin"
)

// 初始化上下文
func InitAppContext(appCtx *context.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		savedCtx := ctx.Request.Context()
		defer func() {
			ctx.Request = ctx.Request.WithContext(savedCtx)
		}()

		newCtx := context.FillAppContext(savedCtx, appCtx)
		ctx.Request = ctx.Request.WithContext(newCtx)

		ctx.Next()
	}
}
