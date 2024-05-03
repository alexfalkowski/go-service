package config_test

import (
	"encoding/base64"
	"errors"
	"os"
	"testing"

	"github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/go-service/config"
	"github.com/alexfalkowski/go-service/security"
	"github.com/alexfalkowski/go-service/server"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestValidEnvConfig(t *testing.T) {
	configs := []string{
		"../test/configs/config.json",
		"../test/configs/config.toml",
		"../test/configs/config.yml",
	}

	for _, f := range configs {
		Convey("Given I have configuration file", t, func() {
			So(os.Setenv("CONFIG_FILE", f), ShouldBeNil)

			c, err := test.NewCmdConfig("env:CONFIG_FILE")
			So(err, ShouldBeNil)

			Convey("When I try to parse the configuration file", func() {
				cfg, err := config.NewConfigurator(&cmd.InputConfig{Config: c})
				So(err, ShouldBeNil)

				Convey("Then I should have a valid configuration", func() {
					verifyConfig(cfg)
				})
			})

			So(os.Unsetenv("CONFIG_FILE"), ShouldBeNil)
		})
	}
}

func TestValidFileConfig(t *testing.T) {
	Convey("Given I have configuration file", t, func() {
		c, err := test.NewCmdConfig("file:../test/configs/config.yml")
		So(err, ShouldBeNil)

		Convey("When I try to parse the configuration file", func() {
			cfg, err := config.NewConfigurator(&cmd.InputConfig{Config: c})
			So(err, ShouldBeNil)

			Convey("Then I should have a valid configuration", func() {
				verifyConfig(cfg)
			})
		})
	})
}

func TestMissingFileConfig(t *testing.T) {
	Convey("Given I have missing configuration file", t, func() {
		c, err := test.NewCmdConfig("file:../test/configs/missing.yml")
		So(err, ShouldBeNil)

		Convey("When I try to parse the configuration file", func() {
			_, err := config.NewConfigurator(&cmd.InputConfig{Config: c})

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
				So(errors.Is(err, os.ErrNotExist), ShouldBeTrue)
			})
		})
	})
}

func TestValidMemConfig(t *testing.T) {
	Convey("Given I have configuration file", t, func() {
		d, err := os.ReadFile("../test/configs/config.yml")
		So(err, ShouldBeNil)

		So(os.Setenv("CONFIG_FILE", "yaml:CONFIG"), ShouldBeNil)
		So(os.Setenv("CONFIG", base64.StdEncoding.EncodeToString(d)), ShouldBeNil)

		c, err := test.NewCmdConfig("env:CONFIG_FILE")
		So(err, ShouldBeNil)

		Convey("When I try to parse the configuration file", func() {
			cfg, err := config.NewConfigurator(&cmd.InputConfig{Config: c})
			So(err, ShouldBeNil)

			Convey("Then I should have a valid configuration", func() {
				verifyConfig(cfg)
			})
		})

		So(os.Unsetenv("CONFIG_FILE"), ShouldBeNil)
		So(os.Unsetenv("CONFIG"), ShouldBeNil)
	})
}

//nolint:funlen
func verifyConfig(cfg config.Configurator) {
	So(string(cfg.EnvironmentConfig()), ShouldEqual, "development")
	So(server.IsEnabled(cfg.DebugConfig().Config), ShouldBeTrue)
	So(cfg.DebugConfig().Port, ShouldEqual, "6060")
	So(security.IsEnabled(cfg.DebugConfig().Config.Security), ShouldBeFalse)
	So(cfg.FeatureConfig().Kind, ShouldEqual, "flipt")
	So(cfg.FeatureConfig().Host, ShouldEqual, "localhost:9000")
	So(cfg.HooksConfig().Secret, ShouldEqual, "YWJjZGUxMjM0NQ==")
	So(cfg.RedisConfig().Compressor, ShouldEqual, "snappy")
	So(cfg.RedisConfig().Marshaller, ShouldEqual, "proto")
	So(cfg.RedisConfig().Addresses, ShouldResemble, map[string]string{"server": "localhost:6379"})
	So(cfg.RistrettoConfig().BufferItems, ShouldEqual, 64)
	So(cfg.RistrettoConfig().MaxCost, ShouldEqual, 100000000)
	So(cfg.RistrettoConfig().NumCounters, ShouldEqual, 10000000)
	So(len(cfg.PGConfig().Masters), ShouldEqual, 1)
	So(cfg.PGConfig().Masters[0].URL, ShouldEqual, "postgres://test:test@localhost:5432/test?sslmode=disable")
	So(len(cfg.PGConfig().Slaves), ShouldEqual, 1)
	So(cfg.PGConfig().Slaves[0].URL, ShouldEqual, "postgres://test:test@localhost:5432/test?sslmode=disable")
	So(cfg.PGConfig().MaxIdleConns, ShouldEqual, 5)
	So(cfg.PGConfig().MaxOpenConns, ShouldEqual, 5)
	So(cfg.PGConfig().ConnMaxLifetime, ShouldEqual, "1h")
	So(cfg.TokenConfig().Kind, ShouldEqual, "none")
	So(cfg.LimiterConfig().Enabled, ShouldBeTrue)
	So(cfg.LimiterConfig().Kind, ShouldEqual, "user-agent")
	So(cfg.LimiterConfig().Pattern, ShouldEqual, "10-S")
	So(cfg.LoggerConfig().Level, ShouldEqual, "info")
	So(cfg.MetricsConfig().Kind, ShouldEqual, "prometheus")
	So(cfg.NTPConfig().Host, ShouldEqual, "0.beevik-ntp.pool.ntp.org")
	So(cfg.NTSConfig().Host, ShouldEqual, "time.cloudflare.com")
	So(cfg.TracerConfig().Host, ShouldEqual, "http://localhost:4318/v1/traces")
	So(cfg.TracerConfig().Kind, ShouldEqual, "otlp")
	So(server.IsEnabled(cfg.GRPCConfig().Config), ShouldBeTrue)
	So(cfg.GRPCConfig().Port, ShouldEqual, "12000")
	So(cfg.GRPCConfig().Retry.Attempts, ShouldEqual, 3)
	So(cfg.GRPCConfig().Retry.Timeout, ShouldEqual, "1s")
	So(cfg.GRPCConfig().UserAgent, ShouldEqual, "Service grpc/1.0")
	So(security.IsEnabled(cfg.GRPCConfig().Config.Security), ShouldBeFalse)
	So(server.IsEnabled(cfg.HTTPConfig().Config), ShouldBeTrue)
	So(cfg.HTTPConfig().Port, ShouldEqual, "11000")
	So(cfg.HTTPConfig().Retry.Attempts, ShouldEqual, 3)
	So(cfg.HTTPConfig().Retry.Timeout, ShouldEqual, "1s")
	So(cfg.HTTPConfig().UserAgent, ShouldEqual, "Service http/1.0")
	So(security.IsEnabled(cfg.HTTPConfig().Config.Security), ShouldBeFalse)
}
