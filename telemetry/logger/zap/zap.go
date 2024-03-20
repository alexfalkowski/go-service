package zap

import (
	"context"

	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/version"
	pretty "github.com/thessem/zap-prettyconsole"
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
	if !IsEnabled(params.Config) {
		return zap.NewNop(), nil
	}

	fields := zap.Fields(
		zap.String("name", os.ExecutableName()),
		zap.String("environment", string(params.Environment)),
		zap.String("version", string(params.Version)),
	)

	var logger *zap.Logger

	if params.Environment.IsDevelopment() {
		logger = pretty.NewLogger(zap.DebugLevel).WithOptions(fields)
	} else {
		l, err := params.ZapConfig.Build(fields)
		if err != nil {
			return nil, err
		}

		logger = l
	}

	params.Lifecycle.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			_ = logger.Sync()

			return nil
		},
	})

	return logger, nil
}
