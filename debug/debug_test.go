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

func TestInsecureDebug(t *testing.T) {
	Convey("When I have a all the servers", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		m := test.NewOTLPMeter(lc)
		p := debug.ServerParams{
			Shutdowner: test.NewShutdowner(),
			Mux:        debug.NewServeMux(),
			Config:     test.NewInsecureDebugConfig(),
			Logger:     logger,
		}

		server, err := debug.NewServer(p)
		So(err, ShouldBeNil)

		debug.RegisterPprof(p.Mux)
		debug.RegisterFgprof(p.Mux)
		debug.RegisterPsutil(p.Mux, marshaller.NewJSON())
		debug.RegisterStatsviz(p.Mux)

		transport.Register(transport.RegisterParams{Lifecycle: lc, Servers: []transport.Server{server}})
		lc.RequireStart()

		Convey("Then all the debug URLs are valid", func() {
			port := p.Config.Port
			cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: test.NewOTLPTracerConfig(), Transport: test.NewInsecureTransportConfig(), Meter: m}
			client := cl.NewHTTP()
			urls := []string{
				url("http", port, "debug/statsviz"),
				url("http", port, "debug/pprof/"),
				url("http", port, "debug/pprof/cmdline"),
				url("http", port, "debug/pprof/symbol"),
				url("http", port, "debug/pprof/trace"),
				url("http", port, "debug/fgprof?seconds=1"),
				url("http", port, "debug/psutil"),
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

func TestSecureDebug(t *testing.T) {
	Convey("When I have a all the servers", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		m := test.NewOTLPMeter(lc)
		p := debug.ServerParams{
			Shutdowner: test.NewShutdowner(),
			Mux:        debug.NewServeMux(),
			Config:     test.NewSecureDebugConfig(),
			Logger:     logger,
		}

		server, err := debug.NewServer(p)
		So(err, ShouldBeNil)

		debug.RegisterPprof(p.Mux)
		debug.RegisterFgprof(p.Mux)
		debug.RegisterPsutil(p.Mux, marshaller.NewJSON())
		debug.RegisterStatsviz(p.Mux)

		transport.Register(transport.RegisterParams{Lifecycle: lc, Servers: []transport.Server{server}})
		lc.RequireStart()

		Convey("Then all the debug URLs are valid", func() {
			port := p.Config.Port
			cl := &test.Client{
				Lifecycle: lc, Logger: logger, Tracer: test.NewOTLPTracerConfig(),
				Transport: test.NewSecureTransportConfig(),
				TLS:       test.NewTLSClientConfig(),
				Meter:     m,
			}

			client := cl.NewHTTP()
			urls := []string{
				url("https", port, "debug/statsviz"),
				url("https", port, "debug/pprof/"),
				url("https", port, "debug/pprof/cmdline"),
				url("https", port, "debug/pprof/symbol"),
				url("https", port, "debug/pprof/trace"),
				url("https", port, "debug/fgprof?seconds=1"),
				url("https", port, "debug/psutil"),
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

func url(scheme, port, path string) string {
	return fmt.Sprintf("%s://localhost:%s/%s", scheme, port, path)
}
