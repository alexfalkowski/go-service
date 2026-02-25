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

// Bool returns an Attr representing a boolean key/value pair.
//
// It is an alias of `slog.Bool`.
func Bool(key string, v bool) Attr {
	return slog.Bool(key, v)
}

// Int returns an Attr representing an integer key/value pair.
//
// It is an alias of `slog.Int`.
func Int(key string, value int) Attr {
	return slog.Int(key, value)
}

// LogError logs an error message using the process-wide default slog logger.
//
// It is an alias of `slog.ErrorContext`.
func LogError(ctx context.Context, msg string, args ...any) {
	slog.ErrorContext(ctx, msg, args...)
}

// String returns an Attr representing a string key/value pair.
//
// It is an alias of `slog.String`.
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

	// ID is typically attached as a static attribute (for example "id") to log records.
	ID env.ID

	// Name is typically attached as a static attribute (for example "name") to log records.
	Name env.Name

	// Version is typically attached as a static attribute (for example "version") to log records.
	Version env.Version

	// Environment is typically attached as a static attribute (for example "environment") to log records.
	Environment env.Environment
}

// NewLogger constructs the configured slog logger and installs it as the slog default.
//
// If logging is disabled (`params.Config == nil`), it returns (nil, nil).
// If `Config.Kind` is unknown, it returns ErrNotFound.
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
// When logging via `Log`/`LogAttrs`, the logger appends:
//
//   - metadata attributes extracted from the provided context (via Meta)
//   - an "error" attribute when the Message contains a non-nil error (via Error)
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
