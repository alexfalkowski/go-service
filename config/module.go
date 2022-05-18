package config

import (
	"go.uber.org/fx"
)

var (
	// Module for fx.
	Module = fx.Options(ConfiguratorModule, UnmarshalModule, ConfigModule, WatchModule)

	// ConfiguratorModule for fx.
	ConfiguratorModule = fx.Provide(NewConfigurator)

	// UnmarshalModule for fx.
	UnmarshalModule = fx.Invoke(UnmarshalFromFile)

	// WatchModule for fx.
	WatchModule = fx.Invoke(Watch)

	// ConfigModule for fx.
	ConfigModule = fx.Options(
		fx.Provide(redisConfig), fx.Provide(ristrettoConfig),
		fx.Provide(auth0Config),
		fx.Provide(pgConfig),
		fx.Provide(opentracingConfig),
		fx.Provide(transportConfig),
		fx.Provide(grpcConfig),
		fx.Provide(httpConfig),
		fx.Provide(nsqConfig),
	)
)
