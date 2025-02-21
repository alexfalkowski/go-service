package test

import (
	"github.com/alexfalkowski/go-service/telemetry/logger"
	"go.uber.org/fx"
)

// NewLogger for test.
func NewLogger(lc fx.Lifecycle, config *logger.Config) *logger.Logger {
	return logger.NewLogger(logger.Params{Lifecycle: lc, Config: config, Version: Version})
}
