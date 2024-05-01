package grpc_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/security"
	"github.com/alexfalkowski/go-service/server"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/alexfalkowski/go-service/test"
	"github.com/alexfalkowski/go-service/transport/grpc"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
)

func TestServer(t *testing.T) {
	Convey("Given I have invalid creds", t, func() {
		c := &grpc.Config{
			Config: &server.Config{
				Enabled:  true,
				Security: &security.Config{Enabled: true, CertFile: "bob", KeyFile: "bob"},
			},
		}

		Convey("When I create a server", func() {
			_, err := grpc.NewServer(grpc.ServerParams{Config: c})

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})

	Convey("Given I have secure creds", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)

		c := test.NewInsecureTransportConfig()
		c.GRPC.Security = &security.Config{Enabled: true}

		Convey("When I create a server", func() {
			s := &test.Server{Lifecycle: lc, Logger: logger, Transport: c, Meter: metrics.NewNoopMeter()}
			s.Register()

			Convey("Then I should start the server", func() {
				lc.RequireStart()
			})
		})

		lc.RequireStop()
	})
}
