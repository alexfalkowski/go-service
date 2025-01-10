//nolint:varnamelen
package metrics_test

import (
	"context"
	"testing"

	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
)

func TestOTLP(t *testing.T) {
	Convey("Given I register OTLP metrics", t, func() {
		lc := fxtest.NewLifecycle(t)
		m := test.NewOTLPMeter(lc)

		Convey("When I create a metric", func() {
			c, err := m.Int64Counter("test_otlp")
			So(err, ShouldBeNil)

			lc.RequireStart()
			c.Add(context.Background(), 1)

			Convey("Then I should have a metric", func() {
				lc.RequireStop()
				So(c, ShouldNotBeNil)
			})
		})
	})
}
