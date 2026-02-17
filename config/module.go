package config

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires config decoding, validation, and feature sub-configs into Fx.
var Module = di.Module(
	di.Constructor(NewValidator), di.Constructor(NewDecoder), di.Constructor(NewConfig[Config]),
	di.Constructor(cryptoAESConfig), di.Constructor(cryptoED25519Config),
	di.Constructor(cryptoHMACConfig), di.Constructor(cryptoRSAConfig), di.Constructor(cryptoSSHConfig),
	di.Constructor(environmentConfig), di.Constructor(cacheConfig),
	di.Constructor(debugConfig), di.Constructor(idConfig), di.Constructor(timeConfig),
	di.Constructor(pgConfig), di.Constructor(featureConfig), di.Constructor(hooksConfig),
	di.Constructor(loggerConfig), di.Constructor(tracerConfig), di.Constructor(metricsConfig),
	di.Constructor(grpcConfig), di.Constructor(httpConfig),
)
