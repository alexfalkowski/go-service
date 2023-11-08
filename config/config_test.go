package config_test

import (
	"encoding/base64"
	"os"
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/go-service/config"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestValidEnvConfig(t *testing.T) {
	for _, f := range []string{"../test/config.yml", "../test/config.toml"} {
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
		c, err := test.NewCmdConfig("file:../test/config.yml")
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

func TestValidMemConfig(t *testing.T) {
	Convey("Given I have configuration file", t, func() {
		d, err := os.ReadFile("../test/config.yml")
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

func verifyConfig(cfg config.Configurator) {
	So(string(cfg.EnvironmentConfig()), ShouldEqual, "development")
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
	So(cfg.PGConfig().ConnMaxLifetime, ShouldEqual, time.Hour)
	So(cfg.TracerConfig().Host, ShouldEqual, "localhost:4318")
	So(cfg.GRPCConfig().Enabled, ShouldEqual, true)
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
}
