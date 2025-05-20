package meta_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/meta"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSnakeCase(t *testing.T) {
	Convey("Given I have some meta values", t, func() {
		ctx := t.Context()
		ctx = meta.WithAttribute(ctx, "testId", meta.String("1"))
		ctx = meta.WithAttribute(ctx, "see", meta.Ignored("secret"))
		ctx = meta.WithAttribute(ctx, "redacted", meta.Redacted("2"))

		Convey("When I get the strings", func() {
			m := meta.SnakeStrings(ctx, "")

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
			m := meta.CamelStrings(ctx, "")

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
			m := meta.Strings(ctx, "")

			Convey("Then I should have valid strings", func() {
				So(m, ShouldEqual, meta.Map{"testId": "1", "redacted": "*"})
			})
		})
	})
}
