// Package metrics wires up OpenTelemetry metrics for a go-service application.
//
// # Overview
//
// This package provides constructors and Fx/Dig wiring for OpenTelemetry metrics:
//
//   - `NewReader` constructs an `sdk/metric.Reader` based on `Config`.
//   - `NewMeterProvider` constructs and installs a global `sdk/metric.MeterProvider`
//     (via `otel.SetMeterProvider`) and starts runtime/host metric instrumentation.
//   - `NewMeter` (from `metrics.go`) returns a `metric.Meter` scoped to the service name
//     and version, and is nil-safe when the provider is nil.
//
// # Exporters / readers
//
// `Config.Kind` selects the reader/exporter implementation:
//
//   - "otlp": uses an OTLP/HTTP metrics exporter and a periodic reader.
//   - "prometheus": uses the Prometheus exporter/reader with a namespace derived from
//     the service name.
//
// If `Config.Kind` is unknown, `NewReader` returns `ErrNotFound`. If metrics are disabled
// (`*Config` is nil), `NewReader` returns a nil reader and `NewMeterProvider` returns nil.
//
// # Resource attributes
//
// When enabled, the meter provider is configured with a resource that includes standard
// service identity attributes: host ID, service name, service version, and deployment
// environment name.
//
// # Lifecycle behavior
//
// The reader and provider register lifecycle hooks to shut down cleanly on application stop.
// Provider shutdown errors are intentionally ignored to avoid blocking other stop hooks.
package metrics
