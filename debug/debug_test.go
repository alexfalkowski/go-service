package debug_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/alexfalkowski/go-service/debug"
	"github.com/alexfalkowski/go-service/net/http/content"
	"github.com/alexfalkowski/go-service/test"
	"github.com/alexfalkowski/go-service/transport"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
)

func init() {
	content.Register(test.Pool)
}

func TestInsecureDebug(t *testing.T) {
	Convey("When I have a all the servers", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		m := test.NewOTLPMeter(lc)
		p := debug.ServerParams{
			Shutdowner: test.NewShutdowner(),
			Config:     test.NewInsecureDebugConfig(),
			Logger:     logger,
		}

		server, err := debug.NewServer(p)
		So(err, ShouldBeNil)

		debug.RegisterPprof(server)
		debug.RegisterFgprof(server)
		debug.RegisterPsutil(server, test.Content)

		err = debug.RegisterStatsviz(server)
		So(err, ShouldBeNil)

		transport.Register(transport.RegisterParams{Lifecycle: lc, Servers: []transport.Server{server}})
		lc.RequireStart()

		Convey("Then all the debug URLs are valid", func() {
			cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: test.NewOTLPTracerConfig(), Transport: test.NewInsecureTransportConfig(), Meter: m}
			client := cl.NewHTTP()
			urls := []string{
				url("http", p.Config.Address, "debug/statsviz"),
				url("http", p.Config.Address, "debug/pprof/"),
				url("http", p.Config.Address, "debug/pprof/cmdline"),
				url("http", p.Config.Address, "debug/pprof/symbol"),
				url("http", p.Config.Address, "debug/pprof/trace"),
				url("http", p.Config.Address, "debug/fgprof?seconds=1"),
				url("http", p.Config.Address, "debug/psutil"),
			}

			for _, u := range urls {
				req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, u, http.NoBody)
				So(err, ShouldBeNil)

				req.Header.Set("Content-Type", "application/json")

				res, err := client.Do(req)
				So(err, ShouldBeNil)

				defer res.Body.Close()

				So(res.StatusCode, ShouldEqual, 200)
			}

			lc.RequireStop()
		})
	})
}

func TestSecureDebug(t *testing.T) {
	Convey("When I have a all the servers", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		m := test.NewOTLPMeter(lc)
		p := debug.ServerParams{
			Shutdowner: test.NewShutdowner(),
			Config:     test.NewSecureDebugConfig(),
			Logger:     logger,
		}

		server, err := debug.NewServer(p)
		So(err, ShouldBeNil)

		debug.RegisterPprof(server)
		debug.RegisterFgprof(server)
		debug.RegisterPsutil(server, test.Content)

		err = debug.RegisterStatsviz(server)
		So(err, ShouldBeNil)

		transport.Register(transport.RegisterParams{Lifecycle: lc, Servers: []transport.Server{server}})
		lc.RequireStart()

		Convey("Then all the debug URLs are valid", func() {
			cl := &test.Client{
				Lifecycle: lc, Logger: logger, Tracer: test.NewOTLPTracerConfig(),
				Transport: test.NewSecureTransportConfig(),
				TLS:       test.NewTLSClientConfig(),
				Meter:     m,
			}

			client := cl.NewHTTP()
			urls := []string{
				url("https", p.Config.Address, "debug/statsviz"),
				url("https", p.Config.Address, "debug/pprof/"),
				url("https", p.Config.Address, "debug/pprof/cmdline"),
				url("https", p.Config.Address, "debug/pprof/symbol"),
				url("https", p.Config.Address, "debug/pprof/trace"),
				url("https", p.Config.Address, "debug/fgprof?seconds=1"),
				url("https", p.Config.Address, "debug/psutil"),
			}

			for _, u := range urls {
				req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, u, http.NoBody)
				So(err, ShouldBeNil)

				req.Header.Set("Content-Type", "application/json")

				res, err := client.Do(req)
				So(err, ShouldBeNil)

				defer res.Body.Close()

				So(res.StatusCode, ShouldEqual, 200)
			}

			lc.RequireStop()
		})
	})
}

func url(scheme, address, path string) string {
	return fmt.Sprintf("%s://%s/%s", scheme, address, path)
}
