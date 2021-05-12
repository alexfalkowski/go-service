package sql_test

import (
	"os"
	"testing"

	"github.com/alexfalkowski/go-service/pkg/sql"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestSQL(t *testing.T) {
	Convey("Given I have a configuration", t, func() {
		os.Setenv("APP_NAME", "test")
		os.Setenv("POSTGRES_URL", "postgres://test:test@localhost:5432/test?sslmode=disable")

		cfg, err := sql.NewConfig()
		So(err, ShouldBeNil)

		Convey("When I try to get a database", func() {
			lc := fxtest.NewLifecycle(t)

			db, err := sql.NewDB(lc, cfg)
			So(err, ShouldBeNil)

			lc.RequireStart()

			Convey("Then I should have a valid database", func() {
				So(db, ShouldNotBeNil)
			})

			lc.RequireStop()
		})

		So(os.Unsetenv("APP_NAME"), ShouldBeNil)
		So(os.Unsetenv("POSTGRES_URL"), ShouldBeNil)
	})
}

func TestInvalidSQL(t *testing.T) {
	Convey("Given I have an invalid configuration", t, func() {
		cfg := &sql.Config{
			PostgresURL: "invalid url",
		}

		Convey("When I try to get a database", func() {
			lc := fxtest.NewLifecycle(t)
			_, err := sql.NewDB(lc, cfg)

			lc.RequireStart()

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})

			lc.RequireStop()
		})
	})
}
