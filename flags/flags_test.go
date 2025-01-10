//nolint:varnamelen
package flags_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/flags"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"github.com/spf13/cobra"
)

func TestFlags(t *testing.T) {
	Convey("Given I have a var", t, func() {
		name := string(test.Name)
		c := &cobra.Command{
			Use: name, Short: name, Long: name,
			Version: string(test.Version),
		}
		b := flags.Bool()
		s := flags.String()

		Convey("When I add it to the command", func() {
			flags.BoolVar(c, b, "bool", "y", false, "test")
			flags.StringVar(c, s, "string", "z", "", "test")

			Convey("Then I should unset flags", func() {
				So(flags.IsBoolSet(b), ShouldBeFalse)
				So(flags.IsStringSet(s), ShouldBeFalse)
			})
		})
	})
}
