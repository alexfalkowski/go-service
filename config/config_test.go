package config_test

import (
	"testing"

	cache "github.com/alexfalkowski/go-service/v2/cache/config"
	"github.com/alexfalkowski/go-service/v2/config"
	"github.com/alexfalkowski/go-service/v2/crypto"
	"github.com/alexfalkowski/go-service/v2/crypto/aes"
	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/crypto/hmac"
	"github.com/alexfalkowski/go-service/v2/crypto/rsa"
	"github.com/alexfalkowski/go-service/v2/crypto/ssh"
	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	"github.com/alexfalkowski/go-service/v2/debug"
	"github.com/alexfalkowski/go-service/v2/encoding/base64"
	"github.com/alexfalkowski/go-service/v2/feature"
	"github.com/alexfalkowski/go-service/v2/flag"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/server"
	"github.com/alexfalkowski/go-service/v2/token"
	. "github.com/smartystreets/goconvey/convey"
)

func TestValidFileConfig(t *testing.T) {
	files := []string{
		test.FilePath("configs/config.json"),
		test.FilePath("configs/config.toml"),
		test.FilePath("configs/config.yml"),
	}

	for _, file := range files {
		Convey("Given I have configuration file", t, func() {
			set := flag.NewFlagSet("test")
			set.AddInput(file)

			decoder := test.NewDecoder(set)

			Convey("When I try to parse the configuration file", func() {
				config, err := config.NewConfig[config.Config](decoder, test.Validator)
				So(err, ShouldBeNil)

				Convey("Then I should have a valid configuration", func() {
					verifyConfig(config)
				})
			})
		})
	}
}

func TestInvalidFileConfig(t *testing.T) {
	files := []string{
		test.FilePath("configs/invalid.yml"),
		test.FilePath("configs/invalid_trace.yml"),
		test.FilePath("configs/missing.yml"),
		test.FilePath("configs/script.sh"),
		test.FilePath("config.go"),
		"",
		"env:BOB",
	}

	for _, file := range files {
		Convey("Given I have configuration file", t, func() {
			set := flag.NewFlagSet("test")
			set.AddInput(file)

			decoder := test.NewDecoder(set)

			Convey("When I try to parse the configuration file", func() {
				_, err := config.NewConfig[config.Config](decoder, test.Validator)

				Convey("Then I should have an error", func() {
					So(err, ShouldBeError)
				})
			})
		})
	}
}

func TestValidEnvConfig(t *testing.T) {
	Convey("Given I have configuration file", t, func() {
		d, err := test.FS.ReadFile(test.Path("configs/config.yml"))
		So(err, ShouldBeNil)

		t.Setenv("CONFIG", "yaml:"+base64.Encode(d))

		set := flag.NewFlagSet("test")
		set.AddInput("env:CONFIG")

		decoder := test.NewDecoder(set)

		Convey("When I try to parse the configuration file", func() {
			config, err := config.NewConfig[config.Config](decoder, test.Validator)
			So(err, ShouldBeNil)

			Convey("Then I should have a valid configuration", func() {
				verifyConfig(config)
			})
		})
	})
}

func TestInvalidEnvConfig(t *testing.T) {
	Convey("Given I have configuration file", t, func() {
		t.Setenv("CONFIG", "yaml:not_good")

		set := flag.NewFlagSet("test")
		set.AddInput("env:CONFIG")

		decoder := test.NewDecoder(set)

		Convey("When I try to parse the configuration file", func() {
			_, err := config.NewConfig[config.Config](decoder, test.Validator)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}

func TestValidCommonConfig(t *testing.T) {
	Convey("Given I have configuration file", t, func() {
		home := os.UserHomeDir()
		path := test.FS.Join(home, ".config", test.Name.String())

		err := test.FS.MkdirAll(path, 0o777)
		So(err, ShouldBeNil)

		data, err := test.FS.ReadFile(test.Path("configs/config.yml"))
		So(err, ShouldBeNil)

		err = test.FS.WriteFile(test.FS.Join(path, test.Name.String()+".yml"), data, 0o600)
		So(err, ShouldBeNil)

		set := flag.NewFlagSet("test")
		set.AddInput("")

		decoder := test.NewDecoder(set)

		Convey("When I try to parse the configuration file", func() {
			config, err := config.NewConfig[config.Config](decoder, test.Validator)
			So(err, ShouldBeNil)

			Convey("Then I should have a valid configuration", func() {
				verifyConfig(config)
			})
		})

		err = test.FS.RemoveAll(path)
		So(err, ShouldBeNil)
	})
}

func TestInvalidCommonConfig(t *testing.T) {
	Convey("Given I do not have a configuration file", t, func() {
		set := flag.NewFlagSet("test")
		set.AddInput("")

		decoder := test.NewDecoder(set)

		Convey("When I try to parse the configuration file", func() {
			_, err := config.NewConfig[config.Config](decoder, test.Validator)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}

func TestInvalidKindConfig(t *testing.T) {
	Convey("When I try to parse the configuration file", t, func() {
		set := flag.NewFlagSet("test")
		set.AddInput("test:test")

		decoder := test.NewDecoder(set)

		Convey("When I try to parse the configuration file", func() {
			_, err := config.NewConfig[config.Config](decoder, test.Validator)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}

//nolint:funlen
func verifyConfig(config *config.Config) {
	So(config.Environment.String(), ShouldEqual, "development")
	So(cache.IsEnabled(config.Cache), ShouldBeTrue)
	So(config.Cache.Kind, ShouldEqual, "redis")
	So(config.Cache.Compressor, ShouldEqual, "snappy")
	So(config.Cache.Encoder, ShouldEqual, "proto")
	So(config.Cache.Options["url"], ShouldEqual, "file:../test/secrets/redis")
	So(crypto.IsEnabled(config.Crypto), ShouldBeTrue)
	So(aes.IsEnabled(config.Crypto.AES), ShouldBeTrue)
	So(config.Crypto.AES.Key, ShouldNotBeBlank)
	So(ed25519.IsEnabled(config.Crypto.Ed25519), ShouldBeTrue)
	So(config.Crypto.Ed25519.Public, ShouldNotBeBlank)
	So(config.Crypto.Ed25519.Private, ShouldNotBeBlank)
	So(hmac.IsEnabled(config.Crypto.HMAC), ShouldBeTrue)
	So(config.Crypto.HMAC.Key, ShouldNotBeBlank)
	So(rsa.IsEnabled(config.Crypto.RSA), ShouldBeTrue)
	So(config.Crypto.RSA.Public, ShouldNotBeBlank)
	So(config.Crypto.RSA.Private, ShouldNotBeBlank)
	So(ssh.IsEnabled(config.Crypto.SSH), ShouldBeTrue)
	So(config.Crypto.SSH.Public, ShouldNotBeBlank)
	So(config.Crypto.SSH.Private, ShouldNotBeBlank)
	So(debug.IsEnabled(config.Debug), ShouldBeTrue)
	So(config.Debug.Address, ShouldEqual, ":6060")
	So(tls.IsEnabled(config.Debug.TLS), ShouldBeFalse)
	So(feature.IsEnabled(config.Feature), ShouldBeTrue)
	So(config.Feature.Address, ShouldEqual, "localhost:9000")
	So(config.ID.Kind, ShouldEqual, "uuid")
	So(config.Hooks.Secret, ShouldEqual, "file:../test/secrets/hooks")
	So(len(config.SQL.PG.Masters), ShouldEqual, 1)
	So(config.SQL.PG.Masters[0].URL, ShouldEqual, "file:../test/secrets/pg")
	So(len(config.SQL.PG.Slaves), ShouldEqual, 1)
	So(config.SQL.PG.Slaves[0].URL, ShouldEqual, "file:../test/secrets/pg")
	So(config.SQL.PG.MaxIdleConns, ShouldEqual, 5)
	So(config.SQL.PG.MaxOpenConns, ShouldEqual, 5)
	So(config.SQL.PG.ConnMaxLifetime, ShouldEqual, "1h")
	So(config.Limiter.Kind, ShouldEqual, "user-agent")
	So(config.Limiter.Tokens, ShouldEqual, 10)
	So(config.Limiter.Interval, ShouldEqual, "1s")
	So(config.Telemetry.Logger.Kind, ShouldEqual, "text")
	So(config.Telemetry.Logger.Level, ShouldEqual, "info")
	So(config.Telemetry.Metrics.Kind, ShouldEqual, "prometheus")
	So(config.Time.Kind, ShouldEqual, "nts")
	So(config.Time.Address, ShouldEqual, "time.cloudflare.com")
	So(config.Telemetry.Tracer.URL, ShouldEqual, "http://localhost:4318/v1/traces")
	So(config.Telemetry.Tracer.Kind, ShouldEqual, "otlp")
	So(token.IsEnabled(config.Token), ShouldBeTrue)
	So(config.Token.Kind, ShouldEqual, "jwt")
	So(config.Token.JWT.Expiration, ShouldEqual, "1h")
	So(config.Token.JWT.Issuer, ShouldEqual, "iss")
	So(config.Token.JWT.KeyID, ShouldEqual, "1234567890")
	So(server.IsEnabled(config.Transport.GRPC.Config), ShouldBeTrue)
	So(config.Transport.GRPC.Address, ShouldEqual, ":12000")
	So(config.Transport.GRPC.Retry.Attempts, ShouldEqual, 3)
	So(config.Transport.GRPC.Retry.Timeout, ShouldEqual, "1s")
	So(tls.IsEnabled(config.Transport.GRPC.TLS), ShouldBeFalse)
	So(server.IsEnabled(config.Transport.HTTP.Config), ShouldBeTrue)
	So(config.Transport.HTTP.Address, ShouldEqual, ":11000")
	So(config.Transport.HTTP.Retry.Attempts, ShouldEqual, 3)
	So(config.Transport.HTTP.Retry.Timeout, ShouldEqual, "1s")
	So(tls.IsEnabled(config.Transport.HTTP.TLS), ShouldBeFalse)
}
