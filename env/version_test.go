package env_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/strings"
	. "github.com/smartystreets/goconvey/convey"
)

func TestVersion(t *testing.T) {
	Convey("Given I have a system version", t, func() {
		v := env.NewVersion()

		Convey("When I get the string representation", func() {
			s := v.String()

			Convey("Then I have a valid version", func() {
				So(s, ShouldEqual, "(devel)")
			})
		})
	})

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
		v := env.Version(strings.Empty)

		Convey("When I get the string representation", func() {
			s := v.String()

			Convey("Then I have an empty string", func() {
				So(s, ShouldBeBlank)
			})
		})
	})

	Convey("Given I have a version set via a env variable", t, func() {
		t.Setenv("SERVICE_VERSION", test.Version.String())

		Convey("When I get a version", func() {
			version := env.NewVersion()

			Convey("Then I have a valid version", func() {
				So(version.String(), ShouldEqual, "1.0.0")
			})
		})
	})
}
