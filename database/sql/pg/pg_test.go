package pg_test

import (
	"context"
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/database/sql/config"
	"github.com/alexfalkowski/go-service/database/sql/pg"
	ptracer "github.com/alexfalkowski/go-service/database/sql/pg/telemetry/tracer"
	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
	"go.uber.org/multierr"
)

func init() {
	tracer.Register()
}

func TestSQL(t *testing.T) {
	Convey("Given I have a configuration", t, func() {
		Convey("When I try to get a database", func() {
			lc := fxtest.NewLifecycle(t)
			logger := test.NewLogger(lc)

			tracer, err := ptracer.NewTracer(ptracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
			So(err, ShouldBeNil)

			pg.Register(tracer, logger)

			db, err := pg.Open(pg.OpenParams{Lifecycle: lc, Config: test.NewPGConfig()})
			So(err, ShouldBeNil)
			So(db, ShouldNotBeNil)

			lc.RequireStart()

			Convey("Then I should have a valid database", func() {
				So(multierr.Combine(db.Ping()...), ShouldBeNil)
			})

			lc.RequireStop()
		})
	})
}

func TestDBQuery(t *testing.T) {
	Convey("Given I have a ready database", t, func() {
		m := test.NewMigrator()

		err := m.Up()
		So(err, ShouldBeNil)

		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)

		tracer, err := ptracer.NewTracer(ptracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		pg.Register(tracer, logger)

		db, err := pg.Open(pg.OpenParams{Lifecycle: lc, Config: test.NewPGConfig()})
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I select data with a query", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			ctx = meta.WithAttribute(ctx, "test", "test")

			//nolint:rowserrcheck
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

		lc.RequireStop()

		err = m.Drop()
		So(err, ShouldBeNil)
	})
}

func TestDBExec(t *testing.T) {
	Convey("Given I have a ready database", t, func() {
		m := test.NewMigrator()

		err := m.Up()
		So(err, ShouldBeNil)

		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)

		tracer, err := ptracer.NewTracer(ptracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		pg.Register(tracer, logger)

		db, err := pg.Open(pg.OpenParams{Lifecycle: lc, Config: test.NewPGConfig()})
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I insert data into a table", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			ctx = meta.WithAttribute(ctx, "test", "test")

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

		lc.RequireStop()

		err = m.Drop()
		So(err, ShouldBeNil)
	})
}

func TestDBCommitTransExec(t *testing.T) {
	Convey("Given I have a ready database", t, func() {
		m := test.NewMigrator()

		err := m.Up()
		So(err, ShouldBeNil)

		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)

		tracer, err := ptracer.NewTracer(ptracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		pg.Register(tracer, logger)

		db, err := pg.Open(pg.OpenParams{Lifecycle: lc, Config: test.NewPGConfig()})
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I insert data into a table", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			ctx = meta.WithAttribute(ctx, "test", "test")

			tx, err := db.BeginTx(ctx, nil)
			So(err, ShouldBeNil)

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

		lc.RequireStop()

		err = m.Drop()
		So(err, ShouldBeNil)
	})
}

func TestDBRollbackTransExec(t *testing.T) {
	Convey("Given I have a ready database", t, func() {
		m := test.NewMigrator()

		err := m.Up()
		So(err, ShouldBeNil)

		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)

		tracer, err := ptracer.NewTracer(ptracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		pg.Register(tracer, logger)

		db, err := pg.Open(pg.OpenParams{Lifecycle: lc, Config: test.NewPGConfig()})
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I insert data into a table", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			ctx = meta.WithAttribute(ctx, "test", "test")

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

		lc.RequireStop()

		err = m.Drop()
		So(err, ShouldBeNil)
	})
}

func TestStatementQuery(t *testing.T) {
	Convey("Given I have a ready database", t, func() {
		m := test.NewMigrator()

		err := m.Up()
		So(err, ShouldBeNil)

		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)

		tracer, err := ptracer.NewTracer(ptracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		pg.Register(tracer, logger)

		db, err := pg.Open(pg.OpenParams{Lifecycle: lc, Config: test.NewPGConfig()})
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I select data with a query", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			ctx = meta.WithAttribute(ctx, "test", "test")

			_, stmt, err := db.PrepareContext(ctx, "SELECT table_name FROM information_schema.tables WHERE table_schema = $1")
			So(err, ShouldBeNil)

			defer stmt.Close()

			//nolint:rowserrcheck
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

		lc.RequireStop()

		err = m.Drop()
		So(err, ShouldBeNil)
	})
}

func TestStatementExec(t *testing.T) {
	Convey("Given I have a ready database", t, func() {
		m := test.NewMigrator()

		err := m.Up()
		So(err, ShouldBeNil)

		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)

		tracer, err := ptracer.NewTracer(ptracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		pg.Register(tracer, logger)

		db, err := pg.Open(pg.OpenParams{Lifecycle: lc, Config: test.NewPGConfig()})
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I insert data into a table", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			ctx = meta.WithAttribute(ctx, "test", "test")

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

		lc.RequireStop()

		err = m.Drop()
		So(err, ShouldBeNil)
	})
}

func TestTransStatementExec(t *testing.T) {
	Convey("Given I have a ready database", t, func() {
		m := test.NewMigrator()

		err := m.Up()
		So(err, ShouldBeNil)

		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)

		tracer, err := ptracer.NewTracer(ptracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		pg.Register(tracer, logger)

		db, err := pg.Open(pg.OpenParams{Lifecycle: lc, Config: test.NewPGConfig()})
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I insert data into a table", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			ctx = meta.WithAttribute(ctx, "test", "test")

			tx, err := db.Begin()
			So(err, ShouldBeNil)

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

		lc.RequireStop()

		err = m.Drop()
		So(err, ShouldBeNil)
	})
}

func TestInvalidStatementQuery(t *testing.T) {
	Convey("Given I have a ready database", t, func() {
		m := test.NewMigrator()

		err := m.Up()
		So(err, ShouldBeNil)

		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)

		tracer, err := ptracer.NewTracer(ptracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		pg.Register(tracer, logger)

		db, err := pg.Open(pg.OpenParams{Lifecycle: lc, Config: test.NewPGConfig()})
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I select data with an invalid query", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			ctx = meta.WithAttribute(ctx, "test", "test")

			_, stmt, err := db.PrepareContext(ctx, "SELECT table_name FROM information_schema.tables WHERE table_schema = $1")
			So(err, ShouldBeNil)

			defer stmt.Close()

			//nolint:sqlclosecheck,rowserrcheck
			_, err = stmt.QueryContext(ctx, 1)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})

		lc.RequireStop()

		err = m.Drop()
		So(err, ShouldBeNil)
	})
}

func TestInvalidSQLPort(t *testing.T) {
	Convey("Given I have a configuration", t, func() {
		cfg := &pg.Config{Config: config.Config{
			Masters:         []config.DSN{{URL: "postgres://test:test@localhost:5444/test?sslmode=disable"}},
			Slaves:          []config.DSN{{URL: "postgres://test:test@localhost:5444/test?sslmode=disable"}},
			MaxOpenConns:    5,
			MaxIdleConns:    5,
			ConnMaxLifetime: time.Hour,
		}}

		Convey("When I try to get a database", func() {
			lc := fxtest.NewLifecycle(t)
			logger := test.NewLogger(lc)

			tracer, err := ptracer.NewTracer(ptracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
			So(err, ShouldBeNil)

			pg.Register(tracer, logger)

			db, err := pg.Open(pg.OpenParams{Lifecycle: lc, Config: cfg})
			So(err, ShouldBeNil)

			lc.RequireStart()

			Convey("Then I should have an invalid database", func() {
				So(multierr.Combine(db.Ping()...), ShouldBeError)
			})

			lc.RequireStop()
		})
	})
}
