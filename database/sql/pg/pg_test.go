package pg_test

import (
	"io/fs"
	"testing"

	"github.com/alexfalkowski/go-service/v2/database/sql/config"
	"github.com/alexfalkowski/go-service/v2/database/sql/driver"
	"github.com/alexfalkowski/go-service/v2/database/sql/pg"
	"github.com/alexfalkowski/go-service/v2/internal/test"
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
			name: "invalid writers",
			config: test.NewPGConfigWithDSNs(
				[]config.DSN{{URL: test.FilePath("secrets/none")}},
				[]config.DSN{{URL: test.FilePath("secrets/pg")}},
			),
			wantErr: fs.ErrNotExist,
		},
		{
			name: "invalid readers",
			config: test.NewPGConfigWithDSNs(
				[]config.DSN{{URL: test.FilePath("secrets/pg")}},
				[]config.DSN{{URL: test.FilePath("secrets/none")}},
			),
			wantErr: fs.ErrNotExist,
		},
		{
			name:    "empty dsn configuration",
			config:  test.NewPGConfigWithDSNs(nil, nil),
			wantErr: driver.ErrNoDSNs,
		},
		{
			name: "empty writer dsn",
			config: test.NewPGConfigWithDSNs(
				[]config.DSN{{}},
				[]config.DSN{{URL: test.FilePath("secrets/pg")}},
			),
			wantErr: driver.ErrEmptyDSN,
		},
		{
			name: "empty reader dsn",
			config: test.NewPGConfigWithDSNs(
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

	require.NoError(t, world.DB.Ping())
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

	test.RequireDBSystemName(t, reader, "postgresql")
}

func TestOpenClosesDBsOnStop(t *testing.T) {
	cfg := test.NewPGConfig()
	lc := fxtest.NewLifecycle(t)

	pg.Register()

	db, err := pg.Open(lc, test.FS, cfg)
	require.NoError(t, err)
	require.NotNil(t, db)

	lc.RequireStart()
	require.NoError(t, db.Ping())

	lc.RequireStop()
	require.Error(t, db.Ping())
}

func TestDBQuery(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldPGConfig(nil))

	ctx, cleanup := test.SetupAccounts(t, world.DB)
	defer cleanup()

	reader, err := world.DB.Reader()
	require.NoError(t, err)

	rows, err := reader.QueryContext(ctx, "SELECT table_name FROM information_schema.tables WHERE table_schema='public'")
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

	ctx, cleanup := test.SetupAccounts(t, world.DB)
	defer cleanup()

	writer, err := world.DB.Writer()
	require.NoError(t, err)

	result, err := writer.ExecContext(ctx, "INSERT INTO accounts(created_at) VALUES($1)", time.Now())
	require.NoError(t, err)

	num, err := result.RowsAffected()
	require.NoError(t, err)
	require.Positive(t, num)
}

func TestDBCommitTransExec(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldPGConfig(nil))

	ctx, cleanup := test.SetupAccounts(t, world.DB)
	defer cleanup()

	writer, err := world.DB.Writer()
	require.NoError(t, err)

	tx, err := writer.BeginTx(ctx, nil)
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

	ctx, cleanup := test.SetupAccounts(t, world.DB)
	defer cleanup()

	writer, err := world.DB.Writer()
	require.NoError(t, err)

	tx, err := writer.BeginTx(ctx, nil)
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

	ctx, cleanup := test.SetupAccounts(t, world.DB)
	defer cleanup()

	reader, err := world.DB.Reader()
	require.NoError(t, err)

	stmt, err := reader.PrepareContext(ctx, "SELECT table_name FROM information_schema.tables WHERE table_schema = $1")
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

	ctx, cleanup := test.SetupAccounts(t, world.DB)
	defer cleanup()

	writer, err := world.DB.Writer()
	require.NoError(t, err)

	stmt, err := writer.PrepareContext(ctx, "INSERT INTO accounts(created_at) VALUES($1)")
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

	ctx, cleanup := test.SetupAccounts(t, world.DB)
	defer cleanup()

	writer, err := world.DB.Writer()
	require.NoError(t, err)

	tx, err := writer.Begin()
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

	ctx, cleanup := test.SetupAccounts(t, world.DB)
	defer cleanup()

	reader, err := world.DB.Reader()
	require.NoError(t, err)

	stmt, err := reader.PrepareContext(ctx, "SELECT table_name FROM information_schema.tables WHERE table_schema = $1")
	require.NoError(t, err)

	defer stmt.Close()

	_, err = stmt.QueryContext(ctx, 1)
	require.Error(t, err)
}

func TestInvalidSQLPort(t *testing.T) {
	cfg := &pg.Config{Config: &config.Config{
		Writers:         []config.DSN{{URL: test.FilePath("secrets/pg_invalid")}},
		Readers:         []config.DSN{{URL: test.FilePath("secrets/pg_invalid")}},
		MaxOpenConns:    5,
		MaxIdleConns:    5,
		ConnMaxIdleTime: 30 * time.Minute,
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
	require.Error(t, db.Ping())
	lc.RequireStop()
}
