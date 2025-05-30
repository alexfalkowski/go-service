package test

import (
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"go.uber.org/fx"
)

// NewLogger for test.
func NewLogger(lc fx.Lifecycle, config *logger.Config) *logger.Logger {
	logger, err := logger.NewLogger(logger.Params{Lifecycle: lc, Config: config, Version: Version})
	runtime.Must(err)

	return logger
}

// WithWorldLogger for test.
func WithWorldLogger(logger *logger.Logger) WorldOption {
	return worldOptionFunc(func(o *worldOpts) {
		o.logger = logger
	})
}

// WithWorldLogger for test.
func WithWorldLoggerConfig(config string) WorldOption {
	return worldOptionFunc(func(o *worldOpts) {
		o.loggerConfig = config
	})
}

func createLogger(lc fx.Lifecycle, os *worldOpts) *logger.Logger {
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
