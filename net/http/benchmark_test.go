package http_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func BenchmarkClientTelemetry(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))
	defer server.Close()

	bench := func(name string, setup func(testing.TB)) {
		b.Run(name, func(b *testing.B) {
			resetTelemetry(b)
			setup(b)
			defer resetTelemetry(b)

			b.ReportAllocs()

			client := http.NewClient(http.DefaultTransport, time.Second)
			req, err := http.NewRequestWithContext(b.Context(), http.MethodGet, server.URL, http.NoBody)
			require.NoError(b, err)

			b.ResetTimer()

			for b.Loop() {
				resp, err := client.Do(req)
				require.NoError(b, err)
				resp.Body.Close()
			}

			b.StopTimer()
			client.CloseIdleConnections()
		})
	}

	bench("disabled", func(testing.TB) {})
	bench("metrics", enableMetrics)
	bench("tracer", enableTracer)
}

func resetTelemetry(tb testing.TB) {
	tb.Helper()

	require.NoError(tb, tracer.Register(tracer.TracerParams{Lifecycle: fxtest.NewLifecycle(tb)}))
	metrics.NewMeterProvider(metrics.MeterProviderParams{Lifecycle: fxtest.NewLifecycle(tb)})
}

func enableMetrics(tb testing.TB) {
	tb.Helper()

	metrics.NewMeterProvider(metrics.MeterProviderParams{
		Lifecycle:   fxtest.NewLifecycle(tb),
		Config:      &metrics.Config{},
		Reader:      metrics.NewManualReader(),
		ID:          test.ID,
		Name:        test.Name,
		Version:     test.Version,
		Environment: test.Environment,
	})
}

func enableTracer(tb testing.TB) {
	tb.Helper()

	require.NoError(tb, tracer.Register(tracer.TracerParams{
		Lifecycle:   fxtest.NewLifecycle(tb),
		Config:      &tracer.Config{Kind: "otlp"},
		ID:          test.ID,
		Name:        test.Name,
		Version:     test.Version,
		Environment: test.Environment,
	}))
}
