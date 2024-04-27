package meta_test

import (
	"context"
	"net"
	"testing"

	"github.com/alexfalkowski/go-service/transport/grpc/meta"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func TestIPAddr(t *testing.T) {
	Convey("Given I have a x-forwarded-for header", t, func() {
		md := metadata.MD{"x-forwarded-for": []string{"test"}}

		Convey("When I get the ip address", func() {
			ip := meta.IPAddr(context.Background(), md)

			Convey("Then I should have an ip address", func() {
				So(ip, ShouldEqual, "test")
			})
		})
	})

	Convey("Given I have no context", t, func() {
		md := metadata.MD{}

		Convey("When I get the ip address", func() {
			ip := meta.IPAddr(context.Background(), md)

			Convey("Then I should have an ip address", func() {
				So(ip, ShouldBeBlank)
			})
		})
	})

	Convey("Given I have a peer with no port", t, func() {
		ip, err := net.ResolveIPAddr("ip", "203.0.113.0")
		So(err, ShouldBeNil)

		ctx := peer.NewContext(context.Background(), &peer.Peer{Addr: ip})
		md := metadata.MD{}

		Convey("When I get the ip address", func() {
			ip := meta.IPAddr(ctx, md)

			Convey("Then I should have an ip address", func() {
				So(ip, ShouldEqual, "203.0.113.0")
			})
		})
	})

	Convey("Given I have a peer with port", t, func() {
		ip, err := net.ResolveUDPAddr("udp", "203.0.113.0:53")
		So(err, ShouldBeNil)

		ctx := peer.NewContext(context.Background(), &peer.Peer{Addr: ip})
		md := metadata.MD{}

		Convey("When I get the ip address", func() {
			ip := meta.IPAddr(ctx, md)

			Convey("Then I should have an ip address", func() {
				So(ip, ShouldEqual, "203.0.113.0")
			})
		})
	})
}
