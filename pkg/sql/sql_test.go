package sql_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/pkg/sql"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestSQL(t *testing.T) {
	Convey("Given I have a configuration", t, func() {
		lc := fxtest.NewLifecycle(t)
		cfg := &sql.Config{
			PostgresURL: "postgres://test:test@localhost:5432/test?sslmode=disable",
		}

		Convey("When I try to get a database", func() {
			db, err := sql.NewDB(lc, cfg)
			So(err, ShouldBeNil)

			lc.RequireStart()

			Convey("Then I should have a valid database", func() {
				So(db, ShouldNotBeNil)

				lc.RequireStop()
			})
		})
	})
}

func TestInvalidSQL(t *testing.T) {
	Convey("Given I have an invalid configuration", t, func() {
		lc := fxtest.NewLifecycle(t)
		cfg := &sql.Config{
			PostgresURL: "invalid url",
		}

		lc.RequireStart()

		Convey("When I try to get a database", func() {
			_, err := sql.NewDB(lc, cfg)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)

				lc.RequireStop()
			})
		})
	})
}
