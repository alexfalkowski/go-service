package test

import (
	"github.com/alexfalkowski/go-service/telemetry"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// NewLogger for test.
func NewLogger(lc fx.Lifecycle) *zap.Logger {
	logger, _ := telemetry.NewLogger(telemetry.LoggerParams{Lifecycle: lc, Version: Version})

	return logger
}
