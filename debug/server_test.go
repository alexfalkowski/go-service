package debug_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/debug"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/server"
	. "github.com/smartystreets/goconvey/convey"
)

var paths = []string{
	"debug/statsviz",
	"debug/pprof/",
	"debug/pprof/cmdline",
	"debug/pprof/symbol",
	"debug/pprof/trace",
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
				url := world.NamedDebugURL("http", path)

				res, err := world.ResponseWithNoBody(t.Context(), url, http.MethodGet, header)
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
				url := world.NamedDebugURL("https", path)

				res, err := world.ResponseWithNoBody(t.Context(), url, http.MethodGet, header)
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
		params := debug.ServerParams{
			Shutdowner: test.NewShutdowner(),
			Config:     cfg,
			FS:         test.FS,
		}

		_, err := debug.NewServer(params)

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})
}
