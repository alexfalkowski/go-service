package tracer_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/transport/grpc/telemetry/tracer"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc/metadata"
)

func TestCarrier(t *testing.T) {
	Convey("Given I have a carrier", t, func() {
		md := metadata.MD{}
		carrier := tracer.NewCarrier(md)

		Convey("When I set some keys", func() {
			carrier.Set("test", "test")

			Convey("Then I should have keys", func() {
				So(carrier.Get("test"), ShouldEqual, "test")
				So(carrier.Keys(), ShouldEqual, []string{"test"})
			})
		})
	})
}
