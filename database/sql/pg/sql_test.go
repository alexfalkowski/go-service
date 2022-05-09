package pg_test

import (
	"context"
	"testing"

	"github.com/alexfalkowski/go-service/database/sql/pg"
	"github.com/alexfalkowski/go-service/database/sql/pg/trace/opentracing"
	"github.com/alexfalkowski/go-service/test"
	"github.com/alexfalkowski/go-service/version"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

// nolint:dupl
func TestSQL(t *testing.T) {
	Convey("Given I have a configuration", t, func() {
		cfg := &pg.Config{URL: "postgres://test:test@localhost:5432/test?sslmode=disable"}

		Convey("When I try to get a database", func() {
			lc := fxtest.NewLifecycle(t)
			version := version.Version("1.0.0")

			tracer, err := opentracing.NewTracer(opentracing.TracerParams{Lifecycle: lc, Config: test.NewJaegerConfig(), Version: version})
			So(err, ShouldBeNil)

			ctx := context.Background()
			ctx, span := opentracing.StartSpanFromContext(ctx, tracer, "test", "test")
			defer span.Finish()

			db, err := pg.NewDB(pg.DBParams{Lifecycle: lc, Config: cfg, Version: version})
			So(err, ShouldBeNil)

			lc.RequireStart()

			Convey("Then I should have a valid database", func() {
				So(db, ShouldNotBeNil)

				err = db.PingContext(ctx)
				So(err, ShouldBeNil)
			})

			lc.RequireStop()
		})
	})
}

// nolint:dupl
func TestInvalidSQLPort(t *testing.T) {
	Convey("Given I have a configuration", t, func() {
		cfg := &pg.Config{URL: "postgres://test:test@localhost:5444/test?sslmode=disable"}

		Convey("When I try to get a database", func() {
			lc := fxtest.NewLifecycle(t)
			version := version.Version("1.0.0")

			tracer, err := opentracing.NewTracer(opentracing.TracerParams{Lifecycle: lc, Config: test.NewDatadogConfig(), Version: version})
			So(err, ShouldBeNil)

			ctx := context.Background()
			ctx, span := opentracing.StartSpanFromContext(ctx, tracer, "test", "test")
			defer span.Finish()

			db, err := pg.NewDB(pg.DBParams{Lifecycle: lc, Config: cfg, Version: version})
			So(err, ShouldBeNil)

			lc.RequireStart()

			Convey("Then I should have an invalid database", func() {
				So(db, ShouldNotBeNil)

				err = db.PingContext(ctx)
				So(err, ShouldNotBeNil)
			})

			lc.RequireStop()
		})
	})
}

func TestInvalidSQL(t *testing.T) {
	Convey("Given I have an invalid configuration", t, func() {
		cfg := &pg.Config{URL: "invalid url"}

		Convey("When I try to get a database", func() {
			lc := fxtest.NewLifecycle(t)
			_, err := pg.NewDB(pg.DBParams{Lifecycle: lc, Config: cfg, Version: version.Version("1.0.0")})

			lc.RequireStart()

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})

			lc.RequireStop()
		})
	})
}
