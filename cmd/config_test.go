package cmd_test

import (
	"path/filepath"
	"testing"

	"github.com/alexfalkowski/go-service/bytes"
	"github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/go-service/encoding/base64"
	"github.com/alexfalkowski/go-service/internal/test"
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/strings"
	. "github.com/smartystreets/goconvey/convey"
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
			err := output.Write(strings.Bytes("test"), os.ModeAppend)

			Convey("Then I should have a valid configuration", func() {
				So(err, ShouldBeError)
			})
		})
	}
}

//nolint:funlen
func TestReadValidConfigFile(t *testing.T) {
	Convey("Given I have configuration file", t, func() {
		t.Setenv("CONFIG_FILE", test.Path("configs/config.yml"))

		set := cmd.NewFlagSet("test")
		set.AddInput("env:CONFIG_FILE")

		Convey("When I read the config", func() {
			input := test.NewInputConfig(set)

			Convey("Then I should have a valid configuration", func() {
				So(input.Kind(), ShouldEqual, "yml")

				var d map[string]any

				So(input.Decode(&d), ShouldBeNil)
			})
		})
	})

	Convey("Given I have configuration file", t, func() {
		Convey("When I read the config", func() {
			set := cmd.NewFlagSet("test")
			set.AddInput("file:" + test.Path("configs/config.yml"))

			input := test.NewInputConfig(set)

			Convey("Then I should have a valid configuration", func() {
				So(input.Kind(), ShouldEqual, "yml")

				var d map[string]any

				So(input.Decode(&d), ShouldBeNil)
			})
		})
	})

	Convey("Given I have configuration file", t, func() {
		home := os.UserHomeDir()
		path := filepath.Join(home, ".config", test.Name.String())

		err := os.MkdirAll(path, 0o777)
		So(err, ShouldBeNil)

		data, err := os.ReadFile(test.Path("configs/config.yml"))
		So(err, ShouldBeNil)

		err = os.WriteFile(filepath.Join(path, test.Name.String()+".yml"), data, 0o600)
		So(err, ShouldBeNil)

		Convey("When I read the config", func() {
			set := cmd.NewFlagSet("test")
			set.AddInput("")

			input := test.NewInputConfig(set)

			Convey("Then I should have a valid configuration", func() {
				So(input.Kind(), ShouldEqual, "yml")

				var d map[string]any

				So(input.Decode(&d), ShouldBeNil)
			})
		})

		err = os.RemoveAll(path)
		So(err, ShouldBeNil)
	})
}

//nolint:funlen
func TestWriteValidConfigFile(t *testing.T) {
	Convey("Given I have configuration file", t, func() {
		file := test.Path("configs/new_config.yml")

		So(os.WriteFile(file, strings.Bytes("environment: development"), 0o600), ShouldBeNil)
		t.Setenv("CONFIG_FILE", file)

		set := cmd.NewFlagSet("test")
		set.AddOutput("env:CONFIG_FILE")

		Convey("When I write the config", func() {
			input := test.NewOutputConfig(set)

			err := input.Write(strings.Bytes("test"), os.ModeAppend)
			So(err, ShouldBeNil)

			Convey("Then I should have a valid configuration", func() {
				d, err := os.ReadFile(file)
				So(err, ShouldBeNil)

				So(bytes.String(d), ShouldEqual, "test")
			})

			So(os.Remove(file), ShouldBeNil)
		})
	})

	Convey("Given I have configuration file", t, func() {
		file := test.Path("configs/new_config.yml")

		So(os.WriteFile(file, strings.Bytes("environment: development"), 0o600), ShouldBeNil)

		set := cmd.NewFlagSet("test")
		set.AddOutput("file:" + file)

		Convey("When I write the config", func() {
			output := test.NewOutputConfig(set)

			err := output.Write(strings.Bytes("test"), os.ModeAppend)
			So(err, ShouldBeNil)

			Convey("Then I should have a valid configuration", func() {
				d, err := os.ReadFile(file)
				So(err, ShouldBeNil)

				So(bytes.String(d), ShouldEqual, "test")
			})

			So(os.Remove(file), ShouldBeNil)
		})
	})

	Convey("Given I have configuration file", t, func() {
		home := os.UserHomeDir()
		path := filepath.Join(home, ".config", test.Name.String())

		err := os.MkdirAll(path, 0o777)
		So(err, ShouldBeNil)

		data, err := os.ReadFile(test.Path("configs/config.yml"))
		So(err, ShouldBeNil)

		file := filepath.Join(path, test.Name.String()+".yml")

		err = os.WriteFile(file, data, 0o600)
		So(err, ShouldBeNil)

		set := cmd.NewFlagSet("test")
		set.AddOutput("")

		Convey("When I write the config", func() {
			output := test.NewOutputConfig(set)

			err := output.Write(strings.Bytes("test"), os.ModeAppend)
			So(err, ShouldBeNil)

			Convey("Then I should have a valid configuration", func() {
				d, err := os.ReadFile(file)
				So(err, ShouldBeNil)

				So(bytes.String(d), ShouldEqual, "test")
			})

			So(os.Remove(file), ShouldBeNil)
		})

		err = os.RemoveAll(path)
		So(err, ShouldBeNil)
	})
}

func TestValidConfigEnv(t *testing.T) {
	Convey("Given I have configuration file", t, func() {
		t.Setenv("CONFIG_FILE", "yaml:CONFIG")
		t.Setenv("CONFIG", "ZW52aXJvbm1lbnQ6IGRldmVsb3BtZW50Cg==")

		set := cmd.NewFlagSet("test")
		set.AddInput("env:CONFIG_FILE")

		Convey("When I read the config", func() {
			input := test.NewInputConfig(set)

			Convey("Then I should have a valid configuration", func() {
				So(input.Kind(), ShouldEqual, "yaml")
			})
		})

		Convey("When I write the config", func() {
			output := test.NewInputConfig(set)

			err := output.Write(strings.Bytes("test"), os.ModeAppend)
			So(err, ShouldBeNil)

			Convey("Then I should have a valid configuration", func() {
				So(os.Getenv("CONFIG"), ShouldEqual, base64.Encode(strings.Bytes("test")))
			})
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
		t.Setenv("CONFIG_FILE", "../bob")

		set := cmd.NewFlagSet("test")
		set.AddInput("env:CONFIG_FILE")

		Convey("When I try to parse the configuration file", func() {
			input := test.NewInputConfig(set)

			Convey("Then I should have a valid configuration", func() {
				So(input.Kind(), ShouldBeEmpty)
			})
		})
	})
}

func TestInvalidKindConfig(t *testing.T) {
	Convey("When I try to parse the configuration file", t, func() {
		set := cmd.NewFlagSet("test")
		set.AddInput("test:test")

		input := test.NewInputConfig(set)

		Convey("Then I should have a valid configuration", func() {
			So(input.Kind(), ShouldEqual, "yaml")
		})
	})
}

func TestInvalidConfig(t *testing.T) {
	Convey("Given I have invalid kind configuration file", t, func() {
		t.Setenv("CONFIG_FILE", "config.go")

		set := cmd.NewFlagSet("test")
		set.AddInput("env:CONFIG_FILE")

		Convey("When I try to parse the configuration file", func() {
			input := test.NewInputConfig(set)

			Convey("Then I should have an error when decoding", func() {
				So(input.Decode(nil), ShouldBeError)
			})
		})
	})
}
