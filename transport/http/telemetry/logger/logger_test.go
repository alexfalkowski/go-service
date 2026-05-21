package logger_test

import (
	"bytes"
	"errors"
	"log/slog"
	"net/http"
	"testing"

	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	httplogger "github.com/alexfalkowski/go-service/v2/transport/http/telemetry/logger"
	"github.com/stretchr/testify/require"
)

func TestRoundTripperLogsTransportError(t *testing.T) {
	var logs bytes.Buffer
	transportErr := errors.New("dial failed")
	base := errorRoundTripper{err: transportErr}
	slogLogger := slog.New(slog.NewJSONHandler(&logs, &slog.HandlerOptions{}))
	rt := httplogger.NewRoundTripper(&logger.Logger{Logger: slogLogger}, base)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "https://example.com/users", http.NoBody)
	require.NoError(t, err)

	res, err := rt.RoundTrip(req)
	require.ErrorIs(t, err, transportErr)
	require.Nil(t, res)
	require.Contains(t, logs.String(), `"level":"ERROR"`)
	require.Contains(t, logs.String(), `"error":"dial failed"`)
}

type errorRoundTripper struct {
	err error
}

func (r errorRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, r.err
}
