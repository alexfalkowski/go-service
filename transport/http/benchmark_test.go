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
	addr := test.Address()

	mux.HandleFunc("GET /hello", func(_ http.ResponseWriter, _ *http.Request) {})

	server := &http.Server{
		Handler:           mux,
		Addr:              addr,
		ReadHeaderTimeout: time.Second,
	}
	defer server.Close()

	//nolint:errcheck
	go server.ListenAndServe()

	client := &http.Client{Transport: http.DefaultTransport}
	url := fmt.Sprintf("http://%s/hello", addr)

	b.ResetTimer()

	b.Run("std", func(b *testing.B) {
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

	client := &http.Client{Transport: http.DefaultTransport}
	url := fmt.Sprintf("http://%s/hello", cfg.HTTP.Address)

	lc.RequireStart()
	b.ResetTimer()

	b.Run("none", func(b *testing.B) {
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

	client := &http.Client{Transport: http.DefaultTransport}
	url := fmt.Sprintf("http://%s/hello", cfg.HTTP.Address)

	lc.RequireStart()
	b.ResetTimer()

	b.Run("log", func(b *testing.B) {
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

	client := &http.Client{Transport: http.DefaultTransport}
	url := fmt.Sprintf("http://%s/hello", cfg.HTTP.Address)

	lc.RequireStart()
	b.ResetTimer()

	b.Run("trace", func(b *testing.B) {
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

	mux := http.NewServeMux()
	lc := fxtest.NewLifecycle(b)
	logger := zap.NewNop()
	cfg := test.NewInsecureTransportConfig()
	tc := test.NewOTLPTracerConfig()
	m := test.NewOTLPMeter(lc)

	s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, TransportConfig: cfg, Meter: m, Mux: mux}
	s.Register()

	cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}

	v := mvc.NewViews(mvc.ViewsParams{FS: &test.Views, Patterns: mvc.Patterns{"views/*.tmpl"}})
	mvc.Register(mux, v)

	mvc.Route("GET /hello", func(_ context.Context) (mvc.View, *test.PageData, error) {
		return mvc.View("hello.tmpl"), &test.Model, nil
	})

	client := cl.NewHTTP()

	lc.RequireStart()
	b.ResetTimer()

	b.Run("html", func(b *testing.B) {
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, fmt.Sprintf("http://%s/hello", cfg.HTTP.Address), http.NoBody)
		runtime.Must(err)

		req.Header.Set("Content-Type", "text/html")

		for range b.N {
			_, _ = client.Do(req)
		}
	})

	b.StopTimer()
	lc.RequireStop()
}

func BenchmarkRPC(b *testing.B) {
	b.ReportAllocs()

	mux := http.NewServeMux()
	lc := fxtest.NewLifecycle(b)
	logger := zap.NewNop()

	cfg := test.NewInsecureTransportConfig()
	tc := test.NewOTLPTracerConfig()
	m := test.NewOTLPMeter(lc)

	s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, TransportConfig: cfg, Meter: m, Mux: mux}
	s.Register()

	cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}
	t := cl.NewHTTP().Transport

	rpc.Register(mux, test.Content, test.Pool)
	rpc.Route("/hello", test.SuccessSayHello)

	url := fmt.Sprintf("http://%s/hello", cfg.HTTP.Address)

	lc.RequireStart()
	b.ResetTimer()

	for _, mt := range []string{"json", "yaml", "yml", "toml", "gob"} {
		client := rpc.NewClient[test.Request, test.Response](url,
			rpc.WithClientContentType("application/"+mt),
			rpc.WithClientRoundTripper(t))

		b.Run(mt, func(b *testing.B) {
			for range b.N {
				_, _ = client.Invoke(context.Background(), &test.Request{Name: "Bob"})
			}
		})
	}

	b.StopTimer()
	lc.RequireStop()
}

func BenchmarkProtobuf(b *testing.B) {
	b.ReportAllocs()

	mux := http.NewServeMux()
	lc := fxtest.NewLifecycle(b)
	logger := zap.NewNop()

	cfg := test.NewInsecureTransportConfig()
	tc := test.NewOTLPTracerConfig()
	m := test.NewOTLPMeter(lc)

	s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, TransportConfig: cfg, Meter: m, Mux: mux}
	s.Register()

	cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}
	t := cl.NewHTTP().Transport

	rpc.Register(mux, test.Content, test.Pool)
	rpc.Route("/hello", test.SuccessProtobufSayHello)

	url := fmt.Sprintf("http://%s/hello", cfg.HTTP.Address)

	lc.RequireStart()
	b.ResetTimer()

	for _, mt := range []string{"proto", "protobuf", "prototext", "protojson"} {
		client := rpc.NewClient[v1.SayHelloRequest, v1.SayHelloResponse](url,
			rpc.WithClientContentType("application/"+mt),
			rpc.WithClientRoundTripper(t))

		b.Run(mt, func(b *testing.B) {
			for range b.N {
				_, _ = client.Invoke(context.Background(), &v1.SayHelloRequest{Name: "Bob"})
			}
		})
	}

	b.StopTimer()
	lc.RequireStop()
}
