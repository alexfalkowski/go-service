package grpc_test

import (
	"net/http"
	"testing"

	"github.com/alexfalkowski/go-service/crypto/tls"
	"github.com/alexfalkowski/go-service/server"
	"github.com/alexfalkowski/go-service/test"
	"github.com/alexfalkowski/go-service/transport/grpc"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.opentelemetry.io/otel/metric/noop"
	"go.uber.org/fx/fxtest"
)

func TestServer(t *testing.T) {
	Convey("Given I have invalid creds", t, func() {
		config := &grpc.Config{
			Config: &server.Config{
				TLS: &tls.Config{Cert: "bob", Key: "bob"},
			},
		}

		Convey("When I create a server", func() {
			_, err := grpc.NewServer(grpc.ServerParams{Config: config})

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})

	Convey("Given I have secure creds", t, func() {
		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)

		c := test.NewInsecureTransportConfig()
		c.GRPC.TLS = &tls.Config{}

		Convey("When I create a server", func() {
			s := &test.Server{Lifecycle: lc, Logger: logger, Transport: c, Meter: noop.Meter{}, Mux: mux}
			s.Register()

			Convey("Then I should start the server", func() {
				lc.RequireStart()
			})
		})

		lc.RequireStop()
	})
}
