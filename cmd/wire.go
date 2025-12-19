//go:build wireinject
// +build wireinject

package cmd

import (
	"intehub/internal/app/context"

	"intehub/internal/app/router"

	"intehub/internal/app/api/v1/auth"
	"intehub/internal/app/api/v1/health"

	sa "intehub/internal/app/service/auth"

	"github.com/google/wire"

	"intehub/internal/app"
)

var wireSet = wire.NewSet(
	// Config & DB
	MustProvideConfig,
	MustProvidePostgreSQLDB,
	MustProvideModel,

	router.NewHttpRouter,

	// routers
	health.New,
	auth.New,

	// Services
	sa.New,

	ProvideAuthService,
	ProvideAppService,
	ProvidePushService,
	ProvideFieldService,
	ProvideSystemService,
	ProvideScheduleService,

	// Handlers
	// ProvideAuthHandler,
	ProvideAppHandler,
	ProvidePushHandler,
	ProvideFieldHandler,
	ProvideSystemHandler,
	ProvideScheduleHandler,
)

func NewServer() (*app.Server, error) {
	wire.Build(
		wireSet,
		// Context & Server
		wire.Struct(new(context.AppContext), "*"),
		app.NewServerWire,
	)
	return nil, nil
}
