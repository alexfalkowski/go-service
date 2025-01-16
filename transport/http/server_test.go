package http_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/server"
	"github.com/alexfalkowski/go-service/test"
	"github.com/alexfalkowski/go-service/transport/http"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestInvalidServer(t *testing.T) {
	Convey("When I try to create a server with invalid tls configuration", t, func() {
		cfg := &http.Config{
			Config: &server.Config{
				Timeout: "5s",
				TLS:     test.NewTLSConfig("certs/client-cert.pem", "secrets/none"),
			},
		}
		p := http.ServerParams{
			Shutdowner: test.NewShutdowner(),
			Config:     cfg,
		}

		_, err := http.NewServer(p)

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})
}
