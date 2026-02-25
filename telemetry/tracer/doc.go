// Package tracer wires up OpenTelemetry tracing for a go-service application.
//
// # What it does
//
// When enabled via `Config`, `Register` configures and installs a global OpenTelemetry
// `trace.TracerProvider` (via `otel.SetTracerProvider`) backed by an OTLP/HTTP exporter.
// The provider is configured with a resource describing the running service (host ID,
// service name, version, and deployment environment).
//
// # Lifecycle behavior
//
// `Register` appends hooks to the provided lifecycle:
//   - OnStart: starts the OTLP exporter
//   - OnStop: shuts down the tracer provider and exporter
//
// Shutdown errors are intentionally ignored to avoid blocking other lifecycle stop hooks.
//
// # Configuration
//
// The exporter request headers are provided by `Config.Headers`. Header values may be
// configured as "source strings" (for example `env:NAME`, `file:/path`, or a literal value)
// and are resolved by the telemetry header helpers used by `Config` consumers.
package tracer
