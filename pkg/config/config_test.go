package config_test

import (
	"os"
	"testing"

	"github.com/alexfalkowski/go-service/pkg/config"
	. "github.com/smartystreets/goconvey/convey"
)

func TestConfig(t *testing.T) {
	Convey("Given I have environment variable set", t, func() {
		os.Setenv("APP_NAME", "test")
		os.Setenv("HTTP_PORT", "8000")
		os.Setenv("GRPC_PORT", "9000")
		os.Setenv("DATABASE_URL", "postgres://test:test@localhost:5432/test?sslmode=disable")

		Convey("When I get the config", func() {
			cfg, err := config.NewConfig()
			So(err, ShouldBeNil)

			Convey("Then I should have valid config", func() {
				So(cfg.AppName, ShouldEqual, "test")
				So(cfg.HTTPPort, ShouldEqual, "8000")
				So(cfg.GRPCPort, ShouldEqual, "9000")
				So(cfg.DatabaseURL, ShouldEqual, "postgres://test:test@localhost:5432/test?sslmode=disable")
			})

			So(os.Unsetenv("APP_NAME"), ShouldBeNil)
			So(os.Unsetenv("HTTP_PORT"), ShouldBeNil)
			So(os.Unsetenv("GRPC_PORT"), ShouldBeNil)
			So(os.Unsetenv("DATABASE_URL"), ShouldBeNil)
		})
	})
}

func TestMissingConfig(t *testing.T) {
	Convey("Given I have some environment variable set", t, func() {
		os.Setenv("HTTP_PORT", "8000")
		os.Setenv("GRPC_PORT", "9000")

		Convey("When I get the config", func() {
			_, err := config.NewConfig()

			Convey("Then I should have an error getting config", func() {
				So(err, ShouldBeError)
			})
		})

		So(os.Unsetenv("HTTP_PORT"), ShouldBeNil)
		So(os.Unsetenv("GRPC_PORT"), ShouldBeNil)
	})
}
