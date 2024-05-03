package nts_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/time/nts"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestService(t *testing.T) {
	Convey("Given I have NTP setup correctly", t, func() {
		c := &nts.Config{Host: "time.cloudflare.com"}
		n, err := nts.NewService(c)
		So(err, ShouldBeNil)

		Convey("When I query", func() {
			_, err := n.Query()

			Convey("I should not have an error", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("When I create an invalid service", t, func() {
		c := &nts.Config{}
		_, err := nts.NewService(c)

		Convey("I should not have an error", func() {
			So(err, ShouldBeError)
		})
	})
}
