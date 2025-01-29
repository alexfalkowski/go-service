package proxy_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/proxy"
	"github.com/alexfalkowski/go-service/server"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestInvalidServer(t *testing.T) {
	Convey("When I try to create a server with invalid tls configuration", t, func() {
		cfg := &proxy.Config{
			Config: &server.Config{
				Timeout: "5s",
				TLS:     test.NewTLSConfig("certs/client-cert.pem", "secrets/none"),
			},
		}
		p := proxy.ServerParams{
			Shutdowner: test.NewShutdowner(),
			Config:     cfg,
		}

		_, err := proxy.NewServer(p)

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})
}
