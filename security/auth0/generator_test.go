package auth0_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/security/auth0"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/test"
	htracer "github.com/alexfalkowski/go-service/transport/http/telemetry/tracer"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func init() {
	tracer.Register()
}

func TestGenerate(t *testing.T) {
	Convey("Given I have a generator", t, func() {
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

		m, err := metrics.NewMeter(lc)
		So(err, ShouldBeNil)

		cache := test.NewRistrettoCache(lc, m)

		tracer, err := htracer.NewTracer(htracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		params := auth0.GeneratorParams{Config: acfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		gen, err := auth0.NewGenerator(params)
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I generate a token", func() {
			ctx := context.Background()

			token, err := gen.Generate(ctx)
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
		acfg := &auth0.Config{
			URL:           os.Getenv("AUTH0_URL"),
			ClientID:      os.Getenv("AUTH0_CLIENT_ID"),
			ClientSecret:  "invalid-secret",
			Audience:      os.Getenv("AUTH0_AUDIENCE"),
			Issuer:        os.Getenv("AUTH0_ISSUER"),
			Algorithm:     os.Getenv("AUTH0_ALGORITHM"),
			JSONWebKeySet: os.Getenv("AUTH0_JSON_WEB_KEY_SET"),
		}
		logger := test.NewLogger(lc)

		m, err := metrics.NewMeter(lc)
		So(err, ShouldBeNil)

		cache := test.NewRistrettoCache(lc, m)

		tracer, err := htracer.NewTracer(htracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		params := auth0.GeneratorParams{Config: acfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		gen, err := auth0.NewGenerator(params)
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I generate a token", func() {
			ctx := context.Background()
			_, err := gen.Generate(ctx)

			Convey("Then I should have an error", func() {
				So(err, ShouldEqual, auth0.ErrInvalidResponse)
			})
		})

		lc.RequireStop()
	})
}

func TestInvalidURLGenerate(t *testing.T) {
	Convey("Given I have an invalid generator", t, func() {
		lc := fxtest.NewLifecycle(t)
		acfg := &auth0.Config{
			URL:           "not a valid URL",
			ClientID:      os.Getenv("AUTH0_CLIENT_ID"),
			ClientSecret:  "invalid-secret",
			Audience:      os.Getenv("AUTH0_AUDIENCE"),
			Issuer:        os.Getenv("AUTH0_ISSUER"),
			Algorithm:     os.Getenv("AUTH0_ALGORITHM"),
			JSONWebKeySet: os.Getenv("AUTH0_JSON_WEB_KEY_SET"),
		}
		logger := test.NewLogger(lc)

		m, err := metrics.NewMeter(lc)
		So(err, ShouldBeNil)

		cache := test.NewRistrettoCache(lc, m)

		tracer, err := htracer.NewTracer(htracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		params := auth0.GeneratorParams{Config: acfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		gen, err := auth0.NewGenerator(params)
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I generate a token", func() {
			ctx := context.Background()
			_, err := gen.Generate(ctx)

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
		acfg := &auth0.Config{
			URL:           string([]byte{0x7f}),
			ClientID:      os.Getenv("AUTH0_CLIENT_ID"),
			ClientSecret:  "invalid-secret",
			Audience:      os.Getenv("AUTH0_AUDIENCE"),
			Issuer:        os.Getenv("AUTH0_ISSUER"),
			Algorithm:     os.Getenv("AUTH0_ALGORITHM"),
			JSONWebKeySet: os.Getenv("AUTH0_JSON_WEB_KEY_SET"),
		}
		logger := test.NewLogger(lc)

		m, err := metrics.NewMeter(lc)
		So(err, ShouldBeNil)

		cache := test.NewRistrettoCache(lc, m)

		tracer, err := htracer.NewTracer(htracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		params := auth0.GeneratorParams{Config: acfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		gen, err := auth0.NewGenerator(params)
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I generate a token", func() {
			ctx := context.Background()
			_, err := gen.Generate(ctx)

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

		m, err := metrics.NewMeter(lc)
		So(err, ShouldBeNil)

		cache := test.NewRistrettoCache(lc, m)

		tracer, err := htracer.NewTracer(htracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		params := auth0.GeneratorParams{Config: acfg, HTTPConfig: &test.NewTransportConfig().HTTP, Cache: cache, Logger: logger, Tracer: tracer}
		gen, err := auth0.NewGenerator(params)
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I generate a token twice", func() {
			ctx := context.Background()

			_, err = gen.Generate(ctx)
			So(err, ShouldBeNil)

			time.Sleep(1 * time.Second)

			token, err := gen.Generate(ctx)
			So(err, ShouldBeNil)

			Convey("Then I should have a valid cached token", func() {
				So(token, ShouldNotBeEmpty)
			})
		})

		lc.RequireStop()
	})
}
