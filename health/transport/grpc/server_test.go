package grpc_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/health/transport/grpc"
	. "github.com/smartystreets/goconvey/convey"
)

func TestServer(t *testing.T) {
	Convey("When I try to create a server with no observer", t, func() {
		server := grpc.NewServer(grpc.ServerParams{})

		Convey("Then I should have a nil server", func() {
			So(server, ShouldBeNil)
		})
	})
}
