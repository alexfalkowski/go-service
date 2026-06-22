package driver_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/database/sql"
	"github.com/alexfalkowski/go-service/v2/database/sql/config"
	"github.com/alexfalkowski/go-service/v2/database/sql/driver"
	"github.com/alexfalkowski/go-service/v2/database/sql/telemetry"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/telemetry/attributes"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func TestOpenUnregistersDBStatsMetrics(t *testing.T) {
	reader := test.EnableMetricsReader(t)
	driverName := test.RegisterBenchmarkSQLDriver(t, "test-sql-")

	lc := fxtest.NewLifecycle(t)
	db, err := driver.Open(lc, driverName, test.FS, &config.Config{
		Writers: []config.DSN{{URL: "benchmark"}},
	})
	require.NoError(t, err)
	require.NotNil(t, db)

	lc.RequireStart()

	test.RequireDBStatsMetrics(t, reader)

	lc.RequireStop()

	test.RequireNoDBStatsMetrics(t, reader)
	require.Error(t, db.Ping())
}

func TestOpenReturnsConnectError(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	db, err := driver.Open(lc, "missing-dsns", test.FS, &config.Config{})

	require.Nil(t, db)
	require.ErrorIs(t, err, driver.ErrNoDSNs)
}

func TestConnectDestroyUnregistersDBStatsMetrics(t *testing.T) {
	reader := test.EnableMetricsReader(t)
	driverName := test.RegisterBenchmarkSQLDriver(t, "test-sql-")

	db, err := driver.Connect(driverName, test.FS, &config.Config{
		Writers: []config.DSN{{URL: "benchmark"}},
	})
	require.NoError(t, err)
	require.NotNil(t, db)

	test.RequireDBStatsMetrics(t, reader)

	require.NoError(t, db.Destroy())

	test.RequireNoDBStatsMetrics(t, reader)
}

func TestConnectUsesTelemetryOptionsForDBStatsMetrics(t *testing.T) {
	reader := test.EnableMetricsReader(t)
	driverName := test.RegisterBenchmarkSQLDriver(t, "test-sql-")

	db, err := driver.Connect(driverName, test.FS, &config.Config{
		Writers: []config.DSN{{URL: "benchmark"}},
	}, telemetry.WithAttributes(attributes.DBSystemNamePostgreSQL))
	require.NoError(t, err)
	require.NotNil(t, db)

	test.RequireDBStatsMetrics(t, reader)
	test.RequireDBSystemName(t, reader, "postgresql")

	require.NoError(t, db.Destroy())
}

func TestRegisterReturnsDuplicateRegistrationError(t *testing.T) {
	driverName := test.RegisterBenchmarkSQLDriver(t, "test-sql-")

	err := driver.Register(driverName, test.BenchmarkSQLDriver{})
	require.Error(t, err)
}

func TestRegisterWrapsDriverWhenTelemetryIsEnabled(t *testing.T) {
	test.EnableMetricsReader(t)
	exporter := test.EnableSpanExporter(t)

	driverName := test.RegisterBenchmarkSQLDriver(t, "test-sql-")

	db, err := sql.Open(driverName, "benchmark")
	require.NoError(t, err)
	defer db.Close()

	rows, err := db.QueryContext(t.Context(), "SELECT register-telemetry")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, rows.Close())
	}()

	require.NoError(t, rows.Err())

	spans := exporter.Spans()
	require.NotEmpty(t, spans)
}

func TestConnectUnregistersReaderDBStatsMetrics(t *testing.T) {
	reader := test.EnableMetricsReader(t)
	driverName := test.RegisterBenchmarkSQLDriver(t, "test-sql-")

	db, err := driver.Connect(driverName, test.FS, &config.Config{
		Readers: []config.DSN{{URL: "benchmark"}},
	})
	require.NoError(t, err)
	require.NotNil(t, db)

	test.RequireDBStatsMetrics(t, reader)

	require.NoError(t, db.Destroy())

	test.RequireNoDBStatsMetrics(t, reader)
}

func TestConnectAppliesPoolSettings(t *testing.T) {
	driverName := test.RegisterBenchmarkSQLDriver(t, "test-sql-")

	db, err := driver.Connect(driverName, test.FS, &config.Config{
		Writers:      []config.DSN{{URL: "benchmark"}},
		Readers:      []config.DSN{{URL: "benchmark"}},
		MaxOpenConns: 3,
	})
	require.NoError(t, err)
	require.NotNil(t, db)
	defer func() {
		require.NoError(t, db.Destroy())
	}()

	writers := db.Writers()
	require.Len(t, writers, 1)
	require.Equal(t, 3, writers[0].Stats().MaxOpenConnections)

	readers := db.Readers()
	require.Len(t, readers, 1)
	require.Equal(t, 3, readers[0].Stats().MaxOpenConnections)
}

func TestConnectWritersReadersReturnsErrors(t *testing.T) {
	db, errs := driver.ConnectWritersReaders("missing-driver", []string{"benchmark"}, nil)

	require.Nil(t, db)
	require.Error(t, errors.Join(errs...))
}
