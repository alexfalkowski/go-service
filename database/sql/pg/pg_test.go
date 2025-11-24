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
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func up(db *mssqlx.DBs) error {
	ctx, cancel := test.Timeout()
	defer cancel()

	query := `CREATE TABLE accounts (
		user_id serial PRIMARY KEY,
		created_at TIMESTAMP NOT NULL
	);`

	_, err := db.ExecContext(ctx, query)

	return err
}

func down(db *mssqlx.DBs) error {
	ctx, cancel := test.Timeout()
	defer cancel()

	query := "DROP TABLE accounts;"

	_, err := db.ExecContext(ctx, query)

	return err
}

func TestOpen(t *testing.T) {
	Convey("Given I have a configuration", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldLoggerConfig("text"))

		Convey("When I try open the database", func() {
			_, err := world.OpenDatabase()

			world.RequireStart()

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})

			world.RequireStop()
		})
	})
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
		Convey("Given I have an invalid config", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldPGConfig(config), test.WithWorldLoggerConfig("json"))
			world.Register()

			Convey("When I try open the database", func() {
				_, err := world.OpenDatabase()

				world.RequireStart()

				Convey("Then I should have an error", func() {
					So(err, ShouldBeError)
				})

				world.RequireStop()
			})
		})
	}
}

func TestSQL(t *testing.T) {
	Convey("Given I have a configuration", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldLoggerConfig("otlp"))
		world.Register()

		Convey("When I try to get a database", func() {
			db, err := world.OpenDatabase()
			So(err, ShouldBeNil)

			world.RequireStart()

			Convey("Then I should have a valid database", func() {
				So(errors.Join(db.Ping()...), ShouldBeNil)
			})

			world.RequireStop()
		})
	})
}

func TestDBQuery(t *testing.T) {
	Convey("Given I have a ready database", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"))
		world.Register()

		db, err := world.OpenDatabase()
		So(err, ShouldBeNil)

		world.RequireStart()

		err = up(db)
		So(err, ShouldBeNil)

		Convey("When I select data with a query", func() {
			ctx, cancel := test.Timeout()
			defer cancel()

			ctx = meta.WithAttribute(ctx, "test", meta.String("test"))
			rows, err := db.QueryContext(ctx, "SELECT table_name FROM information_schema.tables WHERE table_schema='public'")

			Convey("Then I should have valid data", func() {
				So(err, ShouldBeNil)

				count := 0

				for rows.Next() {
					count++
				}

				So(count, ShouldBeGreaterThan, 0)
				So(rows.Err(), ShouldBeNil)
				So(rows.Close(), ShouldBeNil)
			})
		})

		err = down(db)
		So(err, ShouldBeNil)

		world.RequireStop()
	})
}

func TestDBExec(t *testing.T) {
	Convey("Given I have a ready database", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"))
		world.Register()

		db, err := world.OpenDatabase()
		So(err, ShouldBeNil)

		world.RequireStart()

		err = up(db)
		So(err, ShouldBeNil)

		Convey("When I insert data into a table", func() {
			ctx, cancel := test.Timeout()
			defer cancel()

			ctx = meta.WithAttribute(ctx, "test", meta.String("test"))

			result, err := db.ExecContext(ctx, "INSERT INTO accounts(created_at) VALUES($1)", time.Now())
			So(err, ShouldBeNil)

			Convey("Then I should have successfully inserted data", func() {
				_, err := result.LastInsertId()
				So(err, ShouldBeError)

				num, err := result.RowsAffected()
				So(err, ShouldBeNil)
				So(num, ShouldBeGreaterThan, 0)
			})
		})

		err = down(db)
		So(err, ShouldBeNil)

		world.RequireStop()
	})
}

func TestDBCommitTransExec(t *testing.T) {
	Convey("Given I have a ready database", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"))
		world.Register()

		db, err := world.OpenDatabase()
		So(err, ShouldBeNil)

		world.RequireStart()

		err = up(db)
		So(err, ShouldBeNil)

		Convey("When I insert data into a table", func() {
			ctx, cancel := test.Timeout()
			defer cancel()

			ctx = meta.WithAttribute(ctx, "test", meta.String("test"))

			tx, err := db.BeginTx(ctx, nil)
			So(err, ShouldBeNil)

			//nolint:errcheck
			defer tx.Rollback()

			result, err := tx.ExecContext(ctx, "INSERT INTO accounts(created_at) VALUES($1)", time.Now())
			So(err, ShouldBeNil)

			err = tx.Commit()
			So(err, ShouldBeNil)

			Convey("Then I should have successfully inserted data", func() {
				_, err := result.LastInsertId()
				So(err, ShouldBeError)

				num, err := result.RowsAffected()
				So(err, ShouldBeNil)
				So(num, ShouldBeGreaterThan, 0)
			})
		})

		err = down(db)
		So(err, ShouldBeNil)

		world.RequireStop()
	})
}

func TestDBRollbackTransExec(t *testing.T) {
	Convey("Given I have a ready database", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"))
		world.Register()

		db, err := world.OpenDatabase()
		So(err, ShouldBeNil)

		world.RequireStart()

		err = up(db)
		So(err, ShouldBeNil)

		Convey("When I insert data into a table", func() {
			ctx, cancel := test.Timeout()
			defer cancel()

			ctx = meta.WithAttribute(ctx, "test", meta.String("test"))

			tx, err := db.BeginTx(ctx, nil)
			So(err, ShouldBeNil)

			result, err := tx.ExecContext(ctx, "INSERT INTO accounts(created_at) VALUES($1)", time.Now())
			So(err, ShouldBeNil)

			err = tx.Rollback()
			So(err, ShouldBeNil)

			Convey("Then I should have successfully inserted data", func() {
				_, err := result.LastInsertId()
				So(err, ShouldBeError)

				num, err := result.RowsAffected()
				So(err, ShouldBeNil)
				So(num, ShouldBeGreaterThan, 0)
			})
		})

		err = down(db)
		So(err, ShouldBeNil)

		world.RequireStop()
	})
}

func TestStatementQuery(t *testing.T) {
	Convey("Given I have a ready database", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"))
		world.Register()

		db, err := world.OpenDatabase()
		So(err, ShouldBeNil)

		world.RequireStart()

		err = up(db)
		So(err, ShouldBeNil)

		Convey("When I select data with a query", func() {
			ctx, cancel := test.Timeout()
			defer cancel()

			ctx = meta.WithAttribute(ctx, "test", meta.String("test"))

			_, stmt, err := db.PrepareContext(ctx, "SELECT table_name FROM information_schema.tables WHERE table_schema = $1")
			So(err, ShouldBeNil)

			defer stmt.Close()

			rows, err := stmt.QueryContext(ctx, "public")

			Convey("Then I should have valid data", func() {
				So(err, ShouldBeNil)

				count := 0

				for rows.Next() {
					count++
				}

				So(count, ShouldBeGreaterThan, 0)
				So(rows.Err(), ShouldBeNil)
				So(rows.Close(), ShouldBeNil)
			})
		})

		err = down(db)
		So(err, ShouldBeNil)

		world.RequireStop()
	})
}

func TestStatementExec(t *testing.T) {
	Convey("Given I have a ready database", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"))
		world.Register()

		db, err := world.OpenDatabase()
		So(err, ShouldBeNil)

		world.RequireStart()

		err = up(db)
		So(err, ShouldBeNil)

		Convey("When I insert data into a table", func() {
			ctx, cancel := test.Timeout()
			defer cancel()

			ctx = meta.WithAttribute(ctx, "test", meta.String("test"))

			_, stmt, err := db.PrepareContext(ctx, "INSERT INTO accounts(created_at) VALUES($1)")
			So(err, ShouldBeNil)

			defer stmt.Close()

			result, err := stmt.ExecContext(ctx, time.Now())
			So(err, ShouldBeNil)

			Convey("Then I should have successfully inserted data", func() {
				_, err := result.LastInsertId()
				So(err, ShouldBeError)

				num, err := result.RowsAffected()
				So(err, ShouldBeNil)
				So(num, ShouldBeGreaterThan, 0)
			})
		})

		err = down(db)
		So(err, ShouldBeNil)

		world.RequireStop()
	})
}

func TestTransStatementExec(t *testing.T) {
	Convey("Given I have a ready database", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"))
		world.Register()

		db, err := world.OpenDatabase()
		So(err, ShouldBeNil)

		world.RequireStart()

		err = up(db)
		So(err, ShouldBeNil)

		Convey("When I insert data into a table", func() {
			ctx, cancel := test.Timeout()
			defer cancel()

			ctx = meta.WithAttribute(ctx, "test", meta.String("test"))

			tx, err := db.Begin()
			So(err, ShouldBeNil)

			//nolint:errcheck
			defer tx.Rollback()

			stmt, err := tx.PrepareContext(ctx, "INSERT INTO accounts(created_at) VALUES($1)")
			So(err, ShouldBeNil)

			defer stmt.Close()

			result, err := stmt.ExecContext(ctx, time.Now())
			So(err, ShouldBeNil)

			err = tx.Commit()
			So(err, ShouldBeNil)

			Convey("Then I should have successfully inserted data", func() {
				_, err := result.LastInsertId()
				So(err, ShouldBeError)

				num, err := result.RowsAffected()
				So(err, ShouldBeNil)
				So(num, ShouldBeGreaterThan, 0)
			})
		})

		err = down(db)
		So(err, ShouldBeNil)

		world.RequireStop()
	})
}

func TestInvalidStatementQuery(t *testing.T) {
	Convey("Given I have a ready database", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldLoggerConfig("tint"))
		world.Register()

		db, err := world.OpenDatabase()
		So(err, ShouldBeNil)

		world.RequireStart()

		err = up(db)
		So(err, ShouldBeNil)

		Convey("When I select data with an invalid query", func() {
			ctx, cancel := test.Timeout()
			defer cancel()

			ctx = meta.WithAttribute(ctx, "test", meta.String("test"))

			_, stmt, err := db.PrepareContext(ctx, "SELECT table_name FROM information_schema.tables WHERE table_schema = $1")
			So(err, ShouldBeNil)

			defer stmt.Close()

			_, err = stmt.QueryContext(ctx, 1)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})

		err = down(db)
		So(err, ShouldBeNil)

		world.RequireStop()
	})
}

func TestInvalidSQLPort(t *testing.T) {
	Convey("Given I have a configuration", t, func() {
		cfg := &pg.Config{Config: &config.Config{
			Masters:         []config.DSN{{URL: test.FilePath("secrets/pg_invalid")}},
			Slaves:          []config.DSN{{URL: test.FilePath("secrets/pg_invalid")}},
			MaxOpenConns:    5,
			MaxIdleConns:    5,
			ConnMaxLifetime: time.Hour.String(),
		}}

		Convey("When I try to get a database", func() {
			lc := fxtest.NewLifecycle(t)
			_ = test.NewLogger(lc, test.NewTextLoggerConfig())
			tc := test.NewOTLPTracerConfig()
			_ = test.NewTracer(lc, tc)

			pg.Register()

			db, err := pg.Open(lc, test.FS, cfg)
			So(err, ShouldBeNil)

			lc.RequireStart()

			Convey("Then I should have an invalid database", func() {
				So(errors.Join(db.Ping()...), ShouldBeError)
			})

			lc.RequireStop()
		})
	})
}
