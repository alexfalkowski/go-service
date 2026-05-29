package telemetry_test

import (
	"strconv"
	"testing"

	"github.com/alexfalkowski/go-service/v2/database/sql"
	"github.com/alexfalkowski/go-service/v2/database/sql/driver"
	"github.com/alexfalkowski/go-service/v2/database/sql/telemetry"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-sync"
	"github.com/stretchr/testify/require"
)

var driverID sync.Uint64

func TestOpen(t *testing.T) {
	db, err := telemetry.Open("missing", "dsn")
	require.Nil(t, db)
	require.Error(t, err)
	require.ErrorContains(t, err, `sql: unknown driver "missing"`)
}

func TestOpenDisablesRawQueryTextByDefault(t *testing.T) {
	exporter, shutdown := setupSpans(t)
	defer shutdown()

	driverName := registerDriver(t, test.BenchmarkSQLDriver{})
	const query = "SELECT open-secret-query"

	db, err := telemetry.Open(driverName, "benchmark")
	require.NoError(t, err)
	defer db.Close()

	queryRows(t, db, query)

	spans := exporter.Spans()
	require.NotEmpty(t, spans)
	requireNoRawQueryText(t, spans, query)
}

func TestWrapDriverDisablesRawQueryTextByDefault(t *testing.T) {
	exporter, shutdown := setupSpans(t)
	defer shutdown()

	driverName := registerDriver(t, telemetry.WrapDriver(test.BenchmarkSQLDriver{}))
	const query = "SELECT wrap-secret-query"

	db, err := sql.Open(driverName, "benchmark")
	require.NoError(t, err)
	defer db.Close()

	queryRows(t, db, query)

	spans := exporter.Spans()
	require.NotEmpty(t, spans)
	requireNoRawQueryText(t, spans, query)
}

func TestWithSpanOptionsCanEnableRawQueryText(t *testing.T) {
	exporter, shutdown := setupSpans(t)
	defer shutdown()

	driverName := registerDriver(t, test.BenchmarkSQLDriver{})
	const query = "SELECT opt-in-query"

	db, err := telemetry.Open(driverName, "benchmark", telemetry.WithSpanOptions(telemetry.SpanOptions{}))
	require.NoError(t, err)
	defer db.Close()

	queryRows(t, db, query)

	spans := exporter.Spans()
	require.NotEmpty(t, spans)
	requireRawQueryText(t, spans, query)
}

func registerDriver(t *testing.T, sqlDriver driver.Driver) string {
	t.Helper()

	driverName := "test-sql-telemetry-" + strconv.FormatUint(driverID.Add(1), 10)
	require.NoError(t, driver.Register(driverName, sqlDriver))

	return driverName
}

func queryRows(t *testing.T, db *sql.DB, query string) {
	t.Helper()

	rows, err := db.QueryContext(t.Context(), query)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, rows.Close())
	}()

	require.NoError(t, rows.Err())
}

func setupSpans(t *testing.T) (*test.SpanExporter, func()) {
	t.Helper()

	test.ResetTelemetry(t)
	t.Cleanup(func() {
		test.ResetTelemetry(t)
	})

	exporter := &test.SpanExporter{}
	provider := tracer.NewProvider(tracer.WithSyncer(exporter))
	tracer.SetProvider(provider)

	return exporter, func() {
		require.NoError(t, provider.Shutdown(t.Context()))
		tracer.SetProvider(tracer.NewNoopProvider())
	}
}

func requireNoRawQueryText(t *testing.T, spans []tracer.ReadOnlySpan, query string) {
	t.Helper()

	for _, span := range spans {
		for _, attr := range span.Attributes() {
			require.NotEqual(t, "db.statement", string(attr.Key))
			require.NotEqual(t, "db.query.text", string(attr.Key))
			require.False(t, strings.Contains(attr.Value.AsString(), query), "sql trace attribute leaked query text")
		}
	}
}

func requireRawQueryText(t *testing.T, spans []tracer.ReadOnlySpan, query string) {
	t.Helper()

	for _, span := range spans {
		for _, attr := range span.Attributes() {
			if strings.Contains(attr.Value.AsString(), query) {
				return
			}
		}
	}

	require.Failf(t, "missing raw query text", "expected a SQL span attribute to contain %q", query)
}
