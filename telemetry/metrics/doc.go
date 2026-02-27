// Package metrics wires OpenTelemetry metrics into a go-service application.
//
// # Overview
//
// This package provides constructors and Fx/Dig wiring for OpenTelemetry metrics:
//
//   - NewReader constructs an OpenTelemetry SDK metric reader/exporter based on Config.
//   - NewMeterProvider constructs an SDK MeterProvider, installs it globally via
//     otel.SetMeterProvider, and starts runtime/host metric instrumentation.
//   - NewMeter returns a metric.Meter scoped to the service name and instrumentation
//     version and is nil-safe when the provider is nil.
//
// The intention is to make metrics wiring consistent across services while keeping
// implementation details (exporters/readers, resource attributes, lifecycle hooks)
// in one place.
//
// # Enablement model
//
// Metrics are enabled by presence:
//
//   - If *Config is nil, metrics are treated as disabled.
//   - If the configured Reader is nil, metrics are treated as disabled even when
//     Config is non-nil (this allows DI to short-circuit metrics when reader
//     construction fails or is intentionally omitted).
//
// When disabled, NewReader returns (nil, nil) and NewMeterProvider returns nil.
//
// # Exporters / readers
//
// Config.Kind selects the reader/exporter implementation. This package typically supports:
//
//   - "otlp": uses an OTLP/HTTP metrics exporter and a periodic reader.
//   - "prometheus": uses the Prometheus exporter/reader with a namespace derived from
//     the service name.
//
// If Config.Kind is unknown, NewReader returns ErrNotFound.
//
// # Global provider installation
//
// When enabled, NewMeterProvider installs the constructed provider as the process-wide
// default via otel.SetMeterProvider. Instrumentation that uses the global provider
// (directly or indirectly) will emit metrics through this provider.
//
// This package does not configure propagation. Propagation concerns apply primarily to
// tracing context (and are handled by the top-level telemetry package).
//
// # Resource attributes
//
// When enabled, the meter provider is configured with a resource that includes standard
// service identity attributes:
//
//   - host.id
//   - service.name
//   - service.version
//   - deployment.environment.name
//
// These values come from go-service env types provided via DI.
//
// # Lifecycle behavior
//
// The reader and provider register lifecycle hooks to shut down cleanly on application
// stop. Provider shutdown errors are intentionally ignored to avoid blocking other stop
// hooks.
package metrics
