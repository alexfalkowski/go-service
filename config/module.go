package config

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires the configuration subsystem into Fx/Dig.
//
// It provides the core configuration components:
//   - Validator (NewValidator) used to validate decoded config structs.
//   - Decoder (NewDecoder) that dispatches config loading based on the "-i" flag.
//   - *config.Config (NewConfig[Config]) as the standard top-level configuration.
//
// It also provides a set of small "projection" constructors that extract commonly-used sub-configs from
// the top-level Config. These projections are used by other modules so they can depend directly on the
// sub-config they need, without having to understand the full top-level shape.
//
// Projection constructors are nil-safe by convention: if the parent feature is disabled (i.e. the parent
// config pointer is nil), the projection returns nil so downstream modules treat that subsystem as disabled.
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
