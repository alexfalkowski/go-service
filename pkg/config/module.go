package config

import (
	"go.uber.org/fx"
)

var (
	// Module for fx.
	Module = fx.Options(ConfiguratorModule, UnmarshalModule, ConfigModule)

	// ConfiguratorModule for fx.
	ConfiguratorModule = fx.Provide(NewConfigurator)

	// UnmarshalModule for fx.
	UnmarshalModule = fx.Invoke(Unmarshal)

	// ConfigModule for fx.
	ConfigModule = fx.Options(
		fx.Provide(redisConfig), fx.Provide(ristrettoConfig),
		fx.Provide(auth0Config),
		fx.Provide(pgConfig),
		fx.Provide(datadogConfig), fx.Provide(jaegerConfig),
		fx.Provide(grpcConfig), fx.Provide(httpConfig), fx.Provide(nsqConfig),
	)
)
