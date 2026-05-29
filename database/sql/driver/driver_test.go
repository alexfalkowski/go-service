package driver_test

import (
	"strconv"
	"testing"

	"github.com/alexfalkowski/go-service/v2/database/sql"
	"github.com/alexfalkowski/go-service/v2/database/sql/config"
	"github.com/alexfalkowski/go-service/v2/database/sql/driver"
	"github.com/alexfalkowski/go-service/v2/database/sql/telemetry"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/telemetry/attributes"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/fx/fxtest"
)

func TestOpenUnregistersDBStatsMetrics(t *testing.T) {
	reader := test.EnableMetricsReader(t)
	driverName := registerDriver(t)

	lc := fxtest.NewLifecycle(t)
	db, err := driver.Open(lc, driverName, test.FS, &config.Config{
		Masters: []config.DSN{{URL: "benchmark"}},
	})
	require.NoError(t, err)
	require.NotNil(t, db)

	lc.RequireStart()

	requireDBStatsMetrics(t, reader)

	lc.RequireStop()

	requireNoDBStatsMetrics(t, reader)
	require.Error(t, errors.Join(db.Ping()...))
}

func TestOpenReturnsConnectError(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	db, err := driver.Open(lc, "missing-dsns", test.FS, &config.Config{})

	require.Nil(t, db)
	require.ErrorIs(t, err, driver.ErrNoDSNs)
}

func TestConnectDestroyUnregistersDBStatsMetrics(t *testing.T) {
	reader := test.EnableMetricsReader(t)
	driverName := registerDriver(t)

	db, err := driver.Connect(driverName, test.FS, &config.Config{
		Masters: []config.DSN{{URL: "benchmark"}},
	})
	require.NoError(t, err)
	require.NotNil(t, db)

	requireDBStatsMetrics(t, reader)

	require.NoError(t, db.Destroy())

	requireNoDBStatsMetrics(t, reader)
}

func TestConnectUsesTelemetryOptionsForDBStatsMetrics(t *testing.T) {
	reader := test.EnableMetricsReader(t)
	driverName := registerDriver(t)

	db, err := driver.Connect(driverName, test.FS, &config.Config{
		Masters: []config.DSN{{URL: "benchmark"}},
	}, telemetry.WithAttributes(attributes.DBSystemNamePostgreSQL))
	require.NoError(t, err)
	require.NotNil(t, db)

	requireDBStatsMetrics(t, reader)
	requireDBSystemName(t, reader, "postgresql")

	require.NoError(t, db.Destroy())
}

func TestRegisterReturnsDuplicateRegistrationError(t *testing.T) {
	driverName := registerDriver(t)

	err := driver.Register(driverName, test.BenchmarkSQLDriver{})
	require.Error(t, err)
}

func TestRegisterWrapsDriverWhenTelemetryIsEnabled(t *testing.T) {
	test.EnableMetricsReader(t)
	exporter, shutdown := setupSpans(t)
	defer shutdown()

	driverName := registerDriver(t)

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

func TestConnectUnregistersSlaveDBStatsMetrics(t *testing.T) {
	reader := test.EnableMetricsReader(t)
	driverName := registerDriver(t)

	db, err := driver.Connect(driverName, test.FS, &config.Config{
		Slaves: []config.DSN{{URL: "benchmark"}},
	})
	require.NoError(t, err)
	require.NotNil(t, db)

	requireDBStatsMetrics(t, reader)

	require.NoError(t, db.Destroy())

	requireNoDBStatsMetrics(t, reader)
}

func TestConnectAppliesPoolSettings(t *testing.T) {
	driverName := registerDriver(t)

	db, err := driver.Connect(driverName, test.FS, &config.Config{
		Masters:      []config.DSN{{URL: "benchmark"}},
		Slaves:       []config.DSN{{URL: "benchmark"}},
		MaxOpenConns: 3,
	})
	require.NoError(t, err)
	require.NotNil(t, db)
	defer func() {
		require.NoError(t, db.Destroy())
	}()

	masters, _ := db.GetAllMasters()
	require.Len(t, masters, 1)
	require.Equal(t, 3, masters[0].Stats().MaxOpenConnections)

	slaves, _ := db.GetAllSlaves()
	require.Len(t, slaves, 1)
	require.Equal(t, 3, slaves[0].Stats().MaxOpenConnections)
}

func TestConnectMasterSlavesReturnsErrors(t *testing.T) {
	db, errs := driver.ConnectMasterSlaves("missing-driver", []string{"benchmark"}, nil)

	require.Nil(t, db)
	require.Error(t, errors.Join(errs...))
}

func registerDriver(t *testing.T) string {
	t.Helper()

	driverName := "test-sql-" + strconv.FormatUint(driverID.Add(1), 10)
	require.NoError(t, driver.Register(driverName, test.BenchmarkSQLDriver{}))

	return driverName
}

func requireDBStatsMetrics(t *testing.T, reader metrics.Reader) {
	t.Helper()

	got := &metrics.ResourceMetrics{}
	require.NoError(t, reader.Collect(t.Context(), got))
	require.Len(t, got.ScopeMetrics, 1)
	require.Len(t, got.ScopeMetrics[0].Metrics, 7)
}

func requireNoDBStatsMetrics(t *testing.T, reader metrics.Reader) {
	t.Helper()

	got := &metrics.ResourceMetrics{}
	require.NoError(t, reader.Collect(t.Context(), got))
	require.Empty(t, got.ScopeMetrics)
}

func requireDBSystemName(t *testing.T, reader metrics.Reader, name string) {
	t.Helper()

	got := &metrics.ResourceMetrics{}
	require.NoError(t, reader.Collect(t.Context(), got))

	for _, scope := range got.ScopeMetrics {
		for _, metric := range scope.Metrics {
			if hasDBSystemName(metric, name) {
				return
			}
		}
	}

	require.Failf(t, "missing db system name", "expected %q in DB stats metrics", name)
}

func setupSpans(t *testing.T) (*test.SpanExporter, func()) {
	t.Helper()

	exporter := &test.SpanExporter{}
	provider := trace.NewTracerProvider(trace.WithSyncer(exporter))
	otel.SetTracerProvider(provider)

	return exporter, func() {
		require.NoError(t, provider.Shutdown(t.Context()))
		otel.SetTracerProvider(noop.NewTracerProvider())
	}
}

func hasDBSystemName(metric metrics.Metrics, name string) bool {
	switch data := metric.Data.(type) {
	case metrics.Gauge[int64]:
		return hasDBSystemNameDataPoint(data.DataPoints, name)
	case metrics.Sum[int64]:
		return hasDBSystemNameDataPoint(data.DataPoints, name)
	default:
		return false
	}
}

func hasDBSystemNameDataPoint(points []metrics.DataPoint[int64], name string) bool {
	for _, point := range points {
		value, ok := point.Attributes.Value(attributes.DBSystemNameKey)
		if ok && value.AsString() == name {
			return true
		}
	}

	return false
}
