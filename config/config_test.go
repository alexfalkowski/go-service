package config_test

import (
	"encoding/base64"
	"errors"
	"os"
	"testing"

	"github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/go-service/config"
	"github.com/alexfalkowski/go-service/crypto"
	"github.com/alexfalkowski/go-service/crypto/aes"
	"github.com/alexfalkowski/go-service/crypto/ed25519"
	"github.com/alexfalkowski/go-service/crypto/hmac"
	"github.com/alexfalkowski/go-service/crypto/tls"
	"github.com/alexfalkowski/go-service/debug"
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
				cfg, err := config.NewConfig(&cmd.InputConfig{Config: c})
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
			cfg, err := config.NewConfig(&cmd.InputConfig{Config: c})
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
			_, err := config.NewConfig(&cmd.InputConfig{Config: c})

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
			cfg, err := config.NewConfig(&cmd.InputConfig{Config: c})
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
func verifyConfig(cfg *config.Config) {
	So(string(cfg.Environment), ShouldEqual, "development")
	So(crypto.IsEnabled(cfg.Crypto), ShouldBeTrue)
	So(aes.IsEnabled(cfg.Crypto.AES), ShouldBeTrue)
	So(cfg.Crypto.AES.Key, ShouldNotBeBlank)
	So(ed25519.IsEnabled(cfg.Crypto.Ed25519), ShouldBeTrue)
	So(cfg.Crypto.Ed25519.Public, ShouldNotBeBlank)
	So(cfg.Crypto.Ed25519.GetPrivate(), ShouldNotBeBlank)
	So(hmac.IsEnabled(cfg.Crypto.HMAC), ShouldBeTrue)
	So(cfg.Crypto.HMAC.Key, ShouldNotBeBlank)
	So(cfg.Crypto.RSA.Public, ShouldNotBeBlank)
	So(cfg.Crypto.RSA.GetPrivate(), ShouldNotBeBlank)
	So(cfg.Debug.Port, ShouldEqual, "6060")
	So(debug.IsEnabled(cfg.Debug), ShouldBeTrue)
	So(cfg.Debug.Port, ShouldEqual, "6060")
	So(tls.IsEnabled(cfg.Debug.TLS), ShouldBeFalse)
	So(cfg.Feature.Kind, ShouldEqual, "flipt")
	So(cfg.Feature.Host, ShouldEqual, "localhost:9000")
	So(cfg.Hooks.Secret, ShouldEqual, "YWJjZGUxMjM0NQ==")
	So(cfg.Cache.Redis.Compressor, ShouldEqual, "snappy")
	So(cfg.Cache.Redis.Marshaller, ShouldEqual, "proto")
	So(cfg.Cache.Redis.Addresses, ShouldResemble, map[string]string{"server": "localhost:6379"})
	So(cfg.Cache.Ristretto.BufferItems, ShouldEqual, 64)
	So(cfg.Cache.Ristretto.MaxCost, ShouldEqual, 100000000)
	So(cfg.Cache.Ristretto.NumCounters, ShouldEqual, 10000000)
	So(len(cfg.SQL.PG.Masters), ShouldEqual, 1)
	So(cfg.SQL.PG.Masters[0].URL, ShouldEqual, "postgres://test:test@localhost:5432/test?sslmode=disable")
	So(len(cfg.SQL.PG.Slaves), ShouldEqual, 1)
	So(cfg.SQL.PG.Slaves[0].URL, ShouldEqual, "postgres://test:test@localhost:5432/test?sslmode=disable")
	So(cfg.SQL.PG.MaxIdleConns, ShouldEqual, 5)
	So(cfg.SQL.PG.MaxOpenConns, ShouldEqual, 5)
	So(cfg.SQL.PG.ConnMaxLifetime, ShouldEqual, "1h")
	So(cfg.Token.Kind, ShouldEqual, "none")
	So(cfg.Limiter.Kind, ShouldEqual, "user-agent")
	So(cfg.Limiter.Pattern, ShouldEqual, "10-S")
	So(cfg.Telemetry.Logger.Level, ShouldEqual, "info")
	So(cfg.Telemetry.Metrics.Kind, ShouldEqual, "prometheus")
	So(cfg.Time.Kind, ShouldEqual, "nts")
	So(cfg.Time.Host, ShouldEqual, "time.cloudflare.com")
	So(cfg.Telemetry.Tracer.Host, ShouldEqual, "http://localhost:4318/v1/traces")
	So(cfg.Telemetry.Tracer.Kind, ShouldEqual, "otlp")
	So(server.IsEnabled(cfg.Transport.GRPC.Config), ShouldBeTrue)
	So(cfg.Transport.GRPC.Port, ShouldEqual, "12000")
	So(cfg.Transport.GRPC.Retry.Attempts, ShouldEqual, 3)
	So(cfg.Transport.GRPC.Retry.Timeout, ShouldEqual, "1s")
	So(cfg.Transport.GRPC.UserAgent, ShouldEqual, "Service grpc/1.0")
	So(tls.IsEnabled(cfg.Transport.GRPC.TLS), ShouldBeFalse)
	So(server.IsEnabled(cfg.Transport.HTTP.Config), ShouldBeTrue)
	So(cfg.Transport.HTTP.Port, ShouldEqual, "11000")
	So(cfg.Transport.HTTP.Retry.Attempts, ShouldEqual, 3)
	So(cfg.Transport.HTTP.Retry.Timeout, ShouldEqual, "1s")
	So(cfg.Transport.HTTP.UserAgent, ShouldEqual, "Service http/1.0")
	So(tls.IsEnabled(cfg.Transport.HTTP.TLS), ShouldBeFalse)
}
