package logger

import (
	"log/slog"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/errors"
)

const (
	// LevelError is the error log level.
	//
	// It is an alias of slog.LevelError, provided so callers can depend on the
	// go-service logger package while using standard slog levels.
	LevelError = slog.LevelError

	// LevelInfo is the info log level.
	//
	// It is an alias of slog.LevelInfo, provided so callers can depend on the
	// go-service logger package while using standard slog levels.
	LevelInfo = slog.LevelInfo

	// LevelWarn is the warning log level.
	//
	// It is an alias of slog.LevelWarn, provided so callers can depend on the
	// go-service logger package while using standard slog levels.
	LevelWarn = slog.LevelWarn
)

// Attr is a structured logging attribute.
type Attr = slog.Attr

// Level is the slog logging level type.
type Level = slog.Level

// ErrNotFound is returned when Config.Kind is unknown.
var ErrNotFound = errors.New("logger: not found")

// Bool returns an Attr representing a boolean key/value pair.
//
// This is a thin wrapper around slog.Bool and does not change semantics.
func Bool(key string, v bool) Attr {
	return slog.Bool(key, v)
}

// Int returns an Attr representing an integer key/value pair.
//
// This is a thin wrapper around slog.Int and does not change semantics.
func Int(key string, value int) Attr {
	return slog.Int(key, value)
}

// LogError logs an error message using the process-wide default slog logger.
//
// This is a thin wrapper around slog.ErrorContext and does not change semantics.
// It is useful in code that prefers importing go-service packages rather than
// log/slog directly.
func LogError(ctx context.Context, msg string, args ...any) {
	slog.ErrorContext(ctx, msg, args...)
}

// String returns an Attr representing a string key/value pair.
//
// This is a thin wrapper around slog.String and does not change semantics.
func String(key, value string) Attr {
	return slog.String(key, value)
}

// LoggerParams declares the dependencies required by NewLogger.
//
// It is intended for Fx/Dig injection and includes service identity fields that
// are attached as static attributes by most logger implementations.
type LoggerParams struct {
	di.In

	// Lifecycle is used by some logger kinds (for example OTLP) to shut down
	// exporters/providers on application stop.
	Lifecycle di.Lifecycle

	// Config enables logging when non-nil and selects the logger kind.
	Config *Config

	// ID is the host identifier typically attached as a static attribute (for example "id").
	ID env.ID

	// Name is the service name typically attached as a static attribute (for example "name").
	Name env.Name

	// Version is the service version typically attached as a static attribute (for example "version").
	Version env.Version

	// Environment is the deployment environment typically attached as a static attribute (for example "environment").
	Environment env.Environment
}

// NewLogger constructs the configured slog logger and installs it as the process-wide default.
//
// When logging is enabled (params.Config != nil), NewLogger builds the configured logger,
// installs it as the global default via slog.SetDefault, and returns a *Logger wrapper.
//
// If logging is disabled (params.Config == nil), NewLogger returns (nil, nil).
// If Config.Kind is unknown, NewLogger returns ErrNotFound.
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

// Logger wraps slog.Logger and adds helpers that standardize contextual metadata and errors.
//
// When logging via Log / LogAttrs, the logger appends:
//
//   - metadata attributes extracted from the provided context (via Meta), and
//   - an "error" attribute when the Message contains a non-nil error (via Error).
//
// This keeps log records consistent across handlers/exporters.
type Logger struct {
	*slog.Logger
}

// Log logs the Message at its derived level and appends context metadata.
//
// The level is derived from `msg.Level()`. Use `LogAttrs` when you need to override
// the level explicitly.
func (l *Logger) Log(ctx context.Context, msg Message, attrs ...slog.Attr) {
	l.LogAttrs(ctx, msg.Level(), msg, attrs...)
}

// LogAttrs logs the Message at level and appends context metadata.
//
// It appends `Meta(ctx)` and `Error(msg.Error)` to the provided attrs before
// delegating to the underlying slog.Logger.
func (l *Logger) LogAttrs(ctx context.Context, level slog.Level, msg Message, attrs ...slog.Attr) {
	attrs = append(attrs, Meta(ctx)...)
	attrs = append(attrs, Error(msg.Error))

	l.Logger.LogAttrs(ctx, level, msg.Text, attrs...)
}

// GetLogger returns the underlying slog.Logger.
//
// It returns nil if the receiver is nil, making it safe to call when logging is
// disabled and a *Logger dependency was not provided.
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
