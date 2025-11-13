package grpc_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/server"
	"github.com/alexfalkowski/go-service/v2/transport/grpc"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func init() {
	grpc.Register(test.Name, test.FS)
}

func TestServer(t *testing.T) {
	Convey("Given I have secure credentials", t, func() {
		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc, test.NewTextLoggerConfig())
		meter := test.NewPrometheusMeter(lc)

		c := test.NewInsecureTransportConfig()
		c.GRPC.TLS = &tls.Config{}

		Convey("When I create a server", func() {
			s := &test.Server{Lifecycle: lc, Logger: logger, TransportConfig: c, Meter: meter, Mux: mux}
			s.Register()

			Convey("Then I should start the server", func() {
				lc.RequireStart()
			})
		})

		lc.RequireStop()
	})
}

func TestInvalidServer(t *testing.T) {
	Convey("When I try to create a server with invalid tls configuration", t, func() {
		cfg := &grpc.Config{
			Config: &server.Config{
				Timeout: "5s",
				TLS:     test.NewTLSConfig("certs/client-cert.pem", "secrets/none"),
			},
		}
		params := grpc.ServerParams{
			Shutdowner: test.NewShutdowner(),
			Config:     cfg,
		}

		_, err := grpc.NewServer(params)

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})
}
