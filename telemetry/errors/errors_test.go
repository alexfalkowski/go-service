package errors_test

import (
	"log/slog"
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
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

func TestHandleLogsErrors(t *testing.T) {
	stdout := captureStdout(t)
	handler := errors.NewHandler()
	originalDefault := slog.Default()
	t.Cleanup(func() {
		slog.SetDefault(originalDefault)
	})
	capture := &test.CaptureHandler{}
	slog.SetDefault(slog.New(capture))

	handler.Handle(context.Canceled)
	handler.Handle(context.DeadlineExceeded)

	require.Empty(t, capture.Records)
	_, err := stdout.Seek(0, 0)
	require.NoError(t, err)
	data, reader, err := io.ReadAll(stdout)
	require.NoError(t, err)
	require.NoError(t, reader.Close())
	logs := string(data)
	require.Equal(t, 2, strings.Count(logs, "\n"))
	require.Equal(t, 2, strings.Count(logs, `"msg":"telemetry: global error"`))
	require.Contains(t, logs, `"error":"context canceled"`)
	require.Contains(t, logs, `"error":"context deadline exceeded"`)
}

func captureStdout(t *testing.T) *os.File {
	t.Helper()

	original := os.Stdout
	file, err := os.CreateTemp(t.TempDir(), "telemetry-errors-*.log")
	require.NoError(t, err)
	os.Stdout = file
	t.Cleanup(func() {
		os.Stdout = original
		require.NoError(t, file.Close())
	})

	return file
}
