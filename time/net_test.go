package time_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/time"
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
