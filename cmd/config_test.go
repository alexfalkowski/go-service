package cmd_test

import (
	"encoding/base64"
	"testing"

	"github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/go-service/internal/test"
	"github.com/alexfalkowski/go-service/os"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestNoneConfig(t *testing.T) {
	for _, flag := range []string{"", "env:BOB"} {
		Convey("When I read the config", t, func() {
			set := cmd.NewFlagSet("test")
			set.AddInput(flag)

			input := test.NewInputConfig(set)

			Convey("Then I should have a valid configuration", func() {
				So(input.Decode(nil), ShouldBeError)
			})
		})

		Convey("When I write the config", t, func() {
			set := cmd.NewFlagSet("test")
			set.AddOutput(flag)

			output := test.NewOutputConfig(set)
			err := output.Write([]byte("test"), os.ModeAppend)

			Convey("Then I should have a valid configuration", func() {
				So(err, ShouldBeError)
			})
		})
	}
}

func TestReadValidConfigFile(t *testing.T) {
	Convey("Given I have configuration file", t, func() {
		So(os.SetVariable("CONFIG_FILE", "config.yml"), ShouldBeNil)

		set := cmd.NewFlagSet("test")
		set.AddInput("env:CONFIG_FILE")

		Convey("When I read the config", func() {
			input := test.NewInputConfig(set)

			Convey("Then I should have a valid configuration", func() {
				So(input.Kind(), ShouldEqual, "yml")
			})

			So(os.UnsetVariable("CONFIG_FILE"), ShouldBeNil)
		})
	})

	Convey("Given I have configuration file", t, func() {
		Convey("When I read the config", func() {
			set := cmd.NewFlagSet("test")
			set.AddInput("file:config.yml")

			input := test.NewInputConfig(set)

			Convey("Then I should have a valid configuration", func() {
				So(input, ShouldNotBeNil)
			})
		})
	})
}

func TestWriteValidConfigFile(t *testing.T) {
	Convey("Given I have configuration file", t, func() {
		file := test.Path("configs/new_config.yml")

		So(os.WriteFile(file, []byte("environment: development"), 0o600), ShouldBeNil)
		So(os.SetVariable("CONFIG_FILE", file), ShouldBeNil)

		set := cmd.NewFlagSet("test")
		set.AddOutput("env:CONFIG_FILE")

		Convey("When I write the config", func() {
			input := test.NewOutputConfig(set)

			err := input.Write([]byte("test"), os.ModeAppend)
			So(err, ShouldBeNil)

			Convey("Then I should have a valid configuration", func() {
				d, err := os.ReadFile(file)
				So(err, ShouldBeNil)

				So(string(d), ShouldEqual, "test")
			})

			So(os.UnsetVariable("CONFIG_FILE"), ShouldBeNil)
			So(os.Remove(file), ShouldBeNil)
		})
	})

	Convey("Given I have configuration file", t, func() {
		file := test.Path("configs/new_config.yml")

		So(os.WriteFile(file, []byte("environment: development"), 0o600), ShouldBeNil)

		set := cmd.NewFlagSet("test")
		set.AddOutput("file:" + file)

		Convey("When I write the config", func() {
			output := test.NewOutputConfig(set)

			err := output.Write([]byte("test"), os.ModeAppend)
			So(err, ShouldBeNil)

			Convey("Then I should have a valid configuration", func() {
				d, err := os.ReadFile(file)
				So(err, ShouldBeNil)

				So(string(d), ShouldEqual, "test")
			})

			So(os.Remove(file), ShouldBeNil)
		})
	})
}

func TestValidConfigEnv(t *testing.T) {
	Convey("Given I have configuration file", t, func() {
		So(os.SetVariable("CONFIG_FILE", "yaml:CONFIG"), ShouldBeNil)
		So(os.SetVariable("CONFIG", "ZW52aXJvbm1lbnQ6IGRldmVsb3BtZW50Cg=="), ShouldBeNil)

		set := cmd.NewFlagSet("test")
		set.AddInput("env:CONFIG_FILE")

		Convey("When I read the config", func() {
			input := test.NewInputConfig(set)

			Convey("Then I should have a valid configuration", func() {
				So(input, ShouldNotBeNil)
			})

			So(os.UnsetVariable("CONFIG"), ShouldBeNil)
			So(os.UnsetVariable("CONFIG_FILE"), ShouldBeNil)
		})

		Convey("When I write the config", func() {
			output := test.NewInputConfig(set)

			err := output.Write([]byte("test"), os.ModeAppend)
			So(err, ShouldBeNil)

			Convey("Then I should have a valid configuration", func() {
				So(os.GetVariable("CONFIG"), ShouldEqual, base64.StdEncoding.EncodeToString([]byte("test")))
			})

			So(os.UnsetVariable("CONFIG"), ShouldBeNil)
			So(os.UnsetVariable("CONFIG_FILE"), ShouldBeNil)
		})
	})
}

func TestMissingConfig(t *testing.T) {
	Convey("When I read the config", t, func() {
		set := cmd.NewFlagSet("test")
		set.AddInput("env:CONFIG_FILE")

		input := test.NewInputConfig(set)

		Convey("Then I should have a valid configuration", func() {
			So(input, ShouldNotBeNil)
		})
	})
}

func TestNonExistentConfig(t *testing.T) {
	Convey("Given I have non existent configuration file", t, func() {
		So(os.SetVariable("CONFIG_FILE", "../bob"), ShouldBeNil)

		set := cmd.NewFlagSet("test")
		set.AddInput("env:CONFIG_FILE")

		Convey("When I try to parse the configuration file", func() {
			input := test.NewInputConfig(set)

			Convey("Then I should have a valid configuration", func() {
				So(input, ShouldNotBeNil)
			})

			So(os.UnsetVariable("CONFIG_FILE"), ShouldBeNil)
		})
	})
}

func TestInvalidKindConfig(t *testing.T) {
	Convey("When I try to parse the configuration file", t, func() {
		set := cmd.NewFlagSet("test")
		set.AddInput("test:test")

		input := test.NewInputConfig(set)

		Convey("Then I should have a valid configuration", func() {
			So(input, ShouldNotBeNil)
		})
	})
}

func TestInvalidConfig(t *testing.T) {
	Convey("Given I have invalid kind configuration file", t, func() {
		So(os.SetVariable("CONFIG_FILE", "config.go"), ShouldBeNil)

		set := cmd.NewFlagSet("test")
		set.AddInput("env:CONFIG_FILE")

		Convey("When I try to parse the configuration file", func() {
			input := test.NewInputConfig(set)

			Convey("Then I should have an error when decoding", func() {
				So(input.Decode(nil), ShouldBeError)
			})

			So(os.UnsetVariable("CONFIG_FILE"), ShouldBeNil)
		})
	})
}
