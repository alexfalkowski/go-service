package telemetry_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/database/sql"
	"github.com/alexfalkowski/go-service/v2/database/sql/telemetry"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/stretchr/testify/require"
)

func TestOpen(t *testing.T) {
	db, err := telemetry.Open("missing", "dsn")
	require.Nil(t, db)
	require.Error(t, err)
	require.ErrorContains(t, err, `sql: unknown driver "missing"`)
}

func TestOpenDisablesRawQueryTextByDefault(t *testing.T) {
	exporter := test.EnableIsolatedSpanExporter(t)

	driverName := test.RegisterBenchmarkSQLDriver(t, "test-sql-telemetry-")
	const query = "SELECT open-secret-query"

	db, err := telemetry.Open(driverName, "benchmark")
	require.NoError(t, err)
	defer db.Close()

	test.RequireSQLQueryRows(t, db, query)

	spans := exporter.Spans()
	require.NotEmpty(t, spans)
	test.RequireNoSQLSpanQueryText(t, spans, query)
}

func TestWrapDriverDisablesRawQueryTextByDefault(t *testing.T) {
	exporter := test.EnableIsolatedSpanExporter(t)

	driverName := test.RegisterSQLDriver(t, "test-sql-telemetry-", telemetry.WrapDriver(test.BenchmarkSQLDriver{}))
	const query = "SELECT wrap-secret-query"

	db, err := sql.Open(driverName, "benchmark")
	require.NoError(t, err)
	defer db.Close()

	test.RequireSQLQueryRows(t, db, query)

	spans := exporter.Spans()
	require.NotEmpty(t, spans)
	test.RequireNoSQLSpanQueryText(t, spans, query)
}

func TestWithSpanOptionsCanEnableRawQueryText(t *testing.T) {
	exporter := test.EnableIsolatedSpanExporter(t)

	driverName := test.RegisterBenchmarkSQLDriver(t, "test-sql-telemetry-")
	const query = "SELECT opt-in-query"

	db, err := telemetry.Open(driverName, "benchmark", telemetry.WithSpanOptions(telemetry.SpanOptions{}))
	require.NoError(t, err)
	defer db.Close()

	test.RequireSQLQueryRows(t, db, query)

	spans := exporter.Spans()
	require.NotEmpty(t, spans)
	test.RequireSQLSpanQueryText(t, spans, query)
}
