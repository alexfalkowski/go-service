package env_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/internal/test"
	"github.com/alexfalkowski/go-service/os"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestName(t *testing.T) {
	Convey("When I get a name", t, func() {
		name := env.NewName()

		Convey("Then I have a valid name", func() {
			So(name.String(), ShouldEqual, "env.test")
		})
	})

	Convey("Given I have a name set via a env variable", t, func() {
		So(os.SetVariable("SERVICE_NAME", test.Name.String()), ShouldBeNil)

		Convey("When I get a name", func() {
			name := env.NewName()

			Convey("Then I have a valid name", func() {
				So(name.String(), ShouldEqual, "test")
			})
		})

		So(os.UnsetVariable("SERVICE_NAME"), ShouldBeNil)
	})
}

func TestUserAgent(t *testing.T) {
	Convey("Given I have a name and version", t, func() {
		version := env.Version("v1.0.0")
		name := env.Name("test")

		Convey("When I get a user agent", func() {
			ua := env.NewUserAgent(name, version)

			Convey("Then I have a valid user agent", func() {
				So(ua.String(), ShouldEqual, "test/1.0.0")
			})
		})
	})

	Convey("Given I have a name and invalid version", t, func() {
		version := env.Version("test")
		name := env.Name("test")

		Convey("When I get a user agent", func() {
			ua := env.NewUserAgent(name, version)

			Convey("Then I have a valid user agent", func() {
				So(ua.String(), ShouldEqual, "test/test")
			})
		})
	})
}
