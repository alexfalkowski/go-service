package config

import (
	"go.uber.org/fx"
)

var (
	// Module for fx.
	// nolint:gochecknoglobals
	Module = fx.Options(ConfiguratorModule, UnmarshalModule, ConfigModule, WatchModule)

	// ConfiguratorModule for fx.
	// nolint:gochecknoglobals
	ConfiguratorModule = fx.Provide(NewConfigurator)

	// UnmarshalModule for fx.
	// nolint:gochecknoglobals
	UnmarshalModule = fx.Invoke(UnmarshalFromFile)

	// WatchModule for fx.
	// nolint:gochecknoglobals
	WatchModule = fx.Invoke(Watch)

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
