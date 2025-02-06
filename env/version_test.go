package env_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/env"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestVersion(t *testing.T) {
	Convey("Given I have a version", t, func() {
		v := env.Version("v1.0.0")

		Convey("When I get the string representation", func() {
			s := v.String()

			Convey("Then I have a valid version", func() {
				So(s, ShouldEqual, "1.0.0")
			})
		})
	})

	Convey("Given I have an invalid version", t, func() {
		v := env.Version("what")

		Convey("When I get the string representation", func() {
			s := v.String()

			Convey("Then I have the same invalid version", func() {
				So(s, ShouldEqual, "what")
			})
		})
	})

	Convey("Given I have an empty version", t, func() {
		v := env.Version("")

		Convey("When I get the string representation", func() {
			s := v.String()

			Convey("Then I have an empty string", func() {
				So(s, ShouldBeBlank)
			})
		})
	})
}
