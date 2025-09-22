package test

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/telemetry/errors"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
)

// NewLogger for test.
func NewLogger(lc di.Lifecycle, config *logger.Config) *logger.Logger {
	logger, err := logger.NewLogger(logger.LoggerParams{Lifecycle: lc, Config: config, Version: Version})
	runtime.Must(err)

	return logger
}

// WithWorldLogger for test.
func WithWorldLogger(logger *logger.Logger) WorldOption {
	return worldOptionFunc(func(o *worldOpts) {
		o.logger = logger
	})
}

// WithWorldLoggerConfig for test.
func WithWorldLoggerConfig(config string) WorldOption {
	return worldOptionFunc(func(o *worldOpts) {
		o.loggerConfig = config
	})
}

func (w *World) registerTelemetry() {
	errors.Register(errors.NewHandler(w.Logger))
}

func createLogger(lc di.Lifecycle, os *worldOpts) *logger.Logger {
	if os.logger != nil {
		return os.logger
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
