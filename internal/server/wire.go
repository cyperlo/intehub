//go:build wireinject
// +build wireinject

package server

import "github.com/google/wire"

func NewServer() (*Server, error) {
	wire.Build(
		// Config & DB
		MustProvideConfig,
		MustProvidePostgreSQLDB,

		// Services
		ProvideAuthService,
		ProvideAppService,
		ProvidePushService,
		ProvideFieldService,
		ProvideSystemService,
		ProvideScheduleService,

		// Handlers
		ProvideAuthHandler,
		ProvideAppHandler,
		ProvidePushHandler,
		ProvideFieldHandler,
		ProvideSystemHandler,
		ProvideScheduleHandler,

		// Context & Server
		wire.Struct(new(AppContext), "*"),
		NewServerWire,
	)
	return nil, nil
}
