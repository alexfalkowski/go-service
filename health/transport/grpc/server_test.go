package grpc_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/health/transport/grpc"
	. "github.com/smartystreets/goconvey/convey"
)

func TestServer(t *testing.T) {
	Convey("When I try to create a server with no observer", t, func() {
		server := grpc.NewServer(nil)

		Convey("Then I should have a nil server", func() {
			So(server, ShouldBeNil)
		})
	})
}
