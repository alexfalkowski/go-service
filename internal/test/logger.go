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
