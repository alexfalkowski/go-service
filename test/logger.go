package test

import (
	"github.com/alexfalkowski/go-service/runtime"
	logger "github.com/alexfalkowski/go-service/telemetry/logger/zap"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// NewLogger for test.
func NewLogger(lc fx.Lifecycle) *zap.Logger {
	c := &logger.Config{Level: "info"}

	cfg, err := logger.NewConfig(Environment, c)
	runtime.Must(err)

	logger, err := logger.NewLogger(logger.LoggerParams{Lifecycle: lc, Config: c, Logger: cfg, Version: Version})
	runtime.Must(err)

	return logger
}
