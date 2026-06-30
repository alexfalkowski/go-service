package driver_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/database/sql"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/stretchr/testify/require"
)

// BenchmarkSQLTelemetry measures database/sql wrapper overhead with telemetry disabled, partial, and fully enabled.
func BenchmarkSQLTelemetry(b *testing.B) {
	bench := func(name string, setup func(testing.TB)) {
		b.Run(name, func(b *testing.B) {
			test.ResetTelemetry(b)
			setup(b)
			defer test.ResetTelemetry(b)

			driverName := test.RegisterBenchmarkSQLDriver(b, "benchmark-sql-"+name+"-")

			db, err := sql.Open(driverName, "benchmark")
			require.NoError(b, err)
			defer db.Close()

			b.ReportAllocs()
			b.ResetTimer()

			for b.Loop() {
				func() {
					rows, err := db.QueryContext(b.Context(), "SELECT 1")
					require.NoError(b, err)
					defer func() {
						require.NoError(b, rows.Close())
					}()

					require.NoError(b, rows.Err())
				}()
			}
		})
	}

	bench("disabled", func(testing.TB) {})
	bench("metrics", test.EnableMetrics)
	bench("tracer", test.EnableTracer)
	bench("enabled", test.EnableTelemetry)
}
