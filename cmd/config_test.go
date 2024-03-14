package cmd_test

import (
	"encoding/base64"
	"os"
	"testing"

	"github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestReadValidConfigFile(t *testing.T) {
	Convey("Given I have configuration file", t, func() {
		So(os.Setenv("CONFIG_FILE", "../test/config.yml"), ShouldBeNil)

		Convey("When I read the config", func() {
			c, err := test.NewCmdConfig("")

			Convey("Then I should have a valid configuration", func() {
				So(err, ShouldBeNil)
				So(c.Kind(), ShouldEqual, "yml")
			})

			So(os.Unsetenv("CONFIG_FILE"), ShouldBeNil)
		})
	})

	Convey("Given I have configuration file", t, func() {
		Convey("When I read the config", func() {
			_, err := test.NewCmdConfig("file:../test/config.yml")

			Convey("Then I should have a valid configuration", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestWriteValidConfigFile(t *testing.T) {
	Convey("Given I have configuration file", t, func() {
		file := "../test/new_config.yml"

		So(os.WriteFile(file, []byte("environment: development"), os.ModePerm), ShouldBeNil)
		So(os.Setenv("CONFIG_FILE", file), ShouldBeNil)

		Convey("When I write the config", func() {
			c, err := test.NewCmdConfig("")
			So(err, ShouldBeNil)

			err = c.Write([]byte("test"), os.ModeAppend)
			So(err, ShouldBeNil)

			Convey("Then I should have a valid configuration", func() {
				b, err := os.ReadFile(file)
				So(err, ShouldBeNil)

				So(string(b), ShouldEqual, "test")
			})

			So(os.Unsetenv("CONFIG_FILE"), ShouldBeNil)
			So(os.Remove(file), ShouldBeNil)
		})
	})

	Convey("Given I have configuration file", t, func() {
		file := "../test/new_config.yml"

		So(os.WriteFile(file, []byte("environment: development"), os.ModePerm), ShouldBeNil)

		Convey("When I write the config", func() {
			c, err := test.NewCmdConfig("file:" + file)
			So(err, ShouldBeNil)

			err = c.Write([]byte("test"), os.ModeAppend)
			So(err, ShouldBeNil)

			Convey("Then I should have a valid configuration", func() {
				b, err := os.ReadFile(file)
				So(err, ShouldBeNil)

				So(string(b), ShouldEqual, "test")
			})

			So(os.Remove(file), ShouldBeNil)
		})
	})
}

func TestValidConfigEnv(t *testing.T) {
	Convey("Given I have configuration file", t, func() {
		So(os.Setenv("CONFIG_FILE", "yaml:CONFIG"), ShouldBeNil)
		So(os.Setenv("CONFIG", "ZW52aXJvbm1lbnQ6IGRldmVsb3BtZW50Cg=="), ShouldBeNil)

		Convey("When I read the config", func() {
			_, err := test.NewCmdConfig("")

			Convey("Then I should have a valid configuration", func() {
				So(err, ShouldBeNil)
			})

			So(os.Unsetenv("CONFIG"), ShouldBeNil)
			So(os.Unsetenv("CONFIG_FILE"), ShouldBeNil)
		})

		Convey("When I write the config", func() {
			c, err := test.NewCmdConfig("")
			So(err, ShouldBeNil)

			err = c.Write([]byte("test"), os.ModeAppend)
			So(err, ShouldBeNil)

			Convey("Then I should have a valid configuration", func() {
				So(os.Getenv("CONFIG"), ShouldEqual, base64.StdEncoding.EncodeToString([]byte("test")))
			})

			So(os.Unsetenv("CONFIG"), ShouldBeNil)
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
		So(os.Setenv("CONFIG_FILE", "../../test/bob"), ShouldBeNil)

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
		So(os.Setenv("CONFIG_FILE", "../test/greet/v1/service.proto"), ShouldBeNil)

		Convey("When I try to parse the configuration file", func() {
			_, err := test.NewCmdConfig("test:test")

			Convey("Then I should have an error of invalid kind config file", func() {
				So(err, ShouldBeError)
			})

			So(os.Unsetenv("CONFIG_FILE"), ShouldBeNil)
		})
	})
}
