package debug_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/debug"
	"github.com/alexfalkowski/go-service/internal/test"
	"github.com/alexfalkowski/go-service/server"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestInvalidServer(t *testing.T) {
	Convey("When I try to create a server with invalid tls configuration", t, func() {
		cfg := &debug.Config{
			Config: &server.Config{
				Timeout: "5s",
				TLS:     test.NewTLSConfig("certs/client-cert.pem", "secrets/none"),
			},
		}
		p := debug.ServerParams{
			Shutdowner: test.NewShutdowner(),
			Config:     cfg,
		}

		_, err := debug.NewServer(p)

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})
}
