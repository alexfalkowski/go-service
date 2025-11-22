package metrics_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestInvalidReader(t *testing.T) {
	Convey("Given I have an invalid configuration", t, func() {
		lc := fxtest.NewLifecycle(t)
		cfg := &metrics.Config{Kind: "wrong"}

		Convey("When I try to get a reader", func() {
			_, err := metrics.NewReader(lc, test.Name, cfg)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}
