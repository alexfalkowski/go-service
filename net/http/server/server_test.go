package server_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/config"
	server "github.com/alexfalkowski/go-service/v2/net/http/server"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

func TestNewServerWithRawAddress(t *testing.T) {
	srv, err := server.NewServer(&http.Server{Handler: http.NewServeMux()}, &config.Config{Address: ":0"})
	require.NoError(t, err)
	require.NotEmpty(t, srv.String())

	client := http.NewClient(http.DefaultTransport, time.Second)
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Serve()
	}()

	conn, err := test.Connect(t.Context(), srv.String())
	require.NoError(t, err)
	require.NoError(t, conn.Close())

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://"+srv.String(), nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())

	require.NoError(t, srv.Shutdown(context.Background()))
	require.NoError(t, <-errCh)
}
