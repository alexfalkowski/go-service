package http_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/net/http/rpc"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/stretchr/testify/require"
)

func TestRPCNoContent(t *testing.T) {
	for _, mt := range []string{"json", "hjson", "yaml", "yml", "toml", "gob"} {
		t.Run(mt, func(t *testing.T) {
			world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())

			rpc.Route("/hello", test.NoContent)
			httpClient, err := world.NewHTTP()
			require.NoError(t, err)

			client := rpc.NewClient(world.ServerURL("http"),
				rpc.WithClientContentType("application/"+mt),
				rpc.WithClientRoundTripper(httpClient.Transport),
				rpc.WithClientTimeout("10s"),
			)
			req := &test.Request{Name: "Bob"}
			res := &test.Response{}

			err = client.Post(t.Context(), "/hello", req, res)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, status.Code(err))
		})
	}
}

func TestRPCWithContent(t *testing.T) {
	for _, mt := range []string{"json", "hjson", "yaml", "yml", "toml", "gob"} {
		t.Run(mt, func(t *testing.T) {
			world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())

			rpc.Route("/hello", test.SuccessSayHello)
			httpClient, err := world.NewHTTP()
			require.NoError(t, err)

			client := rpc.NewClient(world.ServerURL("http"),
				rpc.WithClientContentType("application/"+mt),
				rpc.WithClientRoundTripper(httpClient.Transport),
				rpc.WithClientTimeout("10s"),
			)
			req := &test.Request{Name: "Bob"}
			res := &test.Response{}

			err = client.Post(t.Context(), "/hello", req, res)
			require.NoError(t, err)
			require.Equal(t, "Hello Bob", res.Greeting)
		})
	}
}

func TestSuccessProtobufRPC(t *testing.T) {
	for _, mt := range []string{"proto", "protobuf", "prototext", "protojson"} {
		t.Run(mt, func(t *testing.T) {
			world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())

			rpc.Route("/hello", test.SuccessProtobufSayHello)

			client := rpc.NewClient(world.ServerURL("http"), rpc.WithClientContentType("application/"+mt))
			req := &v1.SayHelloRequest{Name: "Bob"}
			res := &v1.SayHelloResponse{}

			err := client.Post(t.Context(), "/hello", req, res)
			require.NoError(t, err)
			require.Equal(t, "Hello Bob", res.GetMessage())
		})
	}
}

func TestErroneousProtobufRPC(t *testing.T) {
	handlers := []struct {
		handler content.RequestHandler[v1.SayHelloRequest, v1.SayHelloResponse]
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

					rpc.Route("/hello", handler.handler)

					client := rpc.NewClient(world.ServerURL("http"), rpc.WithClientContentType("application/"+mt))
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

func TestErroneousUnmarshalRPC(t *testing.T) {
	for _, mt := range []string{"json", "hjson", "yaml", "yml", "toml", "gob"} {
		t.Run(mt, func(t *testing.T) {
			world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())

			rpc.Route("/hello", test.SuccessSayHello)

			url := world.PathServerURL("http", "hello")

			header := http.Header{}
			header.Set(content.TypeKey, "application/"+mt)

			res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodPost, header, bytes.NewBufferString("an erroneous payload"))
			require.NoError(t, err)
			require.NotEmpty(t, body)
			require.Equal(t, http.StatusBadRequest, res.StatusCode)
		})
	}
}

func TestErrorRPC(t *testing.T) {
	handlers := []struct {
		handler content.RequestHandler[test.Request, test.Response]
		name    string
	}{
		{name: "mapped", handler: test.ErrorSayHello},
		{name: "not-mapped", handler: test.ErrorNotMappedSayHello},
	}

	for _, handler := range handlers {
		t.Run(handler.name, func(t *testing.T) {
			for _, mt := range []string{"json", "hjson", "yaml", "yml", "toml", "gob"} {
				t.Run(mt, func(t *testing.T) {
					world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())

					rpc.Route("/hello", handler.handler)

					header := http.Header{}
					header.Set(content.TypeKey, "application/"+mt)

					enc := test.Encoder.Get(mt)

					b := test.Pool.Get()
					defer test.Pool.Put(b)

					err := enc.Encode(b, test.Request{Name: "Bob"})
					require.NoError(t, err)

					url := world.PathServerURL("http", "hello")

					res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodPost, header, b)
					require.NoError(t, err)
					require.Equal(t, "failed", body)
					require.Equal(t, http.StatusInternalServerError, res.StatusCode)
				})
			}
		})
	}
}

func TestRPCNotFound(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())

	header := http.Header{}
	header.Set(content.TypeKey, media.JSON)

	res, body, err := world.ResponseWithBody(t.Context(), world.PathServerURL("http", "missing"), http.MethodPost, header, http.NoBody)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, res.StatusCode)
	require.Equal(t, media.WithUTF8(media.Error), res.Header.Get(content.TypeKey))
	require.Equal(t, http.StatusText(http.StatusNotFound), body)
}

func TestAllowedRPC(t *testing.T) {
	for _, mt := range []string{"json", "hjson", "yaml", "yml", "toml", "gob"} {
		t.Run(mt, func(t *testing.T) {
			world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())

			rpc.Route("/hello", test.SuccessSayHello)
			httpClient, err := world.NewHTTP()
			require.NoError(t, err)

			client := rpc.NewClient(world.ServerURL("http"),
				rpc.WithClientContentType("application/"+mt),
				rpc.WithClientRoundTripper(httpClient.Transport))
			req := &test.Request{Name: "Bob"}
			res := &test.Response{}

			err = client.Post(t.Context(), "/hello", req, res)
			require.NoError(t, err)
			require.Equal(t, "Hello Bob", res.Greeting)
		})
	}
}

func TestDisallowedRPC(t *testing.T) {
	for _, mt := range []string{media.JSON, media.HJSON, media.YAML, "application/yml", media.TOML, "application/gob", "test"} {
		t.Run(mt, func(t *testing.T) {
			world := test.NewStartedWorld(t,
				test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)),
				test.WithWorldToken(nil, test.NewVerifier("test")), test.WithWorldHTTP(),
			)

			rpc.Route("/hello", test.SuccessSayHello)
			httpClient, err := world.NewHTTP()
			require.NoError(t, err)

			client := rpc.NewClient(world.ServerURL("http"),
				rpc.WithClientContentType(mt),
				rpc.WithClientRoundTripper(httpClient.Transport))
			req := &test.Request{Name: "Bob"}
			res := &test.Response{}

			err = client.Post(t.Context(), "/hello", req, res)
			require.Error(t, err)
			require.True(t, status.IsError(err))
			require.Equal(t, http.StatusUnauthorized, status.Code(err))
			require.Equal(t, "token: invalid match", err.Error())
		})
	}
}

func TestInvalidRPCRequest(t *testing.T) {
	for _, mt := range []string{"gob"} {
		t.Run(mt, func(t *testing.T) {
			world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())

			rpc.Route("/hello", test.SuccessSayHello)
			httpClient, err := world.NewHTTP()
			require.NoError(t, err)

			client := rpc.NewClient(world.ServerURL("http"),
				rpc.WithClientContentType("application/"+mt),
				rpc.WithClientRoundTripper(httpClient.Transport))

			require.Error(t, client.Post(t.Context(), "/hello", nil, &test.Response{}))
		})
	}
}

func TestInvalidRPCResponse(t *testing.T) {
	for _, mt := range []string{"json"} {
		t.Run(mt, func(t *testing.T) {
			world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())

			rpc.Route("/hello", test.SuccessSayHello)
			httpClient, err := world.NewHTTP()
			require.NoError(t, err)

			client := rpc.NewClient(world.ServerURL("http"),
				rpc.WithClientContentType("application/"+mt),
				rpc.WithClientRoundTripper(httpClient.Transport))

			require.Error(t, client.Post(t.Context(), "/hello", &test.Request{Name: "Bob"}, nil))
		})
	}
}
