package http_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/internal/test"
	"github.com/alexfalkowski/go-service/server"
	"github.com/alexfalkowski/go-service/transport/http"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	http.Register(test.FS)
}

func TestInvalidServer(t *testing.T) {
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
			FS:         test.FS,
		}

		_, err := http.NewServer(params)

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})
}
