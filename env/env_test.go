package env_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/id/uuid"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	. "github.com/smartystreets/goconvey/convey"
)

func TestName(t *testing.T) {
	Convey("When I get a name", t, func() {
		name := env.NewName(test.FS)

		Convey("Then I have a valid name", func() {
			So(name.String(), ShouldEqual, "env.test")
		})
	})

	Convey("Given I have a name set via a env variable", t, func() {
		t.Setenv("SERVICE_NAME", test.Name.String())

		Convey("When I get a name", func() {
			name := env.NewName(test.FS)

			Convey("Then I have a valid name", func() {
				So(name.String(), ShouldEqual, "test")
			})
		})
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

func TestID(t *testing.T) {
	generator := uuid.NewGenerator()

	Convey("When I get a id", t, func() {
		id := env.NewID(generator)

		Convey("Then I have a valid id", func() {
			So(id.String(), ShouldNotBeBlank)
		})
	})

	Convey("Given I have a id set via a env variable", t, func() {
		t.Setenv("SERVICE_ID", "new_id")

		Convey("When I get a id", func() {
			id := env.NewID(generator)

			Convey("Then I have a valid id", func() {
				So(id.String(), ShouldEqual, "new_id")
			})
		})
	})
}
