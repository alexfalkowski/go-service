package health_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/transport/grpc/health"
	. "github.com/smartystreets/goconvey/convey"
)

func TestServer(t *testing.T) {
	Convey("When I try to create a server with no observer", t, func() {
		server := health.NewServer(health.ServerParams{})

		Convey("Then I should have a nil server", func() {
			So(server, ShouldBeNil)
		})
	})
}
