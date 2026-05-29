package logger_test

import (
	"bytes"
	"errors"
	"log/slog"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	httplogger "github.com/alexfalkowski/go-service/v2/transport/http/telemetry/logger"
	"github.com/stretchr/testify/require"
)

func TestHandlerLogsResponse(t *testing.T) {
	tests := []struct {
		name  string
		level string
		code  int
	}{
		{name: "success", code: http.StatusOK, level: "INFO"},
		{name: "client error", code: http.StatusBadRequest, level: "WARN"},
		{name: "server error", code: http.StatusInternalServerError, level: "ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var logs bytes.Buffer
			slogLogger := slog.New(slog.NewJSONHandler(&logs, &slog.HandlerOptions{}))
			handler := httplogger.NewHandler(test.Name, &logger.Logger{Logger: slogLogger})
			res := httptest.NewRecorder()
			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/greeter/say-hello", http.NoBody)

			handler.ServeHTTP(res, req, func(res http.ResponseWriter, _ *http.Request) {
				res.WriteHeader(tt.code)
			})

			require.Contains(t, logs.String(), `"level":"`+tt.level+`"`)
			require.Contains(t, logs.String(), `"msg":"http: say-hello greeter"`)
			require.Contains(t, logs.String(), `"system":"http"`)
			require.Contains(t, logs.String(), `"service":"greeter"`)
			require.Contains(t, logs.String(), `"method":"say-hello"`)
			require.Contains(t, logs.String(), `"code":`+strconv.Itoa(tt.code))
		})
	}
}

func TestHandlerSkipsOperationPath(t *testing.T) {
	var logs bytes.Buffer
	slogLogger := slog.New(slog.NewJSONHandler(&logs, &slog.HandlerOptions{}))
	handler := httplogger.NewHandler(test.Name, &logger.Logger{Logger: slogLogger})
	res := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/test/metrics", http.NoBody)
	called := false

	handler.ServeHTTP(res, req, func(res http.ResponseWriter, _ *http.Request) {
		called = true
		res.WriteHeader(http.StatusNoContent)
	})

	require.True(t, called)
	require.Empty(t, logs.String())
}

func TestRoundTripperLogsTransportError(t *testing.T) {
	var logs bytes.Buffer
	transportErr := errors.New("dial failed")
	base := &test.ErrorRoundTripper{Err: transportErr}
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

func TestRoundTripperLogsResponse(t *testing.T) {
	tests := []struct {
		name  string
		level string
		code  int
	}{
		{name: "success", code: http.StatusOK, level: "INFO"},
		{name: "client error", code: http.StatusBadRequest, level: "WARN"},
		{name: "server error", code: http.StatusInternalServerError, level: "ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var logs bytes.Buffer
			base := &test.StatusRoundTripper{Status: tt.code}
			slogLogger := slog.New(slog.NewJSONHandler(&logs, &slog.HandlerOptions{}))
			rt := httplogger.NewRoundTripper(&logger.Logger{Logger: slogLogger}, base)
			req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "https://example.com/users", http.NoBody)
			require.NoError(t, err)

			res, err := rt.RoundTrip(req)
			require.NoError(t, err)
			require.Equal(t, tt.code, res.StatusCode)
			require.Contains(t, logs.String(), `"level":"`+tt.level+`"`)
			require.Contains(t, logs.String(), `"msg":"http: get users"`)
			require.Contains(t, logs.String(), `"system":"http"`)
			require.Contains(t, logs.String(), `"service":"users"`)
			require.Contains(t, logs.String(), `"method":"get"`)
			require.Contains(t, logs.String(), `"code":`+strconv.Itoa(tt.code))
		})
	}
}
