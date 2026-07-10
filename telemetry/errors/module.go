package errors

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires OpenTelemetry global error handling into [go.uber.org/fx].
//
// Including this module in an Fx application provides:
//
//   - NewHandler: constructs a *[Handler] with its own local JSON logger.
//   - Register: installs that handler as the process-wide OpenTelemetry error
//     handler via otel.SetErrorHandler.
//
// This keeps OpenTelemetry exporter/SDK failures visible on local stdout without
// feeding them back into a failing configured OTLP logger.
//
// Note: the OpenTelemetry error handler is global; the last handler registered
// wins.
var Module = di.Module(
	di.Constructor(NewHandler),
	di.Register(Register),
)
