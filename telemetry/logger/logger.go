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

// Attr is an alias of [slog.Attr].
//
// It lets packages expose go-service logging types without importing log/slog
// directly.
type Attr = slog.Attr

// Level is an alias of [slog.Level].
//
// It lets packages refer to slog levels through the go-service logger package
// when that keeps imports or public APIs simpler.
type Level = slog.Level

// ErrNotFound is returned when Config.Kind is unknown.
var ErrNotFound = errors.New("logger: not found")

// ErrInvalidLevel is returned when Config.Level is unknown.
var ErrInvalidLevel = errors.New("logger: invalid level")

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
// It is a thin wrapper around [slog.ErrorContext] and does not change
// semantics. In applications wired with [NewLogger], the default logger is the
// instance installed by this package. If another package replaces the process
// default via [slog.SetDefault], LogError uses that logger instead.
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
// It is intended for Fx/Dig injection. Stdout-oriented logger kinds attach the
// identity fields as static slog attributes, while the OTLP logger kind maps
// them into OpenTelemetry resource attributes.
type LoggerParams struct {
	di.In

	// Lifecycle is used by logger kinds that need shutdown hooks.
	//
	// For example, the OTLP logger appends an OnStop hook that shuts down its
	// OpenTelemetry logger provider and exporter.
	Lifecycle di.Lifecycle

	// Config enables logging when non-nil and selects the logger kind and level.
	//
	// A nil Config disables logging and causes NewLogger to return (nil, nil).
	Config *Config

	// ID identifies the host or instance emitting logs.
	//
	// Stdout-oriented logger kinds attach it as an "id" attribute. The OTLP
	// logger uses it when constructing the OpenTelemetry resource.
	ID env.ID

	// Name identifies the service emitting logs.
	//
	// Stdout-oriented logger kinds attach it as a "name" attribute. The OTLP
	// logger uses it when constructing the OpenTelemetry resource.
	Name env.Name

	// Version identifies the build or release version of the service.
	//
	// Stdout-oriented logger kinds attach it as a "version" attribute. The OTLP
	// logger uses it when constructing the OpenTelemetry resource.
	Version env.Version

	// Environment identifies the deployment environment for emitted logs.
	//
	// Stdout-oriented logger kinds attach it as an "environment" attribute. The
	// OTLP logger uses it when constructing the OpenTelemetry resource.
	Environment env.Environment
}

// NewLogger constructs the configured slog logger, installs it as the
// process-wide default, and returns a wrapper with go-service logging helpers.
//
// When params.Config is nil, logging is disabled and NewLogger returns (nil,
// nil).
//
// When logging is enabled, NewLogger validates Config.Level, builds the
// implementation selected by Config.Kind, installs it via [slog.SetDefault],
// and returns a [Logger].
//
// Stdout-oriented logger kinds attach ID, Name, Version, and Environment as
// static attributes on every log record. The "otlp" logger kind instead maps
// those values to OpenTelemetry resource attributes, installs a global
// OpenTelemetry logger provider, and registers shutdown hooks on
// params.Lifecycle.
//
// NewLogger returns [ErrNotFound] when Config.Kind is unknown and
// [ErrInvalidLevel] when Config.Level is unsupported.
//
// NewLogger may panic if the selected logger kind performs mandatory startup
// wiring that fails. In particular, the "otlp" logger panics if its exporter
// cannot be created.
func NewLogger(params LoggerParams) (*Logger, error) {
	if !params.Config.IsEnabled() {
		return nil, nil
	}
	if err := validateLevel(params.Config); err != nil {
		return nil, err
	}

	logger, err := newLogger(params)
	if err != nil {
		return nil, err
	}

	slog.SetDefault(logger)
	return &Logger{logger}, nil
}

// Logger wraps [slog.Logger] and adds go-service logging helpers.
//
// The embedded slog.Logger remains available for native slog APIs. The helper
// methods on Logger are nil-safe so callers can treat logging as optional when
// Config is nil and dependency injection provides no logger.
//
// When logging via [Logger.Log] or [Logger.LogAttrs], the logger appends:
//
//   - metadata attributes extracted from the provided context (via Meta), and
//   - an "error" attribute when the Message contains a non-nil error (via Error).
//
// The pass-through severity helpers [Logger.Info], [Logger.Warn], and
// [Logger.Error] keep slog semantics unchanged and do not add those extra
// attributes automatically.
type Logger struct {
	*slog.Logger
}

// Info logs at [LevelInfo].
//
// It forwards msg and args to the embedded [slog.Logger] unchanged. It is a
// no-op when l is nil.
func (l *Logger) Info(msg string, args ...any) {
	if l == nil {
		return
	}

	l.Logger.Info(msg, args...)
}

// Error logs at [LevelError].
//
// It forwards msg and args to the embedded [slog.Logger] unchanged. It is a
// no-op when l is nil.
func (l *Logger) Error(msg string, args ...any) {
	if l == nil {
		return
	}

	l.Logger.Error(msg, args...)
}

// Warn logs at [LevelWarn].
//
// It forwards msg and args to the embedded [slog.Logger] unchanged. It is a
// no-op when l is nil.
func (l *Logger) Warn(msg string, args ...any) {
	if l == nil {
		return
	}

	l.Logger.Warn(msg, args...)
}

// Log logs msg.Text at the level derived from [Message.Level].
//
// It preserves the supplied attrs, then appends the output of [Meta] and the
// standardized error attribute from [Error] before delegating to
// [Logger.LogAttrs]. It is a no-op when l is nil.
//
// Use [Logger.LogAttrs] when you need to override the derived level
// explicitly.
func (l *Logger) Log(ctx context.Context, msg Message, attrs ...slog.Attr) {
	if l == nil {
		return
	}

	l.LogAttrs(ctx, msg.Level(), msg, attrs...)
}

// LogAttrs logs msg.Text at level and enriches the record with go-service
// logging conventions.
//
// It preserves the supplied attrs, then appends the output of [Meta] and the
// standardized error attribute from [Error] before delegating to the embedded
// [slog.Logger]. It is a no-op when l is nil.
func (l *Logger) LogAttrs(ctx context.Context, level slog.Level, msg Message, attrs ...slog.Attr) {
	if l == nil {
		return
	}

	attrs = append(attrs, Meta(ctx)...)
	attrs = append(attrs, Error(msg.Error))

	l.Logger.LogAttrs(ctx, level, msg.Text, attrs...)
}

// GetLogger returns the embedded [slog.Logger].
//
// It returns nil when l is nil, making it safe for adapters that accept an
// optional [Logger] but need to work with a raw `*slog.Logger`.
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
