package metrics

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires metrics construction into Fx.
//
// It provides constructors for:
//
//   - `NewReader`, which constructs an OpenTelemetry SDK metric reader/exporter based on `*Config`.
//   - `NewMeterProvider`, which installs a global MeterProvider (via `otel.SetMeterProvider`) when enabled.
//   - `NewMeter`, which returns a Meter scoped to the service name and instrumentation version.
//
// When metrics are disabled (`*Config` is nil), `NewReader` returns a nil reader and
// `NewMeterProvider` returns a nil provider.
var Module = di.Module(
	di.Constructor(NewReader),
	di.Constructor(NewMeterProvider),
	di.Constructor(NewMeter),
)
