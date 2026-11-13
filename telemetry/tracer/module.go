package tracer

import "github.com/alexfalkowski/go-signal/v2"

// Module wires tracer registration into [github.com/alexfalkowski/go-signal/v2].
//
// It registers [Register], which configures and installs a global OpenTelemetry
// TracerProvider when tracing is enabled.
var Module = signal.Module(
	signal.Register(Register),
)
