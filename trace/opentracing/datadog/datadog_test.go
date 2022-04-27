package datadog_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/trace/opentracing/datadog"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestDatadog(t *testing.T) {
	Convey("Given I have a configuration", t, func() {
		cfg := &datadog.Config{Host: "localhost:8126"}
		lc := fxtest.NewLifecycle(t)

		Convey("When I register the trace system", func() {
			t := datadog.NewTracer(datadog.TracerParams{Lifecycle: lc, Config: cfg})

			lc.RequireStart()

			Convey("Then I should have registered successfully", func() {
				So(t, ShouldNotBeNil)
			})

			lc.RequireStop()
		})
	})
}
