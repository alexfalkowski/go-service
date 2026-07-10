package test

import (
	"log/slog"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/telemetry/errors"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
)

// CapturedRecord is a slog record captured by CaptureHandler.
type CapturedRecord struct {
	Attrs   map[string]slog.Value
	Message string
	Level   slog.Level
}

// CaptureHandler records slog records for tests.
type CaptureHandler struct {
	Records []CapturedRecord
}

// Enabled reports whether the handler accepts records at level.
func (h *CaptureHandler) Enabled(context.Context, slog.Level) bool {
	return true
}

// Handle records a slog record.
func (h *CaptureHandler) Handle(_ context.Context, record slog.Record) error {
	attrs := make(map[string]slog.Value)
	record.Attrs(func(attr slog.Attr) bool {
		if attr.Key != "" {
			attrs[attr.Key] = attr.Value.Resolve()
		}
		return true
	})
	h.Records = append(h.Records, CapturedRecord{
		Level:   record.Level,
		Message: record.Message,
		Attrs:   attrs,
	})
	return nil
}

// WithAttrs returns a handler with attrs.
func (h *CaptureHandler) WithAttrs([]slog.Attr) slog.Handler {
	return h
}

// WithGroup returns a handler with group.
func (h *CaptureHandler) WithGroup(string) slog.Handler {
	return h
}

// NewLogger constructs a test logger bound to the supplied lifecycle and logger config.
func NewLogger(lc di.Lifecycle, config *logger.Config) (*logger.Logger, error) {
	return logger.NewLogger(logger.LoggerParams{Lifecycle: lc, Config: config, Version: Version})
}

func (w *World) registerTelemetry() {
	errors.Register(errors.NewHandler(nil))
}

func createLogger(lc di.Lifecycle, os *worldOpts) (*logger.Logger, error) {
	if os.logger != nil {
		return os.logger, nil
	}

	var config *logger.Config
	switch os.loggerConfig {
	case "json":
		config = NewJSONLoggerConfig()
	case "text":
		config = NewTextLoggerConfig()
	case "tint":
		config = NewTintLoggerConfig()
	case "otlp":
		config = NewOTLPLoggerConfig()
	default:
		config = NewOTLPLoggerConfig()
	}

	return NewLogger(lc, config)
}
