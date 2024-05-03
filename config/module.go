package config

import (
	"go.uber.org/fx"
)

var (
	// Module for fx.
	Module = fx.Options(
		ConfiguratorModule,
		ConfigModule,
	)

	// ConfiguratorModule for fx.
	ConfiguratorModule = fx.Provide(NewConfigurator)

	// ConfigModule for fx.
	ConfigModule = fx.Options(
		fx.Provide(environmentConfig), fx.Provide(debugConfig), fx.Provide(tokenConfig),
		fx.Provide(featureConfig), fx.Provide(hooksConfig),
		fx.Provide(ntpConfig), fx.Provide(ntsConfig),
		fx.Provide(pgConfig), fx.Provide(redisConfig), fx.Provide(ristrettoConfig),
		fx.Provide(loggerConfig), fx.Provide(tracerConfig), fx.Provide(metricsConfig),
		fx.Provide(grpcConfig), fx.Provide(httpConfig), fx.Provide(limiterConfig),
	)
)
