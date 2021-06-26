package auth0_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/pkg/cache/ristretto"
	"github.com/alexfalkowski/go-service/pkg/logger/zap"
	"github.com/alexfalkowski/go-service/pkg/security/auth0"
	"github.com/alexfalkowski/go-service/pkg/transport/http"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestVerify(t *testing.T) {
	Convey("Given I have a valid token", t, func() {
		os.Setenv("SERVICE_NAME", "test")

		cfg, err := ristretto.NewConfig()
		So(err, ShouldBeNil)

		lc := fxtest.NewLifecycle(t)

		acfg, err := auth0.NewConfig()
		So(err, ShouldBeNil)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cache, err := ristretto.NewCache(lc, cfg)
		So(err, ShouldBeNil)

		client := http.NewClient(logger)
		gen := auth0.NewGenerator(acfg, client, cache)
		cert := auth0.NewCertificator(acfg, client, cache)
		ver := auth0.NewVerifier(acfg, cert)

		lc.RequireStart()

		Convey("When I verify the token", func() {
			ctx := context.Background()

			token, err := gen.Generate(context.Background())
			So(err, ShouldBeNil)

			_, err = ver.Verify(ctx, token)

			Convey("Then I should have no errors", func() {
				So(err, ShouldBeNil)
			})
		})

		lc.RequireStop()
		So(os.Unsetenv("SERVICE_NAME"), ShouldBeNil)
	})
}

func TestCachedVerify(t *testing.T) {
	Convey("Given I have a valid token", t, func() {
		os.Setenv("SERVICE_NAME", "test")

		cfg, err := ristretto.NewConfig()
		So(err, ShouldBeNil)

		lc := fxtest.NewLifecycle(t)

		acfg, err := auth0.NewConfig()
		So(err, ShouldBeNil)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cache, err := ristretto.NewCache(lc, cfg)
		So(err, ShouldBeNil)

		client := http.NewClient(logger)
		gen := auth0.NewGenerator(acfg, client, cache)
		cert := auth0.NewCertificator(acfg, client, cache)
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
		So(os.Unsetenv("SERVICE_NAME"), ShouldBeNil)
	})
}
