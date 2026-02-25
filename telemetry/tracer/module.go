package tracer

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires tracer registration into Fx.
//
// It registers `Register`, which configures and installs a global OpenTelemetry
// TracerProvider when tracing is enabled.
var Module = di.Module(
	di.Register(Register),
)
