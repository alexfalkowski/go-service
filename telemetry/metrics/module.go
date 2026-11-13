package metrics

import "github.com/alexfalkowski/go-signal/v2"

// Module wires metrics construction into [github.com/alexfalkowski/go-signal/v2].
//
// It provides constructors for:
//
//   - [NewReader], which constructs an OpenTelemetry SDK metric reader/exporter based on *[Config].
//   - [NewMeterProvider], which installs a global MeterProvider (via [go.opentelemetry.io/otel.SetMeterProvider]) when enabled.
//   - [NewMeter], which returns a Meter scoped to the service name and instrumentation version.
//
// When metrics are disabled (*[Config] is nil), [NewReader] returns a nil reader and
// [NewMeterProvider] installs and returns the package noop provider.
var Module = signal.Module(
	signal.Constructor(NewReader),
	signal.Constructor(NewMeterProvider),
	signal.Constructor(NewMeter),
	signal.Register(registerMeterProvider),
)

func registerMeterProvider(_ MeterProvider) {}
