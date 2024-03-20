package test

import (
	szap "github.com/alexfalkowski/go-service/telemetry/logger/zap"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// NewLogger for test.
func NewLogger(lc fx.Lifecycle) *zap.Logger {
	c := &szap.Config{Enabled: true, Level: "info"}
	cfg, _ := szap.NewConfig(DevEnvironment, c)
	logger, _ := szap.NewLogger(szap.LoggerParams{Lifecycle: lc, Config: c, ZapConfig: cfg, Version: Version})

	return logger
}
