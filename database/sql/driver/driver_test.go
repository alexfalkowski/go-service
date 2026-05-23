package driver_test

import (
	"strconv"
	"testing"

	"github.com/alexfalkowski/go-service/v2/database/sql/config"
	"github.com/alexfalkowski/go-service/v2/database/sql/driver"
	"github.com/alexfalkowski/go-service/v2/database/sql/telemetry"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/telemetry/attributes"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.uber.org/fx/fxtest"
)

func TestOpenUnregistersDBStatsMetrics(t *testing.T) {
	reader := setupMetrics(t)
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
}

func TestOpenReturnsConnectError(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	db, err := driver.Open(lc, "missing-dsns", test.FS, &config.Config{})

	require.Nil(t, db)
	require.ErrorIs(t, err, driver.ErrNoDSNs)
}

func TestConnectDestroyUnregistersDBStatsMetrics(t *testing.T) {
	reader := setupMetrics(t)
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
	reader := setupMetrics(t)
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

func TestConnectMasterSlavesReturnsErrors(t *testing.T) {
	db, errs := driver.ConnectMasterSlaves("missing-driver", []string{"benchmark"}, nil)

	require.Nil(t, db)
	require.Error(t, errors.Join(errs...))
}

func setupMetrics(t *testing.T) metrics.Reader {
	t.Helper()

	test.ResetTelemetry(t)
	t.Cleanup(func() {
		test.ResetTelemetry(t)
	})

	reader := metrics.NewManualReader()
	metrics.NewMeterProvider(metrics.MeterProviderParams{
		Lifecycle:   fxtest.NewLifecycle(t),
		Config:      &metrics.Config{},
		Reader:      reader,
		ID:          test.ID,
		Name:        test.Name,
		Version:     test.Version,
		Environment: test.Environment,
	})

	return reader
}

func registerDriver(t *testing.T) string {
	t.Helper()

	driverName := "test-sql-" + strconv.FormatUint(driverID.Add(1), 10)
	require.NoError(t, driver.Register(driverName, test.BenchmarkSQLDriver{}))

	return driverName
}

func requireDBStatsMetrics(t *testing.T, reader metrics.Reader) {
	t.Helper()

	got := &metricdata.ResourceMetrics{}
	require.NoError(t, reader.Collect(t.Context(), got))
	require.Len(t, got.ScopeMetrics, 1)
	require.Len(t, got.ScopeMetrics[0].Metrics, 7)
}

func requireNoDBStatsMetrics(t *testing.T, reader metrics.Reader) {
	t.Helper()

	got := &metricdata.ResourceMetrics{}
	require.NoError(t, reader.Collect(t.Context(), got))
	require.Empty(t, got.ScopeMetrics)
}

func requireDBSystemName(t *testing.T, reader metrics.Reader, name string) {
	t.Helper()

	got := &metricdata.ResourceMetrics{}
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

func hasDBSystemName(metric metricdata.Metrics, name string) bool {
	switch data := metric.Data.(type) {
	case metricdata.Gauge[int64]:
		return hasDBSystemNameDataPoint(data.DataPoints, name)
	case metricdata.Sum[int64]:
		return hasDBSystemNameDataPoint(data.DataPoints, name)
	default:
		return false
	}
}

func hasDBSystemNameDataPoint(points []metricdata.DataPoint[int64], name string) bool {
	for _, point := range points {
		value, ok := point.Attributes.Value(attributes.DBSystemNameKey)
		if ok && value.AsString() == name {
			return true
		}
	}

	return false
}
