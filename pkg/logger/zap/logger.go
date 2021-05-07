package zap

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// NewLogger using zap.
func NewLogger(lc fx.Lifecycle, cfg zap.Config) (*zap.Logger, error) {
	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Sync() // nolint:errcheck

			return nil
		},
	})

	return logger, nil
}
