package test

import (
	"github.com/alexfalkowski/go-service/runtime"
	"github.com/alexfalkowski/go-service/telemetry/logger"
	"go.uber.org/fx"
)

// NewLogger for test.
func NewLogger(lc fx.Lifecycle) *logger.Logger {
	c := &logger.Config{Level: "info"}

	cfg, err := logger.NewConfig(Environment, c)
	runtime.Must(err)

	logger, err := logger.NewLogger(logger.Params{Lifecycle: lc, Config: c, Logger: cfg, Version: Version})
	runtime.Must(err)

	return logger
}
