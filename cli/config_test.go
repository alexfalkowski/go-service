package cli_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/cli/flag"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/os"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNoneConfig(t *testing.T) {
	for _, arg := range []string{"", "env:BOB"} {
		Convey("When I read the config", t, func() {
			set := flag.NewFlagSet("test")
			set.AddInput(arg)

			input := test.NewConfig(set)

			Convey("Then I should have a valid configuration", func() {
				So(input.Decode(nil), ShouldBeError)
			})
		})
	}
}

func TestReadValidConfigFile(t *testing.T) {
	Convey("Given I have configuration file", t, func() {
		Convey("When I read the config", func() {
			set := flag.NewFlagSet("test")
			set.AddInput(test.FilePath("configs/config.yml"))

			input := test.NewConfig(set)

			Convey("Then I should have a valid configuration", func() {
				var d map[string]any

				So(input.Decode(&d), ShouldBeNil)
			})
		})
	})

	Convey("Given I have configuration file", t, func() {
		home := os.UserHomeDir()
		path := test.FS.Join(home, ".config", test.Name.String())
		fs := test.FS

		err := fs.MkdirAll(path, 0o777)
		So(err, ShouldBeNil)

		data, err := fs.ReadFile(test.Path("configs/config.yml"))
		So(err, ShouldBeNil)

		err = fs.WriteFile(test.FS.Join(path, test.Name.String()+".yml"), data, 0o600)
		So(err, ShouldBeNil)

		Convey("When I read the config", func() {
			set := flag.NewFlagSet("test")
			set.AddInput("")

			input := test.NewConfig(set)

			Convey("Then I should have a valid configuration", func() {
				var d map[string]any

				So(input.Decode(&d), ShouldBeNil)
			})
		})

		err = fs.RemoveAll(path)
		So(err, ShouldBeNil)
	})
}

func TestValidEnvConfig(t *testing.T) {
	Convey("Given I have configuration file", t, func() {
		t.Setenv("CONFIG", "yaml:ZW52aXJvbm1lbnQ6IGRldmVsb3BtZW50Cg==")

		set := flag.NewFlagSet("test")
		set.AddInput("env:CONFIG")

		Convey("When I read the config", func() {
			input := test.NewConfig(set)

			Convey("Then I should have a valid configuration", func() {
				var d map[string]any

				So(input.Decode(&d), ShouldBeNil)
			})
		})
	})
}

func TestMissingConfig(t *testing.T) {
	Convey("When I read the config", t, func() {
		set := flag.NewFlagSet("test")
		set.AddInput("")

		input := test.NewConfig(set)

		Convey("Then I should have a valid configuration", func() {
			So(input, ShouldNotBeNil)
		})
	})
}

func TestNonExistentConfig(t *testing.T) {
	Convey("Given I have non existent configuration file", t, func() {
		set := flag.NewFlagSet("test")
		set.AddInput("file:../bob")

		Convey("When I try to parse the configuration file", func() {
			input := test.NewConfig(set)

			Convey("Then I should have a invalid configuration", func() {
				So(input.Decode(nil), ShouldBeError)
			})
		})
	})
}

func TestInvalidKindConfig(t *testing.T) {
	Convey("When I try to parse the configuration file", t, func() {
		set := flag.NewFlagSet("test")
		set.AddInput("test:test")

		input := test.NewConfig(set)

		Convey("Then I should have a invalid configuration", func() {
			So(input.Decode(nil), ShouldBeError)
		})
	})
}

func TestInvalidConfig(t *testing.T) {
	Convey("Given I have invalid kind configuration file", t, func() {
		set := flag.NewFlagSet("test")
		set.AddInput("file:config.go")

		Convey("When I try to parse the configuration file", func() {
			input := test.NewConfig(set)

			Convey("Then I should have an error when decoding", func() {
				So(input.Decode(nil), ShouldBeError)
			})
		})
	})
}
