package router

import (
	"fmt"
	"intehub/internal/app/api/v1/app"
	"intehub/internal/app/api/v1/appstore"
	"intehub/internal/app/api/v1/auth"
	"intehub/internal/app/api/v1/field"
	"intehub/internal/app/api/v1/health"
	"intehub/internal/app/api/v1/push"
	"intehub/internal/app/api/v1/schedule"
	"intehub/internal/app/api/v1/system"
	"intehub/internal/app/api/v1/workflow"
	"intehub/internal/app/context"
	"intehub/internal/app/router/middleware"
	"intehub/internal/utils/http"

	"github.com/gin-gonic/gin"
)

type HttpRouter struct {
	*gin.Engine
	appContext      *context.AppContext
	healthHandler   health.Handler
	authHandler     auth.Handler
	appHandler      *app.Handler
	appstoreHandler *appstore.Handler
	workflowHandler *workflow.Handler
	fieldHandler    *field.Handler
	pushHandler     *push.Handler
	scheduleHandler *schedule.Handler
	systemHandler   *system.Handler
}

func NewHttpRouter(
	appContext *context.AppContext,
	healthHandler health.Handler,
	authHandler auth.Handler,
	appHandler *app.Handler,
	appstoreHandler *appstore.Handler,
	workflowHandler *workflow.Handler,
	fieldHandler *field.Handler,
	pushHandler *push.Handler,
	scheduleHandler *schedule.Handler,
	systemHandler *system.Handler,
) *HttpRouter {
	r := HttpRouter{
		appContext:      appContext,
		healthHandler:   healthHandler,
		authHandler:     authHandler,
		appHandler:      appHandler,
		appstoreHandler: appstoreHandler,
		workflowHandler: workflowHandler,
		fieldHandler:    fieldHandler,
		pushHandler:     pushHandler,
		scheduleHandler: scheduleHandler,
		systemHandler:   systemHandler,
	}

	router := gin.New()

	router.Use(middleware.InitAppContext(appContext))

	routerPrefix := appContext.Config.ApiPrefix

	loginRequireConfig := middleware.LoginRequireConfig{
		ExcludedPath: map[string]struct{}{
			fmt.Sprintf(middleware.LoginRequiredURLFmt, "POST",
				routerPrefix+"/v1/auth/login"): {},
		},
	}

	router.Use(middleware.AuthMiddleware(appContext.Config.JWT.Secret, loginRequireConfig))
	router.HandleMethodNotAllowed = true
	router.NoRoute(http.HandlerNotFound)
	router.NoMethod(http.HandlerMethodNotAllowed)

	api := router.Group(routerPrefix)
	v1 := api.Group("/v1")

	handleAPI(&r, v1)

	r.Engine = router
	return &r
}
