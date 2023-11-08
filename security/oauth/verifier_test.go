package oauth_test

import (
	"context"
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/security/oauth"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/test"
	htracer "github.com/alexfalkowski/go-service/transport/http/telemetry/tracer"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
)

func init() {
	tracer.Register()
}

func TestVerify(t *testing.T) {
	Convey("Given I have a valid token", t, func() {
		lc := fxtest.NewLifecycle(t)
		cfg := test.NewOAuthConfig()
		logger := test.NewLogger(lc)

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		cache := test.NewRistrettoCache(lc, m)

		tracer, err := htracer.NewTracer(htracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		gp := oauth.GeneratorParams{Config: cfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		gen, err := oauth.NewGenerator(gp)
		So(err, ShouldBeNil)

		cp := oauth.CertificatorParams{Config: cfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		cert, err := oauth.NewCertificator(cp)
		So(err, ShouldBeNil)

		ver := oauth.NewVerifier(cfg, cert)

		lc.RequireStart()

		Convey("When I verify the token", func() {
			ctx := context.Background()

			_, token, err := gen.Generate(ctx)
			So(err, ShouldBeNil)

			_, err = ver.Verify(ctx, token)

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
		cfg := test.NewOAuthConfig()
		logger := test.NewLogger(lc)

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		cache := test.NewRistrettoCache(lc, m)

		tracer, err := htracer.NewTracer(htracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		gp := oauth.GeneratorParams{Config: cfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		gen, err := oauth.NewGenerator(gp)
		So(err, ShouldBeNil)

		cp := oauth.CertificatorParams{Config: cfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		cert, err := oauth.NewCertificator(cp)
		So(err, ShouldBeNil)

		ver := oauth.NewVerifier(cfg, cert)

		lc.RequireStart()

		Convey("When I verify the token twice", func() {
			ctx := context.Background()

			_, token, err := gen.Generate(ctx)
			So(err, ShouldBeNil)

			_, err = ver.Verify(ctx, token)
			So(err, ShouldBeNil)

			time.Sleep(1 * time.Second)

			_, err = ver.Verify(ctx, token)

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
		cfg := test.NewOAuthConfig()
		logger := test.NewLogger(lc)

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		cache := test.NewRistrettoCache(lc, m)

		tracer, err := htracer.NewTracer(htracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		gp := oauth.GeneratorParams{Config: cfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		gen, err := oauth.NewGenerator(gp)
		So(err, ShouldBeNil)

		cp := oauth.CertificatorParams{Config: cfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		cert, err := oauth.NewCertificator(cp)
		So(err, ShouldBeNil)

		ver := oauth.NewVerifier(cfg, cert)

		lc.RequireStart()

		Convey("When I verify the token", func() {
			ctx := context.Background()

			_, token, err := gen.Generate(ctx)
			So(err, ShouldBeNil)

			cfg.Algorithm = "Algorithm"

			_, err = ver.Verify(ctx, token)

			Convey("Then I should have an invalid algorithm", func() {
				So(err, ShouldBeError)
				So(err, ShouldEqual, oauth.ErrInvalidAlgorithm)
			})
		})

		lc.RequireStop()
	})
}

//nolint:dupl
func TestVerifyInvalidIssuer(t *testing.T) {
	Convey("Given I have a valid token", t, func() {
		lc := fxtest.NewLifecycle(t)
		cfg := test.NewOAuthConfig()
		logger := test.NewLogger(lc)

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		cache := test.NewRistrettoCache(lc, m)

		tracer, err := htracer.NewTracer(htracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		gp := oauth.GeneratorParams{Config: cfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		gen, err := oauth.NewGenerator(gp)
		So(err, ShouldBeNil)

		cp := oauth.CertificatorParams{Config: cfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		cert, err := oauth.NewCertificator(cp)
		So(err, ShouldBeNil)

		ver := oauth.NewVerifier(cfg, cert)

		lc.RequireStart()

		Convey("When I verify the token", func() {
			ctx := context.Background()

			_, token, err := gen.Generate(ctx)
			So(err, ShouldBeNil)

			cfg.Issuer = "Issuer"

			_, err = ver.Verify(ctx, token)

			Convey("Then I should have an invalid issuer", func() {
				So(err, ShouldBeError)
				So(err, ShouldEqual, oauth.ErrInvalidIssuer)
			})
		})

		lc.RequireStop()
	})
}

//nolint:dupl
func TestVerifyInvalidAudience(t *testing.T) {
	Convey("Given I have a valid token", t, func() {
		lc := fxtest.NewLifecycle(t)
		cfg := test.NewOAuthConfig()
		logger := test.NewLogger(lc)

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		cache := test.NewRistrettoCache(lc, m)

		tracer, err := htracer.NewTracer(htracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		gp := oauth.GeneratorParams{Config: cfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		gen, err := oauth.NewGenerator(gp)
		So(err, ShouldBeNil)

		cp := oauth.CertificatorParams{Config: cfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		cert, err := oauth.NewCertificator(cp)
		So(err, ShouldBeNil)

		ver := oauth.NewVerifier(cfg, cert)

		lc.RequireStart()

		Convey("When I verify the token", func() {
			ctx := context.Background()

			_, token, err := gen.Generate(ctx)
			So(err, ShouldBeNil)

			cfg.Audience = "Audience"

			_, err = ver.Verify(ctx, token)

			Convey("Then I should have an invalid audience", func() {
				So(err, ShouldBeError)
				So(err, ShouldEqual, oauth.ErrInvalidAudience)
			})
		})

		lc.RequireStop()
	})
}
