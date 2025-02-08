package meta_test

import (
	"context"
	"testing"

	"github.com/alexfalkowski/go-service/meta"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestSnakeCase(t *testing.T) {
	Convey("Given I have some meta values", t, func() {
		ctx := context.Background()
		ctx = meta.WithAttribute(ctx, "testId", meta.String("1"))

		Convey("When I get the strings", func() {
			m := meta.SnakeStrings(ctx, "")

			Convey("Then I should have valid map", func() {
				So(m, ShouldEqual, meta.Map{"test_id": "1"})
			})
		})
	})
}

func TestCamelCase(t *testing.T) {
	Convey("Given I have some meta values", t, func() {
		ctx := context.Background()
		ctx = meta.WithAttribute(ctx, "test_id", meta.String("1"))

		Convey("When I get the strings", func() {
			m := meta.CamelStrings(ctx, "")

			Convey("Then I should have valid strings", func() {
				So(m, ShouldEqual, meta.Map{"testId": "1"})
			})
		})
	})
}

func TestNoneCase(t *testing.T) {
	Convey("Given I have some meta values", t, func() {
		ctx := context.Background()
		ctx = meta.WithAttribute(ctx, "testId", meta.String("1"))

		Convey("When I get the strings", func() {
			m := meta.Strings(ctx, "")

			Convey("Then I should have valid strings", func() {
				So(m, ShouldEqual, meta.Map{"testId": "1"})
			})
		})
	})
}

func TestBlank(t *testing.T) {
	Convey("When I have a blank value", t, func() {
		v := meta.String("")

		Convey("Then it should be blank", func() {
			So(v.IsBlank(), ShouldBeTrue)
		})
	})

	Convey("When I have a blank value", t, func() {
		v := meta.Redacted("")

		Convey("Then it should be blank", func() {
			So(v.IsBlank(), ShouldBeTrue)
		})
	})

	Convey("When I have a blank value", t, func() {
		v := meta.Ignored("")

		Convey("Then it should be blank", func() {
			So(v.IsBlank(), ShouldBeTrue)
		})
	})

	Convey("When I have a ignored value", t, func() {
		v := meta.Ignored("redacted")

		Convey("Then it should be blank", func() {
			So(v.IsBlank(), ShouldBeTrue)
		})
	})

	Convey("When I have a blank value", t, func() {
		v := meta.Blank()

		Convey("Then it should be blank", func() {
			So(v.IsBlank(), ShouldBeTrue)
		})
	})
}
