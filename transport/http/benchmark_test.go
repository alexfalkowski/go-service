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
	"github.com/alexfalkowski/go-service/v2/net/server"
	"github.com/alexfalkowski/go-service/v2/telemetry/errors"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/time"
	transporthttp "github.com/alexfalkowski/go-service/v2/transport/http"
	"github.com/stretchr/testify/require"
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
		require.NoError(b, err)

		for b.Loop() {
			_, err = client.Do(req)
			if err != nil {
				require.NoError(b, err)
			}
		}

		b.StopTimer()
	})

	b.Run("none", func(b *testing.B) {
		b.ReportAllocs()

		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(b)
		cfg := test.NewInsecureTransportConfig()

		h, err := transporthttp.NewServer(transporthttp.ServerParams{
			Shutdowner: test.NewShutdowner(),
			Mux:        mux,
			Config:     cfg.HTTP,
			UserAgent:  test.UserAgent,
			Version:    test.Version,
			ID:         uuid.NewGenerator(),
		})
		require.NoError(b, err)

		server.Register(lc, []*server.Service{h.GetService()})

		lc.RequireStart()

		_, addr, _ := net.SplitNetworkAddress(cfg.HTTP.Address)
		client := &http.Client{Transport: http.DefaultTransport}
		url := fmt.Sprintf("http://%s/hello", addr)

		b.ResetTimer()

		req, err := http.NewRequestWithContext(b.Context(), http.MethodGet, url, http.NoBody)
		require.NoError(b, err)

		for b.Loop() {
			_, err = client.Do(req)
			if err != nil {
				require.NoError(b, err)
			}
		}

		b.StopTimer()
		lc.RequireStop()
	})

	b.Run("log", func(b *testing.B) {
		b.ReportAllocs()

		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(b)
		logger, err := logger.NewLogger(logger.LoggerParams{})
		require.NoError(b, err)
		cfg := test.NewInsecureTransportConfig()

		h, err := transporthttp.NewServer(transporthttp.ServerParams{
			Shutdowner: test.NewShutdowner(),
			Mux:        mux,
			Config:     cfg.HTTP,
			Logger:     logger,
			UserAgent:  test.UserAgent,
			Version:    test.Version,
			ID:         uuid.NewGenerator(),
		})
		require.NoError(b, err)

		server.Register(lc, []*server.Service{h.GetService()})
		errors.Register(errors.NewHandler(logger))

		lc.RequireStart()

		_, addr, _ := net.SplitNetworkAddress(cfg.HTTP.Address)
		client := &http.Client{Transport: http.DefaultTransport}
		url := fmt.Sprintf("http://%s/hello", addr)

		b.ResetTimer()

		req, err := http.NewRequestWithContext(b.Context(), http.MethodGet, url, http.NoBody)
		require.NoError(b, err)

		for b.Loop() {
			_, err = client.Do(req)
			if err != nil {
				require.NoError(b, err)
			}
		}

		b.StopTimer()
		lc.RequireStop()
	})

	b.Run("trace", func(b *testing.B) {
		b.ReportAllocs()

		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(b)
		logger, err := logger.NewLogger(logger.LoggerParams{})
		require.NoError(b, err)
		cfg := test.NewInsecureTransportConfig()

		test.RegisterTracer(lc, nil)

		h, err := transporthttp.NewServer(transporthttp.ServerParams{
			Shutdowner: test.NewShutdowner(),
			Mux:        mux,
			Config:     cfg.HTTP,
			Logger:     logger,
			UserAgent:  test.UserAgent,
			Version:    test.Version,
			ID:         uuid.NewGenerator(),
		})
		require.NoError(b, err)

		server.Register(lc, []*server.Service{h.GetService()})
		errors.Register(errors.NewHandler(logger))

		lc.RequireStart()

		_, addr, _ := net.SplitNetworkAddress(cfg.HTTP.Address)
		client := &http.Client{Transport: http.DefaultTransport}
		url := fmt.Sprintf("http://%s/hello", addr)

		b.ResetTimer()

		req, err := http.NewRequestWithContext(b.Context(), http.MethodGet, url, http.NoBody)
		require.NoError(b, err)

		for b.Loop() {
			_, err = client.Do(req)
			if err != nil {
				require.NoError(b, err)
			}
		}

		b.StopTimer()
		lc.RequireStop()
	})
}

func BenchmarkMVC(b *testing.B) {
	b.ReportAllocs()

	logger, err := logger.NewLogger(logger.LoggerParams{})
	require.NoError(b, err)

	world := test.NewStartedWorld(b, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP(), test.WithWorldLogger(logger))

	view := mvc.NewFullView("views/hello.tmpl")

	mvc.Get("/hello", func(_ context.Context) (*mvc.View, *test.Page, error) {
		return view, &test.Model, nil
	})

	b.ResetTimer()

	b.Run("html", func(b *testing.B) {
		client, err := world.NewHTTP()
		require.NoError(b, err)
		url := world.PathServerURL("http", "hello")

		req, err := http.NewRequestWithContext(b.Context(), http.MethodGet, url, http.NoBody)
		require.NoError(b, err)

		req.Header.Set(content.TypeKey, mime.HTMLMediaType)

		for b.Loop() {
			_, err = client.Do(req)
			if err != nil {
				require.NoError(b, err)
			}
		}
	})

	b.StopTimer()
}

//nolint:funlen
func BenchmarkRPC(b *testing.B) {
	b.Run("text", func(b *testing.B) {
		b.ReportAllocs()

		logger, err := logger.NewLogger(logger.LoggerParams{})
		require.NoError(b, err)

		world := test.NewStartedWorld(b, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP(), test.WithWorldLogger(logger))

		rpc.Route("/hello", test.SuccessSayHello)

		b.ResetTimer()

		for _, mt := range []string{"json", "hjson", "yaml", "yml", "toml", "gob"} {
			cl, err := world.NewHTTP()
			require.NoError(b, err)
			client := rpc.NewClient(world.ServerURL("http"),
				rpc.WithClientContentType("application/"+mt),
				rpc.WithClientRoundTripper(cl.Transport),
			)

			b.Run(mt, func(b *testing.B) {
				for b.Loop() {
					req := &test.Request{Name: "Bob"}
					res := &test.Response{}

					err := client.Post(b.Context(), "/hello", req, res)
					if err != nil {
						require.NoError(b, err)
					}
				}
			})
		}

		b.StopTimer()
	})

	b.Run("proto", func(b *testing.B) {
		b.ReportAllocs()

		logger, err := logger.NewLogger(logger.LoggerParams{})
		require.NoError(b, err)

		world := test.NewStartedWorld(b, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP(), test.WithWorldLogger(logger))

		rpc.Route("/hello", test.SuccessProtobufSayHello)

		b.ResetTimer()

		for _, mt := range []string{"proto", "protobuf", "prototext", "protojson"} {
			cl, err := world.NewHTTP()
			require.NoError(b, err)
			client := rpc.NewClient(world.ServerURL("http"),
				rpc.WithClientContentType("application/"+mt),
				rpc.WithClientRoundTripper(cl.Transport))

			b.Run(mt, func(b *testing.B) {
				for b.Loop() {
					req := &v1.SayHelloRequest{Name: "Bob"}
					res := &v1.SayHelloResponse{}

					err := client.Post(b.Context(), "/hello", req, res)
					if err != nil {
						require.NoError(b, err)
					}
				}
			})
		}

		b.StopTimer()
	})
}

//nolint:funlen
func BenchmarkRest(b *testing.B) {
	b.Run("text", func(b *testing.B) {
		b.ReportAllocs()

		logger, err := logger.NewLogger(logger.LoggerParams{})
		require.NoError(b, err)

		world := test.NewStartedWorld(b, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP(), test.WithWorldLogger(logger))

		test.RegisterRequestHandlers("/hello", test.RestRequestContent)
		mvc.StaticFile("/robots.txt", "static/robots.txt")

		b.ResetTimer()

		for _, mt := range []string{"json", "hjson", "yaml", "yml", "toml", "gob"} {
			cl, err := world.NewHTTP()
			require.NoError(b, err)
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
					if err != nil {
						require.NoError(b, err)
					}
				}
			})
		}

		b.Run("static", func(b *testing.B) {
			cl, err := world.NewHTTP()
			require.NoError(b, err)
			url := world.PathServerURL("http", "robots.txt")
			client := rest.NewClient(rest.WithClientRoundTripper(cl.Transport))

			for b.Loop() {
				buffer := test.Pool.Get()
				opts := &rest.Options{
					ContentType: mime.TextMediaType,
					Response:    buffer,
				}

				err := client.Get(b.Context(), url, opts)
				if err != nil {
					require.NoError(b, err)
				}

				test.Pool.Put(buffer)
			}
		})

		b.StopTimer()
	})

	b.Run("proto", func(b *testing.B) {
		b.ReportAllocs()

		logger, err := logger.NewLogger(logger.LoggerParams{})
		require.NoError(b, err)

		world := test.NewStartedWorld(b, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP(), test.WithWorldLogger(logger))

		test.RegisterRequestHandlers("/hello", test.RestRequestProtobuf)

		b.ResetTimer()

		for _, mt := range []string{"proto", "protobuf", "prototext", "protojson"} {
			cl, err := world.NewHTTP()
			require.NoError(b, err)
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
					if err != nil {
						require.NoError(b, err)
					}
				}
			})
		}

		b.StopTimer()
	})
}
