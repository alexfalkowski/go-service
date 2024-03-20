package zap_test

import (
	"testing"

	lzap "github.com/alexfalkowski/go-service/telemetry/logger/zap"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

func TestDevLogger(t *testing.T) {
	Convey("Given I have an invalid zap config", t, func() {
		lc := fxtest.NewLifecycle(t)
		cfg := lzap.Config{Enabled: true}
		c := zap.Config{}

		Convey("When I try to get a logger", func() {
			p := lzap.LoggerParams{Lifecycle: lc, Config: &cfg, ZapConfig: c, Environment: test.DevEnvironment, Version: test.Version}
			_, err := lzap.NewLogger(p)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestProdLogger(t *testing.T) {
	Convey("Given I have an invalid zap config", t, func() {
		lc := fxtest.NewLifecycle(t)
		cfg := lzap.Config{Enabled: true}
		c := zap.Config{}

		Convey("When I try to get a logger", func() {
			p := lzap.LoggerParams{Lifecycle: lc, Config: &cfg, ZapConfig: c, Environment: test.ProdEnvironment, Version: test.Version}
			_, err := lzap.NewLogger(p)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}
