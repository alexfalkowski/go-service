package telemetry_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/telemetry"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestLogger(t *testing.T) {
	Convey("Given I have an valid zap config", t, func() {
		lc := fxtest.NewLifecycle(t)

		Convey("When I try to get a logger", func() {
			_, err := telemetry.NewLogger(telemetry.LoggerParams{Lifecycle: lc, Version: test.Version})

			Convey("Then I should not have an error", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}
