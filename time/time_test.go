package time_test

import (
	"testing"
	"time"

	stime "github.com/alexfalkowski/go-service/time"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRandomWaitTime(t *testing.T) {
	Convey("Given I time setup", t, func() {
		Convey("When I try to get random time", func() {
			t := stime.RandomWaitTime()

			Convey("Then I should have a cached item", func() {
				So(t, ShouldBeBetween, 0, 15*time.Second)
			})
		})
	})
}
