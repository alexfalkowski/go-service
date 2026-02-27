// Package attributes provides small, stable helpers and aliases for OpenTelemetry
// semantic convention attributes used by go-service.
//
// This package is intentionally thin. It primarily re-exports selected identifiers
// from the OpenTelemetry semantic conventions package:
//
//	go.opentelemetry.io/otel/semconv/v1.39.0
//
// The goal is to centralize common resource/telemetry attributes that are used
// across go-service telemetry wiring (logging, metrics, tracing) without requiring
// every package to import semconv directly.
//
// # Resource attributes
//
// go-service telemetry providers commonly attach a Resource describing the running
// service instance. This package exposes helper constructors for common resource
// fields, such as:
//
//   - HostID
//   - ServiceName
//   - ServiceVersion
//   - DeploymentEnvironmentName
//
// These helpers return attribute.KeyValue values that can be passed to
// resource.NewWithAttributes.
//
// # Schema URL
//
// SchemaURL is re-exported so resource creation can consistently specify the
// semantic conventions schema used by the attributes.
//
// # Protocol/system identifiers
//
// Some instrumentation requires standard system identifiers (for example RPC system
// names). This package re-exports the gRPC RPC system name constant
// (RPCSystemNameGRPC) from semconv.
//
// # Notes
//
// This package does not define new semantics or attribute keys. It only re-exports
// upstream OpenTelemetry semantic convention identifiers behind go-service imports.
// For authoritative definitions and version-specific details, consult the upstream
// OpenTelemetry semantic conventions documentation for the version vendored by this
// repository.
package attributes
