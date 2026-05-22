package driver_test

import (
	"strconv"
	"testing"

	"github.com/alexfalkowski/go-service/v2/database/sql"
	"github.com/alexfalkowski/go-service/v2/database/sql/driver"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-sync"
	"github.com/stretchr/testify/require"
)

var driverID sync.Uint64

func BenchmarkSQLTelemetry(b *testing.B) {
	bench := func(name string, setup func(testing.TB)) {
		b.Run(name, func(b *testing.B) {
			test.ResetTelemetry(b)
			setup(b)
			defer test.ResetTelemetry(b)

			driverName := "benchmark-sql-" + name + "-" + strconv.FormatUint(driverID.Add(1), 10)
			require.NoError(b, driver.Register(driverName, test.BenchmarkSQLDriver{}))

			db, err := sql.Open(driverName, "benchmark")
			require.NoError(b, err)
			defer db.Close()

			b.ReportAllocs()
			b.ResetTimer()

			for b.Loop() {
				rows, err := db.QueryContext(b.Context(), "SELECT 1")
				require.NoError(b, err)
				require.NoError(b, rows.Err())
				require.NoError(b, rows.Close())
			}
		})
	}

	bench("disabled", func(testing.TB) {})
	bench("metrics", test.EnableMetrics)
	bench("tracer", test.EnableTracer)
	bench("enabled", test.EnableTelemetry)
}
