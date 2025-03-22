package time_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/time"
	. "github.com/smartystreets/goconvey/convey"
)

func TestInvalid(t *testing.T) {
	configs := []*time.Config{nil, {}}

	for _, config := range configs {
		Convey("When I try to create a network", t, func() {
			net, err := time.NewNetwork(config)
			So(err, ShouldBeNil)

			Convey("Then I should not have a network", func() {
				So(net, ShouldBeNil)
			})
		})
	}

	Convey("When I try to create a network", t, func() {
		_, err := time.NewNetwork(&time.Config{Kind: "invalid"})

		Convey("Then I should not have an error", func() {
			So(err, ShouldBeError)
		})
	})
}

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

func TestNTS(t *testing.T) {
	Convey("Given I have NTS setup correctly", t, func() {
		c := &time.Config{Kind: "nts", Address: "time.cloudflare.com"}

		n, err := time.NewNetwork(c)
		So(err, ShouldBeNil)

		Convey("When I get the time", func() {
			_, err := n.Now()

			Convey("I should not have an error", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given I have NTS setup incorrectly", t, func() {
		c := &time.Config{Kind: "nts"}

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
