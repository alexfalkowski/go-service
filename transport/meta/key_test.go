package meta_test

import (
	"context"
	"testing"

	"github.com/alexfalkowski/go-service/meta"
	tm "github.com/alexfalkowski/go-service/transport/meta"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestKeys(t *testing.T) {
	Convey("When I set user-agent", t, func() {
		ctx := context.Background()
		ctx = tm.WithGeolocation(ctx, meta.String("user-agent"))

		Convey("Then I should have user-agent", func() {
			So(tm.Geolocation(ctx), ShouldEqual, meta.String("user-agent"))
		})
	})

	Convey("When I set ip address", t, func() {
		ctx := context.Background()
		ctx = tm.WithGeolocation(ctx, meta.String("127.0.0.1"))

		Convey("Then I should have an ip address", func() {
			So(tm.Geolocation(ctx), ShouldEqual, meta.String("127.0.0.1"))
		})
	})

	Convey("When I set geolocation", t, func() {
		ctx := context.Background()
		ctx = tm.WithGeolocation(ctx, meta.String("geo:47,11"))

		Convey("Then I should have geolocation", func() {
			So(tm.Geolocation(ctx), ShouldEqual, meta.String("geo:47,11"))
		})
	})
}
