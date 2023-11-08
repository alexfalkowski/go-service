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

func TestGenerate(t *testing.T) {
	Convey("Given I have a generator", t, func() {
		lc := fxtest.NewLifecycle(t)
		cfg := test.NewOAuthConfig()
		logger := test.NewLogger(lc)

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		cache := test.NewRistrettoCache(lc, m)

		tracer, err := htracer.NewTracer(htracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		params := oauth.GeneratorParams{Config: cfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		gen, err := oauth.NewGenerator(params)
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I generate a token", func() {
			ctx := context.Background()

			_, token, err := gen.Generate(ctx)
			So(err, ShouldBeNil)

			Convey("Then I should have a valid token", func() {
				So(token, ShouldNotBeEmpty)
			})
		})

		lc.RequireStop()
	})
}

func TestInvalidResponseGenerate(t *testing.T) {
	Convey("Given I have an invalid generator", t, func() {
		lc := fxtest.NewLifecycle(t)

		cfg := test.NewOAuthConfig()
		cfg.ClientSecret = "invalid-secret"

		logger := test.NewLogger(lc)

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		cache := test.NewRistrettoCache(lc, m)

		tracer, err := htracer.NewTracer(htracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		params := oauth.GeneratorParams{Config: cfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		gen, err := oauth.NewGenerator(params)
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I generate a token", func() {
			ctx := context.Background()
			_, _, err := gen.Generate(ctx)

			Convey("Then I should have an error", func() {
				So(err, ShouldEqual, oauth.ErrInvalidResponse)
			})
		})

		lc.RequireStop()
	})
}

func TestInvalidURLGenerate(t *testing.T) {
	Convey("Given I have an invalid generator", t, func() {
		lc := fxtest.NewLifecycle(t)

		cfg := test.NewOAuthConfig()
		cfg.URL = "not a valid URL"

		logger := test.NewLogger(lc)

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		cache := test.NewRistrettoCache(lc, m)

		tracer, err := htracer.NewTracer(htracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		params := oauth.GeneratorParams{Config: cfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		gen, err := oauth.NewGenerator(params)
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I generate a token", func() {
			ctx := context.Background()
			_, _, err := gen.Generate(ctx)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})

		lc.RequireStop()
	})
}

func TestMalformedURLGenerate(t *testing.T) {
	Convey("Given I have an invalid generator", t, func() {
		lc := fxtest.NewLifecycle(t)

		cfg := test.NewOAuthConfig()
		cfg.URL = string([]byte{0x7f})

		logger := test.NewLogger(lc)

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		cache := test.NewRistrettoCache(lc, m)

		tracer, err := htracer.NewTracer(htracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		params := oauth.GeneratorParams{Config: cfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		gen, err := oauth.NewGenerator(params)
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I generate a token", func() {
			ctx := context.Background()
			_, _, err := gen.Generate(ctx)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})

		lc.RequireStop()
	})
}

func TestCachedGenerate(t *testing.T) {
	Convey("Given I have a generator", t, func() {
		lc := fxtest.NewLifecycle(t)
		cfg := test.NewOAuthConfig()
		logger := test.NewLogger(lc)

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		cache := test.NewRistrettoCache(lc, m)

		tracer, err := htracer.NewTracer(htracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		params := oauth.GeneratorParams{Config: cfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		gen, err := oauth.NewGenerator(params)
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I generate a token twice", func() {
			ctx := context.Background()

			_, _, err = gen.Generate(ctx)
			So(err, ShouldBeNil)

			time.Sleep(1 * time.Second)

			_, token, err := gen.Generate(ctx)
			So(err, ShouldBeNil)

			Convey("Then I should have a valid cached token", func() {
				So(token, ShouldNotBeEmpty)
			})
		})

		lc.RequireStop()
	})
}
