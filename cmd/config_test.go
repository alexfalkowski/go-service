package cmd_test

import (
	"os"
	"testing"

	"github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey"
)

func TestValidConfig(t *testing.T) {
	Convey("Given I have configuration file", t, func() {
		os.Setenv("CONFIG_FILE", "../test/config.yml")

		Convey("When I read the config", func() {
			_, err := test.NewCmdConfig("")
			So(err, ShouldBeNil)

			Convey("Then I should have a valid configuration", func() {
				So(err, ShouldBeNil)
			})

			So(os.Unsetenv("CONFIG_FILE"), ShouldBeNil)
		})
	})
}

func TestMissingConfig(t *testing.T) {
	Convey("Given I don't have configuration file", t, func() {
		Convey("When I read the config", func() {
			_, err := test.NewCmdConfig("")

			Convey("Then I should have an error of missing config file", func() {
				So(err, ShouldBeError)
			})
		})
	})
}

func TestNonExistentConfig(t *testing.T) {
	Convey("Given I have non existent configuration file", t, func() {
		os.Setenv("CONFIG_FILE", "../../test/bob")

		Convey("When I try to parse the configuration file", func() {
			_, err := test.NewCmdConfig("")

			Convey("Then I should have an error of non existent config file", func() {
				So(err, ShouldBeError)
			})

			So(os.Unsetenv("CONFIG_FILE"), ShouldBeNil)
		})
	})
}

func TestInvalidKindConfig(t *testing.T) {
	Convey("Given I have non existent configuration file", t, func() {
		Convey("When I try to parse the configuration file", func() {
			_, err := test.NewCmdConfig("test:test")

			Convey("Then I should have an error of non existent config file", func() {
				So(err, ShouldBeError)
				So(err, ShouldEqual, cmd.ErrInvalidKind)
			})
		})
	})
}

func TestInvalidConfig(t *testing.T) {
	Convey("Given I have invalid kind configuration file", t, func() {
		os.Setenv("CONFIG_FILE", "../test/greet/v1/service.proto")

		Convey("When I try to parse the configuration file", func() {
			_, err := test.NewCmdConfig("test:test")

			Convey("Then I should have an error of invalid kind config file", func() {
				So(err, ShouldBeError)
			})

			So(os.Unsetenv("CONFIG_FILE"), ShouldBeNil)
		})
	})
}
