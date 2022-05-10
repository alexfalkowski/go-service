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
			_ = logger.Sync()

			return nil
		},
	})

	return logger, nil
}
