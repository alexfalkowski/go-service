package logger

import (
	"context"
	"log/slog"

	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/errors"
	"go.uber.org/fx"
)

// ErrNotFound for logger.
var ErrNotFound = errors.New("logger: not found")

// Params for logger.
type Params struct {
	fx.In

	Lifecycle   fx.Lifecycle
	Config      *Config
	ID          env.ID
	Name        env.Name
	Version     env.Version
	Environment env.Environment
}

// NewLogger using slog.
func NewLogger(params Params) (*Logger, error) {
	logger, err := logger(params)
	if err != nil {
		return nil, err
	}

	if logger == nil {
		return nil, nil
	}

	slog.SetDefault(logger)

	return &Logger{logger}, nil
}

// Logger allows to pass a function to log.
type Logger struct {
	*slog.Logger
}

// Log attrs for logger.
func (l *Logger) Log(ctx context.Context, msg Message, attrs ...slog.Attr) {
	l.LogAttrs(ctx, msg.Level(), msg, attrs...)
}

// LogAttrs for logger.
func (l *Logger) LogAttrs(ctx context.Context, level slog.Level, msg Message, attrs ...slog.Attr) {
	attrs = append(attrs, Meta(ctx)...)
	attrs = append(attrs, Error(msg.Error))

	l.Logger.LogAttrs(ctx, level, msg.Text, attrs...)
}

// GetLogger if defined.
func (l *Logger) GetLogger() *slog.Logger {
	if l == nil {
		return nil
	}

	return l.Logger
}

func logger(params Params) (*slog.Logger, error) {
	switch {
	case !IsEnabled(params.Config):
		return nil, nil
	case params.Config.IsOTLP():
		return newOtlpLogger(params), nil
	case params.Config.IsJSON():
		return newJSONLogger(params), nil
	case params.Config.IsText():
		return newTextLogger(params), nil
	case params.Config.IsTint():
		return newTintLogger(params), nil
	}

	return nil, ErrNotFound
}

func provide(logger *Logger) *slog.Logger {
	return logger.GetLogger()
}
