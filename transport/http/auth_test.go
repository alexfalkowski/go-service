package http_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/id/uuid"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/mime"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/net/http/rpc"
	"github.com/alexfalkowski/go-service/v2/token"
	"github.com/alexfalkowski/go-service/v2/transport/strings"
	"github.com/stretchr/testify/require"
)

func TestTokenAuthUnary(t *testing.T) {
	for _, kind := range []string{"jwt", "paseto", "ssh"} {
		cfg := test.NewToken(kind)
		ec := test.NewEd25519()
		signer, _ := ed25519.NewSigner(test.PEM, ec)
		verifier, _ := ed25519.NewVerifier(test.PEM, ec)
		gen := uuid.NewGenerator()
		tkn := token.NewToken(test.Name, cfg, test.FS, signer, verifier, gen)

		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(tkn, tkn), test.WithWorldHTTP())
		world.Register()
		world.RequireStart()

		rpc.Route("/hello", test.SuccessSayHello)

		header := http.Header{}
		header.Set(content.TypeKey, mime.JSONMediaType)
		header.Set("Request-Id", "test")
		header.Set("X-Forwarded-For", "127.0.0.1")
		header.Set("Geolocation", "geo:47,11")

		url := world.PathServerURL("http", "hello")

		res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodPost, header, bytes.NewBufferString(`{"name":"test"}`))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)
		require.NotEmpty(t, body)

		world.RequireStop()
	}
}

func TestValidAuthUnary(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(test.NewGenerator("test", nil), test.NewVerifier("test")), test.WithWorldHTTP())
	world.Register()
	world.RequireStart()

	rpc.Route("/hello", test.SuccessSayHello)

	header := http.Header{}
	header.Set(content.TypeKey, mime.JSONMediaType)
	header.Set("Request-Id", "test")
	header.Set("X-Forwarded-For", "127.0.0.1")

	url := world.PathServerURL("http", "hello")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodPost, header, bytes.NewBufferString(`{"name":"test"}`))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.NotEmpty(t, body)

	world.RequireStop()
}

func TestInvalidAuthUnary(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(test.NewGenerator("bob", nil), test.NewVerifier("test")), test.WithWorldHTTP())
	world.Register()
	world.RequireStart()

	rpc.Route("/hello", test.SuccessSayHello)

	header := http.Header{}
	header.Set(content.TypeKey, mime.JSONMediaType)
	header.Set("Request-Id", "test")

	url := world.PathServerURL("http", "hello")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodPost, header, bytes.NewBufferString(`{"name":"test"}`))
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, res.StatusCode)
	require.Contains(t, body, `token: invalid match`)

	world.RequireStop()
}

func TestAuthUnaryWithAppend(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())
	world.Register()
	world.RequireStart()

	rpc.Route("/hello", test.SuccessSayHello)

	header := http.Header{}
	header.Set(content.TypeKey, mime.JSONMediaType)
	header.Set("Request-Id", "test")
	header.Set("Authorization", "What Invalid")

	url := world.PathServerURL("http", "hello")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodPost, header, bytes.NewBufferString(`{"name":"test"}`))
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, res.StatusCode)
	require.NotEmpty(t, body)

	world.RequireStop()
}

func TestMissingAuthUnary(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(nil, test.NewVerifier("test")), test.WithWorldHTTP())
	world.Register()
	world.RequireStart()

	rpc.Route("/hello", test.SuccessSayHello)

	header := http.Header{}
	header.Set(content.TypeKey, mime.JSONMediaType)
	header.Set("Request-Id", "test")

	url := world.PathServerURL("http", "hello")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodPost, header, bytes.NewBufferString(`{"name":"test"}`))
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, res.StatusCode)
	require.Contains(t, body, "invalid match")

	world.RequireStop()
}

func TestEmptyAuthUnary(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(test.NewGenerator(strings.Empty, nil), test.NewVerifier("test")), test.WithWorldHTTP())
	world.Register()
	world.RequireStart()

	rpc.Route("/hello", test.SuccessSayHello)

	header := http.Header{}
	header.Set(content.TypeKey, mime.JSONMediaType)
	header.Set("Request-Id", "test")

	url := world.PathServerURL("http", "hello")

	_, _, err := world.ResponseWithBody(t.Context(), url, http.MethodPost, header, bytes.NewBufferString(`{"name":"test"}`))
	require.Error(t, err)
	require.Contains(t, err.Error(), "authorization is invalid")

	world.RequireStop()
}

func TestMissingClientAuthUnary(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(nil, test.NewVerifier("test")), test.WithWorldHTTP())
	world.Register()
	world.RequireStart()

	rpc.Route("/hello", test.SuccessSayHello)

	header := http.Header{}
	header.Set(content.TypeKey, mime.JSONMediaType)
	header.Set("Request-Id", "test")

	url := world.PathServerURL("http", "hello")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodPost, header, bytes.NewBufferString(`{"name":"test"}`))
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, res.StatusCode)
	require.Contains(t, body, "invalid match")

	world.RequireStop()
}

func TestTokenErrorAuthUnary(t *testing.T) {
	world := test.NewWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldToken(test.NewGenerator(strings.Empty, test.ErrGenerate), test.NewVerifier("test")),
		test.WithWorldHTTP())
	world.Register()
	world.RequireStart()

	rpc.Route("/hello", test.SuccessSayHello)

	header := http.Header{}
	header.Set(content.TypeKey, mime.JSONMediaType)
	header.Set("Request-Id", "test")

	url := world.PathServerURL("http", "hello")

	_, _, err := world.ResponseWithBody(t.Context(), url, http.MethodPost, header, bytes.NewBufferString(`{"name":"test"}`))
	require.Error(t, err)
	require.Contains(t, err.Error(), "token: generation issue")

	world.RequireStop()
}

func TestBreakerAuthUnary(t *testing.T) {
	world := test.NewWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldToken(test.NewGenerator(strings.Empty, test.ErrGenerate), test.NewVerifier("test")),
		test.WithWorldHTTP(),
	)
	world.Register()
	world.RequireStart()

	var err error
	url := world.PathServerURL("http", "hello")

	for range 10 {
		header := http.Header{}
		header.Set(content.TypeKey, mime.JSONMediaType)
		header.Set("Request-Id", "test")

		_, _, err = world.ResponseWithBody(t.Context(), url, http.MethodPost, header, bytes.NewBufferString(`{"name":"test"}`))
	}
	require.Error(t, err)

	world.RequireStop()
}
