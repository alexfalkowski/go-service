package http_test

import (
	"fmt"
	"io"
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

		listener, err := net.Listen(b.Context(), "tcp", "localhost:0")
		require.NoError(b, err)
		defer listener.Close()

		server := &http.Server{
			Handler:           mux,
			ReadHeaderTimeout: time.Second.Duration(),
		}
		defer server.Close()

		//nolint:errcheck
		go server.Serve(listener)

		client := &http.Client{Transport: http.DefaultTransport}
		url := fmt.Sprintf("http://%s/hello", listener.Addr().String())

		b.ResetTimer()

		req, err := http.NewRequestWithContext(b.Context(), http.MethodGet, url, http.NoBody)
		require.NoError(b, err)

		for b.Loop() {
			resp, err := client.Do(req)
			if err != nil {
				require.NoError(b, err)
			}
			closeResponse(b, resp)
		}

		b.StopTimer()
		client.CloseIdleConnections()
	})

	b.Run("none", func(b *testing.B) {
		b.ReportAllocs()

		mux := transporthttp.NewServeMux()
		mux.HandleFunc("GET /hello", func(_ http.ResponseWriter, _ *http.Request) {})
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
		cfg.HTTP.Address = test.BoundAddress(cfg.HTTP.Address, h.GetService().String())

		server.Register(lc, []*server.Service{h.GetService()})

		lc.RequireStart()

		_, addr, _ := net.SplitNetworkAddress(cfg.HTTP.Address)
		client := &transporthttp.Client{Transport: transporthttp.DefaultTransport}
		url := fmt.Sprintf("http://%s/hello", addr)

		b.ResetTimer()

		req, err := transporthttp.NewRequestWithContext(b.Context(), transporthttp.MethodGet, url, transporthttp.NoBody)
		require.NoError(b, err)

		for b.Loop() {
			resp, err := client.Do(req)
			if err != nil {
				require.NoError(b, err)
			}
			closeResponse(b, resp)
		}

		b.StopTimer()
		client.CloseIdleConnections()
		lc.RequireStop()
	})

	b.Run("log", func(b *testing.B) {
		b.ReportAllocs()

		mux := transporthttp.NewServeMux()
		mux.HandleFunc("GET /hello", func(_ http.ResponseWriter, _ *http.Request) {})
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
		cfg.HTTP.Address = test.BoundAddress(cfg.HTTP.Address, h.GetService().String())

		server.Register(lc, []*server.Service{h.GetService()})
		errors.Register(errors.NewHandler(logger))

		lc.RequireStart()

		_, addr, _ := net.SplitNetworkAddress(cfg.HTTP.Address)
		client := &transporthttp.Client{Transport: transporthttp.DefaultTransport}
		url := fmt.Sprintf("http://%s/hello", addr)

		b.ResetTimer()

		req, err := transporthttp.NewRequestWithContext(b.Context(), transporthttp.MethodGet, url, transporthttp.NoBody)
		require.NoError(b, err)

		for b.Loop() {
			resp, err := client.Do(req)
			if err != nil {
				require.NoError(b, err)
			}
			closeResponse(b, resp)
		}

		b.StopTimer()
		client.CloseIdleConnections()
		lc.RequireStop()
	})

	b.Run("trace", func(b *testing.B) {
		b.ReportAllocs()

		mux := transporthttp.NewServeMux()
		mux.HandleFunc("GET /hello", func(_ http.ResponseWriter, _ *http.Request) {})
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
		cfg.HTTP.Address = test.BoundAddress(cfg.HTTP.Address, h.GetService().String())

		server.Register(lc, []*server.Service{h.GetService()})
		errors.Register(errors.NewHandler(logger))

		lc.RequireStart()

		_, addr, _ := net.SplitNetworkAddress(cfg.HTTP.Address)
		client := &transporthttp.Client{Transport: transporthttp.DefaultTransport}
		url := fmt.Sprintf("http://%s/hello", addr)

		b.ResetTimer()

		req, err := transporthttp.NewRequestWithContext(b.Context(), transporthttp.MethodGet, url, transporthttp.NoBody)
		require.NoError(b, err)

		for b.Loop() {
			resp, err := client.Do(req)
			if err != nil {
				require.NoError(b, err)
			}
			closeResponse(b, resp)
		}

		b.StopTimer()
		client.CloseIdleConnections()
		lc.RequireStop()
	})
}

func BenchmarkMVC(b *testing.B) {
	b.Run("html", func(b *testing.B) {
		b.ReportAllocs()

		world := newHTTPBenchmarkWorld(b)
		view := mvc.NewFullView("views/hello.tmpl")

		mvc.Get("/hello", func(_ context.Context) (*mvc.View, *test.Page, error) {
			return view, &test.Model, nil
		})

		startHTTPBenchmarkWorld(b, world)

		client, err := world.NewHTTP()
		require.NoError(b, err)
		url := world.PathServerURL("http", "hello")

		req, err := transporthttp.NewRequestWithContext(b.Context(), transporthttp.MethodGet, url, transporthttp.NoBody)
		require.NoError(b, err)

		req.Header.Set(content.TypeKey, mime.HTMLMediaType)

		b.ResetTimer()

		for b.Loop() {
			resp, err := client.Do(req)
			if err != nil {
				require.NoError(b, err)
			}
			closeResponse(b, resp)
		}

		client.CloseIdleConnections()
		world.RequireStop()
	})
}

func closeResponse(b *testing.B, resp *http.Response) {
	b.Helper()

	_, err := io.Copy(io.Discard, resp.Body)
	require.NoError(b, err)
	require.NoError(b, resp.Body.Close())
}

func newHTTPBenchmarkWorld(b *testing.B) *test.World {
	b.Helper()

	logger, err := logger.NewLogger(logger.LoggerParams{})
	require.NoError(b, err)

	return test.NewWorld(b, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP(), test.WithWorldLogger(logger))
}

func startHTTPBenchmarkWorld(b *testing.B, world *test.World) {
	b.Helper()

	world.RequireStart()
	conn, err := test.Connect(b.Context(), world.TransportConfig.HTTP.Address)
	require.NoError(b, err)
	require.NoError(b, conn.Close())
}

//nolint:funlen
func BenchmarkRPC(b *testing.B) {
	b.Run("text", func(b *testing.B) {
		b.ReportAllocs()

		world := newHTTPBenchmarkWorld(b)
		rpc.Route("/hello", test.SuccessSayHello)
		startHTTPBenchmarkWorld(b, world)

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
			cl.CloseIdleConnections()
		}

		b.StopTimer()
		world.RequireStop()
	})

	b.Run("proto", func(b *testing.B) {
		b.ReportAllocs()

		world := newHTTPBenchmarkWorld(b)
		rpc.Route("/hello", test.SuccessProtobufSayHello)
		startHTTPBenchmarkWorld(b, world)

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
			cl.CloseIdleConnections()
		}

		b.StopTimer()
		world.RequireStop()
	})
}

//nolint:funlen
func BenchmarkRest(b *testing.B) {
	b.Run("text", func(b *testing.B) {
		b.ReportAllocs()

		world := newHTTPBenchmarkWorld(b)
		test.RegisterRequestHandlers("/hello", test.RestRequestContent)
		mvc.StaticFile("/robots.txt", "static/robots.txt")
		startHTTPBenchmarkWorld(b, world)

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
			cl.CloseIdleConnections()
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

			cl.CloseIdleConnections()
		})

		b.StopTimer()
		world.RequireStop()
	})

	b.Run("proto", func(b *testing.B) {
		b.ReportAllocs()

		world := newHTTPBenchmarkWorld(b)
		test.RegisterRequestHandlers("/hello", test.RestRequestProtobuf)
		startHTTPBenchmarkWorld(b, world)

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
			cl.CloseIdleConnections()
		}

		b.StopTimer()
		world.RequireStop()
	})
}
