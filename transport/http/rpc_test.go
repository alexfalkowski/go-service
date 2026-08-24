package http_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/config/options"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/client"
	"github.com/alexfalkowski/go-service/v2/net/http/content/stream"
	"github.com/alexfalkowski/go-service/v2/net/http/content/unary"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/net/http/rpc"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

func TestRPCReturnsNoContent(t *testing.T) {
	for _, mt := range test.MessageMediaTypes() {
		t.Run(mt.Name, func(t *testing.T) {
			world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())

			world.RPCServer.Route("/hello", test.NoContent)
			httpClient, err := world.NewHTTP()
			require.NoError(t, err)

			client := rpc.NewClient(world.ServerURL("http"), test.NewContentClient(client.WithRoundTripper(httpClient.Transport)),
				rpc.WithClientContentType(mt.ContentType),
			)
			req := &test.Request{Name: "Bob"}
			res := &test.Response{}

			err = client.Post(t.Context(), "/hello", req, res)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, status.Code(err))
		})
	}
}

func TestRPCWritesContentResponse(t *testing.T) {
	for _, mt := range test.MessageMediaTypes() {
		t.Run(mt.Name, func(t *testing.T) {
			requireSuccessfulRPCPost(t, mt.ContentType)
		})
	}
}

func TestRPCWritesProtobufResponse(t *testing.T) {
	for _, mt := range []string{"proto", "protobuf", "prototext", "protojson"} {
		t.Run(mt, func(t *testing.T) {
			world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())

			world.RPCServer.Route("/hello", test.SuccessProtobufSayHello)

			client := rpc.NewClient(world.ServerURL("http"), test.NewContentClient(), rpc.WithClientContentType("application/"+mt))
			req := &v1.SayHelloRequest{Name: "Bob"}
			res := &v1.SayHelloResponse{}

			err := client.Post(t.Context(), "/hello", req, res)
			require.NoError(t, err)
			require.Equal(t, "Hello Bob", res.GetMessage())
		})
	}
}

func TestRPCPropagatesProtobufErrors(t *testing.T) {
	handlers := []struct {
		handler unary.RequestHandler[v1.SayHelloRequest, v1.SayHelloResponse]
		name    string
	}{
		{name: "mapped", handler: test.ErrorsProtobufSayHello},
		{name: "not-mapped", handler: test.ErrorsNotMappedProtobufSayHello},
		{name: "internal", handler: test.ErrorsInternalProtobufSayHello},
	}

	for _, handler := range handlers {
		t.Run(handler.name, func(t *testing.T) {
			for _, mt := range []string{"proto", "protobuf", "prototext", "protojson"} {
				t.Run(mt, func(t *testing.T) {
					world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())

					world.RPCServer.Route("/hello", handler.handler)

					client := rpc.NewClient(world.ServerURL("http"), test.NewContentClient(), rpc.WithClientContentType("application/"+mt))
					req := &v1.SayHelloRequest{Name: "Bob"}
					res := &v1.SayHelloResponse{}

					err := client.Post(t.Context(), "/hello", req, res)
					require.Error(t, err)
					require.True(t, status.IsError(err))
					require.Equal(t, http.StatusInternalServerError, status.Code(err))
				})
			}
		})
	}
}

func TestRPCReportsUnmarshalErrors(t *testing.T) {
	for _, mt := range test.MessageMediaTypes() {
		t.Run(mt.Name, func(t *testing.T) {
			world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())

			world.RPCServer.Route("/hello", test.SuccessSayHello)

			url := world.PathServerURL("http", "hello")

			header := http.Header{}
			header.Set(http.ContentTypeKey, mt.ContentType)

			res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodPost, header, bytes.NewBufferString("an erroneous payload"))
			require.NoError(t, err)
			require.NotEmpty(t, body)
			require.Equal(t, http.StatusBadRequest, res.StatusCode)
		})
	}
}

func TestRPCReturnsStatusError(t *testing.T) {
	handlers := []struct {
		handler unary.RequestHandler[test.Request, test.Response]
		name    string
	}{
		{name: "mapped", handler: test.ErrorSayHello},
		{name: "not-mapped", handler: test.ErrorNotMappedSayHello},
	}

	for _, handler := range handlers {
		t.Run(handler.name, func(t *testing.T) {
			for _, mt := range test.MessageMediaTypes() {
				t.Run(mt.Name, func(t *testing.T) {
					world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())

					world.RPCServer.Route("/hello", handler.handler)

					header := http.Header{}
					header.Set(http.ContentTypeKey, mt.ContentType)

					enc := test.Encoder.Get(mt.Kind)

					b := test.Pool.Get()
					defer test.Pool.Put(b)
					err := enc.Encode(b, test.Request{Name: "Bob"})
					require.NoError(t, err)
					payload := bytes.Clone(b.Bytes())

					url := world.PathServerURL("http", "hello")

					res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodPost, header, bytes.NewReader(payload))
					require.NoError(t, err)
					if handler.name == "mapped" {
						require.Equal(t, "failed", body)
					} else {
						require.Equal(t, "http: internal server error", body)
					}
					require.Equal(t, http.StatusInternalServerError, res.StatusCode)
				})
			}
		})
	}
}

func TestRPCReturnsNotFound(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())

	header := http.Header{}
	header.Set(http.ContentTypeKey, media.JSON)

	res, body, err := world.ResponseWithBody(t.Context(), world.PathServerURL("http", "missing"), http.MethodPost, header, http.NoBody)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, res.StatusCode)
	require.Equal(t, "text/error; charset=utf-8", res.Header.Get(http.ContentTypeKey))
	require.Equal(t, "http: not found", body)
}

func TestRPCAllowsConfiguredOperation(t *testing.T) {
	for _, mt := range test.MessageMediaTypes() {
		t.Run(mt.Name, func(t *testing.T) {
			requireSuccessfulRPCPost(t, mt.ContentType)
		})
	}
}

func TestRPCRejectsDisallowedMediaType(t *testing.T) {
	for _, mt := range []string{media.JSON, media.HumanJSON, media.YAML, media.TOML, "application/gob", media.MessagePack, "test"} {
		t.Run(mt, func(t *testing.T) {
			world := test.NewStartedWorld(t,
				test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)),
				test.WithWorldToken(nil, test.NewVerifier("test")), test.WithWorldHTTP(),
			)

			world.RPCServer.Route("/hello", test.SuccessSayHello)
			httpClient, err := world.NewHTTP()
			require.NoError(t, err)

			client := rpc.NewClient(world.ServerURL("http"), test.NewContentClient(client.WithRoundTripper(httpClient.Transport)),
				rpc.WithClientContentType(mt),
			)
			req := &test.Request{Name: "Bob"}
			res := &test.Response{}

			err = client.Post(t.Context(), "/hello", req, res)
			require.Error(t, err)
			require.True(t, status.IsError(err))
			require.Equal(t, http.StatusUnauthorized, status.Code(err))
			require.Equal(t, "http: unauthorized", err.Error())
		})
	}
}

func TestRejectsInvalidRPCRequest(t *testing.T) {
	for _, mt := range []string{"gob", "toml", "yml"} {
		t.Run(mt, func(t *testing.T) {
			world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())

			world.RPCServer.Route("/hello", test.SuccessSayHello)
			header := http.Header{}
			header.Set(http.ContentTypeKey, "application/"+mt)
			res, body, err := world.ResponseWithBody(t.Context(), world.PathServerURL("http", "hello"), http.MethodPost, header, http.NoBody)

			require.NoError(t, err)
			require.Equal(t, http.StatusUnsupportedMediaType, res.StatusCode)
			require.Equal(t, "http: unsupported media type", body)
		})
	}
}

func TestRejectsInvalidRPCResponse(t *testing.T) {
	for _, mt := range []string{"json"} {
		t.Run(mt, func(t *testing.T) {
			world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())

			world.RPCServer.Route("/hello", test.SuccessSayHello)
			httpClient, err := world.NewHTTP()
			require.NoError(t, err)

			client := rpc.NewClient(world.ServerURL("http"), test.NewContentClient(client.WithRoundTripper(httpClient.Transport)),
				rpc.WithClientContentType("application/"+mt),
			)

			require.Error(t, client.Post(t.Context(), "/hello", &test.Request{Name: "Bob"}, nil))
		})
	}
}

func TestRPCStreamRouteRefreshesReadDeadlineOverHTTP2(t *testing.T) {
	config := test.NewSecureTransportConfig()
	config.HTTP.Timeout = 100 * time.Millisecond
	config.HTTP.Options = options.Map{"read_timeout": "100ms", "write_timeout": "100ms"}
	world := test.NewWorld(t, test.WithWorldTransportConfig(config), test.WithWorldSecure(), test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())
	rpcServer := rpc.NewServer(world.Router, test.UnaryContent, test.StreamContent, stream.Options{
		ReadTimeout: config.HTTP.GetReadTimeout(), WriteTimeout: config.HTTP.GetWriteTimeout(), MaxReceiveSize: config.HTTP.GetMaxReceiveSize(),
	})
	rpcServer.StreamRoute("/hello", echoStreamServer)
	world.Start()

	httpClient, err := world.NewHTTP()
	require.NoError(t, err)

	c := client.NewClient(test.UnaryContent, test.StreamContent, test.Pool, client.WithRoundTripper(httpClient.Transport))

	url := world.PathServerURL("https", "hello")

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	var greetings []string
	err = c.StreamPost(ctx, url, client.Options{ContentType: media.NDJSON, Accept: media.NDJSON}, staleClientAfterFourEchoes(&greetings))
	require.Error(t, err)
	require.Equal(t, []string{"Hello Bob", "Hello Alice", "Hello Carol", "Hello Dan"}, greetings)
}

func requireSuccessfulRPCPost(t *testing.T, contentType string) {
	t.Helper()

	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())

	world.RPCServer.Route("/hello", test.SuccessSayHello)
	httpClient, err := world.NewHTTP()
	require.NoError(t, err)

	client := rpc.NewClient(world.ServerURL("http"), test.NewContentClient(client.WithRoundTripper(httpClient.Transport)),
		rpc.WithClientContentType(contentType),
	)
	req := &test.Request{Name: "Bob"}
	res := &test.Response{}

	err = client.Post(t.Context(), "/hello", req, res)
	require.NoError(t, err)
	require.Equal(t, "Hello Bob", res.Greeting)
}
