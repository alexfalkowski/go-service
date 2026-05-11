package driver_test

import (
	"strconv"
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/database/sql"
	"github.com/alexfalkowski/go-service/v2/database/sql/driver"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-sync"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

var driverID sync.Uint64

func BenchmarkSQLTelemetry(b *testing.B) {
	bench := func(name string, setup func(testing.TB)) {
		b.Run(name, func(b *testing.B) {
			resetTelemetry(b)
			setup(b)
			defer resetTelemetry(b)

			driverName := "benchmark-sql-" + name + "-" + strconv.FormatUint(driverID.Add(1), 10)
			require.NoError(b, driver.Register(driverName, benchmarkDriver{}))

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
	bench("metrics", enableMetrics)
	bench("tracer", enableTracer)
	bench("enabled", enableTelemetry)
}

type benchmarkDriver struct{}

func (benchmarkDriver) Open(string) (driver.Conn, error) {
	return benchmarkConn{}, nil
}

type benchmarkConn struct{}

func (benchmarkConn) Prepare(string) (driver.Stmt, error) {
	return nil, driver.ErrSkip
}

func (benchmarkConn) Close() error {
	return nil
}

func (benchmarkConn) Begin() (driver.Tx, error) {
	return nil, driver.ErrSkip
}

func (benchmarkConn) QueryContext(context.Context, string, []driver.NamedValue) (driver.Rows, error) {
	return benchmarkRows{}, nil
}

type benchmarkRows struct{}

func (benchmarkRows) Columns() []string {
	return []string{"value"}
}

func (benchmarkRows) Close() error {
	return nil
}

func (benchmarkRows) Next([]driver.Value) error {
	return io.EOF
}

func resetTelemetry(tb testing.TB) {
	tb.Helper()

	tracer.Register(tracer.TracerParams{Lifecycle: fxtest.NewLifecycle(tb)})
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

	tracer.Register(tracer.TracerParams{
		Lifecycle:   fxtest.NewLifecycle(tb),
		Config:      &tracer.Config{},
		ID:          test.ID,
		Name:        test.Name,
		Version:     test.Version,
		Environment: test.Environment,
	})
}

func enableTelemetry(tb testing.TB) {
	tb.Helper()

	enableMetrics(tb)
	enableTracer(tb)
}
