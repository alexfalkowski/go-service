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
	value := strings.Repeat("a", 1023) + "é" + "z"
	ctx := meta.WithAttributes(
		t.Context(),
		meta.WithRequestID(meta.String(value)),
	)

	attrs := logger.Meta(ctx)
	truncated := attrs[0].Value.String()

	require.Len(t, attrs, 1)
	require.True(t, utf8.ValidString(truncated))
	require.Equal(t, strings.Repeat("a", 1023), truncated)
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
	require.Error(t, err)
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
