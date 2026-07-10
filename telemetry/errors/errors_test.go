package errors_test

import (
	"log/slog"
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/encoding/json"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/telemetry/errors"
	"github.com/stretchr/testify/require"
)

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
	file, err := os.CreateTemp(t.TempDir(), "otel-*.log")
	require.NoError(t, err)

	stdout := os.Stdout
	os.Stdout = file
	t.Cleanup(func() {
		os.Stdout = stdout
		_ = file.Close()
	})

	handler := errors.NewHandler(nil)
	require.NotNil(t, handler)

	// Replace the process-wide default logger to prove the handler writes to its
	// own independent sink rather than slog.Default.
	original := slog.Default()
	capture := &test.CaptureHandler{}
	slog.SetDefault(slog.New(capture))
	t.Cleanup(func() { slog.SetDefault(original) })

	handler.Handle(context.Canceled)

	content, err := test.FS.ReadFile(file.Name())
	require.NoError(t, err)

	record := map[string]any{}
	require.NoError(t, json.Unmarshal(content, &record))
	require.Equal(t, "error", record["level"])
	require.Equal(t, "telemetry: global error", record["msg"])
	require.Equal(t, context.Canceled.Error(), record["error"])

	require.Empty(t, capture.Records)
}
