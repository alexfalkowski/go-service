package errors

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires OpenTelemetry global error handling into Fx.
//
// Including this module in an Fx application provides:
//
//   - NewHandler: constructs a *Handler that logs OpenTelemetry internal/SDK errors
//     through the go-service telemetry logger when logging is enabled. If no
//     go-service logger is available, NewHandler returns nil.
//   - Register: installs that handler as the process-wide OpenTelemetry error
//     handler via otel.SetErrorHandler. If the constructed handler is nil,
//     Register leaves the current global handler unchanged.
//
// This surfaces OpenTelemetry exporter/SDK errors in service logs when logging
// is configured, while preserving the OpenTelemetry default error handling when
// it is not.
//
// Note: the OpenTelemetry error handler is global; the last handler registered
// wins.
var Module = di.Module(
	di.Constructor(NewHandler),
	di.Register(Register),
)
