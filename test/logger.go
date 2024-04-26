package test

import (
	logger "github.com/alexfalkowski/go-service/telemetry/logger/zap"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// NewLogger for test.
func NewLogger(lc fx.Lifecycle) *zap.Logger {
	c := &logger.Config{Enabled: true, Level: "info"}
	cfg, _ := logger.NewConfig(Environment, c)
	logger, _ := logger.NewLogger(logger.LoggerParams{Lifecycle: lc, Config: c, ZapConfig: cfg, Version: Version})

	return logger
}
