package config

import "github.com/alexfalkowski/go-service/v2/di"

// Module for fx.
var Module = di.Module(
	di.Constructor(NewValidator), di.Constructor(NewDecoder), di.Constructor(NewConfig[Config]),
	di.Constructor(cryptoAESConfig), di.Constructor(cryptoED25519Config),
	di.Constructor(cryptoHMACConfig), di.Constructor(cryptoRSAConfig), di.Constructor(cryptoSSHConfig),
	di.Constructor(tokenConfig), di.Constructor(tokenAccessConfig), di.Constructor(tokenJWTConfig),
	di.Constructor(tokenPasetoConfig), di.Constructor(tokenSSHConfig),
	di.Constructor(environmentConfig), di.Constructor(cacheConfig),
	di.Constructor(debugConfig), di.Constructor(idConfig), di.Constructor(timeConfig),
	di.Constructor(pgConfig), di.Constructor(featureConfig), di.Constructor(hooksConfig),
	di.Constructor(loggerConfig), di.Constructor(tracerConfig), di.Constructor(metricsConfig),
	di.Constructor(grpcConfig), di.Constructor(httpConfig), di.Constructor(limiterConfig),
)
