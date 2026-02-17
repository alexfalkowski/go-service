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

// Attr is an alias for slog.Attr.
type Attr = slog.Attr

// Level is an alias for slog.Level.
type Level = slog.Level

// ErrNotFound is returned when the configured logger kind is unknown.
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

// LoggerParams defines dependencies used to construct a Logger.
type LoggerParams struct {
	di.In
	Lifecycle   di.Lifecycle
	Config      *Config
	ID          env.ID
	Name        env.Name
	Version     env.Version
	Environment env.Environment
}

// NewLogger constructs the configured logger and installs it as the slog default.
func NewLogger(params LoggerParams) (*Logger, error) {
	if !params.Config.IsEnabled() {
		return nil, nil
	}

	logger, err := newLogger(params)
	if err != nil {
		return nil, err
	}

	slog.SetDefault(logger)
	return &Logger{logger}, nil
}

// Logger wraps slog.Logger and adds meta/error context helpers.
type Logger struct {
	*slog.Logger
}

// Log logs msg at its level with attrs and context metadata.
func (l *Logger) Log(ctx context.Context, msg Message, attrs ...slog.Attr) {
	l.LogAttrs(ctx, msg.Level(), msg, attrs...)
}

// LogAttrs logs msg at level with attrs and context metadata.
func (l *Logger) LogAttrs(ctx context.Context, level slog.Level, msg Message, attrs ...slog.Attr) {
	attrs = append(attrs, Meta(ctx)...)
	attrs = append(attrs, Error(msg.Error))

	l.Logger.LogAttrs(ctx, level, msg.Text, attrs...)
}

// GetLogger returns the underlying slog.Logger, or nil if Logger is nil.
func (l *Logger) GetLogger() *slog.Logger {
	if l == nil {
		return nil
	}

	return l.Logger
}

func newLogger(params LoggerParams) (*slog.Logger, error) {
	switch params.Config.Kind {
	case "otlp":
		return newOtlpLogger(params), nil
	case "json":
		return newJSONLogger(params), nil
	case "text":
		return newTextLogger(params), nil
	case "tint":
		return newTintLogger(params), nil
	default:
		return nil, ErrNotFound
	}
}

func convertLogger(logger *Logger) *slog.Logger {
	return logger.GetLogger()
}
