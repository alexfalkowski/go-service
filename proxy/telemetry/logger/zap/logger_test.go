package zap_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/proxy/telemetry/logger/zap"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
)

func TestLogger(t *testing.T) {
	Convey("Given I have a logger", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		l := zap.NewLogger(logger)

		lc.RequireStart()

		Convey("When I try log", func() {
			f := func() { l.Printf("%s", "test") }

			Convey("Then I should gave a logged message", func() {
				So(f, ShouldNotPanic)
			})
		})

		lc.RequireStop()
	})
}
