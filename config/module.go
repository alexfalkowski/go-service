package config

import (
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	fx.Provide(NewConfig),
	fx.Provide(environmentConfig), fx.Provide(debugConfig), fx.Provide(tokenConfig),
	fx.Provide(featureConfig), fx.Provide(hooksConfig), fx.Provide(timeConfig),
	fx.Provide(pgConfig), fx.Provide(redisConfig), fx.Provide(ristrettoConfig),
	fx.Provide(loggerConfig), fx.Provide(tracerConfig), fx.Provide(metricsConfig),
	fx.Provide(grpcConfig), fx.Provide(httpConfig), fx.Provide(limiterConfig),
)
