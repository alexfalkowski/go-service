package logger

import (
	"context"

	"github.com/alexfalkowski/go-service/env"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Params for logger.
type Params struct {
	fx.In

	Lifecycle   fx.Lifecycle
	Config      *Config
	Logger      *zap.Config
	Environment env.Environment
	Version     env.Version
	Name        env.Name
}

// NewLogger using logger.
func NewLogger(params Params) (*zap.Logger, error) {
	if !IsEnabled(params.Config) {
		return zap.NewNop(), nil
	}

	fields := zap.Fields(
		zap.Stringer("name", params.Name),
		zap.Stringer("environment", params.Environment),
		zap.Stringer("version", params.Version),
	)

	logger, err := params.Logger.Build(fields)
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

// LogWithLogger for logger.
func LogWithLogger(logger *zap.Logger, msg string, err error, fields ...zapcore.Field) {
	var fn LogFunc

	if err != nil {
		fn = logger.Error
	} else {
		fn = logger.Info
	}

	LogWithFunc(fn, msg, err, fields...)
}

// LogFunc for logger.
type LogFunc func(msg string, fields ...zapcore.Field)

// LogWithFunc for logger.
func LogWithFunc(fn LogFunc, msg string, err error, fields ...zapcore.Field) {
	fields = append(fields, zap.Error(err))

	fn(msg, fields...)
}
