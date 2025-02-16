package flags_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/flags"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestFlags(t *testing.T) {
	Convey("Then I should have a flag set", t, func() {
		So(flags.NewFlagSet("test"), ShouldNotBeNil)
	})

	Convey("Then I should unset flags", t, func() {
		So(flags.IsBoolSet(nil), ShouldBeFalse)
		So(flags.IsStringSet(nil), ShouldBeFalse)
	})

	Convey("Then I should unset flags", t, func() {
		So(flags.IsBoolSet(flags.Bool()), ShouldBeFalse)
		So(flags.IsStringSet(flags.String()), ShouldBeFalse)
	})
}
