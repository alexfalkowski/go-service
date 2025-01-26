package metrics_test

import (
	"context"
	"testing"

	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
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
			counter.Add(context.Background(), 1)

			Convey("Then I should have a metric", func() {
				lc.RequireStop()
				So(counter, ShouldNotBeNil)
			})
		})
	})
}

func TestInvalidReader(t *testing.T) {
	Convey("When I try to create a reader with an invalid reader", t, func() {
		_, err := metrics.NewReader(&test.ErrFS{}, test.NewOTLPMetricsConfig())

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})

	Convey("When I try to create a reader with an invalid configuration", t, func() {
		_, err := metrics.NewReader(test.FS, &metrics.Config{Kind: "invalid"})

		Convey("Then I should have no error", func() {
			So(err, ShouldBeNil)
		})
	})
}
