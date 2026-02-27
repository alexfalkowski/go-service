// Package header provides helpers for configuring telemetry exporter/request headers.
//
// Many telemetry exporters (for example OTLP exporters for logs, metrics, and traces)
// support attaching arbitrary headers to outbound HTTP requests. These headers are
// commonly used for:
//
//   - authentication (for example "Authorization: Bearer <token>"),
//   - routing/multi-tenant metadata (for example "X-Scope-OrgID"),
//   - or other collector-specific requirements.
//
// This package defines Map, a small convenience type for representing such headers
// in configuration, along with helpers for resolving secret values at runtime.
//
// # Source string convention and secret resolution
//
// Header values are often sensitive and should not be committed to configuration
// files. go-service supports a “source string” convention that allows a configured
// value to be read from an alternate source at runtime.
//
// Map.Secrets traverses the map and resolves each value using os.FS.ReadSource,
// which supports these forms:
//
//   - "env:NAME"    reads the value of environment variable NAME.
//   - "file:/path"  reads bytes from the file at /path (including path cleaning and trimming).
//   - otherwise     treats the value as a literal string.
//
// After resolving, the map is updated in place so each header value becomes the
// final literal value that exporters should send.
//
// Map.MustSecrets behaves like Map.Secrets but panics on any resolution error via
// runtime.Must. This is intended for strict startup paths where missing secrets
// should abort service startup.
//
// # Mutability and usage notes
//
// Map.Secrets and Map.MustSecrets mutate the map in place. If you need to preserve
// the original configured source strings, copy the map before resolving.
//
// This package does not itself attach headers to exporters; it provides a consistent
// configuration representation and resolution mechanism used by telemetry packages
// (for example telemetry/logger, telemetry/metrics, telemetry/tracer).
package header
