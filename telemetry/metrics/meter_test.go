package metrics_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/internal/test"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestInvalidMetrics(t *testing.T) {
	Convey("When I try to create a counter", t, func() {
		f := func() { metrics.MustInt64ObservableCounter(test.InvalidMeter{}, "", "") }

		Convey("Then it should panic", func() {
			So(f, ShouldPanic)
		})
	})

	Convey("When I try to create a counter", t, func() {
		f := func() { metrics.MustInt64Counter(test.InvalidMeter{}, "", "") }

		Convey("Then it should panic", func() {
			So(f, ShouldPanic)
		})
	})

	Convey("When I try to create a counter", t, func() {
		f := func() { metrics.MustFloat64ObservableCounter(test.InvalidMeter{}, "", "") }

		Convey("Then it should panic", func() {
			So(f, ShouldPanic)
		})
	})

	Convey("When I try to create a histogram", t, func() {
		f := func() { metrics.MustFloat64Histogram(test.InvalidMeter{}, "", "") }

		Convey("Then it should panic", func() {
			So(f, ShouldPanic)
		})
	})

	Convey("When I try to create a gauge", t, func() {
		f := func() { metrics.MustInt64ObservableGauge(test.InvalidMeter{}, "", "") }

		Convey("Then it should panic", func() {
			So(f, ShouldPanic)
		})
	})
}
