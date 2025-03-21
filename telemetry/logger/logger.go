package logger

import (
	"context"
	"errors"
	"log/slog"

	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/os"
	"go.uber.org/fx"
)

// ErrNotFound for logger.
var ErrNotFound = errors.New("logger: not found")

// Params for logger.
type Params struct {
	fx.In

	Lifecycle   fx.Lifecycle
	FileSystem  os.FileSystem
	Config      *Config
	Environment env.Environment
	Version     env.Version
	Name        env.Name
}

// NewLogger using zap.
func NewLogger(params Params) (*Logger, error) {
	switch {
	case !IsEnabled(params.Config):
		return nil, nil
	case params.Config.IsOTLP():
		logger, err := newOtlpLogger(params)

		return &Logger{logger}, prefix(err)
	case params.Config.IsJSON():
		return &Logger{newJSONLogger(params)}, nil
	case params.Config.IsText():
		return &Logger{newTextLogger(params)}, nil
	}

	return nil, ErrNotFound
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

func provide(logger *Logger) *slog.Logger {
	return logger.GetLogger()
}
