package pg_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/database/sql/config"
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

func TestOpen(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldLoggerConfig("text"))

	_, err := world.OpenDatabase()
	world.RequireStart()
	require.Error(t, err)

	world.RequireStop()
}

func TestInvalidOpen(t *testing.T) {
	configs := []*pg.Config{
		{
			Config: &config.Config{
				Masters:         []config.DSN{{URL: test.FilePath("secrets/none")}},
				Slaves:          []config.DSN{{URL: test.FilePath("secrets/pg")}},
				MaxOpenConns:    5,
				MaxIdleConns:    5,
				ConnMaxLifetime: time.Hour.String(),
			},
		},
		{
			Config: &config.Config{
				Masters:         []config.DSN{{URL: test.FilePath("secrets/pg")}},
				Slaves:          []config.DSN{{URL: test.FilePath("secrets/none")}},
				MaxOpenConns:    5,
				MaxIdleConns:    5,
				ConnMaxLifetime: time.Hour.String(),
			},
		},
	}

	for _, config := range configs {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldPGConfig(config), test.WithWorldLoggerConfig("json"))
		world.Register()

		_, err := world.OpenDatabase()
		world.RequireStart()
		require.Error(t, err)

		world.RequireStop()
	}
}

func TestSQL(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldLoggerConfig("otlp"))
	world.Register()

	db, err := world.OpenDatabase()
	require.NoError(t, err)

	world.RequireStart()
	require.NoError(t, errors.Join(db.Ping()...))

	world.RequireStop()
}

func TestDBQuery(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"))
	world.Register()

	db, err := world.OpenDatabase()
	require.NoError(t, err)

	world.RequireStart()
	require.NoError(t, up(db))

	ctx, cancel := test.Timeout()
	defer cancel()

	ctx = meta.WithAttribute(ctx, "test", meta.String("test"))
	rows, err := db.QueryContext(ctx, "SELECT table_name FROM information_schema.tables WHERE table_schema='public'")
	require.NoError(t, err)

	var count int
	for rows.Next() {
		count++
	}
	require.Positive(t, count)
	require.NoError(t, rows.Err())
	require.NoError(t, rows.Close())

	require.NoError(t, down(db))
	world.RequireStop()
}

func TestDBExec(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"))
	world.Register()

	db, err := world.OpenDatabase()
	require.NoError(t, err)

	world.RequireStart()
	require.NoError(t, up(db))

	ctx, cancel := test.Timeout()
	defer cancel()

	ctx = meta.WithAttribute(ctx, "test", meta.String("test"))

	result, err := db.ExecContext(ctx, "INSERT INTO accounts(created_at) VALUES($1)", time.Now())
	require.NoError(t, err)

	num, err := result.RowsAffected()
	require.NoError(t, err)
	require.Positive(t, num)

	require.NoError(t, down(db))
	world.RequireStop()
}

func TestDBCommitTransExec(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"))
	world.Register()

	db, err := world.OpenDatabase()
	require.NoError(t, err)

	world.RequireStart()
	require.NoError(t, up(db))

	ctx, cancel := test.Timeout()
	defer cancel()

	ctx = meta.WithAttribute(ctx, "test", meta.String("test"))

	tx, err := db.BeginTx(ctx, nil)
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

	require.NoError(t, down(db))
	world.RequireStop()
}

func TestDBRollbackTransExec(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"))
	world.Register()

	db, err := world.OpenDatabase()
	require.NoError(t, err)

	world.RequireStart()
	require.NoError(t, up(db))

	ctx, cancel := test.Timeout()
	defer cancel()

	ctx = meta.WithAttribute(ctx, "test", meta.String("test"))

	tx, err := db.BeginTx(ctx, nil)
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

	require.NoError(t, down(db))
	world.RequireStop()
}

func TestStatementQuery(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"))
	world.Register()

	db, err := world.OpenDatabase()
	require.NoError(t, err)

	world.RequireStart()
	require.NoError(t, up(db))

	ctx, cancel := test.Timeout()
	defer cancel()

	ctx = meta.WithAttribute(ctx, "test", meta.String("test"))

	_, stmt, err := db.PrepareContext(ctx, "SELECT table_name FROM information_schema.tables WHERE table_schema = $1")
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

	require.NoError(t, down(db))
	world.RequireStop()
}

func TestStatementExec(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"))
	world.Register()

	db, err := world.OpenDatabase()
	require.NoError(t, err)

	world.RequireStart()
	require.NoError(t, up(db))

	ctx, cancel := test.Timeout()
	defer cancel()

	ctx = meta.WithAttribute(ctx, "test", meta.String("test"))

	_, stmt, err := db.PrepareContext(ctx, "INSERT INTO accounts(created_at) VALUES($1)")
	require.NoError(t, err)

	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, time.Now())
	require.NoError(t, err)

	_, err = result.LastInsertId()
	require.Error(t, err)

	num, err := result.RowsAffected()
	require.NoError(t, err)
	require.Positive(t, num)

	require.NoError(t, down(db))
	world.RequireStop()
}

func TestTransStatementExec(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"))
	world.Register()

	db, err := world.OpenDatabase()
	require.NoError(t, err)

	world.RequireStart()
	require.NoError(t, up(db))

	ctx, cancel := test.Timeout()
	defer cancel()

	ctx = meta.WithAttribute(ctx, "test", meta.String("test"))

	tx, err := db.Begin()
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

	require.NoError(t, down(db))
	world.RequireStop()
}

func TestInvalidStatementQuery(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldLoggerConfig("tint"))
	world.Register()

	db, err := world.OpenDatabase()
	require.NoError(t, err)

	world.RequireStart()
	require.NoError(t, up(db))

	ctx, cancel := test.Timeout()
	defer cancel()

	ctx = meta.WithAttribute(ctx, "test", meta.String("test"))

	_, stmt, err := db.PrepareContext(ctx, "SELECT table_name FROM information_schema.tables WHERE table_schema = $1")
	require.NoError(t, err)

	defer stmt.Close()

	_, err = stmt.QueryContext(ctx, 1)
	require.Error(t, err)

	require.NoError(t, down(db))
	world.RequireStop()
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
