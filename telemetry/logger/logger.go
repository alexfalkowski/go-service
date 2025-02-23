package logger

import (
	"context"
	"log/slog"

	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/os"
	"go.uber.org/fx"
)

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
	var logger *slog.Logger

	switch {
	case !IsEnabled(params.Config):
		logger = newNoopLogger()
	case params.Config.IsOTLP():
		l, err := newOtlpLogger(params)
		if err != nil {
			return nil, prefix(err)
		}

		logger = l
	case params.Config.IsStdout():
		logger = newStdoutLogger(params)
	}

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
