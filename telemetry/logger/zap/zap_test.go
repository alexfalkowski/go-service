package zap_test

import (
	"testing"

	logger "github.com/alexfalkowski/go-service/telemetry/logger/zap"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

func TestLogger(t *testing.T) {
	Convey("Given I have an invalid zap config", t, func() {
		lc := fxtest.NewLifecycle(t)
		cfg := &logger.Config{}
		c := &zap.Config{}

		Convey("When I try to get a logger", func() {
			p := logger.LoggerParams{Lifecycle: lc, Config: cfg, Logger: c, Environment: test.Environment, Version: test.Version}
			_, err := logger.NewLogger(p)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})

	Convey("Given I have an invalid zap config", t, func() {
		cfg := logger.Config{Level: "bob"}

		Convey("When I try to build a logger config", func() {
			_, err := logger.NewConfig(test.Environment, &cfg)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}
