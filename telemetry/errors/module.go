package errors

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires OpenTelemetry global error handling into Fx.
//
// Including this module in an Fx application provides:
//
//   - NewHandler: constructs a *Handler that logs OpenTelemetry internal/SDK errors
//     through the go-service telemetry logger.
//   - Register: installs that handler as the process-wide OpenTelemetry error
//     handler via otel.SetErrorHandler.
//
// This ensures OpenTelemetry exporter/SDK errors are surfaced in service logs.
//
// Note: the OpenTelemetry error handler is global; the last handler registered
// wins.
var Module = di.Module(
	di.Constructor(NewHandler),
	di.Register(Register),
)
