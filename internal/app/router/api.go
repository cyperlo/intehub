package router

import "github.com/gin-gonic/gin"

func handleAPI(r *HttpRouter, rg *gin.RouterGroup, lrMw gin.HandlerFunc,
	expireMw gin.HandlerFunc, authKeyMw gin.HandlerFunc, sourceTypeMw gin.HandlerFunc,
) {
	{
		// 健康检查
		healthGroup := rg.Group("/health")
		r.healthHandler.HandlerHealthAPI(healthGroup)
	}

	// {
	// 	// v1
	// 	v1 := rg.Group("v1")

	// }
}
