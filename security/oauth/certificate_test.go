package oauth_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/alexfalkowski/go-service/security/oauth"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/test"
	htracer "github.com/alexfalkowski/go-service/transport/http/telemetry/tracer"
	"github.com/golang-jwt/jwt/v4"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
)

func init() {
	tracer.Register()
}

func TestInvalidJSONWebKeySet(t *testing.T) {
	Convey("Given I have a corrupt token", t, func() {
		lc := fxtest.NewLifecycle(t)

		cfg := test.NewOAuthConfig()
		cfg.JSONWebKeySet = "not a valid URL"

		logger := test.NewLogger(lc)

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		cache := test.NewRistrettoCache(lc, m)

		tracer, err := htracer.NewTracer(htracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		params := oauth.CertificatorParams{Config: cfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		cert, err := oauth.NewCertificator(params)
		So(err, ShouldBeNil)

		ver := oauth.NewVerifier(cfg, cert)

		lc.RequireStart()

		Convey("When I try to verify the token", func() {
			claims := jwt.MapClaims{"aud": cfg.Audience, "iss": cfg.Issuer, "kid": "none"}
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

		cfg := test.NewOAuthConfig()
		cfg.JSONWebKeySet = "http://localhost:6000/v1/status/400"

		logger := test.NewLogger(lc)

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		cache := test.NewRistrettoCache(lc, m)

		tracer, err := htracer.NewTracer(htracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		params := oauth.CertificatorParams{Config: cfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		cert, err := oauth.NewCertificator(params)
		So(err, ShouldBeNil)

		ver := oauth.NewVerifier(cfg, cert)

		lc.RequireStart()

		Convey("When I try to verify the token", func() {
			claims := jwt.MapClaims{"aud": cfg.Audience, "iss": cfg.Issuer, "kid": "none"}
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

		cfg := test.NewOAuthConfig()
		cfg.JSONWebKeySet = "http://localhost:6000/v1/status/200"

		logger := test.NewLogger(lc)

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		cache := test.NewRistrettoCache(lc, m)

		tracer, err := htracer.NewTracer(htracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		params := oauth.CertificatorParams{Config: cfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		cert, err := oauth.NewCertificator(params)
		So(err, ShouldBeNil)

		ver := oauth.NewVerifier(cfg, cert)

		lc.RequireStart()

		Convey("When I try to verify the token", func() {
			claims := jwt.MapClaims{"aud": cfg.Audience, "iss": cfg.Issuer, "kid": "none"}
			token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

			key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
			So(err, ShouldBeNil)

			tkn, err := token.SignedString(key)
			So(err, ShouldBeNil)

			ctx := context.Background()
			_, _, err = ver.Verify(ctx, []byte(tkn))

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, "json: cannot unmarshal number into Go value of type oauth.jwksResponse")
			})
		})

		lc.RequireStop()
	})
}

func TestCorruptToken(t *testing.T) {
	Convey("Given I have a corrupt token", t, func() {
		lc := fxtest.NewLifecycle(t)
		cfg := test.NewOAuthConfig()
		logger := test.NewLogger(lc)

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		cache := test.NewRistrettoCache(lc, m)

		tracer, err := htracer.NewTracer(htracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		params := oauth.CertificatorParams{Config: cfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		cert, err := oauth.NewCertificator(params)
		So(err, ShouldBeNil)

		ver := oauth.NewVerifier(cfg, cert)

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
		cfg := test.NewOAuthConfig()
		logger := test.NewLogger(lc)

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		cache := test.NewRistrettoCache(lc, m)

		tracer, err := htracer.NewTracer(htracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		params := oauth.CertificatorParams{Config: cfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		cert, err := oauth.NewCertificator(params)
		So(err, ShouldBeNil)

		ver := oauth.NewVerifier(cfg, cert)

		lc.RequireStart()

		Convey("When I try to verify the token", func() {
			token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims{})

			_, pri, err := ed25519.GenerateKey(rand.Reader)
			So(err, ShouldBeNil)

			tkn, err := token.SignedString(pri)
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
		cfg := test.NewOAuthConfig()
		logger := test.NewLogger(lc)

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		cache := test.NewRistrettoCache(lc, m)

		tracer, err := htracer.NewTracer(htracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		params := oauth.CertificatorParams{Config: cfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		cert, err := oauth.NewCertificator(params)
		So(err, ShouldBeNil)

		ver := oauth.NewVerifier(cfg, cert)

		lc.RequireStart()

		Convey("When I try to verify the token", func() {
			claims := jwt.MapClaims{"aud": cfg.Audience}
			token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)

			_, pri, err := ed25519.GenerateKey(rand.Reader)
			So(err, ShouldBeNil)

			tkn, err := token.SignedString(pri)
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
		cfg := test.NewOAuthConfig()
		logger := test.NewLogger(lc)

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		cache := test.NewRistrettoCache(lc, m)

		tracer, err := htracer.NewTracer(htracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		params := oauth.CertificatorParams{Config: cfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		cert, err := oauth.NewCertificator(params)
		So(err, ShouldBeNil)

		ver := oauth.NewVerifier(cfg, cert)

		lc.RequireStart()

		Convey("When I try to verify the token", func() {
			claims := jwt.MapClaims{"aud": cfg.Audience, "iss": cfg.Issuer}
			token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)

			_, pri, err := ed25519.GenerateKey(rand.Reader)
			So(err, ShouldBeNil)

			tkn, err := token.SignedString(pri)
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
		cfg := test.NewOAuthConfig()
		logger := test.NewLogger(lc)

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		cache := test.NewRistrettoCache(lc, m)

		tracer, err := htracer.NewTracer(htracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		params := oauth.CertificatorParams{Config: cfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		cert, err := oauth.NewCertificator(params)
		So(err, ShouldBeNil)

		ver := oauth.NewVerifier(cfg, cert)

		lc.RequireStart()

		Convey("When I try to verify the token", func() {
			claims := jwt.MapClaims{"aud": cfg.Audience, "iss": cfg.Issuer, "kid": "none"}
			token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)

			_, pri, err := ed25519.GenerateKey(rand.Reader)
			So(err, ShouldBeNil)

			tkn, err := token.SignedString(pri)
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
