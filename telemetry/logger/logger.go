package logger

import (
	"log/slog"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/errors"
)

const (
	// LevelError is an alias of slog.LevelError.
	LevelError = slog.LevelError

	// LevelInfo is an alias of slog.LevelInfo.
	LevelInfo = slog.LevelInfo

	// LevelWarn is an alias of slog.LevelWarn.
	LevelWarn = slog.LevelWarn
)

type (
	// Attr is an alias of slog.Attr.
	Attr = slog.Attr

	// Level is an alias of slog.Level.
	Level = slog.Level
)

// ErrNotFound for logger.
var ErrNotFound = errors.New("logger: not found")

// Bool is an alias of slog.Bool.
func Bool(key string, v bool) Attr {
	return slog.Bool(key, v)
}

// Int is an alias of slog.Int.
func Int(key string, value int) Attr {
	return slog.Int(key, value)
}

// LogError is an alias of slog.ErrorContext.
func LogError(ctx context.Context, msg string, args ...any) {
	slog.ErrorContext(ctx, msg, args...)
}

// String is an alias of slog.String.
func String(key, value string) Attr {
	return slog.String(key, value)
}

// LoggerParams for logger.
type LoggerParams struct {
	di.In
	Lifecycle   di.Lifecycle
	Config      *Config
	ID          env.ID
	Name        env.Name
	Version     env.Version
	Environment env.Environment
}

// NewLogger using slog.
func NewLogger(params LoggerParams) (*Logger, error) {
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

func logger(params LoggerParams) (*slog.Logger, error) {
	switch {
	case !params.Config.IsEnabled():
		return nil, nil
	case params.Config.IsOTLP():
		return newOtlpLogger(params), nil
	case params.Config.IsJSON():
		return newJSONLogger(params), nil
	case params.Config.IsText():
		return newTextLogger(params), nil
	case params.Config.IsTint():
		return newTintLogger(params), nil
	default:
		return nil, ErrNotFound
	}
}

func newLogger(logger *Logger) *slog.Logger {
	return logger.GetLogger()
}
