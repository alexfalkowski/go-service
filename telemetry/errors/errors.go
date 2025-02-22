package errors

import (
	"github.com/alexfalkowski/go-service/telemetry/logger"
	"go.opentelemetry.io/otel"
)

// Register the error handler.
func Register(handler *Handler) {
	otel.SetErrorHandler(handler)
}

// NewHandler creates a new error handler.
func NewHandler(logger *logger.Logger) *Handler {
	return &Handler{logger: logger}
}

// Handler is the error handler.
type Handler struct {
	logger *logger.Logger
}

// Handle handles the error.
func (e *Handler) Handle(err error) {
	e.logger.Error("telemetry: global error", logger.Error(err))
}
