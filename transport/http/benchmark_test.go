//nolint:varnamelen
package http_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/alexfalkowski/go-service/id"
	"github.com/alexfalkowski/go-service/internal/test"
	v1 "github.com/alexfalkowski/go-service/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/net/http/mvc"
	"github.com/alexfalkowski/go-service/net/http/rpc"
	"github.com/alexfalkowski/go-service/runtime"
	"github.com/alexfalkowski/go-service/server"
	"github.com/alexfalkowski/go-service/telemetry/errors"
	"github.com/alexfalkowski/go-service/telemetry/logger"
	"github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/transport"
	th "github.com/alexfalkowski/go-service/transport/http"
	"go.uber.org/fx/fxtest"
)

func BenchmarkDefaultHTTP(b *testing.B) {
	b.ReportAllocs()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello", func(_ http.ResponseWriter, _ *http.Request) {})

	server := &http.Server{
		Handler:           mux,
		Addr:              test.Address(),
		ReadHeaderTimeout: time.Second,
	}
	defer server.Close()

	//nolint:errcheck
	go server.ListenAndServe()

	b.ResetTimer()

	b.Run("std", func(b *testing.B) {
		client := &http.Client{Transport: http.DefaultTransport}
		url := fmt.Sprintf("http://%s/hello", server.Addr)

		req, err := http.NewRequestWithContext(b.Context(), http.MethodGet, url, http.NoBody)
		runtime.Must(err)

		for b.Loop() {
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
		Shutdowner: test.NewShutdowner(),
		Mux:        mux,
		Config:     cfg.HTTP,
		UserAgent:  test.UserAgent,
		Version:    test.Version,
		ID:         &id.UUID{},
	})
	runtime.Must(err)

	transport.Register(lc, []*server.Service{h.GetServer()})

	lc.RequireStart()
	b.ResetTimer()

	b.Run("none", func(b *testing.B) {
		client := &http.Client{Transport: http.DefaultTransport}
		url := fmt.Sprintf("http://%s/hello", cfg.HTTP.Address)

		req, err := http.NewRequestWithContext(b.Context(), http.MethodGet, url, http.NoBody)
		runtime.Must(err)

		for b.Loop() {
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
	logger, _ := logger.NewLogger(logger.Params{})
	cfg := test.NewInsecureTransportConfig()

	h, err := th.NewServer(th.ServerParams{
		Shutdowner: test.NewShutdowner(),
		Mux:        mux,
		Config:     cfg.HTTP,
		Logger:     logger,
		UserAgent:  test.UserAgent,
		Version:    test.Version,
		ID:         &id.UUID{},
	})
	runtime.Must(err)

	transport.Register(lc, []*server.Service{h.GetServer()})
	errors.Register(errors.NewHandler(logger))

	lc.RequireStart()
	b.ResetTimer()

	b.Run("log", func(b *testing.B) {
		client := &http.Client{Transport: http.DefaultTransport}
		url := fmt.Sprintf("http://%s/hello", cfg.HTTP.Address)

		req, err := http.NewRequestWithContext(b.Context(), http.MethodGet, url, http.NoBody)
		runtime.Must(err)

		for b.Loop() {
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
	logger, _ := logger.NewLogger(logger.Params{})
	tracer := test.NewTracer(lc, nil)
	cfg := test.NewInsecureTransportConfig()

	h, err := th.NewServer(th.ServerParams{
		Shutdowner: test.NewShutdowner(),
		Mux:        mux,
		Config:     cfg.HTTP,
		Logger:     logger,
		Tracer:     tracer,
		UserAgent:  test.UserAgent,
		Version:    test.Version,
		ID:         &id.UUID{},
	})
	runtime.Must(err)

	transport.Register(lc, []*server.Service{h.GetServer()})
	errors.Register(errors.NewHandler(logger))

	lc.RequireStart()
	b.ResetTimer()

	b.Run("trace", func(b *testing.B) {
		client := &http.Client{Transport: http.DefaultTransport}
		url := fmt.Sprintf("http://%s/hello", cfg.HTTP.Address)

		req, err := http.NewRequestWithContext(b.Context(), http.MethodGet, url, http.NoBody)
		runtime.Must(err)

		for b.Loop() {
			_, err = client.Do(req)
			runtime.Must(err)
		}
	})

	b.StopTimer()
	lc.RequireStop()
}

func BenchmarkRoute(b *testing.B) {
	b.ReportAllocs()

	logger, _ := logger.NewLogger(logger.Params{})

	world := test.NewWorld(b, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP(), test.WithWorldLogger(logger))
	world.Register()

	world.RequireStart()

	mvc.Route("GET /hello", func(_ context.Context) (mvc.View, *test.PageData, error) {
		return mvc.View("hello.tmpl"), &test.Model, nil
	})

	b.ResetTimer()

	b.Run("html", func(b *testing.B) {
		client := world.NewHTTP()

		req, err := http.NewRequestWithContext(b.Context(), http.MethodGet, fmt.Sprintf("http://%s/hello", world.InsecureServerHost()), http.NoBody)
		runtime.Must(err)

		req.Header.Set("Content-Type", "text/html")

		for b.Loop() {
			_, err := client.Do(req)
			runtime.Must(err)
		}
	})

	b.StopTimer()
	world.RequireStop()
}

func BenchmarkRPC(b *testing.B) {
	b.ReportAllocs()

	logger, _ := logger.NewLogger(logger.Params{})

	world := test.NewWorld(b, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP(), test.WithWorldLogger(logger))
	world.Register()

	world.RequireStart()

	rpc.Route("/hello", test.SuccessSayHello)

	b.ResetTimer()

	for _, mt := range []string{"json", "yaml", "yml", "toml", "gob"} {
		cl := world.NewHTTP()
		url := "http://" + world.InsecureServerHost()
		client := rpc.NewClient(url,
			rpc.WithClientContentType("application/"+mt),
			rpc.WithClientRoundTripper(cl.Transport),
		)

		b.Run(mt, func(b *testing.B) {
			for b.Loop() {
				req := &test.Request{Name: "Bob"}
				res := &test.Response{}

				err := client.Invoke(b.Context(), "/hello", req, res)
				runtime.Must(err)
			}
		})
	}

	b.StopTimer()
	world.RequireStop()
}

func BenchmarkProtobuf(b *testing.B) {
	b.ReportAllocs()

	logger, _ := logger.NewLogger(logger.Params{})

	world := test.NewWorld(b, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP(), test.WithWorldLogger(logger))
	world.Register()

	world.RequireStart()

	rpc.Route("/hello", test.SuccessProtobufSayHello)

	b.ResetTimer()

	for _, mt := range []string{"proto", "protobuf", "prototext", "protojson"} {
		cl := world.NewHTTP()
		url := "http://" + world.InsecureServerHost()
		client := rpc.NewClient(url,
			rpc.WithClientContentType("application/"+mt),
			rpc.WithClientRoundTripper(cl.Transport))

		b.Run(mt, func(b *testing.B) {
			for b.Loop() {
				req := &v1.SayHelloRequest{Name: "Bob"}
				res := &v1.SayHelloResponse{}

				err := client.Invoke(b.Context(), "/hello", req, res)
				runtime.Must(err)
			}
		})
	}

	b.StopTimer()
	world.RequireStop()
}
