package logger

import (
	"context"
	"log/slog"
	"strings"

	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/io"
	"github.com/alexfalkowski/go-service/os"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	otlp "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	"go.uber.org/fx"
)

var levels = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

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
		logger = noopLogger()
	case params.Config.IsOTLP():
		l, err := otlpLogger(params)
		if err != nil {
			return nil, err
		}

		logger = l
	case params.Config.IsStdout():
		logger = stdoutLogger(params)
	}

	return &Logger{logger}, nil
}

// Logger allows to pass a function to log.
type Logger struct {
	*slog.Logger
}

// Log attrs for logger.
func (l *Logger) Log(ctx context.Context, msg string, err error, attrs ...slog.Attr) {
	var level slog.Level

	if err != nil {
		level = slog.LevelError
	} else {
		level = slog.LevelInfo
	}

	l.LogAttrs(ctx, level, msg, err, attrs...)
}

// LogAttrs for logger.
func (l *Logger) LogAttrs(ctx context.Context, level slog.Level, msg string, err error, attrs ...slog.Attr) {
	attrs = append(attrs, Meta(ctx)...)

	if err != nil {
		attrs = append(attrs, slog.String("error", err.Error()))
	}

	l.Logger.LogAttrs(ctx, level, msg, attrs...)
}

func noopLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(&io.NoopWriter{}, nil))
}

func stdoutLogger(params Params) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: levels[params.Config.Level],
		ReplaceAttr: func(_ []string, attr slog.Attr) slog.Attr {
			if attr.Key == slog.LevelKey {
				level := attr.Value.Any().(slog.Level)
				attr.Value = slog.StringValue(strings.ToLower(level.String()))
			}

			return attr
		},
	}

	if params.Environment.IsDevelopment() {
		return slog.New(slog.NewTextHandler(os.Stdout, opts))
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, opts))
}

func otlpLogger(params Params) (*slog.Logger, error) {
	if err := params.Config.Headers.Secrets(params.FileSystem); err != nil {
		return nil, errors.Prefix("tracer", err)
	}

	client, _ := otlp.New(context.Background(), otlp.WithEndpointURL(params.Config.URL), otlp.WithHeaders(params.Config.Headers))
	attrs := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(params.Name.String()),
		semconv.ServiceVersion(params.Version.String()),
		semconv.DeploymentEnvironmentName(params.Environment.String()),
	)

	provider := log.NewLoggerProvider(log.WithProcessor(log.NewBatchProcessor(client)), log.WithResource(attrs))

	global.SetLoggerProvider(provider)
	slog.SetLogLoggerLevel(levels[params.Config.Level])

	params.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			_ = provider.Shutdown(ctx)
			_ = client.Shutdown(ctx)

			return nil
		},
	})

	return otelslog.NewLogger(params.Name.String(), otelslog.WithLoggerProvider(provider)), nil
}
