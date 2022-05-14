package test

import (
	szap "github.com/alexfalkowski/go-service/logger/zap"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// NewLogger for test.
func NewLogger(lc fx.Lifecycle) *zap.Logger {
	logger, _ := szap.NewLogger(szap.LoggerParams{Lifecycle: lc, Config: szap.NewConfig(), Version: Version})

	return logger
}
