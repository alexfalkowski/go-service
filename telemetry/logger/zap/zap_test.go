package zap_test

import (
	"testing"

	lzap "github.com/alexfalkowski/go-service/telemetry/logger/zap"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

func TestLogger(t *testing.T) {
	Convey("Given I have an invalid zap config", t, func() {
		lc := fxtest.NewLifecycle(t)
		cfg := zap.Config{}

		Convey("When I try to get a logger", func() {
			_, err := lzap.NewLogger(lzap.LoggerParams{Lifecycle: lc, Config: cfg, Environment: test.Environment, Version: test.Version})

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}
