package server_test

import (
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/config"
	server "github.com/alexfalkowski/go-service/v2/net/http/server"
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

	resp, err := request(client, t.Context(), srv.String())
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())

	require.NoError(t, srv.Shutdown(context.Background()))
	require.NoError(t, <-errCh)
}

func request(client *http.Client, ctx context.Context, address string) (*http.Response, error) {
	deadline := time.Now().Add(time.Second)
	var err error

	for time.Now().Before(deadline) {
		req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+address, nil)
		if reqErr != nil {
			return nil, reqErr
		}

		resp, respErr := client.Do(req)
		if respErr == nil {
			return resp, nil
		}

		err = respErr
		time.Sleep(10 * time.Millisecond)
	}

	return nil, err
}
