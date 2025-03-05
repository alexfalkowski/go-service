package net_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/net"
	. "github.com/smartystreets/goconvey/convey"
)

func TestHost(t *testing.T) {
	Convey("When I try to get the host of an invalid address", t, func() {
		host := net.Host("none")

		Convey("Then I should just get the address back", func() {
			So(host, ShouldEqual, "none")
		})
	})
}
