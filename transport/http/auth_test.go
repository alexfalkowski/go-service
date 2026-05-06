package http_test

import (
	"strconv"
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/id/uuid"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/mime"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/net/http/rpc"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/token"
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

			server := test.NewHTTPTestServer(t, test.WithHTTPTestVerifier(tkn))
			client := test.NewHTTPTestClient(t, server, test.WithHTTPClientGenerator(tkn))

			rpc.Route("/hello", test.SuccessSayHello)

			header := http.Header{}
			header.Set(content.TypeKey, mime.JSONMediaType)
			header.Set("Request-Id", "test")
			header.Set("X-Forwarded-For", "127.0.0.1")
			header.Set("Geolocation", "geo:47,11")

			url := server.URL + "/hello"

			res, body, err := test.HTTPClientResponseWithBody(t, client, http.MethodPost, url, header, bytes.NewBufferString(`{"name":"test"}`))
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, res.StatusCode)
			require.NotEmpty(t, body)
		})
	}
}

func TestUnknownTokenKindAuthUnary(t *testing.T) {
	cfg := test.NewToken("none")
	tkn := token.NewToken(test.Name, cfg, test.FS, nil, nil, nil)

	server := test.NewHTTPTestServer(t, test.WithHTTPTestVerifier(tkn))
	client := test.NewHTTPTestClient(t, server, test.WithHTTPClientGenerator(test.NewGenerator("test", nil)))

	rpc.Route("/hello", test.SuccessSayHello)

	header := http.Header{}
	header.Set(content.TypeKey, mime.JSONMediaType)
	header.Set("Request-Id", "test")

	url := server.URL + "/hello"

	res, body, err := test.HTTPClientResponseWithBody(t, client, http.MethodPost, url, header, bytes.NewBufferString(`{"name":"test"}`))
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, res.StatusCode)
	require.Contains(t, body, "token: invalid config")
}

func TestValidAuthUnary(t *testing.T) {
	server := test.NewHTTPTestServer(t, test.WithHTTPTestVerifier(test.NewVerifier("test")))
	client := test.NewHTTPTestClient(t, server, test.WithHTTPClientGenerator(test.NewGenerator("test", nil)))

	rpc.Route("/hello", test.SuccessSayHello)

	header := http.Header{}
	header.Set(content.TypeKey, mime.JSONMediaType)
	header.Set("Request-Id", "test")
	header.Set("X-Forwarded-For", "127.0.0.1")

	url := server.URL + "/hello"

	res, body, err := test.HTTPClientResponseWithBody(t, client, http.MethodPost, url, header, bytes.NewBufferString(`{"name":"test"}`))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.NotEmpty(t, body)
}

func TestInvalidAuthUnary(t *testing.T) {
	server := test.NewHTTPTestServer(t, test.WithHTTPTestVerifier(test.NewVerifier("test")))
	client := test.NewHTTPTestClient(t, server, test.WithHTTPClientGenerator(test.NewGenerator("bob", nil)))

	rpc.Route("/hello", test.SuccessSayHello)

	header := http.Header{}
	header.Set(content.TypeKey, mime.JSONMediaType)
	header.Set("Request-Id", "test")

	url := server.URL + "/hello"

	res, body, err := test.HTTPClientResponseWithBody(t, client, http.MethodPost, url, header, bytes.NewBufferString(`{"name":"test"}`))
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, res.StatusCode)
	require.Contains(t, body, `token: invalid match`)
}

func TestAuthUnaryWithAppend(t *testing.T) {
	server := test.NewHTTPTestServer(t)

	rpc.Route("/hello", test.SuccessSayHello)

	header := http.Header{}
	header.Set(content.TypeKey, mime.JSONMediaType)
	header.Set("Request-Id", "test")
	header.Set("Authorization", "What Invalid")

	url := server.URL + "/hello"

	res, body := test.HTTPResponseWithBody(t, server, http.MethodPost, url, header, bytes.NewBufferString(`{"name":"test"}`))
	require.Equal(t, http.StatusBadRequest, res.StatusCode)
	require.NotEmpty(t, body)
}

func TestAuthUnaryWithLowercaseBearer(t *testing.T) {
	server := test.NewHTTPTestServer(t, test.WithHTTPTestVerifier(test.NewVerifier("test")))

	rpc.Route("/hello", test.SuccessSayHello)

	header := http.Header{}
	header.Set(content.TypeKey, mime.JSONMediaType)
	header.Set("Request-Id", "test")
	header.Set("Authorization", "bearer test")

	url := server.URL + "/hello"

	res, body := test.HTTPResponseWithBody(t, server, http.MethodPost, url, header, bytes.NewBufferString(`{"name":"test"}`))
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.NotEmpty(t, body)
}

func TestMissingAuthUnary(t *testing.T) {
	server := test.NewHTTPTestServer(t, test.WithHTTPTestVerifier(test.NewVerifier("test")))

	rpc.Route("/hello", test.SuccessSayHello)

	header := http.Header{}
	header.Set(content.TypeKey, mime.JSONMediaType)
	header.Set("Request-Id", "test")

	url := server.URL + "/hello"

	res, body := test.HTTPResponseWithBody(t, server, http.MethodPost, url, header, bytes.NewBufferString(`{"name":"test"}`))
	require.Equal(t, http.StatusUnauthorized, res.StatusCode)
	require.Contains(t, body, "invalid match")
}

func TestEmptyAuthUnary(t *testing.T) {
	server := test.NewHTTPTestServer(t, test.WithHTTPTestVerifier(test.NewVerifier("test")))
	client := test.NewHTTPTestClient(t, server, test.WithHTTPClientGenerator(test.NewGenerator(strings.Empty, nil)))

	rpc.Route("/hello", test.SuccessSayHello)

	header := http.Header{}
	header.Set(content.TypeKey, mime.JSONMediaType)
	header.Set("Request-Id", "test")

	url := server.URL + "/hello"

	_, _, err := test.HTTPClientResponseWithBody(t, client, http.MethodPost, url, header, bytes.NewBufferString(`{"name":"test"}`))
	require.Error(t, err)
	require.Contains(t, err.Error(), "authorization is invalid")
}

func TestMissingClientAuthUnary(t *testing.T) {
	server := test.NewHTTPTestServer(t, test.WithHTTPTestVerifier(test.NewVerifier("test")))

	rpc.Route("/hello", test.SuccessSayHello)

	header := http.Header{}
	header.Set(content.TypeKey, mime.JSONMediaType)
	header.Set("Request-Id", "test")

	url := server.URL + "/hello"

	res, body := test.HTTPResponseWithBody(t, server, http.MethodPost, url, header, bytes.NewBufferString(`{"name":"test"}`))
	require.Equal(t, http.StatusUnauthorized, res.StatusCode)
	require.Contains(t, body, "invalid match")
}

func TestTokenErrorAuthUnary(t *testing.T) {
	server := test.NewHTTPTestServer(t, test.WithHTTPTestVerifier(test.NewVerifier("test")))
	client := test.NewHTTPTestClient(t, server, test.WithHTTPClientGenerator(test.NewGenerator(strings.Empty, test.ErrGenerate)))

	rpc.Route("/hello", test.SuccessSayHello)

	header := http.Header{}
	header.Set(content.TypeKey, mime.JSONMediaType)
	header.Set("Request-Id", "test")

	url := server.URL + "/hello"

	_, _, err := test.HTTPClientResponseWithBody(t, client, http.MethodPost, url, header, bytes.NewBufferString(`{"name":"test"}`))
	require.Error(t, err)
	require.Contains(t, err.Error(), "token: generation issue")
}

func TestBreakerAuthUnary(t *testing.T) {
	server := test.NewHTTPTestServer(t, test.WithHTTPTestVerifier(test.NewVerifier("test")))
	client := test.NewHTTPTestClient(t, server, test.WithHTTPClientGenerator(test.NewGenerator(strings.Empty, test.ErrGenerate)))

	var err error
	url := server.URL + "/hello"

	for i := range 10 {
		t.Run("attempt-"+strconv.Itoa(i+1), func(t *testing.T) {
			header := http.Header{}
			header.Set(content.TypeKey, mime.JSONMediaType)
			header.Set("Request-Id", "test")

			_, _, err = test.HTTPClientResponseWithBody(t, client, http.MethodPost, url, header, bytes.NewBufferString(`{"name":"test"}`))
		})
	}
	require.Error(t, err)
}
