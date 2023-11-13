package debug_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/debug"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
)

func TestDebug(t *testing.T) {
	Convey("When I have a all the servers", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)

		p := debug.RegisterParams{
			Lifecycle: lc,
			Config:    test.NewDebugConfig(),
			Env:       test.Environment,
			Logger:    logger,
		}

		port := p.Config.Port

		debug.Register(p)

		cfg := test.NewTransportConfig()

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		lc.RequireStart()
		time.Sleep(1 * time.Second)

		Convey("Then all the debug URLs are valid", func() {
			client := test.NewHTTPClient(lc, logger, test.NewTracerConfig(), cfg, m)
			urls := []string{
				url(port, "debug/statsviz"),
				url(port, "debug/pprof/"),
				url(port, "debug/pprof/cmdline"),
				url(port, "debug/pprof/symbol"),
				url(port, "debug/pprof/trace"),
				url(port, "debug/psutil"),
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

func url(port, path string) string {
	return fmt.Sprintf("http://localhost:%s/%s", port, path)
}
