package logger_test

import (
	"log/slog"
	"testing"
	"unicode/utf8"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/telemetry/header"
	"github.com/alexfalkowski/go-service/v2/telemetry/internal/otlp"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func TestLogger(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	log, err := test.NewLogger(lc, test.NewTextLoggerConfig())
	require.NoError(t, err)

	require.NotPanics(t, func() {
		log.Log(t.Context(), logger.NewText("test"), logger.Bool("yes", true))
		log.Log(t.Context(), logger.NewMessage("test", context.Canceled), logger.Bool("yes", true))
		log.LogAttrs(t.Context(), logger.LevelInfo, logger.NewMessage("test", context.Canceled), logger.Bool("yes", true))
		log.Info("hello")
		log.Warn("hello")
		log.Error("hello")
	})
}

func TestLoggerAllowsUserLevelAttr(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	log, err := test.NewLogger(lc, test.NewJSONLoggerConfig())
	require.NoError(t, err)

	require.NotPanics(t, func() {
		log.Log(t.Context(), logger.NewText("test"), logger.String(slog.LevelKey, "user-level"))
	})
}

func TestMetaTruncatesLongValues(t *testing.T) {
	value := strings.Repeat("a", 2048)
	ctx := meta.WithAttributes(
		t.Context(),
		meta.WithRequestID(meta.String(value)),
	)

	attrs := logger.Meta(ctx)

	require.Len(t, attrs, 1)
	require.Equal(t, meta.RequestIDKey, attrs[0].Key)
	require.Len(t, attrs[0].Value.String(), 1024)
}

func TestMetaTruncatesLongValuesAtUTF8Boundary(t *testing.T) {
	maxLength := metaValueLength(t, strings.Repeat("a", 2048))
	prefix := strings.Repeat("a", maxLength-1)

	for _, tt := range []struct {
		name  string
		value string
	}{
		{name: "followed by another rune", value: prefix + "é" + "z"},
		{name: "terminal rune", value: prefix + "é"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctx := meta.WithAttributes(
				t.Context(),
				meta.WithRequestID(meta.String(tt.value)),
			)

			attrs := logger.Meta(ctx)
			truncated := attrs[0].Value.String()

			require.Len(t, attrs, 1)
			require.True(t, utf8.ValidString(truncated))
			require.Equal(t, prefix, truncated)
		})
	}
}

func TestInvalidLogger(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	cfg := &logger.Config{Kind: "wrong", Level: "debug"}
	params := logger.LoggerParams{
		Lifecycle:   lc,
		Config:      cfg,
		ID:          test.ID,
		Name:        test.Name,
		Version:     test.Version,
		Environment: test.Environment,
	}

	log, err := logger.NewLogger(params)
	require.ErrorIs(t, err, logger.ErrNotFound)
	require.Nil(t, log)

	require.NotPanics(t, func() {
		log.Log(t.Context(), logger.NewText("test"), logger.Bool("yes", true))
		log.Log(t.Context(), logger.NewMessage("test", context.Canceled), logger.Bool("yes", true))
		log.LogAttrs(t.Context(), logger.LevelInfo, logger.NewMessage("test", context.Canceled), logger.Bool("yes", true))
		log.Info("hello")
		log.Warn("hello")
		log.Error("hello")
	})
}

func TestDisabledLogger(t *testing.T) {
	original := slog.Default()
	t.Cleanup(func() {
		slog.SetDefault(original)
	})
	replacement := slog.New(&test.CaptureHandler{})
	slog.SetDefault(replacement)

	log, err := logger.NewLogger(logger.LoggerParams{})

	require.NoError(t, err)
	require.Nil(t, log)
	require.Same(t, replacement, slog.Default())
}

func TestLogAddsMetadataAndError(t *testing.T) {
	handler := &test.CaptureHandler{}
	log := &logger.Logger{Logger: slog.New(handler)}
	ctx := meta.WithAttributes(t.Context(), meta.WithRequestID(meta.String("request-id")))

	log.Log(ctx, logger.NewText("plain"), logger.String("component", "test"))
	log.LogAttrs(ctx, logger.LevelWarn, logger.NewMessage("failed", context.Canceled), logger.String("component", "test"))

	require.Len(t, handler.Records, 2)
	require.Equal(t, slog.LevelInfo, handler.Records[0].Level)
	require.Equal(t, "plain", handler.Records[0].Message)
	require.Equal(t, "test", handler.Records[0].Attrs["component"].String())
	require.Equal(t, "request-id", handler.Records[0].Attrs[meta.RequestIDKey].String())

	require.Equal(t, slog.LevelWarn, handler.Records[1].Level)
	require.Equal(t, "failed", handler.Records[1].Message)
	require.Equal(t, "test", handler.Records[1].Attrs["component"].String())
	require.Equal(t, "request-id", handler.Records[1].Attrs[meta.RequestIDKey].String())
	require.Equal(t, context.Canceled, handler.Records[1].Attrs["error"].Any())
}

func TestInvalidLevel(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	cfg := &logger.Config{Kind: "text", Level: "bob"}
	params := logger.LoggerParams{
		Lifecycle:   lc,
		Config:      cfg,
		ID:          test.ID,
		Name:        test.Name,
		Version:     test.Version,
		Environment: test.Environment,
	}

	_, err := logger.NewLogger(params)
	require.ErrorIs(t, err, logger.ErrInvalidLevel)
}

func TestInvalidOTLPEndpoint(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	cfg := &logger.Config{
		Kind: "otlp",
		URL:  "http://collector.example.com/v1/logs",
		Headers: header.Map{
			"Authorization": "Bearer token",
		},
	}
	params := logger.LoggerParams{
		Lifecycle:   lc,
		Config:      cfg,
		ID:          test.ID,
		Name:        test.Name,
		Version:     test.Version,
		Environment: test.Environment,
	}

	_, err := logger.NewLogger(params)
	require.ErrorIs(t, err, otlp.ErrInsecureEndpoint)
}

func TestMissingOTLPEndpointIgnoresEnv(t *testing.T) {
	t.Setenv("OTEL_EXPORTER_OTLP_LOGS_ENDPOINT", "http://collector.example.com/v1/logs")
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "https://collector.example.com")

	lc := fxtest.NewLifecycle(t)
	cfg := &logger.Config{
		Kind: "otlp",
		Headers: header.Map{
			"Authorization": "Bearer token",
		},
	}
	params := logger.LoggerParams{
		Lifecycle:   lc,
		Config:      cfg,
		ID:          test.ID,
		Name:        test.Name,
		Version:     test.Version,
		Environment: test.Environment,
	}

	_, err := logger.NewLogger(params)
	require.ErrorIs(t, err, otlp.ErrMissingEndpoint)
}

func TestOTLPLoggerUsesConfiguredLevel(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	cfg := &logger.Config{
		Kind:  "otlp",
		Level: "error",
		URL:   "https://localhost:4318/v1/logs",
	}
	params := logger.LoggerParams{
		Lifecycle:   lc,
		Config:      cfg,
		ID:          test.ID,
		Name:        test.Name,
		Version:     test.Version,
		Environment: test.Environment,
	}

	log, err := logger.NewLogger(params)
	require.NoError(t, err)

	ctx := t.Context()
	require.False(t, log.Enabled(ctx, slog.LevelInfo))
	require.False(t, log.Enabled(ctx, slog.LevelWarn))
	require.True(t, log.Enabled(ctx, slog.LevelError))

	child := log.With("component", "test").WithGroup("otlp")
	require.False(t, child.Enabled(ctx, slog.LevelWarn))
	require.True(t, child.Enabled(ctx, slog.LevelError))
	require.NotPanics(t, func() {
		child.ErrorContext(ctx, "exported")
	})
}

func metaValueLength(t *testing.T, value string) int {
	t.Helper()

	ctx := meta.WithAttributes(
		t.Context(),
		meta.WithRequestID(meta.String(value)),
	)
	attrs := logger.Meta(ctx)

	require.Len(t, attrs, 1)
	return len(attrs[0].Value.String())
}
