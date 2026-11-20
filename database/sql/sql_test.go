package sql_test

import (
	"strconv"
	"testing"

	"github.com/alexfalkowski/go-service/v3/database/sql"
	"github.com/alexfalkowski/go-service/v3/database/sql/driver"
	"github.com/alexfalkowski/go-service/v3/internal/test"
	"github.com/alexfalkowski/go-service/v3/runtime"
	"github.com/alexfalkowski/go-sync"
	"github.com/stretchr/testify/require"
)

var sqlDriverID sync.Uint64

func TestOpenUnknownDriver(t *testing.T) {
	db, err := sql.Open("go-service-missing-driver", "benchmark")

	require.Error(t, err)
	require.Nil(t, db)
}

func TestConfigIsEnabled(t *testing.T) {
	require.False(t, (*sql.Config)(nil).IsEnabled())
	require.True(t, (&sql.Config{}).IsEnabled())
}

func TestConvertAssignConvertsStringToInteger(t *testing.T) {
	var value int64

	require.NoError(t, sql.ConvertAssign(driver.ScanContext{}, &value, "42"))
	require.Equal(t, int64(42), value)

	require.Error(t, sql.ConvertAssign(driver.ScanContext{}, &value, "invalid"))
}

func TestNamedPreservesNameAndValue(t *testing.T) {
	require.Equal(t, sql.NamedArg{Name: "id", Value: 42}, sql.Named("id", 42))
}

func TestRegisterMakesDriverAvailable(t *testing.T) {
	exporter := test.EnableIsolatedSpanExporter(t)
	name := t.Name() + strconv.FormatUint(sqlDriverID.Add(1), 10)
	require.NoError(t, sql.Register(name, test.BenchmarkSQLDriver{}))

	require.Contains(t, sql.Drivers(), name)

	db, err := sql.Open(name, "benchmark")
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, db.Close())
	})

	test.RequireSQLQueryRows(t, db, "SELECT registered")
	require.NotEmpty(t, exporter.Spans())

	require.ErrorIs(t, sql.Register(name, test.BenchmarkSQLDriver{}), runtime.ErrRecovered)
}

func TestOpenWithConnector(t *testing.T) {
	db := sql.OpenDB(test.BenchmarkSQLConnector{})
	t.Cleanup(func() {
		require.NoError(t, db.Close())
	})

	test.RequireSQLQueryRows(t, db, "SELECT connector")
}
