package http_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/config/server"
	"github.com/alexfalkowski/go-service/v2/context"
	tls "github.com/alexfalkowski/go-service/v2/crypto/tls/config"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
	transporthttp "github.com/alexfalkowski/go-service/v2/transport/http"
	"github.com/stretchr/testify/require"
)

func TestInvalidServer(t *testing.T) {
	transporthttp.Register(test.FS)

	cfg := &transporthttp.Config{
		Config: &server.Config{
			Timeout: 5 * time.Second,
			TLS:     test.NewTLSConfig("certs/client-cert.pem", "secrets/none"),
		},
	}
	params := transporthttp.ServerParams{
		Shutdowner: test.NewShutdowner(),
		Config:     cfg,
	}

	_, err := transporthttp.NewServer(params)
	require.Error(t, err)
}

func TestServerRejectsCAOnlyTLS(t *testing.T) {
	transporthttp.Register(test.FS)

	cfg := &transporthttp.Config{
		Config: &server.Config{
			Timeout: 5 * time.Second,
			TLS:     &tls.Config{CA: test.FilePath("certs/rootCA.pem")},
		},
	}
	params := transporthttp.ServerParams{
		Shutdowner: test.NewShutdowner(),
		Config:     cfg,
	}

	_, err := transporthttp.NewServer(params)
	require.ErrorIs(t, err, server.ErrMissingKeyPair)
}

func TestServerMaxReceiveSize(t *testing.T) {
	cfg := test.NewInsecureTransportConfig()
	cfg.HTTP.MaxReceiveSize = 64

	world := test.NewWorld(t, test.WithWorldTransportConfig(cfg), test.WithWorldHTTP())
	world.Handle("POST /hello", content.NewRequestHandler(test.Content, func(_ context.Context, _ *test.Request) (*test.Response, error) {
		return &test.Response{Greeting: "hello"}, nil
	}))
	world.Start()

	header := http.Header{}
	header.Set(content.TypeKey, media.JSON)

	res, body, err := world.PostBody(
		t.Context(),
		world.PathServerURL("http", "hello"),
		header,
		strings.NewReader(`{"name":"`+strings.Repeat("a", 256)+`"}`),
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusRequestEntityTooLarge, res.StatusCode)
	require.Equal(t, "http: request entity too large", body)
}

func TestServerMaxReceiveSizeWithUnknownLength(t *testing.T) {
	cfg := test.NewInsecureTransportConfig()
	cfg.HTTP.MaxReceiveSize = 64

	world := test.NewWorld(t, test.WithWorldTransportConfig(cfg), test.WithWorldHTTP())
	world.Handle("POST /hello", content.NewRequestHandler(test.Content, func(_ context.Context, _ *test.Request) (*test.Response, error) {
		return &test.Response{Greeting: "hello"}, nil
	}))
	world.Start()

	header := http.Header{}
	header.Set(content.TypeKey, media.JSON)

	res, body, err := world.PostBody(
		t.Context(),
		world.PathServerURL("http", "hello"),
		header,
		&test.UnknownLengthReader{Reader: strings.NewReader(`{"name":"` + strings.Repeat("a", 256) + `"}`)},
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusRequestEntityTooLarge, res.StatusCode)
	require.Equal(t, "http: request entity too large", body)
}

func TestServerRecoversPanic(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldHTTP())
	world.Handle("GET /panic", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic("test panic")
	}))
	world.Handle("GET /panic-after-informational", http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		res.WriteHeader(http.StatusEarlyHints)
		panic("test panic after informational response")
	}))
	world.Handle("GET /panic-after-write", http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		res.WriteHeader(http.StatusOK)
		_, _ = res.Write([]byte("partial"))
		res.(interface{ Flush() }).Flush()
		res.WriteHeader(http.StatusEarlyHints)
		panic("test panic after commit")
	}))
	world.HandleHello()
	world.Start()

	res, body, err := world.GetBody(t.Context(), world.PathServerURL("http", "panic"), http.Header{})
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, res.StatusCode)
	require.Equal(t, "http: internal server error", body)

	res, body, err = world.GetBody(t.Context(), world.PathServerURL("http", "panic-after-informational"), http.Header{})
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, res.StatusCode)
	require.Equal(t, "http: internal server error", body)

	_, _, err = world.GetBody(t.Context(), world.PathServerURL("http", "panic-after-write"), http.Header{"Accept-Encoding": {"identity"}})
	require.Error(t, err)

	res, body, err = world.GetBody(t.Context(), world.PathServerURL("http", "hello"), http.Header{})
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "hello!", body)
}
