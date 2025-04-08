package config

import "go.uber.org/fx"

// Module for fx.
var Module = fx.Options(
	fx.Provide(NewConfig[Config]),
	fx.Provide(cryptoAESConfig), fx.Provide(cryptoED25519Config), fx.Provide(cryptoHMACConfig),
	fx.Provide(cryptoRSAConfig), fx.Provide(cryptoSSHConfig),
	fx.Provide(tokenConfig), fx.Provide(tokenJWTConfig),
	fx.Provide(tokenPasetoConfig), fx.Provide(tokenSSHConfig),
	fx.Provide(environmentConfig), fx.Provide(cacheConfig),
	fx.Provide(debugConfig), fx.Provide(idConfig), fx.Provide(timeConfig),
	fx.Provide(pgConfig), fx.Provide(featureConfig), fx.Provide(hooksConfig),
	fx.Provide(loggerConfig), fx.Provide(tracerConfig), fx.Provide(metricsConfig),
	fx.Provide(grpcConfig), fx.Provide(httpConfig), fx.Provide(limiterConfig),
)
