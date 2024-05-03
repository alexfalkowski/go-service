package ntp_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/time/ntp"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestService(t *testing.T) {
	Convey("Given I have NTP setup correctly", t, func() {
		c := &ntp.Config{Host: "0.beevik-ntp.pool.ntp.org"}
		n := ntp.NewService(c)

		Convey("When I get the time", func() {
			_, err := n.Time()

			Convey("I should not have an error", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When I query", func() {
			_, err := n.Query()

			Convey("I should not have an error", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given I have NTP setup incorrectly", t, func() {
		c := &ntp.Config{}
		n := ntp.NewService(c)

		Convey("When I get the time", func() {
			_, err := n.Time()

			Convey("I should not have an error", func() {
				So(err, ShouldBeError)
			})
		})

		Convey("When I query", func() {
			_, err := n.Query()

			Convey("I should not have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}
