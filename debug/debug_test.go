package debug_test

import (
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/debug"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
)

func TestUDebug(t *testing.T) {
	Convey("When I have a all the servers", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)

		debug.Register(lc, test.Environment, logger)

		cfg := test.NewTransportConfig()

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		lc.RequireStart()
		time.Sleep(1 * time.Second)

		Convey("Then all the debug URLs are valid", func() {
			client := test.NewHTTPClient(lc, logger, test.NewTracerConfig(), cfg, m)
			urls := []string{
				"http://localhost:6060/debug/statsviz",
				"http://localhost:6060/debug/pprof/",
				"http://localhost:6060/debug/pprof/cmdline",
				"http://localhost:6060/debug/pprof/symbol",
				"http://localhost:6060/debug/pprof/trace",
				"http://localhost:6060/debug/psutil",
			}

			for _, u := range urls {
				r, err := client.Get(u)
				So(err, ShouldBeNil)

				defer r.Body.Close()

				So(r.StatusCode, ShouldEqual, 200)
			}

			lc.RequireStop()
		})
	})
}
