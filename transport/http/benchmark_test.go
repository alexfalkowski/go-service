package http_test

import (
	"fmt"
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/id/uuid"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/v2/mime"
	"github.com/alexfalkowski/go-service/v2/net"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/net/http/mvc"
	"github.com/alexfalkowski/go-service/v2/net/http/rest"
	"github.com/alexfalkowski/go-service/v2/net/http/rpc"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/server"
	"github.com/alexfalkowski/go-service/v2/telemetry/errors"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/transport"
	th "github.com/alexfalkowski/go-service/v2/transport/http"
	"go.uber.org/fx/fxtest"
)

//nolint:funlen
func BenchmarkHTTP(b *testing.B) {
	b.Run("std", func(b *testing.B) {
		b.ReportAllocs()

		mux := http.NewServeMux()
		mux.HandleFunc("GET /hello", func(_ http.ResponseWriter, _ *http.Request) {})

		addr := test.RandomHost()

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

		req, err := http.NewRequestWithContext(b.Context(), http.MethodGet, url, http.NoBody)
		runtime.Must(err)

		for b.Loop() {
			_, err = client.Do(req)
			runtime.Must(err)
		}

		b.StopTimer()
	})

	b.Run("none", func(b *testing.B) {
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
			ID:         uuid.NewGenerator(),
		})
		runtime.Must(err)

		transport.Register(lc, []*server.Service{h.GetService()})

		lc.RequireStart()

		_, addr, _ := net.SplitNetworkAddress(cfg.HTTP.Address)
		client := &http.Client{Transport: http.DefaultTransport}
		url := fmt.Sprintf("http://%s/hello", addr)

		b.ResetTimer()

		req, err := http.NewRequestWithContext(b.Context(), http.MethodGet, url, http.NoBody)
		runtime.Must(err)

		for b.Loop() {
			_, err = client.Do(req)
			runtime.Must(err)
		}

		b.StopTimer()
		lc.RequireStop()
	})

	b.Run("log", func(b *testing.B) {
		b.ReportAllocs()

		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(b)
		logger, _ := logger.NewLogger(logger.LoggerParams{})
		cfg := test.NewInsecureTransportConfig()

		h, err := th.NewServer(th.ServerParams{
			Shutdowner: test.NewShutdowner(),
			Mux:        mux,
			Config:     cfg.HTTP,
			Logger:     logger,
			UserAgent:  test.UserAgent,
			Version:    test.Version,
			ID:         uuid.NewGenerator(),
		})
		runtime.Must(err)

		transport.Register(lc, []*server.Service{h.GetService()})
		errors.Register(errors.NewHandler(logger))

		lc.RequireStart()

		_, addr, _ := net.SplitNetworkAddress(cfg.HTTP.Address)
		client := &http.Client{Transport: http.DefaultTransport}
		url := fmt.Sprintf("http://%s/hello", addr)

		b.ResetTimer()

		req, err := http.NewRequestWithContext(b.Context(), http.MethodGet, url, http.NoBody)
		runtime.Must(err)

		for b.Loop() {
			_, err = client.Do(req)
			runtime.Must(err)
		}

		b.StopTimer()
		lc.RequireStop()
	})

	b.Run("trace", func(b *testing.B) {
		b.ReportAllocs()

		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(b)
		logger, _ := logger.NewLogger(logger.LoggerParams{})
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
			ID:         uuid.NewGenerator(),
		})
		runtime.Must(err)

		transport.Register(lc, []*server.Service{h.GetService()})
		errors.Register(errors.NewHandler(logger))

		lc.RequireStart()

		_, addr, _ := net.SplitNetworkAddress(cfg.HTTP.Address)
		client := &http.Client{Transport: http.DefaultTransport}
		url := fmt.Sprintf("http://%s/hello", addr)

		b.ResetTimer()

		req, err := http.NewRequestWithContext(b.Context(), http.MethodGet, url, http.NoBody)
		runtime.Must(err)

		for b.Loop() {
			_, err = client.Do(req)
			runtime.Must(err)
		}

		b.StopTimer()
		lc.RequireStop()
	})
}

func BenchmarkMVC(b *testing.B) {
	b.ReportAllocs()

	logger, _ := logger.NewLogger(logger.LoggerParams{})

	world := test.NewWorld(b, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP(), test.WithWorldLogger(logger))
	world.Register()

	world.RequireStart()

	view := mvc.NewFullView("views/hello.tmpl")

	mvc.Get("/hello", func(_ context.Context) (*mvc.View, *test.Page, error) {
		return view, &test.Model, nil
	})

	b.ResetTimer()

	b.Run("html", func(b *testing.B) {
		client := world.NewHTTP()
		url := world.PathServerURL("http", "hello")

		req, err := http.NewRequestWithContext(b.Context(), http.MethodGet, url, http.NoBody)
		runtime.Must(err)

		req.Header.Set(content.TypeKey, mime.HTMLMediaType)

		for b.Loop() {
			_, err := client.Do(req)
			runtime.Must(err)
		}
	})

	b.StopTimer()
	world.RequireStop()
}

//nolint:funlen
func BenchmarkRPC(b *testing.B) {
	b.Run("text", func(b *testing.B) {
		b.ReportAllocs()

		logger, _ := logger.NewLogger(logger.LoggerParams{})

		world := test.NewWorld(b, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP(), test.WithWorldLogger(logger))
		world.Register()

		world.RequireStart()

		rpc.Route("/hello", test.SuccessSayHello)

		b.ResetTimer()

		for _, mt := range []string{"json", "yaml", "yml", "toml", "gob"} {
			cl := world.NewHTTP()
			client := rpc.NewClient(world.ServerURL("http"),
				rpc.WithClientContentType("application/"+mt),
				rpc.WithClientRoundTripper(cl.Transport),
			)

			b.Run(mt, func(b *testing.B) {
				for b.Loop() {
					req := &test.Request{Name: "Bob"}
					res := &test.Response{}

					err := client.Post(b.Context(), "/hello", req, res)
					runtime.Must(err)
				}
			})
		}

		b.StopTimer()
		world.RequireStop()
	})

	b.Run("proto", func(b *testing.B) {
		b.ReportAllocs()

		logger, _ := logger.NewLogger(logger.LoggerParams{})

		world := test.NewWorld(b, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP(), test.WithWorldLogger(logger))
		world.Register()

		world.RequireStart()

		rpc.Route("/hello", test.SuccessProtobufSayHello)

		b.ResetTimer()

		for _, mt := range []string{"proto", "protobuf", "prototext", "protojson"} {
			cl := world.NewHTTP()
			client := rpc.NewClient(world.ServerURL("http"),
				rpc.WithClientContentType("application/"+mt),
				rpc.WithClientRoundTripper(cl.Transport))

			b.Run(mt, func(b *testing.B) {
				for b.Loop() {
					req := &v1.SayHelloRequest{Name: "Bob"}
					res := &v1.SayHelloResponse{}

					err := client.Post(b.Context(), "/hello", req, res)
					runtime.Must(err)
				}
			})
		}

		b.StopTimer()
		world.RequireStop()
	})
}

//nolint:funlen
func BenchmarkRest(b *testing.B) {
	b.Run("text", func(b *testing.B) {
		b.ReportAllocs()

		logger, _ := logger.NewLogger(logger.LoggerParams{})

		world := test.NewWorld(b, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP(), test.WithWorldLogger(logger))
		world.Register()

		world.RequireStart()

		test.RegisterRequestHandlers("/hello", test.RestRequestContent)
		mvc.StaticFile("/robots.txt", "static/robots.txt")

		b.ResetTimer()

		for _, mt := range []string{"json", "yaml", "yml", "toml", "gob"} {
			cl := world.NewHTTP()
			url := world.NamedServerURL("http", "hello")
			client := rest.NewClient(rest.WithClientRoundTripper(cl.Transport))

			b.Run(mt, func(b *testing.B) {
				for b.Loop() {
					req := &test.Request{Name: "Bob"}
					res := &test.Response{}
					opts := &rest.Options{
						ContentType: "application/" + mt,
						Request:     req,
						Response:    res,
					}

					err := client.Post(b.Context(), url, opts)
					runtime.Must(err)
				}
			})
		}

		b.Run("static", func(b *testing.B) {
			cl := world.NewHTTP()
			url := world.PathServerURL("http", "robots.txt")
			client := rest.NewClient(rest.WithClientRoundTripper(cl.Transport))

			for b.Loop() {
				buffer := test.Pool.Get()
				opts := &rest.Options{
					ContentType: mime.TextMediaType,
					Response:    buffer,
				}

				err := client.Get(b.Context(), url, opts)
				runtime.Must(err)

				test.Pool.Put(buffer)
			}
		})

		b.StopTimer()
		world.RequireStop()
	})

	b.Run("proto", func(b *testing.B) {
		b.ReportAllocs()

		logger, _ := logger.NewLogger(logger.LoggerParams{})

		world := test.NewWorld(b, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP(), test.WithWorldLogger(logger))
		world.Register()

		world.RequireStart()

		test.RegisterRequestHandlers("/hello", test.RestRequestProtobuf)

		b.ResetTimer()

		for _, mt := range []string{"proto", "protobuf", "prototext", "protojson"} {
			cl := world.NewHTTP()
			url := world.NamedServerURL("http", "hello")
			client := rest.NewClient(rest.WithClientRoundTripper(cl.Transport))

			b.Run(mt, func(b *testing.B) {
				for b.Loop() {
					req := &v1.SayHelloRequest{Name: "Bob"}
					res := &v1.SayHelloResponse{}
					opts := &rest.Options{
						ContentType: "application/" + mt,
						Request:     req,
						Response:    res,
					}

					err := client.Post(b.Context(), url, opts)
					runtime.Must(err)
				}
			})
		}

		b.StopTimer()
		world.RequireStop()
	})
}
