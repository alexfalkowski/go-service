package zap

import (
	"context"

	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/version"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// LoggerParams for zap.
type LoggerParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    zap.Config
	Version   version.Version
}

// NewLogger using zap.
func NewLogger(params LoggerParams) (*zap.Logger, error) {
	logger, err := params.Config.Build(zap.Fields(zap.String("name", os.ExecutableName()), zap.String("version", string(params.Version))))
	if err != nil {
		return nil, err
	}

	params.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			_ = logger.Sync()

			return nil
		},
	})

	return logger, nil
}
