package http_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/v2/mime"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/net/http/rpc"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/stretchr/testify/require"
)

func TestRPCNoContent(t *testing.T) {
	for _, mt := range []string{"json", "yaml", "yml", "toml", "gob"} {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())
		world.Register()
		world.RequireStart()

		rpc.Route("/hello", test.NoContent)

		client := rpc.NewClient(world.ServerURL("http"),
			rpc.WithClientContentType("application/"+mt),
			rpc.WithClientRoundTripper(world.NewHTTP().Transport),
			rpc.WithClientTimeout("10s"),
		)
		req := &test.Request{Name: "Bob"}
		res := &test.Response{}

		err := client.Post(t.Context(), "/hello", req, res)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status.Code(err))

		world.RequireStop()
	}
}

func TestRPCWithContent(t *testing.T) {
	for _, mt := range []string{"json", "yaml", "yml", "toml", "gob"} {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())
		world.Register()
		world.RequireStart()

		rpc.Route("/hello", test.SuccessSayHello)

		client := rpc.NewClient(world.ServerURL("http"),
			rpc.WithClientContentType("application/"+mt),
			rpc.WithClientRoundTripper(world.NewHTTP().Transport),
			rpc.WithClientTimeout("10s"),
		)
		req := &test.Request{Name: "Bob"}
		res := &test.Response{}

		err := client.Post(t.Context(), "/hello", req, res)
		require.NoError(t, err)
		require.Equal(t, "Hello Bob", res.Greeting)

		world.RequireStop()
	}
}

func TestSuccessProtobufRPC(t *testing.T) {
	for _, mt := range []string{"proto", "protobuf", "prototext", "protojson"} {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())
		world.Register()
		world.RequireStart()

		rpc.Route("/hello", test.SuccessProtobufSayHello)

		client := rpc.NewClient(world.ServerURL("http"), rpc.WithClientContentType("application/"+mt))
		req := &v1.SayHelloRequest{Name: "Bob"}
		res := &v1.SayHelloResponse{}

		err := client.Post(t.Context(), "/hello", req, res)
		require.NoError(t, err)
		require.Equal(t, "Hello Bob", res.GetMessage())

		world.RequireStop()
	}
}

func TestErroneousProtobufRPC(t *testing.T) {
	handlers := []content.RequestHandler[v1.SayHelloRequest, v1.SayHelloResponse]{
		test.ErrorsProtobufSayHello,
		test.ErrorsNotMappedProtobufSayHello,
		test.ErrorsInternalProtobufSayHello,
	}

	for _, handler := range handlers {
		for _, mt := range []string{"proto", "protobuf", "prototext", "protojson"} {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())
			world.Register()
			world.RequireStart()

			rpc.Route("/hello", handler)

			client := rpc.NewClient(world.ServerURL("http"), rpc.WithClientContentType("application/"+mt))
			req := &v1.SayHelloRequest{Name: "Bob"}
			res := &v1.SayHelloResponse{}

			err := client.Post(t.Context(), "/hello", req, res)
			require.Error(t, err)
			require.True(t, status.IsError(err))
			require.Equal(t, http.StatusInternalServerError, status.Code(err))

			world.RequireStop()
		}
	}
}

func TestErroneousUnmarshalRPC(t *testing.T) {
	for _, mt := range []string{"json", "yaml", "yml", "toml", "gob"} {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())
		world.Register()
		world.RequireStart()

		rpc.Route("/hello", test.SuccessSayHello)

		url := world.PathServerURL("http", "hello")

		header := http.Header{}
		header.Set(content.TypeKey, "application/"+mt)

		res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodPost, header, bytes.NewBufferString("an erroneous payload"))
		require.NoError(t, err)
		require.NotEmpty(t, body)
		require.Equal(t, http.StatusBadRequest, res.StatusCode)

		world.RequireStop()
	}
}

func TestErrorRPC(t *testing.T) {
	handlers := []content.RequestHandler[test.Request, test.Response]{
		test.ErrorSayHello,
		test.ErrorNotMappedSayHello,
	}

	for _, handler := range handlers {
		for _, mt := range []string{"json", "yaml", "yml", "toml", "gob"} {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())
			world.Register()
			world.RequireStart()

			rpc.Route("/hello", handler)

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

			world.RequireStop()
		}
	}
}

func TestAllowedRPC(t *testing.T) {
	for _, mt := range []string{"json", "yaml", "yml", "toml", "gob"} {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())
		world.Register()
		world.RequireStart()

		rpc.Route("/hello", test.SuccessSayHello)

		client := rpc.NewClient(world.ServerURL("http"),
			rpc.WithClientContentType("application/"+mt),
			rpc.WithClientRoundTripper(world.NewHTTP().Transport))
		req := &test.Request{Name: "Bob"}
		res := &test.Response{}

		err := client.Post(t.Context(), "/hello", req, res)
		require.NoError(t, err)
		require.Equal(t, "Hello Bob", res.Greeting)

		world.RequireStop()
	}
}

func TestDisallowedRPC(t *testing.T) {
	for _, mt := range []string{mime.JSONMediaType, mime.YAMLMediaType, "application/yml", mime.TOMLMediaType, "application/gob", "test"} {
		world := test.NewWorld(t,
			test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)),
			test.WithWorldToken(nil, test.NewVerifier("test")), test.WithWorldHTTP(),
		)
		world.Register()
		world.RequireStart()

		rpc.Route("/hello", test.SuccessSayHello)

		client := rpc.NewClient(world.ServerURL("http"),
			rpc.WithClientContentType(mt),
			rpc.WithClientRoundTripper(world.NewHTTP().Transport))
		req := &test.Request{Name: "Bob"}
		res := &test.Response{}

		err := client.Post(t.Context(), "/hello", req, res)
		require.Error(t, err)
		require.True(t, status.IsError(err))
		require.Equal(t, http.StatusUnauthorized, status.Code(err))
		require.Equal(t, "token: invalid match", err.Error())

		world.RequireStop()
	}
}

func TestInvalidRPCRequest(t *testing.T) {
	for _, mt := range []string{"gob"} {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())
		world.Register()
		world.RequireStart()

		rpc.Route("/hello", test.SuccessSayHello)

		client := rpc.NewClient(world.ServerURL("http"),
			rpc.WithClientContentType("application/"+mt),
			rpc.WithClientRoundTripper(world.NewHTTP().Transport))

		require.Error(t, client.Post(t.Context(), "/hello", nil, &test.Response{}))
		world.RequireStop()
	}
}

func TestInvalidRPCResponse(t *testing.T) {
	for _, mt := range []string{"json"} {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())
		world.Register()
		world.RequireStart()

		rpc.Route("/hello", test.SuccessSayHello)

		client := rpc.NewClient(world.ServerURL("http"),
			rpc.WithClientContentType("application/"+mt),
			rpc.WithClientRoundTripper(world.NewHTTP().Transport))

		require.Error(t, client.Post(t.Context(), "/hello", &test.Request{Name: "Bob"}, nil))
		world.RequireStop()
	}
}
