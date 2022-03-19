package pg_test

import (
	"context"
	"testing"

	"github.com/alexfalkowski/go-service/sql/pg"
	sotr "github.com/alexfalkowski/go-service/sql/trace/opentracing"
	totr "github.com/alexfalkowski/go-service/transport/trace/opentracing"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestSQL(t *testing.T) {
	Convey("Given I have a configuration", t, func() {
		cfg := &pg.Config{
			URL: "postgres://test:test@localhost:5432/test?sslmode=disable",
		}

		Convey("When I try to get a database", func() {
			lc := fxtest.NewLifecycle(t)

			ctx := context.Background()
			ctx, span := sotr.StartSpanFromContext(ctx, "test", "test", totr.StartSpanOptions(ctx)...)
			defer span.Finish()

			db, err := pg.NewDB(lc, cfg)
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

func TestInvalidSQL(t *testing.T) {
	Convey("Given I have an invalid configuration", t, func() {
		cfg := &pg.Config{URL: "invalid url"}

		Convey("When I try to get a database", func() {
			lc := fxtest.NewLifecycle(t)
			_, err := pg.NewDB(lc, cfg)

			lc.RequireStart()

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})

			lc.RequireStop()
		})
	})
}
