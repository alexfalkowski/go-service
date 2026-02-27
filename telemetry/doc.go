// Package telemetry provides OpenTelemetry-based telemetry configuration and wiring for go-service.
//
// The telemetry subsystem in this repository is organized as a set of focused packages:
//
//   - telemetry/logger: application/system logging (slog-based) with optional exporters (for example OTLP).
//   - telemetry/metrics: metrics SDK wiring and exporters (for example OTLP or Prometheus), plus helpers
//     for obtaining a Meter from a configured provider.
//   - telemetry/tracer: tracing SDK wiring and OTLP exporter configuration.
//   - telemetry/errors: OpenTelemetry error handler wiring so OTel SDK/internal errors are logged.
//   - telemetry/header: helpers for exporter/request headers, including secret resolution via the
//     go-service “source string” convention (env:/file:/literal).
//   - telemetry/attributes: aliases/helpers for OpenTelemetry semantic convention attributes.
//
// This top-level package primarily provides:
//
//   - Config: a single configuration root that embeds logging/metrics/tracing configuration.
//   - Register: a small initialization hook that configures global OpenTelemetry propagation.
//   - Module: an Fx module that composes telemetry submodules and applies Register.
//
// # Configuration
//
// Config is a convenience root used by services to configure telemetry in one place. It contains pointers
// to per-signal configs:
//
//   - Logger (*telemetry/logger.Config)
//   - Metrics (*telemetry/metrics.Config)
//   - Tracer (*telemetry/tracer.Config)
//
// A nil Config typically means “telemetry disabled” at the top level, while subpackages may also support
// their own enable/disable semantics (for example nil config or empty kind).
//
// # Global propagation (Register)
//
// Register configures the global OpenTelemetry TextMapPropagator to a composite propagator containing:
//
//   - W3C Trace Context (propagation.TraceContext)
//   - W3C Baggage (propagation.Baggage)
//
// This affects context extraction/injection for supported transports (HTTP/gRPC) when instrumentation uses
// the global propagator.
//
// Register is intended to be called once during startup (for example via Module).
//
// # Dependency injection (Module)
//
// Module is an Fx module that wires the telemetry submodules into an application. In particular it
// composes the logger/metrics/tracer/error-handler modules and registers Register so propagation is
// configured as part of application startup.
//
// Module does not itself create spans/metrics/log records; it wires providers/exporters and global
// configuration so instrumentation elsewhere in your service can emit telemetry.
//
// # Notes
//
// Many implementations in the telemetry subtree are thin adapters around upstream OpenTelemetry SDK and
// exporter packages. For exact semantics (for example exporter behavior, supported options, and version-
// specific details), consult the upstream documentation for the versions vendored by this repository.
package telemetry
