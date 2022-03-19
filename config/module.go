package config

import (
	"go.uber.org/fx"
)

var (
	// Module for fx.
	// nolint:gochecknoglobals
	Module = fx.Options(ConfiguratorModule, UnmarshalModule, ConfigModule)

	// ConfiguratorModule for fx.
	// nolint:gochecknoglobals
	ConfiguratorModule = fx.Provide(NewConfigurator)

	// UnmarshalModule for fx.
	// nolint:gochecknoglobals
	UnmarshalModule = fx.Invoke(Unmarshal)

	// ConfigModule for fx.
	// nolint:gochecknoglobals
	ConfigModule = fx.Options(
		fx.Provide(redisConfig), fx.Provide(ristrettoConfig),
		fx.Provide(auth0Config),
		fx.Provide(pgConfig),
		fx.Provide(datadogConfig), fx.Provide(jaegerConfig),
		fx.Provide(grpcConfig), fx.Provide(httpConfig), fx.Provide(nsqConfig),
	)
)
