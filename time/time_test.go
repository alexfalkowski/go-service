package time_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/time"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestMustParseDuration(t *testing.T) {
	t.Parallel()

	Convey("When I try to parse duration", t, func() {
		f := func() { time.MustParseDuration("test") }

		Convey("Then I should have an invalid duration", func() {
			So(f, ShouldPanic)
		})
	})
}
