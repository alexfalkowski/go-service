package driver_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/database/sql"
	"github.com/alexfalkowski/go-service/v2/database/sql/config"
	"github.com/alexfalkowski/go-service/v2/database/sql/driver"
	"github.com/alexfalkowski/go-service/v2/database/sql/telemetry"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/telemetry/attributes"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func TestOpenUnregistersDBStatsMetrics(t *testing.T) {
	reader := test.EnableMetricsReader(t)
	driverName := test.RegisterBenchmarkSQLDriver(t, "test-sql-")

	lc := fxtest.NewLifecycle(t)
	db, err := driver.Open(lc, driverName, test.FS, &config.Config{
		Writer: newPool(1),
	})
	require.NoError(t, err)
	require.NotNil(t, db)

	lc.RequireStart()

	test.RequireDBStatsMetrics(t, reader)

	lc.RequireStop()

	test.RequireNoDBStatsMetrics(t, reader)
	require.Error(t, db.Ping(t.Context()))
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
		Writer: newPool(1),
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
		Writer: newPool(1),
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
	t.Cleanup(func() {
		require.NoError(t, db.Close())
	})

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
		Reader: newPool(1),
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
		Reader: newPool(3),
		Writer: newPool(3),
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

func TestConnectAppliesRolePoolSettings(t *testing.T) {
	driverName := test.RegisterBenchmarkSQLDriver(t, "test-sql-")

	db, err := driver.Connect(driverName, test.FS, &config.Config{
		Reader: newPool(7),
		Writer: newPool(5),
	})
	require.NoError(t, err)
	require.NotNil(t, db)
	defer func() {
		require.NoError(t, db.Destroy())
	}()

	writers := db.Writers()
	require.Len(t, writers, 1)
	require.Equal(t, 5, writers[0].Stats().MaxOpenConnections)

	readers := db.Readers()
	require.Len(t, readers, 1)
	require.Equal(t, 7, readers[0].Stats().MaxOpenConnections)
}

func TestConnectWritersReadersReturnsErrors(t *testing.T) {
	db, errs := driver.ConnectWritersReaders("missing-driver", []string{"benchmark"}, nil)

	require.Nil(t, db)
	require.Error(t, errors.Join(errs...))
}

func TestDBsPingReturnsDeadlineWhenWaitingForConnection(t *testing.T) {
	tests := []struct {
		name    string
		writers []string
		readers []string
		ping    func(*sql.DBs, context.Context) error
	}{
		{
			name:    "all pools",
			writers: []string{"benchmark"},
			ping:    func(db *sql.DBs, ctx context.Context) error { return db.Ping(ctx) },
		},
		{
			name:    "writers",
			writers: []string{"benchmark"},
			ping:    func(db *sql.DBs, ctx context.Context) error { return db.PingWriter(ctx) },
		},
		{
			name:    "readers",
			readers: []string{"benchmark"},
			ping:    func(db *sql.DBs, ctx context.Context) error { return db.PingReader(ctx) },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			driverName := test.RegisterBenchmarkSQLDriver(t, "test-sql-")
			db, errs := sql.ConnectWritersReaders(driverName, tt.writers, tt.readers)
			require.Empty(t, errs)
			t.Cleanup(func() {
				require.NoError(t, db.Destroy())
			})

			pool, err := db.Writer()
			if len(tt.readers) > 0 {
				pool, err = db.Reader()
			}
			require.NoError(t, err)
			pool.SetMaxOpenConns(1)

			conn, err := pool.Conn(t.Context())
			require.NoError(t, err)
			defer func() {
				require.NoError(t, conn.Close())
			}()

			ctx, cancel := context.WithTimeout(t.Context(), time.Millisecond)
			defer cancel()

			require.ErrorIs(t, tt.ping(db, ctx), context.DeadlineExceeded)
		})
	}
}

func newPool(maxOpenConns int) *config.Pool {
	return &config.Pool{
		DSNs: []config.DSN{{URL: "benchmark"}},
		Settings: &config.PoolSettings{
			MaxOpenConns: maxOpenConns,
			MaxIdleConns: maxOpenConns,
		},
	}
}
