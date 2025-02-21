package logger

import (
	"context"

	"github.com/alexfalkowski/go-service/env"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Field for logger.
type Field = zapcore.Field

var (
	// Stringer for logger.
	Stringer = zap.Stringer

	// String for logger.
	String = zap.String

	// Int for logger.
	Int = zap.Int

	// Any for logger.
	Any = zap.Any

	// Err for logger.
	Error = zap.Error

	// NamedError for logger.
	NamedError = zap.NamedError
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

// NewLogger using zap.
func NewLogger(params Params) (*Logger, error) {
	if !IsEnabled(params.Config) {
		return &Logger{zap.NewNop()}, nil
	}

	fields := zap.Fields(
		zap.Stringer("name", params.Name),
		zap.Stringer("environment", params.Environment),
		zap.Stringer("version", params.Version),
	)

	log, err := params.Logger.Build(fields)
	if err != nil {
		return nil, err
	}

	params.Lifecycle.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			_ = log.Sync()

			return nil
		},
	})

	return &Logger{log}, nil
}

type (
	// Logger allows to pass a function to log.
	Logger struct {
		*zap.Logger
	}

	// LogFunc for logger.
	LogFunc func(msg string, fields ...Field)
)

// LogWithLogger for logger.
func (l *Logger) Log(ctx context.Context, msg string, err error, fields ...Field) {
	var fn LogFunc

	if err != nil {
		fn = l.Logger.Error
	} else {
		fn = l.Logger.Info
	}

	l.LogFunc(ctx, fn, msg, err, fields...)
}

// LogWithFunc for logger.
func (l *Logger) LogFunc(ctx context.Context, fn LogFunc, msg string, err error, fields ...Field) {
	fields = append(fields, Meta(ctx)...)
	fields = append(fields, zap.Error(err))

	fn(msg, fields...)
}
