//nolint:varnamelen
package pg_test

import (
	"errors"
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/database/sql/config"
	"github.com/alexfalkowski/go-service/database/sql/pg"
	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/test"
	"github.com/linxGnu/mssqlx"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
)

//nolint:gochecknoinits
func init() {
	tracer.Register()
}

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
		lc := fxtest.NewLifecycle(t)
		c := test.NewPGConfig()

		Convey("When I try open the database", func() {
			_, err := pg.Open(pg.OpenParams{Lifecycle: lc, Config: c})

			lc.RequireStart()

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})

			lc.RequireStop()
		})
	})
}

func TestSQL(t *testing.T) {
	Convey("Given I have a configuration", t, func() {
		Convey("When I try to get a database", func() {
			lc := fxtest.NewLifecycle(t)
			logger := test.NewLogger(lc)
			tc := test.NewOTLPTracerConfig()
			tracer, err := tracer.NewTracer(lc, test.Environment, test.Version, test.Name, tc, logger)
			So(err, ShouldBeNil)

			pg.Register(tracer, logger)

			db, err := pg.Open(pg.OpenParams{Lifecycle: lc, Config: test.NewPGConfig()})
			So(err, ShouldBeNil)
			So(db, ShouldNotBeNil)

			lc.RequireStart()

			Convey("Then I should have a valid database", func() {
				So(errors.Join(db.Ping()...), ShouldBeNil)
			})

			lc.RequireStop()
		})
	})
}

func TestDBQuery(t *testing.T) {
	Convey("Given I have a ready database", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)

		tc := test.NewOTLPTracerConfig()
		tracer, err := tracer.NewTracer(lc, test.Environment, test.Version, test.Name, tc, logger)
		So(err, ShouldBeNil)

		pg.Register(tracer, logger)

		db, err := pg.Open(pg.OpenParams{Lifecycle: lc, Config: test.NewPGConfig()})
		So(err, ShouldBeNil)

		lc.RequireStart()

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

		lc.RequireStop()
	})
}

func TestDBExec(t *testing.T) {
	Convey("Given I have a ready database", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		tc := test.NewOTLPTracerConfig()
		tracer, err := tracer.NewTracer(lc, test.Environment, test.Version, test.Name, tc, logger)
		So(err, ShouldBeNil)

		pg.Register(tracer, logger)

		db, err := pg.Open(pg.OpenParams{Lifecycle: lc, Config: test.NewPGConfig()})
		So(err, ShouldBeNil)

		lc.RequireStart()

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

		lc.RequireStop()
	})
}

func TestDBCommitTransExec(t *testing.T) {
	Convey("Given I have a ready database", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		tc := test.NewOTLPTracerConfig()
		tracer, err := tracer.NewTracer(lc, test.Environment, test.Version, test.Name, tc, logger)
		So(err, ShouldBeNil)

		pg.Register(tracer, logger)

		db, err := pg.Open(pg.OpenParams{Lifecycle: lc, Config: test.NewPGConfig()})
		So(err, ShouldBeNil)

		lc.RequireStart()

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

		lc.RequireStop()
	})
}

func TestDBRollbackTransExec(t *testing.T) {
	Convey("Given I have a ready database", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		tc := test.NewOTLPTracerConfig()
		tracer, err := tracer.NewTracer(lc, test.Environment, test.Version, test.Name, tc, logger)
		So(err, ShouldBeNil)

		pg.Register(tracer, logger)

		db, err := pg.Open(pg.OpenParams{Lifecycle: lc, Config: test.NewPGConfig()})
		So(err, ShouldBeNil)

		lc.RequireStart()

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

		lc.RequireStop()
	})
}

func TestStatementQuery(t *testing.T) {
	Convey("Given I have a ready database", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		tc := test.NewOTLPTracerConfig()
		tracer, err := tracer.NewTracer(lc, test.Environment, test.Version, test.Name, tc, logger)
		So(err, ShouldBeNil)

		pg.Register(tracer, logger)

		db, err := pg.Open(pg.OpenParams{Lifecycle: lc, Config: test.NewPGConfig()})
		So(err, ShouldBeNil)

		lc.RequireStart()

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

		lc.RequireStop()
	})
}

func TestStatementExec(t *testing.T) {
	Convey("Given I have a ready database", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		tc := test.NewOTLPTracerConfig()
		tracer, err := tracer.NewTracer(lc, test.Environment, test.Version, test.Name, tc, logger)
		So(err, ShouldBeNil)

		pg.Register(tracer, logger)

		db, err := pg.Open(pg.OpenParams{Lifecycle: lc, Config: test.NewPGConfig()})
		So(err, ShouldBeNil)

		lc.RequireStart()

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

		lc.RequireStop()
	})
}

func TestTransStatementExec(t *testing.T) {
	Convey("Given I have a ready database", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		tc := test.NewOTLPTracerConfig()
		tracer, err := tracer.NewTracer(lc, test.Environment, test.Version, test.Name, tc, logger)
		So(err, ShouldBeNil)

		pg.Register(tracer, logger)

		db, err := pg.Open(pg.OpenParams{Lifecycle: lc, Config: test.NewPGConfig()})
		So(err, ShouldBeNil)

		lc.RequireStart()

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

		lc.RequireStop()
	})
}

func TestInvalidStatementQuery(t *testing.T) {
	Convey("Given I have a ready database", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		tc := test.NewOTLPTracerConfig()
		tracer, err := tracer.NewTracer(lc, test.Environment, test.Version, test.Name, tc, logger)
		So(err, ShouldBeNil)

		pg.Register(tracer, logger)

		db, err := pg.Open(pg.OpenParams{Lifecycle: lc, Config: test.NewPGConfig()})
		So(err, ShouldBeNil)

		lc.RequireStart()

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

		lc.RequireStop()
	})
}

func TestInvalidSQLPort(t *testing.T) {
	Convey("Given I have a configuration", t, func() {
		cfg := &pg.Config{Config: &config.Config{
			Masters:         []config.DSN{{URL: test.Path("secrets/pg_invalid")}},
			Slaves:          []config.DSN{{URL: test.Path("secrets/pg_invalid")}},
			MaxOpenConns:    5,
			MaxIdleConns:    5,
			ConnMaxLifetime: time.Hour.String(),
		}}

		Convey("When I try to get a database", func() {
			lc := fxtest.NewLifecycle(t)
			logger := test.NewLogger(lc)
			tc := test.NewOTLPTracerConfig()
			tracer, err := tracer.NewTracer(lc, test.Environment, test.Version, test.Name, tc, logger)
			So(err, ShouldBeNil)

			pg.Register(tracer, logger)

			db, err := pg.Open(pg.OpenParams{Lifecycle: lc, Config: cfg})
			So(err, ShouldBeNil)

			lc.RequireStart()

			Convey("Then I should have an invalid database", func() {
				So(errors.Join(db.Ping()...), ShouldBeError)
			})

			lc.RequireStop()
		})
	})
}
