package pg_test

import (
	"io/fs"
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/database/sql"
	"github.com/alexfalkowski/go-service/v2/database/sql/config"
	"github.com/alexfalkowski/go-service/v2/database/sql/driver"
	"github.com/alexfalkowski/go-service/v2/database/sql/pg"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/telemetry/attributes"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func TestConnect(t *testing.T) {
	cfg := test.NewPGConfig()

	_, err := driver.Connect("missing", test.FS, cfg.Config)
	require.Error(t, err)
}

func TestConfigIsEnabled(t *testing.T) {
	tests := []struct {
		config  *pg.Config
		name    string
		enabled bool
	}{
		{name: "nil"},
		{name: "empty", config: &pg.Config{}},
		{name: "configured", config: &pg.Config{Config: &config.Config{}}, enabled: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.enabled, tt.config.IsEnabled())
		})
	}
}

func TestDisabledConfig(t *testing.T) {
	tests := []struct {
		config *pg.Config
		name   string
	}{
		{name: "nil"},
		{name: "empty", config: &pg.Config{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lc := fxtest.NewLifecycle(t)

			db, err := pg.Connect(test.FS, tt.config)
			require.NoError(t, err)
			require.Nil(t, db)

			db, err = pg.Open(lc, test.FS, tt.config)
			require.NoError(t, err)
			require.Nil(t, db)
		})
	}
}

func TestInvalidOpen(t *testing.T) {
	tests := []struct {
		wantErr error
		config  *pg.Config
		name    string
	}{
		{
			name: "invalid masters",
			config: pgConfig(
				[]config.DSN{{URL: test.FilePath("secrets/none")}},
				[]config.DSN{{URL: test.FilePath("secrets/pg")}},
			),
			wantErr: fs.ErrNotExist,
		},
		{
			name: "invalid slaves",
			config: pgConfig(
				[]config.DSN{{URL: test.FilePath("secrets/pg")}},
				[]config.DSN{{URL: test.FilePath("secrets/none")}},
			),
			wantErr: fs.ErrNotExist,
		},
		{
			name:    "empty dsn configuration",
			config:  pgConfig(nil, nil),
			wantErr: driver.ErrNoDSNs,
		},
		{
			name: "empty master dsn",
			config: pgConfig(
				[]config.DSN{{}},
				[]config.DSN{{URL: test.FilePath("secrets/pg")}},
			),
			wantErr: driver.ErrEmptyDSN,
		},
		{
			name: "empty slave dsn",
			config: pgConfig(
				[]config.DSN{{URL: test.FilePath("secrets/pg")}},
				[]config.DSN{{URL: "env:PG_EMPTY_DSN"}},
			),
			wantErr: driver.ErrEmptyDSN,
		},
	}

	t.Setenv("PG_EMPTY_DSN", "")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldPGConfig(tt.config), test.WithWorldLoggerConfig("json"))

			err := world.Lifecycle.Start(t.Context())
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestConfiguredSQLPings(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldPGConfig(nil), test.WithWorldLoggerConfig("otlp"))

	require.NoError(t, errors.Join(world.DB.Ping()...))
}

func TestConnectUsesPostgreSQLTelemetryOptions(t *testing.T) {
	reader := test.EnableMetricsReader(t)
	cfg := test.NewPGConfig()

	pg.Register()

	db, err := pg.Connect(test.FS, cfg)
	require.NoError(t, err)
	require.NotNil(t, db)
	defer func() {
		require.NoError(t, db.Destroy())
	}()

	requireDBSystemName(t, reader, "postgresql")
}

func TestOpenClosesDBsOnStop(t *testing.T) {
	cfg := test.NewPGConfig()
	lc := fxtest.NewLifecycle(t)

	pg.Register()

	db, err := pg.Open(lc, test.FS, cfg)
	require.NoError(t, err)
	require.NotNil(t, db)

	lc.RequireStart()
	require.NoError(t, errors.Join(db.Ping()...))

	lc.RequireStop()
	require.Error(t, errors.Join(db.Ping()...))
}

func TestDBQuery(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldPGConfig(nil))

	ctx, cleanup := setupAccounts(t, world.DB)
	defer cleanup()

	rows, err := world.DB.QueryContext(ctx, "SELECT table_name FROM information_schema.tables WHERE table_schema='public'")
	require.NoError(t, err)

	var count int
	for rows.Next() {
		count++
	}
	require.Positive(t, count)
	require.NoError(t, rows.Err())
	require.NoError(t, rows.Close())
}

func TestDBExec(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldPGConfig(nil))

	ctx, cleanup := setupAccounts(t, world.DB)
	defer cleanup()

	result, err := world.DB.ExecContext(ctx, "INSERT INTO accounts(created_at) VALUES($1)", time.Now())
	require.NoError(t, err)

	num, err := result.RowsAffected()
	require.NoError(t, err)
	require.Positive(t, num)
}

func TestDBCommitTransExec(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldPGConfig(nil))

	ctx, cleanup := setupAccounts(t, world.DB)
	defer cleanup()

	tx, err := world.DB.BeginTx(ctx, nil)
	require.NoError(t, err)

	//nolint:errcheck
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, "INSERT INTO accounts(created_at) VALUES($1)", time.Now())
	require.NoError(t, err)

	err = tx.Commit()
	require.NoError(t, err)

	_, err = result.LastInsertId()
	require.Error(t, err)

	num, err := result.RowsAffected()
	require.NoError(t, err)
	require.Positive(t, num)
}

func TestDBRollbackTransExec(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldPGConfig(nil))

	ctx, cleanup := setupAccounts(t, world.DB)
	defer cleanup()

	tx, err := world.DB.BeginTx(ctx, nil)
	require.NoError(t, err)

	result, err := tx.ExecContext(ctx, "INSERT INTO accounts(created_at) VALUES($1)", time.Now())
	require.NoError(t, err)

	err = tx.Rollback()
	require.NoError(t, err)

	_, err = result.LastInsertId()
	require.Error(t, err)

	num, err := result.RowsAffected()
	require.NoError(t, err)
	require.Positive(t, num)
}

func TestStatementQuery(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldPGConfig(nil))

	ctx, cleanup := setupAccounts(t, world.DB)
	defer cleanup()

	_, stmt, err := world.DB.PrepareContext(ctx, "SELECT table_name FROM information_schema.tables WHERE table_schema = $1")
	require.NoError(t, err)

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, "public")
	require.NoError(t, err)

	var count int
	for rows.Next() {
		count++
	}
	require.Positive(t, count)
	require.NoError(t, rows.Err())
	require.NoError(t, rows.Close())
}

func TestStatementExec(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldPGConfig(nil))

	ctx, cleanup := setupAccounts(t, world.DB)
	defer cleanup()

	_, stmt, err := world.DB.PrepareContext(ctx, "INSERT INTO accounts(created_at) VALUES($1)")
	require.NoError(t, err)

	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, time.Now())
	require.NoError(t, err)

	_, err = result.LastInsertId()
	require.Error(t, err)

	num, err := result.RowsAffected()
	require.NoError(t, err)
	require.Positive(t, num)
}

func TestTransStatementExec(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldPGConfig(nil))

	ctx, cleanup := setupAccounts(t, world.DB)
	defer cleanup()

	tx, err := world.DB.Begin()
	require.NoError(t, err)

	//nolint:errcheck
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO accounts(created_at) VALUES($1)")
	require.NoError(t, err)

	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, time.Now())
	require.NoError(t, err)

	err = tx.Commit()
	require.NoError(t, err)

	_, err = result.LastInsertId()
	require.Error(t, err)

	num, err := result.RowsAffected()
	require.NoError(t, err)
	require.Positive(t, num)
}

func TestInvalidStatementQuery(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldPGConfig(nil), test.WithWorldLoggerConfig("tint"))

	ctx, cleanup := setupAccounts(t, world.DB)
	defer cleanup()

	_, stmt, err := world.DB.PrepareContext(ctx, "SELECT table_name FROM information_schema.tables WHERE table_schema = $1")
	require.NoError(t, err)

	defer stmt.Close()

	_, err = stmt.QueryContext(ctx, 1)
	require.Error(t, err)
}

func TestInvalidSQLPort(t *testing.T) {
	cfg := &pg.Config{Config: &config.Config{
		Masters:         []config.DSN{{URL: test.FilePath("secrets/pg_invalid")}},
		Slaves:          []config.DSN{{URL: test.FilePath("secrets/pg_invalid")}},
		MaxOpenConns:    5,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
	}}

	lc := fxtest.NewLifecycle(t)
	_, err := test.NewLogger(lc, test.NewTextLoggerConfig())
	require.NoError(t, err)
	tc := test.NewOTLPTracerConfig()
	test.RegisterTracer(lc, tc)

	pg.Register()

	db, err := pg.Open(lc, test.FS, cfg)
	require.NoError(t, err)

	lc.RequireStart()
	require.Error(t, errors.Join(db.Ping()...))
	lc.RequireStop()
}

func setupAccounts(t *testing.T, db *sql.DBs) (context.Context, func()) {
	t.Helper()

	require.NoError(t, up(t.Context(), db))

	ctx, cancel := test.Timeout(t.Context())

	return meta.WithAttributes(ctx, test.WithTest(meta.String("test"))), func() {
		require.NoError(t, down(t.Context(), db))
		cancel()
	}
}

func up(ctx context.Context, db *sql.DBs) error {
	ctx, cancel := test.Timeout(ctx)
	defer cancel()

	query := `CREATE TABLE IF NOT EXISTS accounts (
		user_id serial PRIMARY KEY,
		created_at TIMESTAMP NOT NULL
	);`

	_, err := db.ExecContext(ctx, query)

	return err
}

func down(ctx context.Context, db *sql.DBs) error {
	ctx, cancel := test.Timeout(ctx)
	defer cancel()

	query := "DROP TABLE IF EXISTS accounts;"

	_, err := db.ExecContext(ctx, query)

	return err
}

func pgConfig(masters, slaves []config.DSN) *pg.Config {
	return &pg.Config{
		Config: &config.Config{
			Masters:         masters,
			Slaves:          slaves,
			MaxOpenConns:    5,
			MaxIdleConns:    5,
			ConnMaxLifetime: time.Hour,
		},
	}
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
