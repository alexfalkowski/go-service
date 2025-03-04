package time_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/time"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSystem(t *testing.T) {
	configs := []*time.Config{{}, {Kind: "none"}}
	for _, config := range configs {
		Convey("Given I have a config", t, func() {
			n := time.NewNetwork(config)

			Convey("When I get the time", func() {
				_, err := n.Now()

				Convey("I should not have an error", func() {
					So(err, ShouldBeNil)
				})
			})
		})
	}
}

func TestNTP(t *testing.T) {
	Convey("Given I have NTP setup correctly", t, func() {
		c := &time.Config{Kind: "ntp", Address: "0.beevik-ntp.pool.ntp.org"}
		n := time.NewNetwork(c)

		Convey("When I get the time", func() {
			_, err := n.Now()

			Convey("I should not have an error", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given I have NTP setup incorrectly", t, func() {
		c := &time.Config{Kind: "ntp"}
		n := time.NewNetwork(c)

		Convey("When I get the time", func() {
			_, err := n.Now()

			Convey("I should not have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}

func TestNTS(t *testing.T) {
	Convey("Given I have NTS setup correctly", t, func() {
		c := &time.Config{Kind: "nts", Address: "time.cloudflare.com"}
		n := time.NewNetwork(c)

		Convey("When I get the time", func() {
			_, err := n.Now()

			Convey("I should not have an error", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given I have NTS setup incorrectly", t, func() {
		c := &time.Config{Kind: "nts"}
		n := time.NewNetwork(c)

		Convey("When I get the time", func() {
			_, err := n.Now()

			Convey("I should not have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}
