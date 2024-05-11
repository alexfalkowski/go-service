package cmd_test

import (
	"encoding/base64"
	"os"
	"testing"

	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestNoneConfig(t *testing.T) {
	Convey("Given I have no configuration file", t, func() {
		Convey("When I read the config", func() {
			c := test.NewInputConfig("")

			Convey("Then I should have a valid configuration", func() {
				So(c.Unmarshal(nil), ShouldBeNil)
				So(c.Kind(), ShouldEqual, "none")
			})
		})

		Convey("When I write the config", func() {
			c := test.NewOutputConfig("")
			err := c.Write([]byte("test"), os.ModeAppend)

			Convey("Then I should have a valid configuration", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestReadValidConfigFile(t *testing.T) {
	Convey("Given I have configuration file", t, func() {
		So(os.Setenv("CONFIG_FILE", "../test/config.yml"), ShouldBeNil)

		Convey("When I read the config", func() {
			c := test.NewInputConfig("env:CONFIG_FILE")

			Convey("Then I should have a valid configuration", func() {
				So(c.Kind(), ShouldEqual, "yml")
			})

			So(os.Unsetenv("CONFIG_FILE"), ShouldBeNil)
		})
	})

	Convey("Given I have configuration file", t, func() {
		Convey("When I read the config", func() {
			c := test.NewInputConfig("file:../test/config.yml")

			Convey("Then I should have a valid configuration", func() {
				So(c, ShouldNotBeNil)
			})
		})
	})
}

func TestWriteValidConfigFile(t *testing.T) {
	Convey("Given I have configuration file", t, func() {
		file := "../test/configs/new_config.yml"

		So(os.WriteFile(file, []byte("environment: development"), os.ModePerm), ShouldBeNil)
		So(os.Setenv("CONFIG_FILE", file), ShouldBeNil)

		Convey("When I write the config", func() {
			c := test.NewOutputConfig("env:CONFIG_FILE")

			err := c.Write([]byte("test"), os.ModeAppend)
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
		file := "../test/configs/new_config.yml"

		So(os.WriteFile(file, []byte("environment: development"), os.ModePerm), ShouldBeNil)

		Convey("When I write the config", func() {
			c := test.NewOutputConfig("file:" + file)

			err := c.Write([]byte("test"), os.ModeAppend)
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
			c := test.NewInputConfig("env:CONFIG_FILE")

			Convey("Then I should have a valid configuration", func() {
				So(c, ShouldNotBeNil)
			})

			So(os.Unsetenv("CONFIG"), ShouldBeNil)
			So(os.Unsetenv("CONFIG_FILE"), ShouldBeNil)
		})

		Convey("When I write the config", func() {
			c := test.NewInputConfig("env:CONFIG_FILE")

			err := c.Write([]byte("test"), os.ModeAppend)
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
			c := test.NewInputConfig("env:CONFIG_FILE")

			Convey("Then I should have a valid configuration", func() {
				So(c, ShouldNotBeNil)
			})
		})
	})
}

func TestNonExistentConfig(t *testing.T) {
	Convey("Given I have non existent configuration file", t, func() {
		So(os.Setenv("CONFIG_FILE", "../../test/bob"), ShouldBeNil)

		Convey("When I try to parse the configuration file", func() {
			c := test.NewInputConfig("env:CONFIG_FILE")

			Convey("Then I should have a valid configuration", func() {
				So(c, ShouldNotBeNil)
			})

			So(os.Unsetenv("CONFIG_FILE"), ShouldBeNil)
		})
	})
}

func TestInvalidKindConfig(t *testing.T) {
	Convey("Given I have non existent configuration file", t, func() {
		Convey("When I try to parse the configuration file", func() {
			c := test.NewInputConfig("test:test")

			Convey("Then I should have a valid configuration", func() {
				So(c, ShouldNotBeNil)
			})
		})
	})
}

func TestInvalidConfig(t *testing.T) {
	Convey("Given I have invalid kind configuration file", t, func() {
		So(os.Setenv("CONFIG_FILE", "../test/greet/v1/service.proto"), ShouldBeNil)

		Convey("When I try to parse the configuration file", func() {
			c := test.NewInputConfig("test:test")

			Convey("Then I should have a valid configuration", func() {
				So(c, ShouldNotBeNil)
			})

			So(os.Unsetenv("CONFIG_FILE"), ShouldBeNil)
		})
	})
}
