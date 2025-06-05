package meta_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/transport/meta"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUserID(t *testing.T) {
	Convey("When I set user-id", t, func() {
		ctx := meta.WithUserID(t.Context(), meta.String("user-id"))

		Convey("Then I should have user-id", func() {
			So(meta.UserID(ctx), ShouldEqual, meta.String("user-id"))
		})
	})
}

func TestGeolocation(t *testing.T) {
	Convey("When I set user-agent", t, func() {
		ctx := meta.WithGeolocation(t.Context(), meta.String("user-agent"))

		Convey("Then I should have user-agent", func() {
			So(meta.Geolocation(ctx), ShouldEqual, meta.String("user-agent"))
		})
	})

	Convey("When I set ip address", t, func() {
		ctx := meta.WithGeolocation(t.Context(), meta.String("127.0.0.1"))

		Convey("Then I should have an ip address", func() {
			So(meta.Geolocation(ctx), ShouldEqual, meta.String("127.0.0.1"))
		})
	})

	Convey("When I set geolocation", t, func() {
		ctx := meta.WithGeolocation(t.Context(), meta.String("geo:47,11"))

		Convey("Then I should have geolocation", func() {
			So(meta.Geolocation(ctx), ShouldEqual, meta.String("geo:47,11"))
		})
	})
}
