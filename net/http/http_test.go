package http_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/config/options"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/time"
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

func TestNewServerRejectsNegativeTimeoutOption(t *testing.T) {
	handler := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})

	for _, key := range []string{"read_timeout", "write_timeout", "idle_timeout", "read_header_timeout"} {
		t.Run(key, func(t *testing.T) {
			require.Panics(t, func() {
				http.NewServer(options.Map{key: "-1s"}, time.Second, handler)
			})
		})
	}
}

func TestSameOriginRedirect(t *testing.T) {
	tests := []struct {
		want error
		name string
		next string
	}{
		{name: "same origin", next: "https://example.com/next", want: nil},
		{name: "different host", next: "https://other.example.com/next", want: http.ErrUseLastResponse},
		{name: "different scheme", next: "http://example.com/next", want: http.ErrUseLastResponse},
	}

	prev, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "https://example.com/start", http.NoBody)
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			next, err := http.NewRequestWithContext(t.Context(), http.MethodGet, tt.next, http.NoBody)
			require.NoError(t, err)
			err = http.SameOriginRedirect(next, []*http.Request{prev})
			if tt.want == nil {
				require.NoError(t, err)
				return
			}

			require.ErrorIs(t, err, tt.want)
		})
	}
}

func TestIgnoreRedirect(t *testing.T) {
	err := http.IgnoreRedirect(nil, nil)
	require.ErrorIs(t, err, http.ErrUseLastResponse)
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
