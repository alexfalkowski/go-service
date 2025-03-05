package meta_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/meta"
	tm "github.com/alexfalkowski/go-service/transport/meta"
	. "github.com/smartystreets/goconvey/convey"
)

func TestKeys(t *testing.T) {
	Convey("When I set user-agent", t, func() {
		ctx := tm.WithGeolocation(t.Context(), meta.String("user-agent"))

		Convey("Then I should have user-agent", func() {
			So(tm.Geolocation(ctx), ShouldEqual, meta.String("user-agent"))
		})
	})

	Convey("When I set ip address", t, func() {
		ctx := tm.WithGeolocation(t.Context(), meta.String("127.0.0.1"))

		Convey("Then I should have an ip address", func() {
			So(tm.Geolocation(ctx), ShouldEqual, meta.String("127.0.0.1"))
		})
	})

	Convey("When I set geolocation", t, func() {
		ctx := tm.WithGeolocation(t.Context(), meta.String("geo:47,11"))

		Convey("Then I should have geolocation", func() {
			So(tm.Geolocation(ctx), ShouldEqual, meta.String("geo:47,11"))
		})
	})
}
