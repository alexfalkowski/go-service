package flags_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/flags"
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"github.com/spf13/cobra"
)

func TestFlags(t *testing.T) {
	Convey("Given I have a var", t, func() {
		name := os.ExecutableName()
		c := &cobra.Command{
			Use: name, Short: name, Long: name,
			Version: string(test.Version),
		}
		v := flags.Bool()

		Convey("When I add it to the command", func() {
			flags.BoolVar(c, v, "test", "y", true, "test")

			Convey("Then I should have it set to the default value", func() {
				So(*v, ShouldBeTrue)
			})
		})
	})
}
