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

func TestShutdownClosesUnservedListener(t *testing.T) {
	srv, err := server.NewServer(&http.Server{Handler: http.NewServeMux()}, &config.Config{Address: ":0"})
	require.NoError(t, err)

	addr := srv.String()
	require.NoError(t, srv.Shutdown(context.Background()))

	conn, err := test.Connect(t.Context(), addr)
	require.Error(t, err)
	require.Nil(t, conn)
}

func TestShutdownClosesActiveConnectionsWhenContextCanceled(t *testing.T) {
	started := make(chan struct{})

	// Bind loopback rather than the wildcard address: [net/http.ProxyFromEnvironment] bypasses a configured
	// proxy only for localhost and loopback hosts, and a proxy between this client and the server answers
	// the severed connection with its own 502 instead of the transport error asserted below.
	srv, err := server.NewServer(&http.Server{Handler: http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
		close(started)
		<-req.Context().Done()
	})}, &config.Config{Address: test.RandomHost()})
	require.NoError(t, err)

	serveErr := make(chan error, 1)
	go func() {
		serveErr <- srv.Serve()
	}()

	requestCtx, cancelRequest := context.WithCancel(t.Context())
	t.Cleanup(cancelRequest)
	req, err := http.NewRequestWithContext(requestCtx, http.MethodGet, "http://"+srv.String(), nil)
	require.NoError(t, err)
	requestErr := make(chan error, 1)
	go func() {
		res, err := http.NewClient(http.DefaultTransport, time.Second*5).Do(req)
		if res != nil {
			_ = res.Body.Close()
		}
		requestErr <- err
	}()

	select {
	case <-started:
	case <-time.After(time.Second):
		require.FailNow(t, "timed out waiting for request handler to start")
	}

	shutdownCtx, cancelShutdown := context.WithCancel(t.Context())
	cancelShutdown()

	require.ErrorIs(t, srv.Shutdown(shutdownCtx), context.Canceled)

	select {
	case err := <-requestErr:
		require.Error(t, err)
	case <-time.After(time.Second):
		require.FailNow(t, "timed out waiting for shutdown to close the active request")
	}

	require.NoError(t, <-serveErr)
}

func TestNewServerWithInvalidNetwork(t *testing.T) {
	srv, err := server.NewServer(&http.Server{Handler: http.NewServeMux()}, &config.Config{Address: "invalid://:0"})
	require.Error(t, err)
	require.Nil(t, srv)
}
