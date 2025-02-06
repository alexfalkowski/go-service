package env_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/env"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestUserAgent(t *testing.T) {
	Convey("Given I have a name and version", t, func() {
		v := env.Version("v1.0.0")
		n := env.Name("test")

		Convey("When I get a user agent", func() {
			ua := env.NewUserAgent(n, v)

			Convey("Then I have a valid user agent", func() {
				So(string(ua), ShouldEqual, "test/1.0.0")
			})
		})
	})

	Convey("Given I have a name and invalid version", t, func() {
		v := env.Version("test")
		n := env.Name("test")

		Convey("When I get a user agent", func() {
			ua := env.NewUserAgent(n, v)

			Convey("Then I have a valid user agent", func() {
				So(string(ua), ShouldEqual, "test/test")
			})
		})
	})
}
