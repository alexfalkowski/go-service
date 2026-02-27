package feature

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires OpenFeature into Fx/Dig.
//
// It provides:
//   - an OpenFeature client constructor (NewClient), and
//   - lifecycle registration (Register) that optionally installs a FeatureProvider.
//
// # Optional provider behavior
//
// Register is designed to be safe to include even when no OpenFeature provider is supplied by the
// consuming service. The FeatureProvider dependency is optional; if it is not present in the DI graph,
// Register becomes a no-op and OpenFeature uses its default provider semantics.
//
// When a provider is present, Register appends lifecycle hooks to set the provider during startup and
// shut down the OpenFeature SDK during stop. If a metrics provider is available, it also installs
// OpenTelemetry hooks for metrics and traces.
var Module = di.Module(
	di.Constructor(NewClient),
	di.Register(Register),
)
