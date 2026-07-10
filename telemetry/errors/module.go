package errors

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires OpenTelemetry global error handling into [go.uber.org/fx].
//
// Including this module in an Fx application provides:
//
//   - NewHandler: constructs a *[Handler] that logs OpenTelemetry internal/SDK errors
//     through a private stdout logger which mirrors the configured logger format
//     (json, text, or tint). This sink is independent of the configured
//     application logger and its OTLP export pipeline.
//   - Register: installs that handler as the process-wide OpenTelemetry error
//     handler via otel.SetErrorHandler. If the constructed handler is nil,
//     Register leaves the current global handler unchanged.
//
// This surfaces OpenTelemetry exporter/SDK errors on stdout so that export
// failures cannot feed their own diagnostics back into a failing exporter.
//
// Note: the OpenTelemetry error handler is global; the last handler registered
// wins.
var Module = di.Module(
	di.Constructor(NewHandler),
	di.Register(Register),
)
