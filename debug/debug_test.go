package debug_test

import (
	"net/http"
	"testing"

	"github.com/alexfalkowski/go-service/internal/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestInsecureDebug(t *testing.T) {
	paths := []string{
		"debug/statsviz",
		"debug/pprof/",
		"debug/pprof/cmdline",
		"debug/pprof/symbol",
		"debug/pprof/trace",
		"debug/fgprof?seconds=1",
		"debug/psutil",
	}

	for _, path := range paths {
		Convey("When I have a all the servers", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldDebug())
			world.Register()
			world.RequireStart()

			Convey("Then all the debug URLs are valid", func() {
				header := http.Header{}

				res, err := world.ResponseWithNoBody(t.Context(), "http", world.DebugHost(), http.MethodGet, path, header, http.NoBody)
				So(err, ShouldBeNil)

				So(res.StatusCode, ShouldEqual, 200)
			})

			world.RequireStop()
		})
	}
}

func TestSecureDebug(t *testing.T) {
	paths := []string{
		"debug/statsviz",
		"debug/pprof/",
		"debug/pprof/cmdline",
		"debug/pprof/symbol",
		"debug/pprof/trace",
		"debug/fgprof?seconds=1",
		"debug/psutil",
	}

	for _, path := range paths {
		Convey("When I have a all the servers", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldSecure(), test.WithWorldDebug())
			world.Register()
			world.RequireStart()

			Convey("Then all the debug URLs are valid", func() {
				header := http.Header{}

				res, err := world.ResponseWithNoBody(t.Context(), "https", world.DebugHost(), http.MethodGet, path, header, http.NoBody)
				So(err, ShouldBeNil)

				So(res.StatusCode, ShouldEqual, 200)
			})

			world.RequireStop()
		})
	}
}
