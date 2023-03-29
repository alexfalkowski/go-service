package auth0_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/otel"
	"github.com/alexfalkowski/go-service/security/auth0"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func init() {
	otel.Register()
}

func TestVerify(t *testing.T) {
	Convey("Given I have a valid token", t, func() {
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

		gp := auth0.GeneratorParams{Config: acfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		gen := auth0.NewGenerator(gp)
		cp := auth0.CertificatorParams{Config: acfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		cert := auth0.NewCertificator(cp)
		ver := auth0.NewVerifier(acfg, cert)

		lc.RequireStart()

		Convey("When I verify the token", func() {
			ctx := context.Background()

			token, err := gen.Generate(ctx)
			So(err, ShouldBeNil)

			_, _, err = ver.Verify(ctx, token)

			Convey("Then I should have no errors", func() {
				So(err, ShouldBeNil)
			})
		})

		lc.RequireStop()
	})
}

func TestCachedVerify(t *testing.T) {
	Convey("Given I have a valid token", t, func() {
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

		gp := auth0.GeneratorParams{Config: acfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		gen := auth0.NewGenerator(gp)
		cp := auth0.CertificatorParams{Config: acfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		cert := auth0.NewCertificator(cp)
		ver := auth0.NewVerifier(acfg, cert)

		lc.RequireStart()

		Convey("When I verify the token twice", func() {
			ctx := context.Background()

			token, err := gen.Generate(ctx)
			So(err, ShouldBeNil)

			_, _, err = ver.Verify(ctx, token)
			So(err, ShouldBeNil)

			time.Sleep(1 * time.Second)

			_, _, err = ver.Verify(ctx, token)

			Convey("Then I should have no errors", func() {
				So(err, ShouldBeNil)
			})
		})

		lc.RequireStop()
	})
}

//nolint:dupl
func TestVerifyInvalidAlgorithm(t *testing.T) {
	Convey("Given I have a valid token", t, func() {
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

		gp := auth0.GeneratorParams{Config: acfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		gen := auth0.NewGenerator(gp)
		cp := auth0.CertificatorParams{Config: acfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		cert := auth0.NewCertificator(cp)
		ver := auth0.NewVerifier(acfg, cert)

		lc.RequireStart()

		Convey("When I verify the token", func() {
			ctx := context.Background()

			token, err := gen.Generate(ctx)
			So(err, ShouldBeNil)

			acfg.Algorithm = "Algorithm"

			_, _, err = ver.Verify(ctx, token)

			Convey("Then I should have an invalid algorithm", func() {
				So(err, ShouldBeError)
				So(err, ShouldEqual, auth0.ErrInvalidAlgorithm)
			})
		})

		lc.RequireStop()
	})
}

//nolint:dupl
func TestVerifyInvalidIssuer(t *testing.T) {
	Convey("Given I have a valid token", t, func() {
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

		gp := auth0.GeneratorParams{Config: acfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		gen := auth0.NewGenerator(gp)
		cp := auth0.CertificatorParams{Config: acfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		cert := auth0.NewCertificator(cp)
		ver := auth0.NewVerifier(acfg, cert)

		lc.RequireStart()

		Convey("When I verify the token", func() {
			ctx := context.Background()

			token, err := gen.Generate(ctx)
			So(err, ShouldBeNil)

			acfg.Issuer = "Issuer"

			_, _, err = ver.Verify(ctx, token)

			Convey("Then I should have an invalid issuer", func() {
				So(err, ShouldBeError)
				So(err, ShouldEqual, auth0.ErrInvalidIssuer)
			})
		})

		lc.RequireStop()
	})
}

//nolint:dupl
func TestVerifyInvalidAudience(t *testing.T) {
	Convey("Given I have a valid token", t, func() {
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

		gp := auth0.GeneratorParams{Config: acfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		gen := auth0.NewGenerator(gp)
		cp := auth0.CertificatorParams{Config: acfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		cert := auth0.NewCertificator(cp)
		ver := auth0.NewVerifier(acfg, cert)

		lc.RequireStart()

		Convey("When I verify the token", func() {
			ctx := context.Background()

			token, err := gen.Generate(ctx)
			So(err, ShouldBeNil)

			acfg.Audience = "Audience"

			_, _, err = ver.Verify(ctx, token)

			Convey("Then I should have an invalid audience", func() {
				So(err, ShouldBeError)
				So(err, ShouldEqual, auth0.ErrInvalidAudience)
			})
		})

		lc.RequireStop()
	})
}
