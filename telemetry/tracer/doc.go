// Package tracer wires OpenTelemetry tracing into a go-service application.
//
// # What it does
//
// When configured with kind "otlp", Register configures and installs a process-wide
// OpenTelemetry TracerProvider (via otel.SetTracerProvider) backed by an OTLP
// exporter.
//
// The provider is configured with a Resource describing the running service instance,
// including standard service identity attributes (host ID, service name, service version,
// and deployment environment name) plus any configured resource attributes. Fixed
// service identity attributes take precedence on duplicate keys.
//
// # Enablement model
//
// Tracing is enabled by kind: a nil *[Config] or an empty [Config.Kind] indicates tracing
// is not configured. When disabled, Register installs this package's noop provider.
// IsEnabled reports whether the current global provider is not that noop provider.
//
// # Global provider installation
//
// Register installs the configured TracerProvider as the global provider for the process.
// Instrumentation that relies on the global provider (directly or indirectly) will create
// spans using this provider.
//
// # Metadata on spans
//
// The provider installs a span processor that copies request/service metadata
// from the context onto each span as it starts (see [Meta]), so spans created
// by any instrumentation (server, database, cache, ...) carry the same context
// used to correlate them with logs. The server/root span is created before that
// metadata is available, so the transport metadata middleware stamps it
// directly.
//
// This package does not configure propagation. Propagation is configured at the top-level
// telemetry package ([github.com/alexfalkowski/go-service/v2/telemetry.RegisterPropagation]), which sets the global
// TextMapPropagator used for context extraction/injection on supported transports
// (HTTP/gRPC) when instrumentation uses the global propagator.
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
// The supported kind is "otlp". Unknown non-empty kinds cause Register to return
// ErrNotFound.
//
// [Config.Protocol] selects the OTLP transport protocol. The empty value uses
// OTLP/HTTP. Set "grpc" to use OTLP/gRPC with a host:port endpoint.
//
// [Config.Sampler] optionally configures head sampling. When omitted, Register
// preserves the OpenTelemetry SDK default sampler and SDK sampler environment
// handling. When set, it overrides those defaults.
//
// The exporter request headers are provided by [Config.Headers]. Header values may be
// configured as go-service "source strings" (for example "env:NAME", "file:/path", or a
// literal value) and are resolved by [github.com/alexfalkowski/go-service/v2/telemetry/header.Map.Secrets] or
// [github.com/alexfalkowski/go-service/v2/telemetry/header.Map.MustSecrets] by the consumer that projects configuration before
// constructing exporters.
//
// When [Config.Headers] is non-empty, non-loopback "http://" OTLP endpoints are
// rejected to avoid sending credential-bearing headers over cleartext transport.
// Use "https://" for external collectors. Local development collectors on
// "localhost" or loopback IP addresses may use "http://".
// Header-bearing remote OTLP/gRPC endpoints require [Config.TLS]; loopback gRPC
// endpoints may still use cleartext.
package tracer
