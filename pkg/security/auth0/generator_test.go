package auth0_test

import (
	"context"
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/pkg/cache/ristretto"
	"github.com/alexfalkowski/go-service/pkg/config"
	"github.com/alexfalkowski/go-service/pkg/logger/zap"
	"github.com/alexfalkowski/go-service/pkg/security/auth0"
	"github.com/alexfalkowski/go-service/pkg/transport/http"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestGenerate(t *testing.T) {
	Convey("Given I have a generator", t, func() {
		cfg := &config.Config{AppName: "test"}
		lc := fxtest.NewLifecycle(t)

		acfg, err := auth0.NewConfig()
		So(err, ShouldBeNil)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cache, err := ristretto.NewCache(lc, cfg, ristretto.NewConfig())
		So(err, ShouldBeNil)

		client := http.NewClient(logger)
		gen := auth0.NewGenerator(acfg, client, cache)

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

func TestInvalidGenerate(t *testing.T) {
	Convey("Given I have an invalid generator", t, func() {
		cfg := &config.Config{AppName: "test"}
		lc := fxtest.NewLifecycle(t)

		acfg, err := auth0.NewConfig()
		So(err, ShouldBeNil)

		acfg.ClientSecret = "invalid-secret"

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cache, err := ristretto.NewCache(lc, cfg, ristretto.NewConfig())
		So(err, ShouldBeNil)

		client := http.NewClient(logger)
		gen := auth0.NewGenerator(acfg, client, cache)

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

func TestCachedGenerate(t *testing.T) {
	Convey("Given I have a generator", t, func() {
		cfg := &config.Config{AppName: "test"}
		lc := fxtest.NewLifecycle(t)

		acfg, err := auth0.NewConfig()
		So(err, ShouldBeNil)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cache, err := ristretto.NewCache(lc, cfg, ristretto.NewConfig())
		So(err, ShouldBeNil)

		client := http.NewClient(logger)
		gen := auth0.NewGenerator(acfg, client, cache)

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
