package http_test

import (
	"encoding/json"
	"testing"

	"github.com/alexfalkowski/go-service/v2/breaker"
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/id/uuid"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/mime"
	netheader "github.com/alexfalkowski/go-service/v2/net/header"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/net/http/rpc"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/token"
	httpbreaker "github.com/alexfalkowski/go-service/v2/transport/http/breaker"
	"github.com/stretchr/testify/require"
)

func TestTokenAuthUnary(t *testing.T) {
	for _, kind := range []string{"jwt", "paseto", "ssh"} {
		t.Run(kind, func(t *testing.T) {
			cfg := test.NewToken(kind)
			ec := test.NewEd25519()
			signer, _ := ed25519.NewSigner(test.PEM, ec)
			verifier, _ := ed25519.NewVerifier(test.PEM, ec)
			gen := uuid.NewGenerator()
			tkn := token.NewToken(test.Name, cfg, test.FS, signer, verifier, gen)

			world := test.NewStartedHTTPWorld(t, func(*test.World) {
				rpc.Route("/hello", test.SuccessSayHello)
			}, test.WithWorldTelemetry("otlp"), test.WithWorldToken(tkn, tkn))

			header := http.Header{}
			header.Set(content.TypeKey, mime.JSONMediaType)
			header.Set("Request-Id", "test")
			header.Set("X-Forwarded-For", "127.0.0.1")
			header.Set("Geolocation", "geo:47,11")

			url := world.PathServerURL("http", "hello")

			res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodPost, header, bytes.NewBufferString(`{"name":"test"}`))
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, res.StatusCode)
			require.Equal(t, mime.JSONMediaType, res.Header.Get(content.TypeKey))

			var resp greetingResponse
			require.NoError(t, json.Unmarshal([]byte(body), &resp))
			require.Equal(t, "Hello test", resp.Greeting)
		})
	}
}

func TestValidAuthUnary(t *testing.T) {
	world := test.NewStartedHTTPWorld(t, func(*test.World) {
		rpc.Route("/hello", test.SuccessSayHello)
	}, test.WithWorldTelemetry("otlp"), test.WithWorldToken(test.NewGenerator("test", nil), test.NewVerifier("test")))

	header := http.Header{}
	header.Set(content.TypeKey, mime.JSONMediaType)
	header.Set("Request-Id", "test")
	header.Set("X-Forwarded-For", "127.0.0.1")

	url := world.PathServerURL("http", "hello")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodPost, header, bytes.NewBufferString(`{"name":"test"}`))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, mime.JSONMediaType, res.Header.Get(content.TypeKey))

	var resp greetingResponse
	require.NoError(t, json.Unmarshal([]byte(body), &resp))
	require.Equal(t, "Hello test", resp.Greeting)
}

func TestInvalidAuthUnary(t *testing.T) {
	world := test.NewStartedHTTPWorld(t, func(*test.World) {
		rpc.Route("/hello", test.SuccessSayHello)
	}, test.WithWorldTelemetry("otlp"), test.WithWorldToken(test.NewGenerator("bob", nil), test.NewVerifier("test")))

	header := http.Header{}
	header.Set(content.TypeKey, mime.JSONMediaType)
	header.Set("Request-Id", "test")

	url := world.PathServerURL("http", "hello")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodPost, header, bytes.NewBufferString(`{"name":"test"}`))
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, res.StatusCode)
	require.Equal(t, mime.ErrorMediaType, res.Header.Get(content.TypeKey))
	require.Equal(t, test.ErrInvalid.Error(), body)
}

func TestAuthUnaryWithAppend(t *testing.T) {
	world := test.NewStartedHTTPWorld(t, func(*test.World) {
		rpc.Route("/hello", test.SuccessSayHello)
	}, test.WithWorldTelemetry("otlp"))

	header := http.Header{}
	header.Set(content.TypeKey, mime.JSONMediaType)
	header.Set("Request-Id", "test")
	header.Set("Authorization", "What Invalid")

	url := world.PathServerURL("http", "hello")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodPost, header, bytes.NewBufferString(`{"name":"test"}`))
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, res.StatusCode)
	require.Equal(t, mime.ErrorMediaType, res.Header.Get(content.TypeKey))
	require.Equal(t, netheader.ErrNotSupportedAuthorization.Error(), body)
}

func TestMissingAuthUnary(t *testing.T) {
	assertUnauthorizedAuthUnary(t, test.NewStartedHTTPWorld(t, func(*test.World) {
		rpc.Route("/hello", test.SuccessSayHello)
	}, test.WithWorldTelemetry("otlp"), test.WithWorldToken(nil, test.NewVerifier("test"))))
}

func TestEmptyAuthUnary(t *testing.T) {
	world := test.NewStartedHTTPWorld(t, func(world *test.World) {
		world.TransportConfig.HTTP.Retry = nil
		rpc.Route("/hello", test.SuccessSayHello)
	}, test.WithWorldTelemetry("otlp"), test.WithWorldToken(test.NewGenerator(strings.Empty, nil), test.NewVerifier("test")))

	header := http.Header{}
	header.Set(content.TypeKey, mime.JSONMediaType)
	header.Set("Request-Id", "test")

	url := world.PathServerURL("http", "hello")

	_, _, err := world.ResponseWithBody(t.Context(), url, http.MethodPost, header, bytes.NewBufferString(`{"name":"test"}`))
	require.Error(t, err)
	require.Equal(t, http.StatusUnauthorized, status.Code(err))
	require.Equal(t, netheader.ErrInvalidAuthorization.Error(), err.Error())
}

func TestMissingClientAuthUnary(t *testing.T) {
	assertUnauthorizedAuthUnary(t, test.NewStartedHTTPWorld(t, func(*test.World) {
		rpc.Route("/hello", test.SuccessSayHello)
	}, test.WithWorldTelemetry("otlp"), test.WithWorldToken(nil, test.NewVerifier("test"))))
}

func TestTokenErrorAuthUnary(t *testing.T) {
	world := test.NewStartedHTTPWorld(t, func(world *test.World) {
		world.TransportConfig.HTTP.Retry = nil
		rpc.Route("/hello", test.SuccessSayHello)
	},
		test.WithWorldTelemetry("otlp"),
		test.WithWorldToken(test.NewGenerator(strings.Empty, test.ErrGenerate), test.NewVerifier("test")),
	)

	header := http.Header{}
	header.Set(content.TypeKey, mime.JSONMediaType)
	header.Set("Request-Id", "test")

	url := world.PathServerURL("http", "hello")

	_, _, err := world.ResponseWithBody(t.Context(), url, http.MethodPost, header, bytes.NewBufferString(`{"name":"test"}`))
	require.Error(t, err)
	require.Equal(t, http.StatusUnauthorized, status.Code(err))
	require.Equal(t, test.ErrGenerate.Error(), err.Error())
}

func TestBreakerAuthUnary(t *testing.T) {
	world := test.NewStartedHTTPWorld(t, func(world *test.World) {
		world.TransportConfig.HTTP.Retry = nil
		rpc.Route("/hello", test.SuccessSayHello)
	},
		test.WithWorldTelemetry("otlp"),
		test.WithWorldToken(test.NewGenerator("bob", nil), test.NewVerifier("test")),
	)

	url := world.PathServerURL("http", "hello")
	client := world.NewHTTP(
		httpbreaker.WithSettings(httpbreaker.Settings{
			MaxRequests: 1,
			Interval:    0,
			Timeout:     time.Minute,
			ReadyToTrip: func(counts breaker.Counts) bool {
				return counts.ConsecutiveFailures >= 2
			},
		}),
		httpbreaker.WithFailureStatuses(http.StatusUnauthorized),
	)

	req1, err := http.NewRequestWithContext(t.Context(), http.MethodPost, url, bytes.NewBufferString(`{"name":"test"}`))
	require.NoError(t, err)
	req1.Header.Set(content.TypeKey, mime.JSONMediaType)
	req1.Header.Set("Request-Id", "test")

	res, err := client.Do(req1)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, res.StatusCode)
	require.Equal(t, mime.ErrorMediaType, res.Header.Get(content.TypeKey))
	body, _, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Equal(t, test.ErrInvalid.Error()+"\n", string(body))

	req2, err := http.NewRequestWithContext(t.Context(), http.MethodPost, url, bytes.NewBufferString(`{"name":"test"}`))
	require.NoError(t, err)
	req2.Header.Set(content.TypeKey, mime.JSONMediaType)
	req2.Header.Set("Request-Id", "test")

	res, err = client.Do(req2)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, res.StatusCode)
	require.NoError(t, res.Body.Close())

	req3, err := http.NewRequestWithContext(t.Context(), http.MethodPost, url, bytes.NewBufferString(`{"name":"test"}`))
	require.NoError(t, err)
	req3.Header.Set(content.TypeKey, mime.JSONMediaType)
	req3.Header.Set("Request-Id", "test")

	_, err = client.Do(req3)
	require.Error(t, err)
	require.ErrorIs(t, err, breaker.ErrOpenState)
}

type greetingResponse struct {
	Greeting string `json:"greeting"`
}

func assertUnauthorizedAuthUnary(t *testing.T, world *test.World) {
	t.Helper()

	header := http.Header{}
	header.Set(content.TypeKey, mime.JSONMediaType)
	header.Set("Request-Id", "test")

	url := world.PathServerURL("http", "hello")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodPost, header, bytes.NewBufferString(`{"name":"test"}`))
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, res.StatusCode)
	require.Equal(t, mime.ErrorMediaType, res.Header.Get(content.TypeKey))
	require.Equal(t, test.ErrInvalid.Error(), body)
}
