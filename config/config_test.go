package config_test

import (
	"encoding/base64"
	"errors"
	"os"
	"testing"

	"github.com/alexfalkowski/go-service/config"
	"github.com/alexfalkowski/go-service/crypto"
	"github.com/alexfalkowski/go-service/crypto/aes"
	"github.com/alexfalkowski/go-service/crypto/ed25519"
	"github.com/alexfalkowski/go-service/crypto/hmac"
	"github.com/alexfalkowski/go-service/crypto/rsa"
	"github.com/alexfalkowski/go-service/crypto/ssh"
	"github.com/alexfalkowski/go-service/crypto/tls"
	"github.com/alexfalkowski/go-service/debug"
	"github.com/alexfalkowski/go-service/feature"
	"github.com/alexfalkowski/go-service/server"
	"github.com/alexfalkowski/go-service/test"
	"github.com/alexfalkowski/go-service/token"
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

			c := test.NewInputConfig("env:CONFIG_FILE")

			Convey("When I try to parse the configuration file", func() {
				cfg, err := config.NewConfig[config.Config](c)
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
		c := test.NewInputConfig("file:../test/configs/config.yml")

		Convey("When I try to parse the configuration file", func() {
			cfg, err := config.NewConfig[config.Config](c)
			So(err, ShouldBeNil)

			Convey("Then I should have a valid configuration", func() {
				verifyConfig(cfg)
			})
		})
	})
}

func TestMissingFileConfig(t *testing.T) {
	Convey("Given I have missing configuration file", t, func() {
		c := test.NewInputConfig("file:../test/configs/missing.yml")

		Convey("When I try to parse the configuration file", func() {
			_, err := config.NewConfig[*config.Config](c)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
				So(errors.Is(err, os.ErrNotExist), ShouldBeTrue)
			})
		})
	})
}

func TestInvalidFileConfig(t *testing.T) {
	Convey("Given I have configuration file", t, func() {
		c := test.NewInputConfig("file:../test/configs/invalid.yml")

		Convey("When I try to parse the configuration file", func() {
			_, err := config.NewConfig[config.Config](c)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
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

		c := test.NewInputConfig("env:CONFIG_FILE")

		Convey("When I try to parse the configuration file", func() {
			cfg, err := config.NewConfig[config.Config](c)
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
	So(cfg.Crypto.Ed25519.Private, ShouldNotBeBlank)
	So(hmac.IsEnabled(cfg.Crypto.HMAC), ShouldBeTrue)
	So(cfg.Crypto.HMAC.Key, ShouldNotBeBlank)
	So(rsa.IsEnabled(cfg.Crypto.RSA), ShouldBeTrue)
	So(cfg.Crypto.RSA.Public, ShouldNotBeBlank)
	So(cfg.Crypto.RSA.Private, ShouldNotBeBlank)
	So(ssh.IsEnabled(cfg.Crypto.SSH), ShouldBeTrue)
	So(cfg.Crypto.SSH.Public, ShouldNotBeBlank)
	So(cfg.Crypto.SSH.Private, ShouldNotBeBlank)
	So(debug.IsEnabled(cfg.Debug), ShouldBeTrue)
	So(cfg.Debug.Address, ShouldEqual, ":6060")
	So(tls.IsEnabled(cfg.Debug.TLS), ShouldBeFalse)
	So(feature.IsEnabled(cfg.Feature), ShouldBeTrue)
	So(cfg.Feature.Address, ShouldEqual, "localhost:9000")
	So(cfg.Hooks.Secret, ShouldEqual, "../test/secrets/hooks")
	So(cfg.Cache.Redis.Compressor, ShouldEqual, "snappy")
	So(cfg.Cache.Redis.Encoder, ShouldEqual, "proto")
	So(len(cfg.SQL.PG.Masters), ShouldEqual, 1)
	So(cfg.SQL.PG.Masters[0].URL, ShouldEqual, "../test/secrets/pg")
	So(len(cfg.SQL.PG.Slaves), ShouldEqual, 1)
	So(cfg.SQL.PG.Slaves[0].URL, ShouldEqual, "../test/secrets/pg")
	So(cfg.SQL.PG.MaxIdleConns, ShouldEqual, 5)
	So(cfg.SQL.PG.MaxOpenConns, ShouldEqual, 5)
	So(cfg.SQL.PG.ConnMaxLifetime, ShouldEqual, "1h")
	So(cfg.Limiter.Kind, ShouldEqual, "user-agent")
	So(cfg.Limiter.Tokens, ShouldEqual, 10)
	So(cfg.Limiter.Interval, ShouldEqual, "1s")
	So(cfg.Telemetry.Logger.Level, ShouldEqual, "info")
	So(cfg.Telemetry.Metrics.Kind, ShouldEqual, "prometheus")
	So(cfg.Time.Kind, ShouldEqual, "nts")
	So(cfg.Time.Address, ShouldEqual, "time.cloudflare.com")
	So(cfg.Telemetry.Tracer.URL, ShouldEqual, "http://localhost:4318/v1/traces")
	So(cfg.Telemetry.Tracer.Kind, ShouldEqual, "otlp")
	So(token.IsEnabled(cfg.Token), ShouldBeTrue)
	So(cfg.Token.Audience, ShouldEqual, "aud")
	So(cfg.Token.Expiration, ShouldEqual, "1h")
	So(cfg.Token.Issuer, ShouldEqual, "iss")
	So(cfg.Token.Kind, ShouldEqual, "jwt")
	So(cfg.Token.Subject, ShouldEqual, "sub")
	So(server.IsEnabled(cfg.Transport.GRPC.Config), ShouldBeTrue)
	So(cfg.Transport.GRPC.Address, ShouldEqual, ":12000")
	So(cfg.Transport.GRPC.Retry.Attempts, ShouldEqual, 3)
	So(cfg.Transport.GRPC.Retry.Timeout, ShouldEqual, "1s")
	So(tls.IsEnabled(cfg.Transport.GRPC.TLS), ShouldBeFalse)
	So(server.IsEnabled(cfg.Transport.HTTP.Config), ShouldBeTrue)
	So(cfg.Transport.HTTP.Address, ShouldEqual, ":11000")
	So(cfg.Transport.HTTP.Retry.Attempts, ShouldEqual, 3)
	So(cfg.Transport.HTTP.Retry.Timeout, ShouldEqual, "1s")
	So(tls.IsEnabled(cfg.Transport.HTTP.TLS), ShouldBeFalse)
}
