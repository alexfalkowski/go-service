package http_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func TestMaxBytesHandler(t *testing.T) {
	handler := http.MaxBytesHandler(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		_, _, err := io.ReadAll(req.Body)
		var maxBytesError *http.MaxBytesError
		require.ErrorAs(t, err, &maxBytesError)

		_, _ = res.Write([]byte("ok"))
	}), 1)

	res := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", bytes.NewBufferString("too large"))

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
}

func TestHandleWhenTelemetryDisabled(t *testing.T) {
	require.NoError(t, tracer.Register(tracer.TracerParams{Lifecycle: fxtest.NewLifecycle(t)}))
	metrics.NewMeterProvider(metrics.MeterProviderParams{Lifecycle: fxtest.NewLifecycle(t)})

	mux := http.NewServeMux()
	called := false

	http.Handle(mux, "/hello", http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		called = true
		_, _ = res.Write([]byte("hello"))
	}))

	res := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)

	mux.ServeHTTP(res, req)

	require.True(t, called)
	require.Equal(t, http.StatusOK, res.Code)
	require.Equal(t, "hello", res.Body.String())
}

func TestHandleWhenMetricsEnabled(t *testing.T) {
	t.Cleanup(func() {
		require.NoError(t, tracer.Register(tracer.TracerParams{Lifecycle: fxtest.NewLifecycle(t)}))
		metrics.NewMeterProvider(metrics.MeterProviderParams{Lifecycle: fxtest.NewLifecycle(t)})
	})

	require.NoError(t, tracer.Register(tracer.TracerParams{Lifecycle: fxtest.NewLifecycle(t)}))
	metrics.NewMeterProvider(metrics.MeterProviderParams{
		Lifecycle: fxtest.NewLifecycle(t),
		Config:    &metrics.Config{},
		Reader:    metrics.NewManualReader(),
	})

	mux := http.NewServeMux()

	http.Handle(mux, "/hello", http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		_, _ = res.Write([]byte("hello"))
	}))

	res := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)

	mux.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	require.Equal(t, "hello", res.Body.String())
}
