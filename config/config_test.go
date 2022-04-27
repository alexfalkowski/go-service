package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/config"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/yaml.v3"
)

func TestValidConfig(t *testing.T) {
	Convey("Given I have configuration file", t, func() {
		os.Setenv("CONFIG_FILE", "../test/config.yml")

		Convey("When I try to parse the configuration file", func() {
			cfg := config.NewConfigurator()
			err := config.UnmarshalFromFile(cfg)
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
				So(cfg.DatadogConfig().Host, ShouldEqual, "localhost:8126")
				So(cfg.JaegerConfig().Host, ShouldEqual, "localhost:6831")
				So(cfg.GRPCConfig().Port, ShouldEqual, "9090")
				So(cfg.GRPCConfig().Retry.Attempts, ShouldEqual, 3)
				So(cfg.GRPCConfig().Retry.Timeout, ShouldEqual, time.Second)
				So(cfg.HTTPConfig().Port, ShouldEqual, "8080")
				So(cfg.HTTPConfig().Retry.Attempts, ShouldEqual, 3)
				So(cfg.HTTPConfig().Retry.Timeout, ShouldEqual, time.Second)
				So(cfg.NSQConfig().Host, ShouldEqual, "localhost:4150")
				So(cfg.NSQConfig().LookupHost, ShouldEqual, "localhost:4161")
				So(cfg.NSQConfig().Retry.Attempts, ShouldEqual, 3)
				So(cfg.NSQConfig().Retry.Timeout, ShouldEqual, time.Second)
			})

			So(os.Unsetenv("CONFIG_FILE"), ShouldBeNil)
		})
	})
}

func TestMissingConfig(t *testing.T) {
	Convey("Given I have configuration file", t, func() {
		Convey("When I try to parse the configuration file", func() {
			err := config.UnmarshalFromFile(config.NewConfigurator())

			Convey("Then I should have an error of missing config file", func() {
				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, "read .: is a directory")
			})
		})
	})
}

func TestNonExistentConfig(t *testing.T) {
	Convey("Given I have configuration file", t, func() {
		os.Setenv("CONFIG_FILE", "../../test/bob")

		Convey("When I try to parse the configuration file", func() {
			err := config.UnmarshalFromFile(config.NewConfigurator())

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
		os.Setenv("CONFIG_FILE", "../test/greet/v1/service.proto")

		Convey("When I try to parse the configuration file", func() {
			err := config.UnmarshalFromFile(config.NewConfigurator())

			Convey("Then I should have an error of invalid configuration file", func() {
				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, "yaml: line 12: mapping values are not allowed in this context")
			})

			So(os.Unsetenv("CONFIG_FILE"), ShouldBeNil)
		})
	})
}

func TestWriteConfig(t *testing.T) {
	Convey("Given I have configuration file", t, func() {
		file := "../test/new_config.yml"
		os.Setenv("CONFIG_FILE", "../test/config.yml")
		os.Setenv("NEW_CONFIG_FILE", file)

		err := os.Remove(file)
		if !os.IsNotExist(err) {
			So(err, ShouldBeNil)
		}

		Convey("When I try to write the new configuration file", func() {
			cfg := config.NewConfigurator()
			err := config.UnmarshalFromFile(cfg)
			So(err, ShouldBeNil)

			bytes, err := yaml.Marshal(cfg)
			So(err, ShouldBeNil)

			err = config.WriteFileToEnv("NEW_CONFIG_FILE", bytes)
			So(err, ShouldBeNil)

			Convey("Then I should have a new configuration", func() {
				_, err := os.Open(file)
				So(err, ShouldBeNil)
			})

			So(os.Unsetenv("CONFIG_FILE"), ShouldBeNil)
			So(os.Unsetenv("NEW_CONFIG_FILE"), ShouldBeNil)
		})
	})
}

func TestStringMapConfig(t *testing.T) {
	Convey("Given I have configuration file", t, func() {
		file := "../test/config.yml"
		bytes, err := os.ReadFile(file)
		So(err, ShouldBeNil)

		Convey("When I try unmarshal from bytes", func() {
			cfg := config.Map{}
			err = config.UnmarshalFromBytes(bytes, cfg)
			So(err, ShouldBeNil)

			port := cfg.Map("transport").Map("http")["port"]
			So(port, ShouldEqual, 8080)

			Convey("Then I should have a valid configuration", func() {
				bytes, err := config.MarshalToBytes(cfg)
				So(err, ShouldBeNil)
				So(bytes, ShouldNotBeEmpty)
			})
		})
	})
}
