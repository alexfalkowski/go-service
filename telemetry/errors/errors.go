package errors

import (
	"log/slog"

	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
)

// Register installs handler as the global OpenTelemetry error handler.
//
// This function forwards to [go.opentelemetry.io/otel.SetErrorHandler]. The OpenTelemetry error
// handler is process-wide; the last handler registered wins.
//
// Register is typically invoked once during service startup (for example via an
// Fx module) so that OpenTelemetry SDK/internal errors (exporter failures,
// dropped data warnings, etc.) are surfaced through the handler's independent
// diagnostic sink.
//
// If handler is nil, Register leaves the current global OpenTelemetry error
// handler unchanged.
func Register(handler *Handler) {
	if handler == nil {
		return
	}

	SetHandler(handler)
}

// NewHandler constructs a Handler that logs OpenTelemetry internal errors.
//
// The returned Handler implements the OpenTelemetry error handler interface and
// owns a private stdout logger built by [logger.NewDiagnosticLogger]. That sink
// mirrors the configured logger format (json, text, or tint) while remaining
// independent of the configured application logger and its OTLP export pipeline,
// so OpenTelemetry export failures cannot feed their own diagnostics back into a
// failing exporter. A nil cfg, the "otlp" kind, or an unknown kind fall back to
// JSON on stdout.
func NewHandler(cfg *logger.Config) *Handler {
	return &Handler{logger: logger.NewDiagnosticLogger(cfg)}
}

// Handler routes OpenTelemetry SDK/internal errors into a private diagnostic logger.
//
// Handler is intended to be registered via Register so that OpenTelemetry errors
// are visible on stdout independently of the configured application logger. It
// logs a consistent message and includes a standardized "error" attribute
// produced by [logger.Error].
type Handler struct {
	logger *slog.Logger
}

// Handle logs an OpenTelemetry internal error.
//
// This method is called by the OpenTelemetry SDK when it encounters an internal
// error. It logs at error level using the handler's private logger, attaching
// the error under the "error" key.
//
// Handle is nil-safe. If the receiver or its logger is nil, the error is ignored.
func (e *Handler) Handle(err error) {
	if e == nil || e.logger == nil {
		return
	}

	e.logger.Error("telemetry: global error", logger.Error(err))
}
