package auth0_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/pkg/cache/ristretto"
	"github.com/alexfalkowski/go-service/pkg/config"
	"github.com/alexfalkowski/go-service/pkg/logger/zap"
	"github.com/alexfalkowski/go-service/pkg/security/auth0"
	"github.com/alexfalkowski/go-service/pkg/transport/http"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestVerify(t *testing.T) {
	Convey("Given I have a valid token", t, func() {
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
		cert := auth0.NewCertificator(acfg, client, cache)
		ver := auth0.NewVerifier(acfg, cert)

		lc.RequireStart()

		Convey("When I verify the token", func() {
			token, err := gen.Generate()
			So(err, ShouldBeNil)

			err = ver.Verify(token)

			Convey("Then I should have no errors", func() {
				So(err, ShouldBeNil)
			})
		})

		lc.RequireStop()
	})
}

func TestCachedVerify(t *testing.T) {
	Convey("Given I have a valid token", t, func() {
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
		cert := auth0.NewCertificator(acfg, client, cache)
		ver := auth0.NewVerifier(acfg, cert)

		lc.RequireStart()

		Convey("When I verify the token twice", func() {
			token, err := gen.Generate()
			So(err, ShouldBeNil)

			err = ver.Verify(token)
			So(err, ShouldBeNil)

			err = ver.Verify(token)

			Convey("Then I should have no errors", func() {
				So(err, ShouldBeNil)
			})
		})

		lc.RequireStop()
	})
}
