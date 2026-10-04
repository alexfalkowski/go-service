package config

import (
	"github.com/alexfalkowski/go-service/v2/config/validate"
	"github.com/alexfalkowski/go-service/v2/di"
)

// Module wires the configuration subsystem into [go.uber.org/fx]/[go.uber.org/dig].
//
// It provides the core configuration components:
//   - [validate.Validator] ([validate.Module]) used to validate decoded config structs. The
//     `config_size`, `duration_second_precision`, `otlp_cadence`, and `otlp_batch_config` validation
//     tags are contributed by `bytes`, `time`, and `telemetry/otlp`'s own `Module`s, composed via
//     [github.com/alexfalkowski/go-service/v2/module.Library] rather than here, so this package does
//     not need to know about their rules by name.
//   - Decoder (NewDecoder) that dispatches config loading based on the "-config" / "-c" flag.
//   - *[Config] ([NewConfig]) as the standard top-level configuration.
//
// It also provides a set of small "projection" constructors that extract commonly-used sub-configs from
// the top-level Config. These projections are used by other modules so they can depend directly on the
// sub-config they need, without having to understand the full top-level shape.
//
// Projection constructors are nil-safe by convention: if the parent feature is disabled (i.e. the parent
// config pointer is nil), the projection returns nil so downstream modules treat that subsystem as disabled.
//
// Telemetry logger, metrics, and tracer projections also resolve configured header source strings by calling
// [github.com/alexfalkowski/go-service/v2/telemetry/header.Map.MustSecrets]. Missing environment variables or
// unreadable files therefore fail fast during startup projection rather than later exporter construction.
var Module = di.Module(
	validate.Module,
	di.Constructor(NewDecoder), di.Constructor(NewConfig[Config]),
	di.Constructor(cryptoAESConfig), di.Constructor(cryptoED25519Config),
	di.Constructor(cryptoHMACConfig), di.Constructor(cryptoRSAConfig),
	di.Constructor(environmentConfig), di.Constructor(cacheConfig),
	di.Constructor(debugConfig), di.Constructor(idConfig), di.Constructor(timeConfig),
	di.Constructor(pgConfig), di.Constructor(featureConfig), di.Constructor(hooksConfig),
	di.Constructor(attributeMap), di.Constructor(metadataLimit),
	di.Constructor(loggerConfig), di.Constructor(tracerConfig), di.Constructor(metricsConfig),
	di.Constructor(propagationConfig),
	di.Constructor(accessConfig), di.Constructor(grpcConfig), di.Constructor(httpConfig),
)
