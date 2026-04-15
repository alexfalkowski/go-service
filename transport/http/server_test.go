package http_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/config/server"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/mime"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/transport/http"
	"github.com/stretchr/testify/require"
)

func TestInvalidServer(t *testing.T) {
	http.Register(test.FS)

	cfg := &http.Config{
		Config: &server.Config{
			Timeout: 5 * time.Second,
			TLS:     test.NewTLSConfig("certs/client-cert.pem", "secrets/none"),
		},
	}
	params := http.ServerParams{
		Shutdowner: test.NewShutdowner(),
		Config:     cfg,
	}

	_, err := http.NewServer(params)
	require.Error(t, err)
}

func TestServerMaxReceiveSize(t *testing.T) {
	cfg := test.NewInsecureTransportConfig()
	cfg.HTTP.MaxReceiveSize = 64

	world := test.NewWorld(t, test.WithWorldTransportConfig(cfg), test.WithWorldHTTP())
	http.Handle(world.ServeMux, "POST /hello", content.NewRequestHandler(test.Content, func(_ context.Context, _ *test.Request) (*test.Response, error) {
		return &test.Response{Greeting: "hello"}, nil
	}))
	world.Start()

	header := http.Header{}
	header.Set(content.TypeKey, mime.JSONMediaType)

	res, body, err := world.PostBody(
		t.Context(),
		world.PathServerURL("http", "hello"),
		header,
		strings.NewReader(`{"name":"`+strings.Repeat("a", 256)+`"}`),
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusRequestEntityTooLarge, res.StatusCode)
	require.Contains(t, body, "http: request body too large")
}
