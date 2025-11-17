package net_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/net"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDefaultAddress(t *testing.T) {
	Convey("When I try to get a default address", t, func() {
		address := net.DefaultAddress("9000")

		Convey("Then I should have a valid address", func() {
			So(address, ShouldEqual, "tcp://:9000")
		})
	})
}

func TestHost(t *testing.T) {
	Convey("When I try to get the host of an invalid address", t, func() {
		host := net.Host("none")

		Convey("Then I should just get the address back", func() {
			So(host, ShouldEqual, "none")
		})
	})
}

func TestNetworkAddress(t *testing.T) {
	Convey("When I try to get a valid network address", t, func() {
		network, address, ok := net.SplitNetworkAddress("tcp://localhost:9000")

		Convey("Then I should have a valid address", func() {
			So(ok, ShouldBeTrue)
			So(network, ShouldEqual, "tcp")
			So(address, ShouldEqual, "localhost:9000")
		})
	})

	Convey("When I try to get a invalid network address", t, func() {
		network, address, ok := net.SplitNetworkAddress("no:address")

		Convey("Then I should have an invalid address", func() {
			So(ok, ShouldBeFalse)
			So(network, ShouldEqual, "no:address")
			So(address, ShouldBeBlank)
		})
	})
}
