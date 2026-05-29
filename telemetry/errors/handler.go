package errors

import "go.opentelemetry.io/otel"

// ErrorHandler is an alias for otel.ErrorHandler.
type ErrorHandler = otel.ErrorHandler

// GetHandler returns the global OpenTelemetry error handler.
func GetHandler() ErrorHandler {
	return otel.GetErrorHandler()
}

// SetHandler installs the global OpenTelemetry error handler.
func SetHandler(handler ErrorHandler) {
	otel.SetErrorHandler(handler)
}
