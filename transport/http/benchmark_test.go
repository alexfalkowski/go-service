package http_test

import (
	"fmt"
	"testing"

	configclient "github.com/alexfalkowski/go-service/v2/config/client"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/id/uuid"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/net"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/client"
	"github.com/alexfalkowski/go-service/v2/net/http/content/stream"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/net/http/mvc"
	"github.com/alexfalkowski/go-service/v2/net/http/rest"
	"github.com/alexfalkowski/go-service/v2/net/http/rpc"
	"github.com/alexfalkowski/go-service/v2/net/server"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/telemetry/errors"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/time"
	transporthttp "github.com/alexfalkowski/go-service/v2/transport/http"
	"github.com/stretchr/testify/require"
)

// BenchmarkHTTP compares the standard HTTP stack with the supported go-service server stack and telemetry layers.
func BenchmarkHTTP(b *testing.B) {
	b.Run("std", benchmarkStdHTTP)
	benchmarkHTTPLayers(b, benchmarkHTTP)
}

// BenchmarkMVC measures the supported MVC HTML rendering path through the HTTP transport stack.
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

		req.Header.Set(http.ContentTypeKey, media.HTML)

		b.ResetTimer()

		for b.Loop() {
			resp, err := client.Do(req)
			if err != nil {
				require.NoError(b, err)
			}
			closeResponse(b, resp)
		}

		b.StopTimer()
		client.CloseIdleConnections()
		world.RequireStop()
	})
}

// BenchmarkRPC measures RPC client/server body encoding overhead across supported HTTP media types.
func BenchmarkRPC(b *testing.B) {
	b.Run("text", func(b *testing.B) {
		b.ReportAllocs()

		world := newHTTPBenchmarkWorld(b)
		rpc.Route("/hello", test.SuccessSayHello)
		startHTTPBenchmarkWorld(b, world)

		b.ResetTimer()

		for _, mt := range test.MessageMediaTypes() {
			cl, err := world.NewHTTP()
			require.NoError(b, err)
			client := rpc.NewClient(world.ServerURL("http"),
				rpc.WithClientContentType(mt.ContentType),
				rpc.WithClientRoundTripper(cl.Transport),
			)

			b.Run(mt.Name, func(b *testing.B) {
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

// BenchmarkRest measures REST client/server body encoding and static response overhead across supported media types.
//
//nolint:funlen
func BenchmarkRest(b *testing.B) {
	b.Run("text", func(b *testing.B) {
		b.ReportAllocs()

		world := newHTTPBenchmarkWorld(b)
		test.RegisterRequestHandlers("/hello", test.RestRequestContent)
		mvc.StaticFile("/robots.txt", "static/robots.txt")
		startHTTPBenchmarkWorld(b, world)

		b.ResetTimer()

		for _, mt := range test.MessageMediaTypes() {
			cl, err := world.NewHTTP()
			require.NoError(b, err)
			url := world.NamedServerURL("http", "hello")
			client := rest.NewClient(rest.WithClientRoundTripper(cl.Transport))

			b.Run(mt.Name, func(b *testing.B) {
				for b.Loop() {
					req := &test.Request{Name: "Bob"}
					res := &test.Response{}
					opts := rest.Options{
						ContentType: mt.ContentType,
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
				opts := rest.Options{
					ContentType: media.Text,
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
					opts := rest.Options{
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

// BenchmarkRestStream measures REST bidirectional streaming (StreamPost) round-trip overhead.
func BenchmarkRestStream(b *testing.B) {
	b.ReportAllocs()

	world := test.NewWorld(b, test.WithWorldSecure(), test.WithWorldTelemetry("otlp"), test.WithWorldRest(), test.WithWorldHTTP())
	rest.StreamPost("/hello", streamHello)
	startHTTPBenchmarkWorld(b, world)

	url := world.PathServerURL("https", "hello")

	b.ResetTimer()

	for b.Loop() {
		err := world.Rest.StreamPost(b.Context(), url, rest.Options{ContentType: media.NDJSON, Accept: media.NDJSON}, streamHelloClient("Bob", 1))
		if err != nil {
			require.NoError(b, err)
		}
	}

	b.StopTimer()
	world.RequireStop()
}

// BenchmarkHTTPStream measures bidirectional streaming with fixed 1 KiB request payloads through progressively enabled HTTP transport layers and TLS.
func BenchmarkHTTPStream(b *testing.B) {
	for _, messageCount := range []int{1, 10, 100} {
		b.Run(fmt.Sprintf("%d-messages", messageCount), func(b *testing.B) {
			b.Helper()

			benchmarkHTTPLayers(b, func(b *testing.B, log *logger.Logger, trace, tlsEnabled bool) {
				b.Helper()

				benchmarkHTTPStream(b, log, trace, tlsEnabled, messageCount)
			})
		})
	}
}

// BenchmarkRPCStream measures RPC bidirectional streaming (StreamRoute) round-trip overhead.
//
// rpc.Client has no streaming method (only Post), so this benchmarks a raw
// [github.com/alexfalkowski/go-service/v2/net/http/client.Client] directly against the world's server
// URL, matching the correctness test this benchmark mirrors.
func BenchmarkRPCStream(b *testing.B) {
	b.ReportAllocs()

	world := test.NewWorld(b, test.WithWorldSecure(), test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())
	rpc.StreamRoute("/hello", streamHello)
	startHTTPBenchmarkWorld(b, world)

	httpClient, err := world.NewHTTP()
	require.NoError(b, err)

	c := client.NewClient(test.UnaryContent, test.StreamContent, test.Pool, client.WithRoundTripper(httpClient.Transport))
	url := world.PathServerURL("https", "hello")

	b.ResetTimer()

	for b.Loop() {
		err := c.StreamPost(b.Context(), url, client.Options{ContentType: media.NDJSON, Accept: media.NDJSON}, streamHelloClient("Bob", 1))
		if err != nil {
			require.NoError(b, err)
		}
	}

	b.StopTimer()
	httpClient.CloseIdleConnections()
	world.RequireStop()
}

func streamHello(_ context.Context, stream *stream.RequestStream[test.Request, test.Response]) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			if stream.IsFinished(err) {
				return nil
			}

			return err
		}

		if err := stream.Send(&test.Response{Greeting: "Hello " + req.Name}); err != nil {
			return err
		}
	}
}

func streamHelloClient(name string, messageCount int) func(context.Context, *client.RequestResponseStream) error {
	return func(_ context.Context, stream *client.RequestResponseStream) error {
		for range messageCount {
			if err := stream.Send(&test.Request{Name: name}); err != nil {
				return err
			}

			var res test.Response
			if err := stream.Recv(&res); err != nil {
				return err
			}
		}

		return nil
	}
}

func benchmarkHTTPLayers(b *testing.B, benchmark func(*testing.B, *logger.Logger, bool, bool)) {
	b.Helper()

	b.Run("none", func(b *testing.B) {
		benchmark(b, nil, false, false)
	})

	b.Run("log", func(b *testing.B) {
		log, err := logger.NewLogger(logger.LoggerParams{})
		require.NoError(b, err)

		benchmark(b, log, false, false)
	})

	b.Run("trace", func(b *testing.B) {
		log, err := logger.NewLogger(logger.LoggerParams{})
		require.NoError(b, err)

		benchmark(b, log, true, false)
	})

	b.Run("tls", func(b *testing.B) {
		benchmark(b, nil, false, true)
	})
}

func benchmarkStdHTTP(b *testing.B) {
	b.ReportAllocs()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello", func(_ http.ResponseWriter, _ *http.Request) {})

	listener, err := net.Listen(b.Context(), "tcp", "localhost:0")
	require.NoError(b, err)
	defer listener.Close()

	httpServer := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: time.Second.Duration(),
	}
	defer httpServer.Close()

	//nolint:errcheck
	go httpServer.Serve(listener)

	httpClient := &http.Client{Transport: http.DefaultTransport}
	url := fmt.Sprintf("http://%s/hello", listener.Addr().String())
	req, err := http.NewRequestWithContext(b.Context(), http.MethodGet, url, http.NoBody)
	require.NoError(b, err)

	b.ResetTimer()

	for b.Loop() {
		resp, err := httpClient.Do(req)
		if err != nil {
			require.NoError(b, err)
		}
		closeResponse(b, resp)
	}

	b.StopTimer()
	httpClient.CloseIdleConnections()
}

func benchmarkHTTP(b *testing.B, log *logger.Logger, trace, tlsEnabled bool) {
	b.Helper()
	b.ReportAllocs()

	address, cleanup := startBenchmarkHTTPServer(b, log, trace, tlsEnabled, func(router *transporthttp.Router, _ *transporthttp.Config) {
		router.Handle("GET /hello", transporthttp.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))
	})
	httpClient, scheme, closeIdleConnections := newBenchmarkHTTPClient(b, tlsEnabled)
	url := fmt.Sprintf("%s://%s/hello", scheme, address)
	req, err := transporthttp.NewRequestWithContext(b.Context(), transporthttp.MethodGet, url, transporthttp.NoBody)
	require.NoError(b, err)

	b.ResetTimer()

	for b.Loop() {
		resp, err := httpClient.Do(req)
		if err != nil {
			require.NoError(b, err)
		}
		closeResponse(b, resp)
	}

	b.StopTimer()
	closeIdleConnections()
	cleanup()
}

func benchmarkHTTPStream(b *testing.B, log *logger.Logger, trace, tlsEnabled bool, messageCount int) {
	b.Helper()
	b.ReportAllocs()

	address, cleanup := startBenchmarkHTTPServer(b, log, trace, tlsEnabled, func(router *transporthttp.Router, cfg *transporthttp.Config) {
		opts := stream.Options{ReadTimeout: cfg.GetReadTimeout(), WriteTimeout: cfg.GetWriteTimeout(), MaxReceiveSize: cfg.GetMaxReceiveSize()}
		rest.Register(router, test.UnaryContent, test.StreamContent, test.Pool, opts)
		rest.StreamPost("/hello", streamHello)
	})
	httpClient, scheme, closeIdleConnections := newHTTPStreamClient(b, tlsEnabled)
	url := fmt.Sprintf("%s://%s/hello", scheme, address)
	requestStream := streamHelloClient(strings.Repeat("b", 1024), messageCount)

	b.ResetTimer()

	for b.Loop() {
		err := httpClient.StreamPost(b.Context(), url, client.Options{ContentType: media.NDJSON, Accept: media.NDJSON}, requestStream)
		if err != nil {
			require.NoError(b, err)
		}
	}

	b.StopTimer()
	closeIdleConnections()
	cleanup()
}

func startBenchmarkHTTPServer(b *testing.B, log *logger.Logger, trace, tlsEnabled bool, register func(*transporthttp.Router, *transporthttp.Config)) (string, func()) {
	b.Helper()

	mux := transporthttp.NewServeMux()
	policy := transporthttp.NewRoutePolicy()
	router := transporthttp.NewRouter(mux, policy)
	lc := test.QuietLifecycle(b)
	cfg := test.NewInsecureTransportConfig()
	if tlsEnabled {
		cfg = test.NewSecureTransportConfig()
	}

	if trace {
		test.RegisterTracer(lc, nil)
	}

	register(router, cfg.HTTP)

	httpServer, err := transporthttp.NewServer(transporthttp.ServerParams{
		Shutdowner:  test.NewShutdowner(),
		Mux:         mux,
		Pool:        test.Pool,
		Config:      cfg.HTTP,
		Logger:      log,
		UserAgent:   test.UserAgent,
		Version:     test.Version,
		ID:          uuid.NewGenerator(),
		RoutePolicy: policy,
	})
	require.NoError(b, err)
	cfg.HTTP.Address = test.BoundAddress(cfg.HTTP.Address, httpServer.GetService().String())

	if log != nil {
		errors.Register(errors.NewHandler(nil))
	}

	server.Register(server.RegisterParams{Lifecycle: lc, Drain: server.NewDrain(), Services: []*server.Service{httpServer.GetService()}})
	lc.RequireStart()

	_, address, _ := net.SplitNetworkAddress(cfg.HTTP.Address)

	return address, lc.RequireStop
}

func newBenchmarkHTTPClient(b *testing.B, tlsEnabled bool) (*transporthttp.Client, string, func()) {
	b.Helper()

	if !tlsEnabled {
		httpClient := &transporthttp.Client{Transport: transporthttp.DefaultTransport}

		return httpClient, "http", httpClient.CloseIdleConnections
	}

	tlsConfig, err := configclient.NewConfig(test.FS, test.NewTLSClientConfig())
	require.NoError(b, err)

	transport := http.Transport(tlsConfig)
	httpClient := &transporthttp.Client{Transport: transport}

	return httpClient, "https", httpClient.CloseIdleConnections
}

func newHTTPStreamClient(b *testing.B, tlsEnabled bool) (*client.Client, string, func()) {
	b.Helper()

	transport := http.Transport(nil)
	scheme := "http"
	if tlsEnabled {
		tlsConfig, err := configclient.NewConfig(test.FS, test.NewTLSClientConfig())
		require.NoError(b, err)

		transport = http.Transport(tlsConfig)
		scheme = "https"
	} else {
		transport.Protocols.SetHTTP1(false)
	}

	return client.NewClient(test.UnaryContent, test.StreamContent, test.Pool, client.WithRoundTripper(transport)), scheme, transport.CloseIdleConnections
}

func closeResponse(b *testing.B, resp *http.Response) {
	b.Helper()

	_, err := io.Copy(io.Discard, resp.Body)
	require.NoError(b, err)
	require.NoError(b, resp.Body.Close())
}

func newHTTPBenchmarkWorld(b *testing.B) *test.World {
	b.Helper()

	log, err := logger.NewLogger(logger.LoggerParams{})
	require.NoError(b, err)

	return test.NewWorld(b, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP(), test.WithWorldLogger(log))
}

func startHTTPBenchmarkWorld(b *testing.B, world *test.World) {
	b.Helper()

	world.RequireStart()
	conn, err := test.Connect(b.Context(), world.TransportConfig.HTTP.Address)
	require.NoError(b, err)
	require.NoError(b, conn.Close())
}
