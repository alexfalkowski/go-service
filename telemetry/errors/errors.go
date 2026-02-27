package errors

import (
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"go.opentelemetry.io/otel"
)

// Register installs handler as the global OpenTelemetry error handler.
//
// This function forwards to otel.SetErrorHandler(handler). The OpenTelemetry error
// handler is process-wide; the last handler registered wins.
//
// Register is typically invoked once during service startup (for example via an
// Fx module) so that OpenTelemetry SDK/internal errors (exporter failures,
// dropped data warnings, etc.) are routed into application logging.
func Register(handler *Handler) {
	otel.SetErrorHandler(handler)
}

// NewHandler constructs a Handler that logs OpenTelemetry internal errors.
//
// The returned Handler implements the OpenTelemetry error handler interface and
// writes errors using the provided go-service *logger.Logger.
func NewHandler(logger *logger.Logger) *Handler {
	return &Handler{logger: logger}
}

// Handler routes OpenTelemetry SDK/internal errors into a go-service logger.
//
// Handler is intended to be registered via Register so that OpenTelemetry errors
// are visible in service logs. It logs a consistent message and includes a
// standardized "error" attribute produced by logger.Error.
type Handler struct {
	logger *logger.Logger
}

// Handle logs an OpenTelemetry internal error.
//
// This method is called by the OpenTelemetry SDK when it encounters an internal
// error. It logs at error level using the go-service logger, attaching the error
// under the "error" key.
func (e *Handler) Handle(err error) {
	e.logger.Error("telemetry: global error", logger.Error(err))
}
