package debug_test

import (
	"fmt"
	"testing"

	"github.com/alexfalkowski/go-service/debug"
	"github.com/alexfalkowski/go-service/marshaller"
	"github.com/alexfalkowski/go-service/test"
	"github.com/alexfalkowski/go-service/transport"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
)

func TestDebug(t *testing.T) {
	Convey("When I have a all the servers", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		m := test.NewMeter(lc)
		p := debug.ServerParams{
			Shutdowner: test.NewShutdowner(),
			Config:     test.NewDebugConfig(),
			Logger:     logger,
		}

		server, err := debug.NewServer(p)
		So(err, ShouldBeNil)

		debug.RegisterPprof(server)
		debug.RegisterFgprof(server)
		debug.RegisterPsutil(server, marshaller.NewJSON())
		debug.RegisterStatsviz(server)

		transport.Register(transport.RegisterParams{Lifecycle: lc, Servers: []transport.Server{server}})
		lc.RequireStart()

		Convey("Then all the debug URLs are valid", func() {
			port := p.Config.Port
			client := test.NewHTTPClient(lc, logger, test.NewOTLPTracerConfig(), test.NewInsecureTransportConfig(), m)
			urls := []string{
				url(port, "debug/statsviz"),
				url(port, "debug/pprof/"),
				url(port, "debug/pprof/cmdline"),
				url(port, "debug/pprof/symbol"),
				url(port, "debug/pprof/trace"),
				url(port, "debug/fgprof?seconds=1"),
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
