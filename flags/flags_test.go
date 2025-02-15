package flags_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/flags"
	"github.com/alexfalkowski/go-service/types/ptr"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestFlags(t *testing.T) {
	Convey("Then I should unset flags", t, func() {
		So(flags.IsBoolSet(nil), ShouldBeFalse)
		So(flags.IsStringSet(nil), ShouldBeFalse)
	})

	Convey("Then I should set flags", t, func() {
		So(flags.IsBoolSet(ptr.Value(true)), ShouldBeTrue)
		So(flags.IsStringSet(ptr.Value("yes")), ShouldBeTrue)
	})
}
