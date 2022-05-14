package pg_test

import (
	"context"
	"testing"

	"github.com/alexfalkowski/go-service/database/sql/pg"
	"github.com/alexfalkowski/go-service/database/sql/pg/trace/opentracing"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestSQL(t *testing.T) {
	Convey("Given I have a configuration", t, func() {
		cfg := &pg.Config{URL: "postgres://test:test@localhost:5432/test?sslmode=disable"}

		Convey("When I try to get a database", func() {
			lc := fxtest.NewLifecycle(t)
			logger := test.NewLogger(lc)

			tracer, err := opentracing.NewTracer(opentracing.TracerParams{Lifecycle: lc, Config: test.NewJaegerConfig(), Version: test.Version})
			So(err, ShouldBeNil)

			pg.Register(tracer, logger)

			db := pg.Open(pg.DBParams{Lifecycle: lc, Config: cfg, Version: test.Version})

			lc.RequireStart()

			Convey("Then I should have a valid database", func() {
				ctx := context.Background()
				err = db.PingContext(ctx)
				So(err, ShouldBeNil)

				rows, err := db.QueryContext(ctx, "SELECT table_name FROM information_schema.tables WHERE table_schema='public'")
				So(err, ShouldBeNil)

				for rows.Next() {
				}

				So(rows.Err(), ShouldBeNil)
				So(rows.Close(), ShouldBeNil)
			})

			lc.RequireStop()
		})
	})
}

func TestInvalidSQLPort(t *testing.T) {
	Convey("Given I have a configuration", t, func() {
		cfg := &pg.Config{URL: "postgres://test:test@localhost:5444/test?sslmode=disable"}

		Convey("When I try to get a database", func() {
			lc := fxtest.NewLifecycle(t)
			logger := test.NewLogger(lc)

			tracer, err := opentracing.NewTracer(opentracing.TracerParams{Lifecycle: lc, Config: test.NewDatadogConfig(), Version: test.Version})
			So(err, ShouldBeNil)

			pg.Register(tracer, logger)

			db := pg.Open(pg.DBParams{Lifecycle: lc, Config: cfg, Version: test.Version})

			lc.RequireStart()

			Convey("Then I should have an invalid database", func() {
				ctx := context.Background()
				err = db.PingContext(ctx)
				So(err, ShouldNotBeNil)
			})

			lc.RequireStop()
		})
	})
}
