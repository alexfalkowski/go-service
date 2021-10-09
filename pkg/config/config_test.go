package config_test

import (
	"os"
	"testing"

	"github.com/alexfalkowski/go-service/pkg/config"
	. "github.com/smartystreets/goconvey/convey"
)

func TestValidConfig(t *testing.T) {
	Convey("Given I have configuration file", t, func() {
		os.Setenv("CONFIG_FILE", "../../test/config.yml")

		Convey("When I try to parse the configuration file", func() {
			cfg := config.NewConfigurator()
			err := config.Unmarshal(cfg)
			So(err, ShouldBeNil)

			Convey("Then I should have a valid configuration", func() {
				So(cfg.RedisConfig().Host, ShouldEqual, "localhost:6379")
				So(cfg.RistrettoConfig().BufferItems, ShouldEqual, 64)
				So(cfg.RistrettoConfig().MaxCost, ShouldEqual, 100000000)
				So(cfg.RistrettoConfig().NumCounters, ShouldEqual, 10000000)
				So(cfg.Auth0Config().URL, ShouldEqual, "test_url")
				So(cfg.Auth0Config().ClientID, ShouldEqual, "test_client_id")
				So(cfg.Auth0Config().ClientSecret, ShouldEqual, "test_client_secret")
				So(cfg.Auth0Config().Audience, ShouldEqual, "test_audience")
				So(cfg.Auth0Config().Issuer, ShouldEqual, "test_issuer")
				So(cfg.Auth0Config().Algorithm, ShouldEqual, "test_algorithm")
				So(cfg.Auth0Config().JSONWebKeySet, ShouldEqual, "test_json_web_key_set")
				So(cfg.PGConfig().URL, ShouldEqual, "postgres://test:test@localhost:5432/test?sslmode=disable")
				So(cfg.DatadogConfig().Host, ShouldEqual, "localhost:6831")
				So(cfg.JaegerConfig().Host, ShouldEqual, "localhost:6379")
				So(cfg.GRPCConfig().Port, ShouldEqual, "9000")
				So(cfg.HTTPConfig().Port, ShouldEqual, "8000")
				So(cfg.NSQConfig().Host, ShouldEqual, "localhost:4150")
				So(cfg.NSQConfig().LookupHost, ShouldEqual, "localhost:4161")
			})

			So(os.Unsetenv("CONFIG_FILE"), ShouldBeNil)
		})
	})
}

func TestMissingConfig(t *testing.T) {
	Convey("Given I have configuration file", t, func() {
		Convey("When I try to parse the configuration file", func() {
			err := config.Unmarshal(config.NewConfigurator())

			Convey("Then I should have an error of missing config file", func() {
				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, "missing config file")
			})
		})
	})
}

func TestNonExistentConfig(t *testing.T) {
	Convey("Given I have configuration file", t, func() {
		os.Setenv("CONFIG_FILE", "../../test/bob")

		Convey("When I try to parse the configuration file", func() {
			err := config.Unmarshal(config.NewConfigurator())

			Convey("Then I should have an error of non existent config file", func() {
				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, "open ../../test/bob: no such file or directory")
			})

			So(os.Unsetenv("CONFIG_FILE"), ShouldBeNil)
		})
	})
}

func TestInvalidConfig(t *testing.T) {
	Convey("Given I have configuration file", t, func() {
		os.Setenv("CONFIG_FILE", "../../test/greet.proto")

		Convey("When I try to parse the configuration file", func() {
			err := config.Unmarshal(config.NewConfigurator())

			Convey("Then I should have an error of invalid configuration file", func() {
				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, "yaml: line 12: mapping values are not allowed in this context")
			})

			So(os.Unsetenv("CONFIG_FILE"), ShouldBeNil)
		})
	})
}
