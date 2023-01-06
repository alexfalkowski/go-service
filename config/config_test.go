package config_test

import (
	"encoding/base64"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/go-service/config"
	. "github.com/smartystreets/goconvey/convey"
)

//nolint:dupl
func TestValidEnvConfig(t *testing.T) {
	Convey("Given I have configuration file", t, func() {
		os.Setenv("CONFIG_FILE", "../test/config.yml")

		cmd.ConfigFlag = "env:CONFIG_FILE"

		c, err := cmd.NewConfig()
		So(err, ShouldBeNil)

		Convey("When I try to parse the configuration file", func() {
			cfg := config.NewConfigurator()
			err := config.Unmarshal(cfg, c)
			So(err, ShouldBeNil)

			Convey("Then I should have a valid configuration", func() {
				So(cfg.RedisConfig().Addresses, ShouldResemble, map[string]string{"server": "localhost:6379"})
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
				So(len(cfg.PGConfig().Masters), ShouldEqual, 1)
				So(cfg.PGConfig().Masters[0].URL, ShouldEqual, "postgres://test:test@localhost:5432/test?sslmode=disable")
				So(len(cfg.PGConfig().Slaves), ShouldEqual, 1)
				So(cfg.PGConfig().Slaves[0].URL, ShouldEqual, "postgres://test:test@localhost:5432/test?sslmode=disable")
				So(cfg.PGConfig().MaxIdleConns, ShouldEqual, 5)
				So(cfg.PGConfig().MaxOpenConns, ShouldEqual, 5)
				So(cfg.PGConfig().ConnMaxLifetime, ShouldEqual, time.Hour)
				So(cfg.OpentracingConfig().Kind, ShouldEqual, "jaeger")
				So(cfg.OpentracingConfig().Host, ShouldEqual, "localhost:6831")
				So(cfg.TransportConfig().Port, ShouldEqual, "8080")
				So(cfg.GRPCConfig().Retry.Attempts, ShouldEqual, 3)
				So(cfg.GRPCConfig().Retry.Timeout, ShouldEqual, time.Second)
				So(cfg.HTTPConfig().Retry.Attempts, ShouldEqual, 3)
				So(cfg.HTTPConfig().Retry.Timeout, ShouldEqual, time.Second)
				So(cfg.NSQConfig().Host, ShouldEqual, "localhost:4150")
				So(cfg.NSQConfig().LookupHost, ShouldEqual, "localhost:4161")
				So(cfg.NSQConfig().Retry.Attempts, ShouldEqual, 3)
				So(cfg.NSQConfig().Retry.Timeout, ShouldEqual, time.Second)
			})
		})

		So(os.Unsetenv("CONFIG_FILE"), ShouldBeNil)
	})
}

//nolint:dupl
func TestValidFileConfig(t *testing.T) {
	Convey("Given I have configuration file", t, func() {
		cmd.ConfigFlag = "file:../test/config.yml"

		c, err := cmd.NewConfig()
		So(err, ShouldBeNil)

		Convey("When I try to parse the configuration file", func() {
			cfg := config.NewConfigurator()
			err := config.Unmarshal(cfg, c)
			So(err, ShouldBeNil)

			Convey("Then I should have a valid configuration", func() {
				So(cfg.RedisConfig().Addresses, ShouldResemble, map[string]string{"server": "localhost:6379"})
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
				So(len(cfg.PGConfig().Masters), ShouldEqual, 1)
				So(cfg.PGConfig().Masters[0].URL, ShouldEqual, "postgres://test:test@localhost:5432/test?sslmode=disable")
				So(len(cfg.PGConfig().Slaves), ShouldEqual, 1)
				So(cfg.PGConfig().Slaves[0].URL, ShouldEqual, "postgres://test:test@localhost:5432/test?sslmode=disable")
				So(cfg.PGConfig().MaxIdleConns, ShouldEqual, 5)
				So(cfg.PGConfig().MaxOpenConns, ShouldEqual, 5)
				So(cfg.PGConfig().ConnMaxLifetime, ShouldEqual, time.Hour)
				So(cfg.OpentracingConfig().Kind, ShouldEqual, "jaeger")
				So(cfg.OpentracingConfig().Host, ShouldEqual, "localhost:6831")
				So(cfg.TransportConfig().Port, ShouldEqual, "8080")
				So(cfg.GRPCConfig().Retry.Attempts, ShouldEqual, 3)
				So(cfg.GRPCConfig().Retry.Timeout, ShouldEqual, time.Second)
				So(cfg.HTTPConfig().Retry.Attempts, ShouldEqual, 3)
				So(cfg.HTTPConfig().Retry.Timeout, ShouldEqual, time.Second)
				So(cfg.NSQConfig().Host, ShouldEqual, "localhost:4150")
				So(cfg.NSQConfig().LookupHost, ShouldEqual, "localhost:4161")
				So(cfg.NSQConfig().Retry.Attempts, ShouldEqual, 3)
				So(cfg.NSQConfig().Retry.Timeout, ShouldEqual, time.Second)
			})
		})
	})
}

//nolint:dupl
func TestValidMemConfig(t *testing.T) {
	d, _ := os.ReadFile("../test/config.yml")

	Convey("Given I have configuration file", t, func() {
		cmd.ConfigFlag = fmt.Sprintf("mem:%s", base64.StdEncoding.EncodeToString(d))

		c, err := cmd.NewConfig()
		So(err, ShouldBeNil)

		Convey("When I try to parse the configuration file", func() {
			cfg := config.NewConfigurator()
			err := config.Unmarshal(cfg, c)
			So(err, ShouldBeNil)

			Convey("Then I should have a valid configuration", func() {
				So(cfg.RedisConfig().Addresses, ShouldResemble, map[string]string{"server": "localhost:6379"})
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
				So(len(cfg.PGConfig().Masters), ShouldEqual, 1)
				So(cfg.PGConfig().Masters[0].URL, ShouldEqual, "postgres://test:test@localhost:5432/test?sslmode=disable")
				So(len(cfg.PGConfig().Slaves), ShouldEqual, 1)
				So(cfg.PGConfig().Slaves[0].URL, ShouldEqual, "postgres://test:test@localhost:5432/test?sslmode=disable")
				So(cfg.PGConfig().MaxIdleConns, ShouldEqual, 5)
				So(cfg.PGConfig().MaxOpenConns, ShouldEqual, 5)
				So(cfg.PGConfig().ConnMaxLifetime, ShouldEqual, time.Hour)
				So(cfg.OpentracingConfig().Kind, ShouldEqual, "jaeger")
				So(cfg.OpentracingConfig().Host, ShouldEqual, "localhost:6831")
				So(cfg.TransportConfig().Port, ShouldEqual, "8080")
				So(cfg.GRPCConfig().Retry.Attempts, ShouldEqual, 3)
				So(cfg.GRPCConfig().Retry.Timeout, ShouldEqual, time.Second)
				So(cfg.HTTPConfig().Retry.Attempts, ShouldEqual, 3)
				So(cfg.HTTPConfig().Retry.Timeout, ShouldEqual, time.Second)
				So(cfg.NSQConfig().Host, ShouldEqual, "localhost:4150")
				So(cfg.NSQConfig().LookupHost, ShouldEqual, "localhost:4161")
				So(cfg.NSQConfig().Retry.Attempts, ShouldEqual, 3)
				So(cfg.NSQConfig().Retry.Timeout, ShouldEqual, time.Second)
			})
		})
	})
}

func TestInvalidConfig(t *testing.T) {
	Convey("Given I have configuration file", t, func() {
		os.Setenv("CONFIG_FILE", "../test/greet/v1/service.proto")
		cmd.ConfigFlag = ""

		c, err := cmd.NewConfig()
		So(err, ShouldBeNil)

		Convey("When I try to parse the configuration file", func() {
			err := config.Unmarshal(config.NewConfigurator(), c)

			Convey("Then I should have an error of invalid configuration file", func() {
				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, "yaml: line 12: mapping values are not allowed in this context")
			})

			So(os.Unsetenv("CONFIG_FILE"), ShouldBeNil)
		})
	})
}
