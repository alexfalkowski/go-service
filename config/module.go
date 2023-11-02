package config

import (
	"github.com/alexfalkowski/go-service/marshaller"
	"go.uber.org/fx"
)

var (
	// Module for fx.
	Module = fx.Options(
		ConfiguratorModule,
		ConfigModule,
		marshaller.Module,
	)

	// ConfiguratorModule for fx.
	ConfiguratorModule = fx.Provide(NewConfigurator)

	// ConfigModule for fx.
	ConfigModule = fx.Options(
		fx.Provide(environmentConfig),
		fx.Provide(redisConfig), fx.Provide(ristrettoConfig),
		fx.Provide(oauthConfig),
		fx.Provide(pgConfig),
		fx.Provide(loggerConfig), fx.Provide(tracerConfig),
		fx.Provide(transportConfig),
		fx.Provide(grpcConfig), fx.Provide(httpConfig), fx.Provide(nsqConfig),
	)
)
