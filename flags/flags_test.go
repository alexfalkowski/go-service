package flags_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/flags"
	"github.com/alexfalkowski/go-service/internal/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"github.com/spf13/cobra"
)

func TestFlags(t *testing.T) {
	Convey("Given I have a var", t, func() {
		name := string(test.Name)
		command := &cobra.Command{
			Use: name, Short: name, Long: name,
			Version: string(test.Version),
		}
		boolFlag := flags.Bool()
		stringFlag := flags.String()

		Convey("When I add it to the command", func() {
			flags.BoolVar(command, boolFlag, "bool", "y", false, "test")
			flags.StringVar(command, stringFlag, "string", "z", "", "test")

			Convey("Then I should unset flags", func() {
				So(flags.IsBoolSet(boolFlag), ShouldBeFalse)
				So(flags.IsStringSet(stringFlag), ShouldBeFalse)
			})
		})
	})
}
