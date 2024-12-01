package http_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/runtime"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/test"
	"github.com/alexfalkowski/go-service/transport"
	th "github.com/alexfalkowski/go-service/transport/http"
	tm "github.com/alexfalkowski/go-service/transport/meta"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

func init() {
	tm.RegisterKeys()
}

func TestSecure(t *testing.T) {
	Convey("Given I a secure client", t, func() {
		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		tc := test.NewOTLPTracerConfig()
		m := test.NewPrometheusMeter(lc)
		cfg := test.NewSecureTransportConfig()

		s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Mux: mux}
		s.Register()

		cl := &test.Client{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m,
			TLS: test.NewTLSClientConfig(),
		}

		lc.RequireStart()

		Convey("When I query github", func() {
			client := cl.NewHTTP()

			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://github.com/alexfalkowski", http.NoBody)
			So(err, ShouldBeNil)

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			Convey("Then I should have valid response", func() {
				So(resp.StatusCode, ShouldEqual, 200)
			})
		})

		lc.RequireStop()
	})
}

func BenchmarkDefaultHTTP(b *testing.B) {
	b.ReportAllocs()

	mux := http.NewServeMux()
	p := test.Port()

	mux.HandleFunc("GET /hello", func(_ http.ResponseWriter, _ *http.Request) {})

	server := &http.Server{
		Handler:           mux,
		Addr:              ":" + p,
		ReadHeaderTimeout: time.Second,
	}
	defer server.Close()

	//nolint:errcheck
	go server.ListenAndServe()

	client := &http.Client{Transport: http.DefaultTransport}
	url := fmt.Sprintf("http://localhost:%s/hello", p)

	b.ResetTimer()

	b.Run("std", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
			runtime.Must(err)

			_, err = client.Do(req)
			runtime.Must(err)
		}
	})

	b.StopTimer()
}

func BenchmarkHTTP(b *testing.B) {
	b.ReportAllocs()

	mux := http.NewServeMux()
	lc := fxtest.NewLifecycle(b)
	cfg := test.NewInsecureTransportConfig()

	h, err := th.NewServer(th.ServerParams{
		Shutdowner: test.NewShutdowner(), Mux: mux,
		Config:    cfg.HTTP,
		UserAgent: test.UserAgent, Version: test.Version,
	})
	runtime.Must(err)

	transport.Register(transport.RegisterParams{Lifecycle: lc, Servers: []transport.Server{h}})

	client := &http.Client{Transport: http.DefaultTransport}
	url := fmt.Sprintf("http://%s/hello", cfg.HTTP.Address)

	lc.RequireStart()
	b.ResetTimer()

	b.Run("none", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
			runtime.Must(err)

			_, err = client.Do(req)
			runtime.Must(err)
		}
	})

	b.StopTimer()
	lc.RequireStop()
}

func BenchmarkLogHTTP(b *testing.B) {
	b.ReportAllocs()

	mux := http.NewServeMux()
	lc := fxtest.NewLifecycle(b)
	logger := zap.NewNop()
	cfg := test.NewInsecureTransportConfig()

	h, err := th.NewServer(th.ServerParams{
		Shutdowner: test.NewShutdowner(), Mux: mux,
		Config: cfg.HTTP, Logger: logger,
		UserAgent: test.UserAgent, Version: test.Version,
	})
	runtime.Must(err)

	transport.Register(transport.RegisterParams{Lifecycle: lc, Servers: []transport.Server{h}})

	client := &http.Client{Transport: http.DefaultTransport}
	url := fmt.Sprintf("http://%s/hello", cfg.HTTP.Address)

	lc.RequireStart()
	b.ResetTimer()

	b.Run("log", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
			runtime.Must(err)

			_, err = client.Do(req)
			runtime.Must(err)
		}
	})

	b.StopTimer()
	lc.RequireStop()
}

func BenchmarkTraceHTTP(b *testing.B) {
	b.ReportAllocs()

	mux := http.NewServeMux()
	lc := fxtest.NewLifecycle(b)

	tc := test.NewOTLPTracerConfig()
	logger := zap.NewNop()

	tracer, err := tracer.NewTracer(lc, test.Environment, test.Version, test.Name, tc, logger)
	runtime.Must(err)

	cfg := test.NewInsecureTransportConfig()

	h, err := th.NewServer(th.ServerParams{
		Shutdowner: test.NewShutdowner(), Mux: mux,
		Config: cfg.HTTP, Logger: logger, Tracer: tracer,
		UserAgent: test.UserAgent, Version: test.Version,
	})
	runtime.Must(err)

	transport.Register(transport.RegisterParams{Lifecycle: lc, Servers: []transport.Server{h}})

	client := &http.Client{Transport: http.DefaultTransport}
	url := fmt.Sprintf("http://%s/hello", cfg.HTTP.Address)

	lc.RequireStart()
	b.ResetTimer()

	b.Run("trace", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
			runtime.Must(err)

			_, err = client.Do(req)
			runtime.Must(err)
		}
	})

	b.StopTimer()
	lc.RequireStop()
}
