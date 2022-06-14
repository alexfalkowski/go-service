package time_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/time"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRandomWaitTime(t *testing.T) {
	Convey("Given I time setup", t, func() {
		Convey("When I try to get random time", func() {
			t := time.RandomWaitTime()

			Convey("Then I should have a cached item", func() {
				So(t, ShouldBeBetween, 0, time.Timeout)
			})
		})
	})
}
