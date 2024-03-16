package zap

import (
	"context"

	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/version"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// LoggerParams for zap.
type LoggerParams struct {
	fx.In

	Lifecycle   fx.Lifecycle
	Config      *Config
	ZapConfig   zap.Config
	Environment env.Environment
	Version     version.Version
}

// NewLogger using zap.
func NewLogger(params LoggerParams) (*zap.Logger, error) {
	c := params.Config
	if c == nil || !c.Enabled {
		return zap.NewNop(), nil
	}

	fields := zap.Fields(
		zap.String("name", os.ExecutableName()),
		zap.String("environment", string(params.Environment)),
		zap.String("version", string(params.Version)),
	)

	logger, err := params.ZapConfig.Build(fields)
	if err != nil {
		return nil, err
	}

	params.Lifecycle.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			_ = logger.Sync()

			return nil
		},
	})

	return logger, nil
}
