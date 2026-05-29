package errors_test

import (
	"log/slog"
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/telemetry/errors"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/stretchr/testify/require"
)

func TestNewHandlerNilLogger(t *testing.T) {
	require.Nil(t, errors.NewHandler(nil))
}

func TestRegisterNilHandler(t *testing.T) {
	original := errors.GetHandler()
	defer errors.SetHandler(original)

	errors.Register(nil)

	require.Same(t, original, errors.GetHandler())
}

func TestHandleNilHandler(t *testing.T) {
	var handler *errors.Handler

	require.NotPanics(t, func() {
		handler.Handle(context.Canceled)
	})
}

func TestHandleLogsError(t *testing.T) {
	capture := &test.CaptureHandler{}
	handler := errors.NewHandler(&logger.Logger{Logger: slog.New(capture)})
	require.NotNil(t, handler)

	handler.Handle(context.Canceled)

	require.Len(t, capture.Records, 1)
	require.Equal(t, slog.LevelError, capture.Records[0].Level)
	require.Equal(t, "telemetry: global error", capture.Records[0].Message)
	require.Equal(t, context.Canceled, capture.Records[0].Attrs["error"].Any())
}
