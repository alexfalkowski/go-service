package logger_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func TestLogger(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	log := test.NewLogger(lc, test.NewTextLoggerConfig())

	log.Log(t.Context(), logger.NewText("test"), logger.Bool("yes", true))
	log.Log(t.Context(), logger.NewMessage("test", context.Canceled), logger.Bool("yes", true))
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

	_, err := logger.NewLogger(params)
	require.Error(t, err)
}
