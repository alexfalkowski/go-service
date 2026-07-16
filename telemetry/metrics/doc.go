// Package metrics wires OpenTelemetry metrics into a go-service application.
//
// # Overview
//
// This package provides constructors and Fx/Dig wiring for OpenTelemetry metrics:
//
//   - NewReader constructs an OpenTelemetry SDK metric reader/exporter based on Config.
//   - NewMeterProvider constructs an SDK MeterProvider, installs it globally via
//     otel.SetMeterProvider, and starts runtime metric instrumentation. Host metric
//     instrumentation is temporarily disabled until
//     https://github.com/shirou/gopsutil/issues/2115 is fixed.
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
//   - If *[Config] is nil, metrics are treated as disabled.
//   - If the configured Reader is nil, metrics are treated as disabled even when
//     Config is non-nil (this allows DI to short-circuit metrics when reader
//     construction fails or is intentionally omitted).
//
// When disabled, NewReader returns (nil, nil) and NewMeterProvider installs and
// returns this package's noop provider. IsEnabled reports whether the current global
// provider is not that noop provider.
//
// # Exporters / readers
//
// [Config.Kind] selects the reader/exporter implementation. This package typically supports:
//
//   - "otlp": uses an OTLP metrics exporter and a periodic reader.
//   - "prometheus": uses the Prometheus exporter/reader with a namespace derived from
//     the service name.
//
// If [Config.Kind] is unknown, NewReader returns [ErrNotFound].
//
// For "otlp", [Config.Interval] and [Config.Timeout] configure the periodic
// reader export cadence. Zero values preserve the OpenTelemetry SDK defaults.
// [Config.Protocol] selects the OTLP transport protocol. The empty value uses
// OTLP/HTTP. Set "grpc" to use OTLP/gRPC with a host:port endpoint.
//
// For "prometheus", [Config.Prometheus] optionally shapes exporter output by
// dropping unit/counter suffixes, the target_info metric, or the scope-info
// labels. A nil value preserves the default OpenTelemetry-conventional output.
//
// # OTLP endpoint security
//
// When [Config.Headers] is non-empty, non-loopback "http://" OTLP endpoints are
// rejected to avoid sending credential-bearing headers over cleartext transport.
// Use "https://" for external collectors. Local development collectors on
// "localhost" or loopback IP addresses may use "http://".
// Header-bearing remote OTLP/gRPC endpoints require [Config.TLS]; loopback gRPC
// endpoints may still use cleartext.
//
// # Global provider installation
//
// When enabled, NewMeterProvider installs the constructed provider as the process-wide
// default via [go.opentelemetry.io/otel.SetMeterProvider]. Instrumentation that uses the global provider
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
// Additional configured resource attributes are also attached, with the fixed
// identity attributes taking precedence on duplicate keys.
//
// # Lifecycle behavior
//
// The reader and provider register lifecycle hooks to shut down cleanly on application
// stop. Provider shutdown errors are intentionally ignored to avoid blocking other stop
// hooks.
package metrics
