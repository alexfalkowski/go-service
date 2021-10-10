package datadog_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/pkg/trace/opentracing/datadog"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestDatadog(t *testing.T) {
	Convey("Given I have a configuration", t, func() {
		cfg := &datadog.Config{
			Host: "localhost:8126",
		}

		Convey("When I register the trace system", func() {
			lc := fxtest.NewLifecycle(t)
			err := datadog.Register(lc, cfg)

			lc.RequireStart()

			Convey("Then I should have registered successfully", func() {
				So(err, ShouldBeNil)
			})

			lc.RequireStop()
		})
	})
}
