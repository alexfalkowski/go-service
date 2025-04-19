package debug_test

import (
	"net/http"
	"testing"

	"github.com/alexfalkowski/go-service/debug"
	"github.com/alexfalkowski/go-service/internal/test"
	"github.com/alexfalkowski/go-service/server"
	. "github.com/smartystreets/goconvey/convey"
)

var paths = []string{
	"debug/statsviz",
	"debug/pprof/",
	"debug/pprof/cmdline",
	"debug/pprof/symbol",
	"debug/pprof/trace",
	"debug/fgprof?seconds=1",
	"debug/psutil",
}

func TestInsecureDebug(t *testing.T) {
	for _, path := range paths {
		Convey("When I have a all the servers", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldDebug())
			world.Register()
			world.RequireStart()

			Convey("Then all the debug URLs are valid", func() {
				header := http.Header{}

				res, err := world.ResponseWithNoBody(t.Context(), "http", world.InsecureDebugHost(), http.MethodGet, path, header, http.NoBody)
				So(err, ShouldBeNil)

				So(res.StatusCode, ShouldEqual, 200)
			})

			world.RequireStop()
		})
	}
}

func TestSecureDebug(t *testing.T) {
	for _, path := range paths {
		Convey("When I have a all the servers", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldSecure(), test.WithWorldDebug())
			world.Register()
			world.RequireStart()

			Convey("Then all the debug URLs are valid", func() {
				header := http.Header{}

				res, err := world.ResponseWithNoBody(t.Context(), "https", world.SecureDebugHost(), http.MethodGet, path, header, http.NoBody)
				So(err, ShouldBeNil)

				So(res.StatusCode, ShouldEqual, 200)
			})

			world.RequireStop()
		})
	}
}

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
