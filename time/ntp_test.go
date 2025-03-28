package time_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/time"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNTP(t *testing.T) {
	Convey("Given I have NTP setup correctly", t, func() {
		c := &time.Config{Kind: "ntp", Address: "0.beevik-ntp.pool.ntp.org"}

		n, err := time.NewNetwork(c)
		So(err, ShouldBeNil)

		Convey("When I get the time", func() {
			_, err := n.Now()

			Convey("I should not have an error", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given I have NTP setup incorrectly", t, func() {
		c := &time.Config{Kind: "ntp"}

		n, err := time.NewNetwork(c)
		So(err, ShouldBeNil)

		Convey("When I get the time", func() {
			_, err := n.Now()

			Convey("I should not have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}
