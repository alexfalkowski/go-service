package test

import (
	"strconv"
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/database/sql"
	"github.com/alexfalkowski/go-service/v2/database/sql/driver"
	"github.com/alexfalkowski/go-service/v2/database/sql/pg"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/telemetry/attributes"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-sync"
	"github.com/stretchr/testify/require"
)

const dbStatsMetricCount = 7

var sqlDriverID sync.Uint64

// BenchmarkSQLDriver is a database/sql driver test double that returns empty result sets.
type BenchmarkSQLDriver struct{}

// Open implements [driver.Driver].
func (BenchmarkSQLDriver) Open(string) (driver.Conn, error) {
	return BenchmarkSQLConn{}, nil
}

// BenchmarkSQLConn is a database/sql connection test double.
type BenchmarkSQLConn struct{}

// Prepare implements [driver.Conn] and returns [driver.ErrSkip].
func (BenchmarkSQLConn) Prepare(string) (driver.Stmt, error) {
	return nil, driver.ErrSkip
}

// Close implements [driver.Conn] and always succeeds.
func (BenchmarkSQLConn) Close() error {
	return nil
}

// Begin implements [driver.Conn] and returns [driver.ErrSkip].
func (BenchmarkSQLConn) Begin() (driver.Tx, error) {
	return nil, driver.ErrSkip
}

// QueryContext implements [driver.QueryerContext] and returns empty rows.
func (BenchmarkSQLConn) QueryContext(context.Context, string, []driver.NamedValue) (driver.Rows, error) {
	return BenchmarkSQLRows{}, nil
}

// BenchmarkSQLRows is an empty rows test double.
type BenchmarkSQLRows struct{}

// Columns returns the value column.
func (BenchmarkSQLRows) Columns() []string {
	return []string{"value"}
}

// Close implements [driver.Rows] and always succeeds.
func (BenchmarkSQLRows) Close() error {
	return nil
}

// Next implements [driver.Rows] and returns [io.EOF].
func (BenchmarkSQLRows) Next([]driver.Value) error {
	return io.EOF
}

// RegisterSQLDriver registers sqlDriver under a unique test driver name with prefix.
func RegisterSQLDriver(tb testing.TB, prefix string, sqlDriver driver.Driver) string {
	tb.Helper()

	driverName := prefix + strconv.FormatUint(sqlDriverID.Add(1), 10)
	require.NoError(tb, driver.Register(driverName, sqlDriver))

	return driverName
}

// RegisterBenchmarkSQLDriver registers BenchmarkSQLDriver under a unique test driver name.
func RegisterBenchmarkSQLDriver(tb testing.TB, prefix string) string {
	tb.Helper()

	return RegisterSQLDriver(tb, prefix, BenchmarkSQLDriver{})
}

// RequireSQLQueryRows requires query to execute successfully and closes the returned rows.
func RequireSQLQueryRows(tb testing.TB, db *sql.DB, query string) {
	tb.Helper()

	rows, err := db.QueryContext(tb.Context(), query)
	require.NoError(tb, err)
	defer func() {
		require.NoError(tb, rows.Close())
	}()

	require.NoError(tb, rows.Err())
}

// RequireDBStatsMetrics requires the reader to contain the repository DB stats metrics.
func RequireDBStatsMetrics(tb testing.TB, reader metrics.Reader) {
	tb.Helper()

	got := &metrics.ResourceMetrics{}
	require.NoError(tb, reader.Collect(tb.Context(), got))
	require.Len(tb, got.ScopeMetrics, 1)
	require.Len(tb, got.ScopeMetrics[0].Metrics, dbStatsMetricCount)
}

// RequireNoDBStatsMetrics requires the reader to contain no DB stats metrics.
func RequireNoDBStatsMetrics(tb testing.TB, reader metrics.Reader) {
	tb.Helper()

	got := &metrics.ResourceMetrics{}
	require.NoError(tb, reader.Collect(tb.Context(), got))
	require.Empty(tb, got.ScopeMetrics)
}

// RequireDBSystemName requires a DB stats metric with the expected db.system.name attribute.
func RequireDBSystemName(tb testing.TB, reader metrics.Reader, name string) {
	tb.Helper()

	got := &metrics.ResourceMetrics{}
	require.NoError(tb, reader.Collect(tb.Context(), got))

	for _, scope := range got.ScopeMetrics {
		for _, metric := range scope.Metrics {
			if hasDBSystemName(metric, name) {
				return
			}
		}
	}

	require.Failf(tb, "missing db system name", "expected %q in DB stats metrics", name)
}

// RequireNoSQLSpanQueryText requires SQL spans not to contain raw query text.
func RequireNoSQLSpanQueryText(tb testing.TB, spans []tracer.ReadOnlySpan, query string) {
	tb.Helper()

	for _, span := range spans {
		for _, attr := range span.Attributes() {
			require.NotEqual(tb, "db.statement", string(attr.Key))
			require.NotEqual(tb, "db.query.text", string(attr.Key))
			require.False(tb, strings.Contains(attr.Value.AsString(), query), "sql trace attribute leaked query text")
		}
	}
}

// RequireSQLSpanQueryText requires at least one SQL span attribute to contain query.
func RequireSQLSpanQueryText(tb testing.TB, spans []tracer.ReadOnlySpan, query string) {
	tb.Helper()

	for _, span := range spans {
		for _, attr := range span.Attributes() {
			if strings.Contains(attr.Value.AsString(), query) {
				return
			}
		}
	}

	require.Failf(tb, "missing raw query text", "expected a SQL span attribute to contain %q", query)
}

// SetupAccounts creates the shared accounts table fixture and returns a context plus cleanup.
func SetupAccounts(tb testing.TB, db *sql.DBs) (context.Context, func()) {
	tb.Helper()

	require.NoError(tb, upAccounts(tb.Context(), db))

	ctx, cancel := Timeout(tb.Context())

	return meta.WithAttributes(ctx, WithTest(meta.String("test"))), func() {
		require.NoError(tb, downAccounts(tb.Context(), db))
		cancel()
	}
}

func upAccounts(ctx context.Context, db *sql.DBs) error {
	ctx, cancel := Timeout(ctx)
	defer cancel()

	query := `CREATE TABLE IF NOT EXISTS accounts (
		user_id serial PRIMARY KEY,
		created_at TIMESTAMP NOT NULL
	);`

	writer, err := db.Writer()
	if err != nil {
		return err
	}

	_, err = writer.ExecContext(ctx, query)

	return err
}

func downAccounts(ctx context.Context, db *sql.DBs) error {
	ctx, cancel := Timeout(ctx)
	defer cancel()

	query := "DROP TABLE IF EXISTS accounts;"

	writer, err := db.Writer()
	if err != nil {
		return err
	}

	_, err = writer.ExecContext(ctx, query)

	return err
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

func (w *World) registerDatabase() {
	pg.Register()

	w.Append(di.Hook{
		OnStart: func(_ context.Context) error {
			if w.PG == nil || !w.PG.IsEnabled() {
				return nil
			}

			return w.openDatabase()
		},
		OnStop: func(_ context.Context) error {
			if w.DB == nil {
				return nil
			}

			return w.DB.Destroy()
		},
	})
}

func (w *World) openDatabase() error {
	db, err := pg.Connect(FS, w.PG)
	if err != nil {
		return err
	}

	w.DB = db

	return nil
}
