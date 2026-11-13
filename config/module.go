package config

import (
	"github.com/alexfalkowski/go-service/v3/config/validate"
	"github.com/alexfalkowski/go-signal/v2"
)

// Module wires the configuration subsystem into [github.com/alexfalkowski/go-signal/v2].
//
// It provides the core configuration components:
//   - [validate.Validator] ([validate.Module]) used to validate decoded config structs. The
//     `config_size`, `duration_second_precision`, `otlp_cadence`, and `otlp_batch_config` validation
//     tags are contributed by `bytes`, `time`, and `telemetry/otlp`'s own `Module`s, composed via
//     [github.com/alexfalkowski/go-service/v3/module.Library] rather than here, so this package does
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
// [github.com/alexfalkowski/go-service/v3/telemetry/header.Map.MustSecrets]. Missing environment variables or
// unreadable files therefore fail fast during startup projection rather than later exporter construction.
var Module = signal.Module(
	validate.Module,
	signal.Constructor(NewDecoder), signal.Constructor(NewConfig[Config]),
	signal.Constructor(cryptoAESConfig), signal.Constructor(cryptoED25519Config),
	signal.Constructor(cryptoHMACConfig), signal.Constructor(cryptoRSAConfig),
	signal.Constructor(environmentConfig), signal.Constructor(cacheConfig),
	signal.Constructor(debugConfig), signal.Constructor(idConfig), signal.Constructor(timeConfig),
	signal.Constructor(pgConfig), signal.Constructor(featureConfig), signal.Constructor(hooksConfig),
	signal.Constructor(attributeMap), signal.Constructor(metadataLimit),
	signal.Constructor(loggerConfig), signal.Constructor(tracerConfig), signal.Constructor(metricsConfig),
	signal.Constructor(propagationConfig),
	signal.Constructor(accessConfig), signal.Constructor(grpcConfig), signal.Constructor(httpConfig),
)
