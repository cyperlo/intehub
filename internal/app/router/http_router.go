package router

import (
	"intehub/internal/app/api/v1/auth"
	"intehub/internal/app/api/v1/health"
	"intehub/internal/app/context"
	"intehub/internal/app/router/middleware"
	"intehub/internal/utils/http"

	"github.com/gin-gonic/gin"
)

type HttpRouter struct {
	*gin.Engine
	appContext    *context.AppContext
	healthHandler health.Handler
	authHandler   auth.Handler
}

func NewHttpRouter(
	appContext *context.AppContext,
	healthHandler health.Handler,
	authHandler auth.Handler,
) *HttpRouter {
	r := HttpRouter{
		appContext:    appContext,
		healthHandler: healthHandler,
		authHandler:   authHandler,
	}

	router := gin.New()

	router.Use(middleware.InitAppContext(appContext))

	routerPrefix := appContext.Config.ApiPrefix

	// router.Use(middleware.AuthMiddleware("test"))
	router.HandleMethodNotAllowed = true
	router.NoRoute(http.HandlerNotFound)
	router.NoMethod(http.HandlerMethodNotAllowed)

	handleAPI(&r, router.Group(routerPrefix), nil, nil, nil, nil)

	r.Engine = router
	return &r
}
