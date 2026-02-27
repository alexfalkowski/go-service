// Package errors wires OpenTelemetry global error handling into go-service logging.
//
// OpenTelemetry SDKs and instrumentations may emit internal errors (for example exporter
// failures, dropped data warnings, or other SDK/runtime issues). The OpenTelemetry API
// provides a global error handler hook (see go.opentelemetry.io/otel.ErrorHandler) that
// applications can set to control how these errors are surfaced.
//
// This package provides a go-service implementation of the OpenTelemetry error handler
// interface and helpers to register it.
//
// # Handler
//
// Handler implements the OpenTelemetry error handler interface by logging errors through
// a go-service `*telemetry/logger.Logger`. Errors are logged at error level using a
// consistent message and attribute key ("error").
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
// # Dependency injection (Fx)
//
// This package also exports `Module`, which wires:
//   - construction of the Handler (NewHandler), and
//   - registration of the handler (Register)
//
// into an Fx application.
//
// Including `telemetry/errors.Module` (or the top-level `telemetry.Module`) ensures that
// OpenTelemetry internal errors are routed into your service logging.
//
// # Notes
//
// The OpenTelemetry error handler is global and should typically be configured once at
// startup. If you install multiple handlers, the last one set wins.
//
// This package only handles OpenTelemetry internal errors; it does not affect how spans,
// metrics, or logs are exported beyond ensuring SDK errors are visible in logs.
package errors
