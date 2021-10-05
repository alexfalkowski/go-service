package config

import (
	"go.uber.org/fx"
)

var (
	// Module for fx.
	Module = fx.Options(
		fx.Provide(New),
		fx.Provide(redisConfig), fx.Provide(ristrettoConfig),
		fx.Provide(auth0Config),
		fx.Provide(pgConfig),
		fx.Provide(datadogConfig), fx.Provide(jaegerConfig),
		fx.Provide(grpcConfig), fx.Provide(httpConfig), fx.Provide(nsqConfig),
	)
)
