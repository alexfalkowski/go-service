package auth0_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"os"
	"testing"

	"github.com/alexfalkowski/go-service/otel"
	"github.com/alexfalkowski/go-service/security/auth0"
	"github.com/alexfalkowski/go-service/test"
	"github.com/golang-jwt/jwt/v4"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

const (
	algorithm = "ES256"
)

func init() {
	otel.Register()
}

func TestInvalidJSONWebKeySet(t *testing.T) {
	Convey("Given I have a corrupt token", t, func() {
		lc := fxtest.NewLifecycle(t)
		acfg := &auth0.Config{
			URL:           os.Getenv("AUTH0_URL"),
			ClientID:      os.Getenv("AUTH0_CLIENT_ID"),
			ClientSecret:  os.Getenv("AUTH0_CLIENT_SECRET"),
			Audience:      os.Getenv("AUTH0_AUDIENCE"),
			Issuer:        os.Getenv("AUTH0_ISSUER"),
			Algorithm:     os.Getenv("AUTH0_ALGORITHM"),
			JSONWebKeySet: "not a valid URL",
		}
		logger := test.NewLogger(lc)
		cache := test.NewRistrettoCache(lc)

		tracer, err := otel.NewTracer(otel.TracerParams{Lifecycle: lc, Config: test.NewOTELConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		params := auth0.CertificatorParams{Config: acfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		cert := auth0.NewCertificator(params)
		ver := auth0.NewVerifier(acfg, cert)

		lc.RequireStart()

		Convey("When I try to verify the token", func() {
			claims := jwt.MapClaims{"aud": acfg.Audience, "iss": acfg.Issuer, "kid": "none"}
			token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

			key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
			So(err, ShouldBeNil)

			tkn, err := token.SignedString(key)
			So(err, ShouldBeNil)

			ctx := context.Background()
			_, _, err = ver.Verify(ctx, []byte(tkn))

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})

		lc.RequireStop()
	})
}

//nolint:dupl
func TestInvalidResponseJSONWebKeySet(t *testing.T) {
	Convey("Given I have a corrupt token", t, func() {
		lc := fxtest.NewLifecycle(t)
		acfg := &auth0.Config{
			URL:           os.Getenv("AUTH0_URL"),
			ClientID:      os.Getenv("AUTH0_CLIENT_ID"),
			ClientSecret:  os.Getenv("AUTH0_CLIENT_SECRET"),
			Audience:      os.Getenv("AUTH0_AUDIENCE"),
			Issuer:        os.Getenv("AUTH0_ISSUER"),
			Algorithm:     os.Getenv("AUTH0_ALGORITHM"),
			JSONWebKeySet: "https://httpstat.us/400",
		}
		logger := test.NewLogger(lc)
		cache := test.NewRistrettoCache(lc)

		tracer, err := otel.NewTracer(otel.TracerParams{Lifecycle: lc, Config: test.NewOTELConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		params := auth0.CertificatorParams{Config: acfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		cert := auth0.NewCertificator(params)
		ver := auth0.NewVerifier(acfg, cert)

		lc.RequireStart()

		Convey("When I try to verify the token", func() {
			claims := jwt.MapClaims{"aud": acfg.Audience, "iss": acfg.Issuer, "kid": "none"}
			token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

			key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
			So(err, ShouldBeNil)

			tkn, err := token.SignedString(key)
			So(err, ShouldBeNil)

			ctx := context.Background()
			_, _, err = ver.Verify(ctx, []byte(tkn))

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, "invalid response")
			})
		})

		lc.RequireStop()
	})
}

//nolint:dupl
func TestInvalidJSONResponseJSONWebKeySet(t *testing.T) {
	Convey("Given I have a corrupt token", t, func() {
		lc := fxtest.NewLifecycle(t)
		acfg := &auth0.Config{
			URL:           os.Getenv("AUTH0_URL"),
			ClientID:      os.Getenv("AUTH0_CLIENT_ID"),
			ClientSecret:  os.Getenv("AUTH0_CLIENT_SECRET"),
			Audience:      os.Getenv("AUTH0_AUDIENCE"),
			Issuer:        os.Getenv("AUTH0_ISSUER"),
			Algorithm:     os.Getenv("AUTH0_ALGORITHM"),
			JSONWebKeySet: "https://httpstat.us/200",
		}
		logger := test.NewLogger(lc)
		cache := test.NewRistrettoCache(lc)

		tracer, err := otel.NewTracer(otel.TracerParams{Lifecycle: lc, Config: test.NewOTELConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		params := auth0.CertificatorParams{Config: acfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		cert := auth0.NewCertificator(params)
		ver := auth0.NewVerifier(acfg, cert)

		lc.RequireStart()

		Convey("When I try to verify the token", func() {
			claims := jwt.MapClaims{"aud": acfg.Audience, "iss": acfg.Issuer, "kid": "none"}
			token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

			key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
			So(err, ShouldBeNil)

			tkn, err := token.SignedString(key)
			So(err, ShouldBeNil)

			ctx := context.Background()
			_, _, err = ver.Verify(ctx, []byte(tkn))

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, "json: cannot unmarshal number into Go value of type auth0.jwksResponse")
			})
		})

		lc.RequireStop()
	})
}

func TestCorruptToken(t *testing.T) {
	Convey("Given I have a corrupt token", t, func() {
		lc := fxtest.NewLifecycle(t)
		acfg := &auth0.Config{
			URL:           os.Getenv("AUTH0_URL"),
			ClientID:      os.Getenv("AUTH0_CLIENT_ID"),
			ClientSecret:  os.Getenv("AUTH0_CLIENT_SECRET"),
			Audience:      os.Getenv("AUTH0_AUDIENCE"),
			Issuer:        os.Getenv("AUTH0_ISSUER"),
			Algorithm:     os.Getenv("AUTH0_ALGORITHM"),
			JSONWebKeySet: os.Getenv("AUTH0_JSON_WEB_KEY_SET"),
		}
		logger := test.NewLogger(lc)
		cache := test.NewRistrettoCache(lc)

		tracer, err := otel.NewTracer(otel.TracerParams{Lifecycle: lc, Config: test.NewOTELConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		params := auth0.CertificatorParams{Config: acfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		cert := auth0.NewCertificator(params)
		ver := auth0.NewVerifier(acfg, cert)

		lc.RequireStart()

		Convey("When I try to verify the token", func() {
			ctx := context.Background()
			_, _, err = ver.Verify(ctx, []byte("corrupt-token"))

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})

		lc.RequireStop()
	})
}

func TestMissingAudienceToken(t *testing.T) {
	Convey("Given I have a missing audience in token", t, func() {
		lc := fxtest.NewLifecycle(t)
		acfg := &auth0.Config{
			URL:           os.Getenv("AUTH0_URL"),
			ClientID:      os.Getenv("AUTH0_CLIENT_ID"),
			ClientSecret:  os.Getenv("AUTH0_CLIENT_SECRET"),
			Audience:      os.Getenv("AUTH0_AUDIENCE"),
			Issuer:        os.Getenv("AUTH0_ISSUER"),
			Algorithm:     algorithm,
			JSONWebKeySet: os.Getenv("AUTH0_JSON_WEB_KEY_SET"),
		}
		logger := test.NewLogger(lc)
		cache := test.NewRistrettoCache(lc)

		tracer, err := otel.NewTracer(otel.TracerParams{Lifecycle: lc, Config: test.NewOTELConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		params := auth0.CertificatorParams{Config: acfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		cert := auth0.NewCertificator(params)
		ver := auth0.NewVerifier(acfg, cert)

		lc.RequireStart()

		Convey("When I try to verify the token", func() {
			token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{})

			key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
			So(err, ShouldBeNil)

			tkn, err := token.SignedString(key)
			So(err, ShouldBeNil)

			ctx := context.Background()
			_, _, err = ver.Verify(ctx, []byte(tkn))

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})

		lc.RequireStop()
	})
}

func TestMissingIssuerToken(t *testing.T) {
	Convey("Given I have a missing issuer in token", t, func() {
		lc := fxtest.NewLifecycle(t)
		acfg := &auth0.Config{
			URL:           os.Getenv("AUTH0_URL"),
			ClientID:      os.Getenv("AUTH0_CLIENT_ID"),
			ClientSecret:  os.Getenv("AUTH0_CLIENT_SECRET"),
			Audience:      os.Getenv("AUTH0_AUDIENCE"),
			Issuer:        os.Getenv("AUTH0_ISSUER"),
			Algorithm:     algorithm,
			JSONWebKeySet: os.Getenv("AUTH0_JSON_WEB_KEY_SET"),
		}
		logger := test.NewLogger(lc)
		cache := test.NewRistrettoCache(lc)

		tracer, err := otel.NewTracer(otel.TracerParams{Lifecycle: lc, Config: test.NewOTELConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		params := auth0.CertificatorParams{Config: acfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		cert := auth0.NewCertificator(params)
		ver := auth0.NewVerifier(acfg, cert)

		lc.RequireStart()

		Convey("When I try to verify the token", func() {
			claims := jwt.MapClaims{"aud": acfg.Audience}
			token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

			key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
			So(err, ShouldBeNil)

			tkn, err := token.SignedString(key)
			So(err, ShouldBeNil)

			ctx := context.Background()
			_, _, err = ver.Verify(ctx, []byte(tkn))

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})

		lc.RequireStop()
	})
}

func TestInvalidCertificateToken(t *testing.T) {
	Convey("Given I have an invalid jwks endpoint", t, func() {
		lc := fxtest.NewLifecycle(t)
		acfg := &auth0.Config{
			URL:           os.Getenv("AUTH0_URL"),
			ClientID:      os.Getenv("AUTH0_CLIENT_ID"),
			ClientSecret:  os.Getenv("AUTH0_CLIENT_SECRET"),
			Audience:      os.Getenv("AUTH0_AUDIENCE"),
			Issuer:        os.Getenv("AUTH0_ISSUER"),
			Algorithm:     algorithm,
			JSONWebKeySet: "https://non-existent.com/.well-known/jwks.json",
		}
		logger := test.NewLogger(lc)
		cache := test.NewRistrettoCache(lc)

		tracer, err := otel.NewTracer(otel.TracerParams{Lifecycle: lc, Config: test.NewOTELConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		params := auth0.CertificatorParams{Config: acfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		cert := auth0.NewCertificator(params)
		ver := auth0.NewVerifier(acfg, cert)

		lc.RequireStart()

		Convey("When I try to verify the token", func() {
			claims := jwt.MapClaims{"aud": acfg.Audience, "iss": acfg.Issuer}
			token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

			key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
			So(err, ShouldBeNil)

			tkn, err := token.SignedString(key)
			So(err, ShouldBeNil)

			ctx := context.Background()
			_, _, err = ver.Verify(ctx, []byte(tkn))

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})

		lc.RequireStop()
	})
}

func TestMissingKidToken(t *testing.T) {
	Convey("Given I have an invalid jwks endpoint", t, func() {
		lc := fxtest.NewLifecycle(t)
		acfg := &auth0.Config{
			URL:           os.Getenv("AUTH0_URL"),
			ClientID:      os.Getenv("AUTH0_CLIENT_ID"),
			ClientSecret:  os.Getenv("AUTH0_CLIENT_SECRET"),
			Audience:      os.Getenv("AUTH0_AUDIENCE"),
			Issuer:        os.Getenv("AUTH0_ISSUER"),
			Algorithm:     algorithm,
			JSONWebKeySet: os.Getenv("AUTH0_JSON_WEB_KEY_SET"),
		}
		logger := test.NewLogger(lc)
		cache := test.NewRistrettoCache(lc)

		tracer, err := otel.NewTracer(otel.TracerParams{Lifecycle: lc, Config: test.NewOTELConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		params := auth0.CertificatorParams{Config: acfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		cert := auth0.NewCertificator(params)
		ver := auth0.NewVerifier(acfg, cert)

		lc.RequireStart()

		Convey("When I try to verify the token", func() {
			claims := jwt.MapClaims{"aud": acfg.Audience, "iss": acfg.Issuer, "kid": "none"}
			token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

			key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
			So(err, ShouldBeNil)

			tkn, err := token.SignedString(key)
			So(err, ShouldBeNil)

			ctx := context.Background()
			_, _, err = ver.Verify(ctx, []byte(tkn))

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})

		lc.RequireStop()
	})
}
