package pg_test

import (
	"io/fs"
	"testing"

	"github.com/alexfalkowski/go-service/v2/database/sql/config"
	"github.com/alexfalkowski/go-service/v2/database/sql/driver"
	"github.com/alexfalkowski/go-service/v2/database/sql/pg"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/linxGnu/mssqlx"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func up(db *mssqlx.DBs) error {
	ctx, cancel := test.Timeout()
	defer cancel()

	query := `CREATE TABLE IF NOT EXISTS accounts (
		user_id serial PRIMARY KEY,
		created_at TIMESTAMP NOT NULL
	);`

	_, err := db.ExecContext(ctx, query)

	return err
}

func down(db *mssqlx.DBs) error {
	ctx, cancel := test.Timeout()
	defer cancel()

	query := "DROP TABLE IF EXISTS accounts;"

	_, err := db.ExecContext(ctx, query)

	return err
}

func TestConnect(t *testing.T) {
	cfg := test.NewPGConfig()

	_, err := driver.Connect("missing", test.FS, cfg.Config)
	require.Error(t, err)
}

func TestInvalidOpen(t *testing.T) {
	tests := []struct {
		wantErr error
		config  *pg.Config
		name    string
	}{
		{
			name: "invalid masters",
			config: &pg.Config{
				Config: &config.Config{
					Masters:         []config.DSN{{URL: test.FilePath("secrets/none")}},
					Slaves:          []config.DSN{{URL: test.FilePath("secrets/pg")}},
					MaxOpenConns:    5,
					MaxIdleConns:    5,
					ConnMaxLifetime: time.Hour.String(),
				},
			},
			wantErr: fs.ErrNotExist,
		},
		{
			name: "invalid slaves",
			config: &pg.Config{
				Config: &config.Config{
					Masters:         []config.DSN{{URL: test.FilePath("secrets/pg")}},
					Slaves:          []config.DSN{{URL: test.FilePath("secrets/none")}},
					MaxOpenConns:    5,
					MaxIdleConns:    5,
					ConnMaxLifetime: time.Hour.String(),
				},
			},
			wantErr: fs.ErrNotExist,
		},
		{
			name: "empty dsn configuration",
			config: &pg.Config{
				Config: &config.Config{
					MaxOpenConns:    5,
					MaxIdleConns:    5,
					ConnMaxLifetime: time.Hour.String(),
				},
			},
			wantErr: driver.ErrNoDSNs,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldPGConfig(tt.config), test.WithWorldLoggerConfig("json"))

			err := world.Lifecycle.Start(t.Context())
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestSQL(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldPGConfig(nil), test.WithWorldLoggerConfig("otlp"))

	require.NoError(t, errors.Join(world.DB.Ping()...))
}

func TestDBQuery(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldPGConfig(nil))

	require.NoError(t, up(world.DB))

	ctx, cancel := test.Timeout()
	defer cancel()

	ctx = meta.WithAttribute(ctx, "test", meta.String("test"))
	rows, err := world.DB.QueryContext(ctx, "SELECT table_name FROM information_schema.tables WHERE table_schema='public'")
	require.NoError(t, err)

	var count int
	for rows.Next() {
		count++
	}
	require.Positive(t, count)
	require.NoError(t, rows.Err())
	require.NoError(t, rows.Close())

	require.NoError(t, down(world.DB))
}

func TestDBExec(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldPGConfig(nil))

	require.NoError(t, up(world.DB))

	ctx, cancel := test.Timeout()
	defer cancel()

	ctx = meta.WithAttribute(ctx, "test", meta.String("test"))

	result, err := world.DB.ExecContext(ctx, "INSERT INTO accounts(created_at) VALUES($1)", time.Now())
	require.NoError(t, err)

	num, err := result.RowsAffected()
	require.NoError(t, err)
	require.Positive(t, num)

	require.NoError(t, down(world.DB))
}

func TestDBCommitTransExec(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldPGConfig(nil))

	require.NoError(t, up(world.DB))

	ctx, cancel := test.Timeout()
	defer cancel()

	ctx = meta.WithAttribute(ctx, "test", meta.String("test"))

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

	require.NoError(t, down(world.DB))
}

func TestDBRollbackTransExec(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldPGConfig(nil))

	require.NoError(t, up(world.DB))

	ctx, cancel := test.Timeout()
	defer cancel()

	ctx = meta.WithAttribute(ctx, "test", meta.String("test"))

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

	require.NoError(t, down(world.DB))
}

func TestStatementQuery(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldPGConfig(nil))

	require.NoError(t, up(world.DB))

	ctx, cancel := test.Timeout()
	defer cancel()

	ctx = meta.WithAttribute(ctx, "test", meta.String("test"))

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

	require.NoError(t, down(world.DB))
}

func TestStatementExec(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldPGConfig(nil))

	require.NoError(t, up(world.DB))

	ctx, cancel := test.Timeout()
	defer cancel()

	ctx = meta.WithAttribute(ctx, "test", meta.String("test"))

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

	require.NoError(t, down(world.DB))
}

func TestTransStatementExec(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldPGConfig(nil))

	require.NoError(t, up(world.DB))

	ctx, cancel := test.Timeout()
	defer cancel()

	ctx = meta.WithAttribute(ctx, "test", meta.String("test"))

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

	require.NoError(t, down(world.DB))
}

func TestInvalidStatementQuery(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldPGConfig(nil), test.WithWorldLoggerConfig("tint"))

	require.NoError(t, up(world.DB))

	ctx, cancel := test.Timeout()
	defer cancel()

	ctx = meta.WithAttribute(ctx, "test", meta.String("test"))

	_, stmt, err := world.DB.PrepareContext(ctx, "SELECT table_name FROM information_schema.tables WHERE table_schema = $1")
	require.NoError(t, err)

	defer stmt.Close()

	_, err = stmt.QueryContext(ctx, 1)
	require.Error(t, err)

	require.NoError(t, down(world.DB))
}

func TestInvalidSQLPort(t *testing.T) {
	cfg := &pg.Config{Config: &config.Config{
		Masters:         []config.DSN{{URL: test.FilePath("secrets/pg_invalid")}},
		Slaves:          []config.DSN{{URL: test.FilePath("secrets/pg_invalid")}},
		MaxOpenConns:    5,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour.String(),
	}}

	lc := fxtest.NewLifecycle(t)
	_ = test.NewLogger(lc, test.NewTextLoggerConfig())
	tc := test.NewOTLPTracerConfig()
	test.RegisterTracer(lc, tc)

	pg.Register()

	db, err := pg.Open(lc, test.FS, cfg)
	require.NoError(t, err)

	lc.RequireStart()
	require.Error(t, errors.Join(db.Ping()...))
	lc.RequireStop()
}
