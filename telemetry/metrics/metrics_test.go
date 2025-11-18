package metrics_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestOTLP(t *testing.T) {
	Convey("Given I register OTLP metrics", t, func() {
		lc := fxtest.NewLifecycle(t)
		m := test.NewOTLPMeter(lc)

		Convey("When I create a metric", func() {
			counter, err := m.Int64Counter("test_otlp")
			So(err, ShouldBeNil)

			lc.RequireStart()
			counter.Add(t.Context(), 1)

			Convey("Then I should have a metric", func() {
				lc.RequireStop()
				So(counter, ShouldNotBeNil)
			})
		})
	})
}

func TestInvalidReader(t *testing.T) {
	Convey("When I try to create a reader with an invalid configuration", t, func() {
		lc := fxtest.NewLifecycle(t)
		_, err := metrics.NewReader(lc, test.Name, &metrics.Config{Kind: "invalid"})

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})
}
