package tracer_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/transport/grpc/telemetry/tracer"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"google.golang.org/grpc/metadata"
)

func TestCarrier(t *testing.T) {
	Convey("Given I have a carrier", t, func() {
		c := &tracer.Carrier{Metadata: metadata.MD{}}

		Convey("When I set some keys", func() {
			c.Set("test", "test")

			Convey("Then I should have keys", func() {
				So(c.Get("test"), ShouldEqual, "test")
				So(c.Keys(), ShouldEqual, []string{"test"})
			})
		})
	})
}
