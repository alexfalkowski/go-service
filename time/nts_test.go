package time_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/time"
	. "github.com/smartystreets/goconvey/convey"
)

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
