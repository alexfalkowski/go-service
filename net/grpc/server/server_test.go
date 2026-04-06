package server_test

import (
	"net"
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/config"
	"github.com/alexfalkowski/go-service/v2/net/grpc/server"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

func TestNewServerWithRawAddress(t *testing.T) {
	srv, err := server.NewServer(grpc.NewServer(test.ConfigOptions, time.Second), &config.Config{Address: ":0"})
	require.NoError(t, err)
	require.NotEmpty(t, srv.String())

	dialer := &net.Dialer{}
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Serve()
	}()

	conn, err := connect(t.Context(), dialer, srv.String())
	require.NoError(t, err)
	require.NoError(t, conn.Close())

	require.NoError(t, srv.Shutdown(context.Background()))
	require.NoError(t, <-errCh)
}

func connect(ctx context.Context, dialer *net.Dialer, address string) (net.Conn, error) {
	deadline := time.Now().Add(time.Second)
	var err error

	for time.Now().Before(deadline) {
		conn, dialErr := dialer.DialContext(ctx, "tcp", address)
		if dialErr == nil {
			return conn, nil
		}

		err = dialErr
		time.Sleep(10 * time.Millisecond)
	}

	return nil, err
}
