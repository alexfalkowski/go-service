// Package errors wires OpenTelemetry global error handling into go-service logging.
//
// OpenTelemetry SDKs and instrumentations may emit internal errors (for example exporter
// failures, dropped data warnings, or other SDK/runtime issues). The OpenTelemetry API
// provides a global error handler hook (see [go.opentelemetry.io/otel.ErrorHandler]) that
// applications can set to control how these errors are surfaced.
//
// This package provides a go-service implementation of the OpenTelemetry error handler
// interface and helpers to register it.
//
// # Handler
//
// Handler implements the OpenTelemetry error handler interface with its own
// JSON logger on stdout. Errors are logged at error level using a consistent
// message and attribute key ("error").
//
// # Registration
//
// Register installs a provided Handler as the process-wide OpenTelemetry error handler by
// calling:
//
//	otel.SetErrorHandler(handler)
//
// This affects all OpenTelemetry components in the process that report errors via the
// global handler.
//
// If Register is called with nil, it is a no-op and the current global handler is
// preserved.
//
// # Dependency injection (Fx)
//
// This package also exports [Module], which wires:
//   - construction of [Handler] with [NewHandler], and
//   - registration with [Register]
//
// into an Fx application.
//
// Including [github.com/alexfalkowski/go-service/v2/telemetry/errors.Module] (or the top-level [github.com/alexfalkowski/go-service/v2/telemetry.Module]) wires this
// handler into your service. OpenTelemetry internal errors are written as JSON
// to local stdout independently of the configured logger. This prevents an OTLP
// exporter failure from feeding its diagnostic back into the same failed
// exporter.
//
// # Notes
//
// The OpenTelemetry error handler is global and should typically be configured once at
// startup. If you install multiple handlers, the last one set wins.
//
// This package only handles OpenTelemetry internal errors; it does not affect how spans,
// metrics, or logs are exported beyond ensuring SDK errors are visible in logs.
package errors
