package cmd_test

import (
	"os"
	"testing"

	"github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/go-service/marshaller"
	. "github.com/smartystreets/goconvey/convey"
)

func TestValidConfig(t *testing.T) {
	Convey("Given I have configuration file", t, func() {
		os.Setenv("CONFIG_FILE", "../test/config.yml")
		cmd.ConfigFlag = ""

		Convey("When I read the config", func() {
			p := marshaller.FactoryParams{YAML: marshaller.NewYAML()}
			c, err := cmd.NewConfig(marshaller.NewFactory(p))
			So(err, ShouldBeNil)

			Convey("Then I should have a valid configuration", func() {
				So(c.Data, ShouldNotBeEmpty)
			})

			So(os.Unsetenv("CONFIG_FILE"), ShouldBeNil)
		})
	})
}

func TestMissingConfig(t *testing.T) {
	Convey("Given I don't have configuration file", t, func() {
		cmd.ConfigFlag = ""

		Convey("When I read the config", func() {
			p := marshaller.FactoryParams{YAML: marshaller.NewYAML()}
			_, err := cmd.NewConfig(marshaller.NewFactory(p))

			Convey("Then I should have an error of missing config file", func() {
				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, "read .: is a directory")
			})
		})
	})
}

func TestNonExistentConfig(t *testing.T) {
	Convey("Given I have non existent configuration file", t, func() {
		os.Setenv("CONFIG_FILE", "../../test/bob")
		cmd.ConfigFlag = ""

		Convey("When I try to parse the configuration file", func() {
			p := marshaller.FactoryParams{YAML: marshaller.NewYAML()}
			_, err := cmd.NewConfig(marshaller.NewFactory(p))

			Convey("Then I should have an error of non existent config file", func() {
				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, "open ../../test/bob: no such file or directory")
			})

			So(os.Unsetenv("CONFIG_FILE"), ShouldBeNil)
		})
	})
}

func TestInvalidKindConfig(t *testing.T) {
	Convey("Given I have non existent configuration file", t, func() {
		cmd.ConfigFlag = "test:test"

		Convey("When I try to parse the configuration file", func() {
			p := marshaller.FactoryParams{YAML: marshaller.NewYAML()}
			_, err := cmd.NewConfig(marshaller.NewFactory(p))

			Convey("Then I should have an error of non existent config file", func() {
				So(err, ShouldBeError)
				So(err, ShouldEqual, cmd.ErrInvalidKind)
			})
		})
	})
}
