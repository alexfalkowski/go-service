package http_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

func BenchmarkClientTelemetry(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))
	defer server.Close()

	bench := func(name string, setup func(testing.TB)) {
		b.Run(name, func(b *testing.B) {
			test.ResetTelemetry(b)
			setup(b)
			defer test.ResetTelemetry(b)

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
	bench("metrics", test.EnableMetrics)
	bench("tracer", test.EnableTracer)
}
