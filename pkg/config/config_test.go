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
			cfg, err := config.New()
			So(err, ShouldBeNil)

			Convey("Then I should have a valid configuration", func() {
				So(cfg.Cache.Redis.Host, ShouldEqual, "localhost:6379")
				So(cfg.Cache.Redis.Name, ShouldEqual, "test")
				So(cfg.Cache.Ristretto.Name, ShouldEqual, "test")
				So(cfg.Cache.Ristretto.BufferItems, ShouldEqual, 64)
				So(cfg.Cache.Ristretto.MaxCost, ShouldEqual, 100000000)
				So(cfg.Cache.Ristretto.NumCounters, ShouldEqual, 10000000)
				So(cfg.Security.Auth0.URL, ShouldEqual, "test_url")
				So(cfg.Security.Auth0.ClientID, ShouldEqual, "test_client_id")
				So(cfg.Security.Auth0.ClientSecret, ShouldEqual, "test_client_secret")
				So(cfg.Security.Auth0.Audience, ShouldEqual, "test_audience")
				So(cfg.Security.Auth0.Issuer, ShouldEqual, "test_issuer")
				So(cfg.Security.Auth0.Algorithm, ShouldEqual, "test_algorithm")
				So(cfg.Security.Auth0.JSONWebKeySet, ShouldEqual, "test_json_web_key_set")
				So(cfg.SQL.PG.Name, ShouldEqual, "test")
				So(cfg.SQL.PG.URL, ShouldEqual, "postgres://test:test@localhost:5432/test?sslmode=disable")
				So(cfg.Trace.Opentracing.Datadog.Host, ShouldEqual, "localhost:6831")
				So(cfg.Trace.Opentracing.Datadog.Name, ShouldEqual, "test")
				So(cfg.Trace.Opentracing.Jaeger.Host, ShouldEqual, "localhost:6379")
				So(cfg.Trace.Opentracing.Jaeger.Name, ShouldEqual, "test")
				So(cfg.Transport.GRPC.Port, ShouldEqual, "9000")
				So(cfg.Transport.HTTP.Port, ShouldEqual, "8000")
				So(cfg.Transport.NSQ.Host, ShouldEqual, "localhost:4150")
				So(cfg.Transport.NSQ.LookupHost, ShouldEqual, "localhost:4161")
			})

			So(os.Unsetenv("CONFIG_FILE"), ShouldBeNil)
		})
	})
}

func TestMissingConfig(t *testing.T) {
	Convey("Given I have configuration file", t, func() {
		Convey("When I try to parse the configuration file", func() {
			_, err := config.New()

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
			_, err := config.New()

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
			_, err := config.New()

			Convey("Then I should have an error of invalid configuration file", func() {
				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, "yaml: line 12: mapping values are not allowed in this context")
			})

			So(os.Unsetenv("CONFIG_FILE"), ShouldBeNil)
		})
	})
}
