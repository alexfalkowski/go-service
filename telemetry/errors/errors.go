package errors

import (
	"log/slog"

	"github.com/alexfalkowski/go-service/v2/os"
)

// Register installs handler as the global OpenTelemetry error handler.
//
// This function forwards to [go.opentelemetry.io/otel.SetErrorHandler]. The OpenTelemetry error
// handler is process-wide; the last handler registered wins.
//
// Register is typically invoked once during service startup (for example via an
// Fx module) so that OpenTelemetry SDK/internal errors (exporter failures,
// dropped data warnings, etc.) are written to the handler's local logger.
//
// If handler is nil, Register leaves the current global OpenTelemetry error
// handler unchanged.
func Register(handler *Handler) {
	if handler == nil {
		return
	}

	SetHandler(handler)
}

// NewHandler constructs a Handler that logs OpenTelemetry internal errors to an
// independent JSON logger on stdout.
//
// The logger does not use the configured application logger or process-wide
// slog default, so exporter errors cannot feed back into a failing OTLP logger.
func NewHandler() *Handler {
	return &Handler{
		logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}

// Handler routes OpenTelemetry SDK/internal errors to a local JSON logger.
//
// Handler is intended to be registered via Register so that OpenTelemetry errors
// are visible without depending on the application logging pipeline.
type Handler struct {
	logger *slog.Logger
}

// Handle logs an OpenTelemetry internal error.
//
// This method is called by the OpenTelemetry SDK when it encounters an internal
// error. It logs at error level using the handler's local logger, attaching the
// error under the "error" key.
//
// Handle is nil-safe. If the receiver or its logger is nil, the error is ignored.
func (e *Handler) Handle(err error) {
	if e == nil || e.logger == nil {
		return
	}

	e.logger.Error("telemetry: global error", slog.Any("error", err))
}
