package auth0_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/cache/ristretto"
	"github.com/alexfalkowski/go-service/logger/zap"
	"github.com/alexfalkowski/go-service/security/auth0"
	"github.com/alexfalkowski/go-service/test"
	"github.com/alexfalkowski/go-service/transport/http/trace/opentracing"
	"github.com/alexfalkowski/go-service/version"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestVerify(t *testing.T) {
	Convey("Given I have a valid token", t, func() {
		cfg := &ristretto.Config{
			NumCounters: 1e7,
			MaxCost:     1 << 30,
			BufferItems: 64,
		}

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

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		tracer, err := opentracing.NewTracer(lc, test.NewJaegerConfig())
		So(err, ShouldBeNil)

		cache, err := ristretto.NewCache(ristretto.CacheParams{Lifecycle: lc, Config: cfg, Version: version.Version("1.0.0")})
		So(err, ShouldBeNil)

		gp := auth0.GeneratorParams{Config: acfg, HTTPConfig: test.NewHTTPConfig(), Cache: cache, Logger: logger, Tracer: tracer}
		gen := auth0.NewGenerator(gp)
		cp := auth0.CertificatorParams{Config: acfg, HTTPConfig: test.NewHTTPConfig(), Cache: cache, Logger: logger, Tracer: tracer}
		cert := auth0.NewCertificator(cp)
		ver := auth0.NewVerifier(acfg, cert)

		lc.RequireStart()

		Convey("When I verify the token", func() {
			ctx := context.Background()

			token, err := gen.Generate(ctx)
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
		cfg := &ristretto.Config{
			NumCounters: 1e7,
			MaxCost:     1 << 30,
			BufferItems: 64,
		}

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

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		tracer, err := opentracing.NewTracer(lc, test.NewJaegerConfig())
		So(err, ShouldBeNil)

		cache, err := ristretto.NewCache(ristretto.CacheParams{Lifecycle: lc, Config: cfg, Version: version.Version("1.0.0")})
		So(err, ShouldBeNil)

		gp := auth0.GeneratorParams{Config: acfg, HTTPConfig: test.NewHTTPConfig(), Cache: cache, Logger: logger, Tracer: tracer}
		gen := auth0.NewGenerator(gp)
		cp := auth0.CertificatorParams{Config: acfg, HTTPConfig: test.NewHTTPConfig(), Cache: cache, Logger: logger, Tracer: tracer}
		cert := auth0.NewCertificator(cp)
		ver := auth0.NewVerifier(acfg, cert)

		lc.RequireStart()

		Convey("When I verify the token twice", func() {
			ctx := context.Background()

			token, err := gen.Generate(ctx)
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
