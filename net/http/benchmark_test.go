package http_test

import (
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

// BenchmarkClientTelemetry measures client request overhead with telemetry disabled, metrics-only, and tracing enabled.
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

func BenchmarkFirstAcceptItem(b *testing.B) {
	for _, size := range []int{1 << 10, 1 << 16, 1 << 20} {
		b.Run("commas="+strconv.Itoa(size), func(b *testing.B) {
			header := "application/json," + strings.Repeat(",", size)

			b.ReportAllocs()
			b.ResetTimer()

			for b.Loop() {
				_ = http.FirstAcceptItem(header)
			}
		})
	}
}
