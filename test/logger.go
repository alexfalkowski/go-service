package test

import (
	szap "github.com/alexfalkowski/go-service/telemetry/logger/zap"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// NewLogger for test.
func NewLogger(lc fx.Lifecycle) *zap.Logger {
	cfg, _ := szap.NewConfig(&szap.Config{Level: "info"})
	logger, _ := szap.NewLogger(szap.LoggerParams{Lifecycle: lc, Config: cfg, Version: Version})

	return logger
}
