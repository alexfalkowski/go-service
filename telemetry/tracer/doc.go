// Package tracer wires OpenTelemetry tracing into a go-service application.
//
// # What it does
//
// When enabled via Config, Register configures and installs a process-wide OpenTelemetry
// TracerProvider (via otel.SetTracerProvider) backed by an OTLP/HTTP exporter.
//
// The provider is configured with a Resource describing the running service instance,
// including standard service identity attributes (host ID, service name, service version,
// and deployment environment name).
//
// # Enablement model
//
// Tracing is enabled by presence: a nil *Config indicates tracing is disabled.
// When disabled, Register is a no-op.
//
// # Global provider installation
//
// Register installs the configured TracerProvider as the global provider for the process.
// Instrumentation that relies on the global provider (directly or indirectly) will create
// spans using this provider.
//
// This package does not configure propagation. Propagation is configured at the top-level
// telemetry package (telemetry.Register), which sets the global TextMapPropagator used for
// context extraction/injection on supported transports (HTTP/gRPC) when instrumentation
// uses the global propagator.
//
// # Lifecycle behavior
//
// Register appends hooks to the provided lifecycle:
//   - OnStart: starts the OTLP exporter.
//   - OnStop: shuts down the tracer provider and exporter.
//
// Shutdown errors are intentionally ignored to avoid blocking other lifecycle stop hooks.
//
// # Configuration
//
// The exporter request headers are provided by Config.Headers. Header values may be
// configured as go-service “source strings” (for example "env:NAME", "file:/path", or a
// literal value) and are resolved by telemetry/header.Map.Secrets or
// telemetry/header.Map.MustSecrets by the consumer that projects configuration before
// constructing exporters.
package tracer
