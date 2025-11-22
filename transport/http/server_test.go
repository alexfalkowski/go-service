package http_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/server"
	"github.com/alexfalkowski/go-service/v2/transport/http"
	. "github.com/smartystreets/goconvey/convey"
)

func TestInvalidServer(t *testing.T) {
	http.Register(test.FS)

	Convey("When I try to create a server with invalid tls configuration", t, func() {
		cfg := &http.Config{
			Config: &server.Config{
				Timeout: "5s",
				TLS:     test.NewTLSConfig("certs/client-cert.pem", "secrets/none"),
			},
		}
		params := http.ServerParams{
			Shutdowner: test.NewShutdowner(),
			Config:     cfg,
		}

		_, err := http.NewServer(params)

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})
}
