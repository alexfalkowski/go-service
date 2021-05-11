package auth0_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/pkg/logger/zap"
	"github.com/alexfalkowski/go-service/pkg/security/auth0"
	"github.com/alexfalkowski/go-service/pkg/transport/http"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestGenerate(t *testing.T) {
	Convey("Given I have a generator", t, func() {
		lc := fxtest.NewLifecycle(t)

		cfg, err := auth0.NewConfig()
		So(err, ShouldBeNil)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		client := http.NewClient(logger)
		gen := auth0.NewGenerator(cfg, client)

		Convey("When I generate a token", func() {
			token, err := gen.Generate()
			So(err, ShouldBeNil)

			Convey("Then I should have a valid token", func() {
				So(token, ShouldNotBeEmpty)
			})
		})
	})
}

func TestInvalidGenerate(t *testing.T) {
	Convey("Given I have an invalid generator", t, func() {
		lc := fxtest.NewLifecycle(t)

		cfg, err := auth0.NewConfig()
		So(err, ShouldBeNil)

		cfg.ClientSecret = "invalid-secret"

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		client := http.NewClient(logger)
		gen := auth0.NewGenerator(cfg, client)

		Convey("When I generate a token", func() {
			_, err := gen.Generate()

			Convey("Then I should have an error", func() {
				So(err, ShouldEqual, auth0.ErrInvalidResponse)
			})
		})
	})
}
