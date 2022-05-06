package datadog_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/test"
	"github.com/alexfalkowski/go-service/trace/opentracing/datadog"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestDatadog(t *testing.T) {
	Convey("Given I have a configuration", t, func() {
		cfg := test.NewDatadogConfig()
		lc := fxtest.NewLifecycle(t)

		Convey("When I register the trace system", func() {
			t := datadog.NewTracer(datadog.TracerParams{Lifecycle: lc, Host: cfg.Host})

			lc.RequireStart()

			Convey("Then I should have registered successfully", func() {
				So(t, ShouldNotBeNil)
			})

			lc.RequireStop()
		})
	})
}
