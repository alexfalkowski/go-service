//nolint:varnamelen
package http_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/alexfalkowski/go-service/id"
	"github.com/alexfalkowski/go-service/net/http/mvc"
	"github.com/alexfalkowski/go-service/net/http/rpc"
	"github.com/alexfalkowski/go-service/runtime"
	"github.com/alexfalkowski/go-service/test"
	v1 "github.com/alexfalkowski/go-service/test/greet/v1"
	"github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/transport"
	th "github.com/alexfalkowski/go-service/transport/http"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

func BenchmarkDefaultHTTP(b *testing.B) {
	b.ReportAllocs()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello", func(_ http.ResponseWriter, _ *http.Request) {})

	addr := test.Address()

	server := &http.Server{
		Handler:           mux,
		Addr:              addr,
		ReadHeaderTimeout: time.Second,
	}
	defer server.Close()

	//nolint:errcheck
	go server.ListenAndServe()

	b.ResetTimer()

	b.Run("std", func(b *testing.B) {
		client := &http.Client{Transport: http.DefaultTransport}
		url := fmt.Sprintf("http://%s/hello", addr)

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
		runtime.Must(err)

		for range b.N {
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
		ID:         id.Default,
	})
	runtime.Must(err)

	transport.Register(transport.RegisterParams{Lifecycle: lc, Servers: []transport.Server{h}})

	lc.RequireStart()
	b.ResetTimer()

	b.Run("none", func(b *testing.B) {
		client := &http.Client{Transport: http.DefaultTransport}
		url := fmt.Sprintf("http://%s/hello", cfg.HTTP.Address)

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
		runtime.Must(err)

		for range b.N {
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
		Shutdowner: test.NewShutdowner(),
		Mux:        mux,
		Config:     cfg.HTTP,
		Logger:     logger,
		UserAgent:  test.UserAgent,
		Version:    test.Version,
		ID:         id.Default,
	})
	runtime.Must(err)

	transport.Register(transport.RegisterParams{Lifecycle: lc, Servers: []transport.Server{h}})

	lc.RequireStart()
	b.ResetTimer()

	b.Run("log", func(b *testing.B) {
		client := &http.Client{Transport: http.DefaultTransport}
		url := fmt.Sprintf("http://%s/hello", cfg.HTTP.Address)

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
		runtime.Must(err)

		for range b.N {
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
	tracer := test.NewTracer(lc, tc, logger)
	cfg := test.NewInsecureTransportConfig()

	h, err := th.NewServer(th.ServerParams{
		Shutdowner: test.NewShutdowner(),
		Mux:        mux,
		Config:     cfg.HTTP,
		Logger:     logger,
		Tracer:     tracer,
		UserAgent:  test.UserAgent,
		Version:    test.Version,
		ID:         id.Default,
	})
	runtime.Must(err)

	transport.Register(transport.RegisterParams{Lifecycle: lc, Servers: []transport.Server{h}})

	lc.RequireStart()
	b.ResetTimer()

	b.Run("trace", func(b *testing.B) {
		client := &http.Client{Transport: http.DefaultTransport}
		url := fmt.Sprintf("http://%s/hello", cfg.HTTP.Address)

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
		runtime.Must(err)

		for range b.N {
			_, err = client.Do(req)
			runtime.Must(err)
		}
	})

	b.StopTimer()
	lc.RequireStop()
}

func BenchmarkRoute(b *testing.B) {
	b.ReportAllocs()

	world := test.NewWorld(b, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP(), test.WithWorldLogger(zap.NewNop()))
	world.Register()

	world.RequireStart()

	mvc.Route("GET /hello", func(_ context.Context) (mvc.View, *test.PageData, error) {
		return mvc.View("hello.tmpl"), &test.Model, nil
	})

	b.ResetTimer()

	b.Run("html", func(b *testing.B) {
		client := world.NewHTTP()

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, fmt.Sprintf("http://%s/hello", world.ServerHost()), http.NoBody)
		runtime.Must(err)

		req.Header.Set("Content-Type", "text/html")

		for range b.N {
			_, err := client.Do(req)
			runtime.Must(err)
		}
	})

	b.StopTimer()
	world.RequireStop()
}

func BenchmarkRPC(b *testing.B) {
	b.ReportAllocs()

	world := test.NewWorld(b, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP(), test.WithWorldLogger(zap.NewNop()))
	world.Register()

	world.RequireStart()

	rpc.Route("/hello", test.SuccessSayHello)

	b.ResetTimer()

	for _, mt := range []string{"json", "yaml", "yml", "toml", "gob"} {
		cl := world.NewHTTP()
		url := fmt.Sprintf("http://%s/hello", world.ServerHost())
		client := rpc.NewClient[test.Request, test.Response](url,
			rpc.WithClientContentType("application/"+mt),
			rpc.WithClientRoundTripper(cl.Transport),
		)

		b.Run(mt, func(b *testing.B) {
			for range b.N {
				_, err := client.Invoke(context.Background(), &test.Request{Name: "Bob"})
				runtime.Must(err)
			}
		})
	}

	b.StopTimer()
	world.RequireStop()
}

func BenchmarkProtobuf(b *testing.B) {
	b.ReportAllocs()

	world := test.NewWorld(b, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP(), test.WithWorldLogger(zap.NewNop()))
	world.Register()

	world.RequireStart()

	rpc.Route("/hello", test.SuccessProtobufSayHello)

	b.ResetTimer()

	for _, mt := range []string{"proto", "protobuf", "prototext", "protojson"} {
		cl := world.NewHTTP()
		url := fmt.Sprintf("http://%s/hello", world.ServerHost())
		client := rpc.NewClient[v1.SayHelloRequest, v1.SayHelloResponse](url,
			rpc.WithClientContentType("application/"+mt),
			rpc.WithClientRoundTripper(cl.Transport))

		b.Run(mt, func(b *testing.B) {
			for range b.N {
				_, err := client.Invoke(context.Background(), &v1.SayHelloRequest{Name: "Bob"})
				runtime.Must(err)
			}
		})
	}

	b.StopTimer()
	world.RequireStop()
}
