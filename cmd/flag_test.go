package cmd_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"github.com/spf13/cobra"
)

type tuple [2]string

func TestSplitFlag(t *testing.T) {
	tuples := []tuple{{"file", "file.yaml"}, {"file", "file.test.yaml"}, {"file", "test/.config/existing.client.yaml"}}

	for _, tu := range tuples {
		Convey("Given I have a flag", t, func() {
			Convey("When I split the flag", func() {
				k, l := cmd.SplitFlag(tu[0] + ":" + tu[1])

				Convey("Then I should have a valid split", func() {
					So(k, ShouldEqual, tu[0])
					So(l, ShouldEqual, tu[1])
				})
			})
		})
	}
}

func TestReadWriter(t *testing.T) {
	tuples := []tuple{{"file", "file.yaml"}, {"file", "file.test.yaml"}, {"file", "test/.config/existing.client.yaml"}}

	for _, tu := range tuples {
		Convey("Given I have a flag", t, func() {
			Convey("When I get a read writer", func() {
				rw := cmd.NewReadWriter(tu[0], tu[1])

				Convey("Then tI should have a valid split", func() {
					So(rw.Kind(), ShouldEqual, "yaml")
				})
			})
		})
	}
}

func TestVar(t *testing.T) {
	Convey("Given I have a var", t, func() {
		name := os.ExecutableName()
		c := &cobra.Command{
			Use: name, Short: name, Long: name,
			Version: string(test.Version),
		}
		v := cmd.Bool()

		Convey("When I add it to the command", func() {
			cmd.BoolVar(c, v, "test", "y", true, "test")

			Convey("Then I should have it set to the default value", func() {
				So(*v, ShouldBeTrue)
			})
		})
	})
}
