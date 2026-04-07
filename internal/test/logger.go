package test

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/telemetry/errors"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
)

// NewLogger constructs a test logger bound to the supplied lifecycle and logger config.
func NewLogger(lc di.Lifecycle, config *logger.Config) (*logger.Logger, error) {
	return logger.NewLogger(logger.LoggerParams{Lifecycle: lc, Config: config, Version: Version})
}

func (w *World) registerTelemetry() {
	errors.Register(errors.NewHandler(w.Logger))
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
