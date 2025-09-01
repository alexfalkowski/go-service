package config_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/config"
	"github.com/alexfalkowski/go-service/v2/encoding/base64"
	"github.com/alexfalkowski/go-service/v2/flag"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/os"
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
	So(config.Debug.Address, ShouldEqual, "tcp://localhost:6060")
	So(config.Debug.TLS.IsEnabled(), ShouldBeFalse)
	So(config.Cache.IsEnabled(), ShouldBeTrue)
	So(config.Cache.Kind, ShouldEqual, "redis")
	So(config.Cache.Compressor, ShouldEqual, "snappy")
	So(config.Cache.Encoder, ShouldEqual, "proto")
	So(config.Cache.Options["url"], ShouldEqual, "file:../test/secrets/redis")
	So(config.Crypto.IsEnabled(), ShouldBeTrue)
	So(config.Crypto.AES.IsEnabled(), ShouldBeTrue)
	So(config.Crypto.AES.Key, ShouldNotBeBlank)
	So(config.Crypto.Ed25519.IsEnabled(), ShouldBeTrue)
	So(config.Crypto.Ed25519.Public, ShouldNotBeBlank)
	So(config.Crypto.Ed25519.Private, ShouldNotBeBlank)
	So(config.Crypto.HMAC.IsEnabled(), ShouldBeTrue)
	So(config.Crypto.HMAC.Key, ShouldNotBeBlank)
	So(config.Crypto.RSA.IsEnabled(), ShouldBeTrue)
	So(config.Crypto.RSA.Public, ShouldNotBeBlank)
	So(config.Crypto.RSA.Private, ShouldNotBeBlank)
	So(config.Crypto.SSH.IsEnabled(), ShouldBeTrue)
	So(config.Crypto.SSH.Public, ShouldNotBeBlank)
	So(config.Crypto.SSH.Private, ShouldNotBeBlank)
	So(config.Debug.IsEnabled(), ShouldBeTrue)
	So(config.Environment.String(), ShouldEqual, "development")
	So(config.Feature.IsEnabled(), ShouldBeTrue)
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
	So(config.Telemetry.Logger.Kind, ShouldEqual, "text")
	So(config.Telemetry.Logger.Level, ShouldEqual, "info")
	So(config.Telemetry.Metrics.Kind, ShouldEqual, "prometheus")
	So(config.Time.Kind, ShouldEqual, "nts")
	So(config.Time.Address, ShouldEqual, "time.cloudflare.com")
	So(config.Telemetry.Tracer.URL, ShouldEqual, "http://localhost:4318/v1/traces")
	So(config.Telemetry.Tracer.Kind, ShouldEqual, "otlp")
	So(config.Transport.GRPC.Token.IsEnabled(), ShouldBeTrue)
	So(config.Transport.GRPC.Token.Access.Policy, ShouldEqual, "../test/configs/rbac.csv")
	So(config.Transport.GRPC.Token.Kind, ShouldEqual, "jwt")
	So(config.Transport.GRPC.Token.JWT.Expiration, ShouldEqual, "1h")
	So(config.Transport.GRPC.Token.JWT.Issuer, ShouldEqual, "iss")
	So(config.Transport.GRPC.Token.JWT.KeyID, ShouldEqual, "1234567890")
	So(config.Transport.GRPC.Config.IsEnabled(), ShouldBeTrue)
	So(config.Transport.GRPC.Address, ShouldEqual, "tcp://localhost:12000")
	So(config.Transport.GRPC.Limiter.Kind, ShouldEqual, "user-agent")
	So(config.Transport.GRPC.Limiter.Tokens, ShouldEqual, 10)
	So(config.Transport.GRPC.Limiter.Interval, ShouldEqual, "1s")
	So(config.Transport.GRPC.Retry.Attempts, ShouldEqual, 3)
	So(config.Transport.GRPC.Retry.Timeout, ShouldEqual, "1s")
	So(config.Transport.GRPC.TLS.IsEnabled(), ShouldBeFalse)
	So(config.Transport.HTTP.Token.IsEnabled(), ShouldBeTrue)
	So(config.Transport.HTTP.Token.Access.Policy, ShouldEqual, "../test/configs/rbac.csv")
	So(config.Transport.HTTP.Token.Kind, ShouldEqual, "jwt")
	So(config.Transport.HTTP.Token.JWT.Expiration, ShouldEqual, "1h")
	So(config.Transport.HTTP.Token.JWT.Issuer, ShouldEqual, "iss")
	So(config.Transport.HTTP.Token.JWT.KeyID, ShouldEqual, "1234567890")
	So(config.Transport.HTTP.Config.IsEnabled(), ShouldBeTrue)
	So(config.Transport.HTTP.Address, ShouldEqual, "tcp://localhost:11000")
	So(config.Transport.HTTP.Limiter.Kind, ShouldEqual, "user-agent")
	So(config.Transport.HTTP.Limiter.Tokens, ShouldEqual, 10)
	So(config.Transport.HTTP.Limiter.Interval, ShouldEqual, "1s")
	So(config.Transport.HTTP.Retry.Attempts, ShouldEqual, 3)
	So(config.Transport.HTTP.Retry.Timeout, ShouldEqual, "1s")
	So(config.Transport.HTTP.TLS.IsEnabled(), ShouldBeFalse)
}
