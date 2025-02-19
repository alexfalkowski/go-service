package zap

import (
	"context"

	"github.com/alexfalkowski/go-service/env"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LoggerParams for zap.
type LoggerParams struct {
	fx.In

	Lifecycle   fx.Lifecycle
	Config      *Config
	Logger      *zap.Config
	Environment env.Environment
	Version     env.Version
	Name        env.Name
}

// NewLogger using zap.
func NewLogger(params LoggerParams) (*zap.Logger, error) {
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

// LogWithLogger for zap.
func LogWithLogger(msg string, err error, logger *zap.Logger, fields ...zapcore.Field) {
	var fn LogFunc

	if err != nil {
		fn = logger.Error
	} else {
		fn = logger.Info
	}

	LogWithFunc(msg, err, fn, fields...)
}

// LogFunc for zap.
type LogFunc func(msg string, fields ...zapcore.Field)

// LogWithFunc for zap.
func LogWithFunc(msg string, err error, fn LogFunc, fields ...zapcore.Field) {
	if err != nil {
		fields = append(fields, zap.Error(err))
	}

	fn(msg, fields...)
}
