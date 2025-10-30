package meta_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/strings"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSnakeCase(t *testing.T) {
	Convey("Given I have some meta values", t, func() {
		ctx := t.Context()
		ctx = meta.WithAttribute(ctx, "testId", meta.String("1"))
		ctx = meta.WithAttribute(ctx, "see", meta.Ignored("secret"))
		ctx = meta.WithAttribute(ctx, "redacted", meta.Redacted("2"))

		Convey("When I get the strings", func() {
			m := meta.SnakeStrings(ctx, strings.Empty)

			Convey("Then I should have valid map", func() {
				So(m, ShouldEqual, meta.Map{"test_id": "1", "redacted": "*"})
			})
		})
	})
}

func TestCamelCase(t *testing.T) {
	Convey("Given I have some meta values", t, func() {
		ctx := t.Context()
		ctx = meta.WithAttribute(ctx, "test_id", meta.String("1"))
		ctx = meta.WithAttribute(ctx, "see", meta.Ignored("secret"))
		ctx = meta.WithAttribute(ctx, "redacted", meta.Redacted("2"))

		Convey("When I get the strings", func() {
			m := meta.CamelStrings(ctx, strings.Empty)

			Convey("Then I should have valid strings", func() {
				So(m, ShouldEqual, meta.Map{"testId": "1", "redacted": "*"})
			})
		})
	})
}

func TestNoneCase(t *testing.T) {
	Convey("Given I have some meta values", t, func() {
		ctx := t.Context()
		ctx = meta.WithAttribute(ctx, "testId", meta.String("1"))
		ctx = meta.WithAttribute(ctx, "see", meta.Ignored("secret"))
		ctx = meta.WithAttribute(ctx, "redacted", meta.Redacted("2"))

		Convey("When I get the strings", func() {
			m := meta.Strings(ctx, strings.Empty)

			Convey("Then I should have valid strings", func() {
				So(m, ShouldEqual, meta.Map{"testId": "1", "redacted": "*"})
			})
		})
	})
}

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
