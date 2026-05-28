package server_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
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

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Serve()
	}()

	conn, err := test.Connect(t.Context(), srv.String())
	require.NoError(t, err)
	require.NoError(t, conn.Close())

	require.NoError(t, srv.Shutdown(context.Background()))
	require.NoError(t, <-errCh)
}

func TestShutdownClosesUnservedListener(t *testing.T) {
	srv, err := server.NewServer(grpc.NewServer(test.ConfigOptions, time.Second), &config.Config{Address: ":0"})
	require.NoError(t, err)

	addr := srv.String()
	require.NoError(t, srv.Shutdown(context.Background()))

	conn, err := test.Connect(t.Context(), addr)
	require.Error(t, err)
	require.Nil(t, conn)
}

func TestShutdownStopsServerWhenContextCanceled(t *testing.T) {
	srv, blocking := newServerWithBlockingGreeter(t)
	serveErrCh := serve(srv)
	rpcErrCh := startBlockingSayHello(t, srv, blocking)

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	require.ErrorIs(t, srv.Shutdown(ctx), context.Canceled)
	require.Error(t, <-rpcErrCh)
	require.NoError(t, <-serveErrCh)
}

func TestNewServerWithInvalidNetwork(t *testing.T) {
	srv, err := server.NewServer(grpc.NewServer(test.ConfigOptions, time.Second), &config.Config{Address: "invalid://:0"})
	require.Error(t, err)
	require.Nil(t, srv)
}

func newServerWithBlockingGreeter(t *testing.T) (*server.Server, *blockingService) {
	t.Helper()

	grpcServer := grpc.NewServer(test.ConfigOptions, time.Second)
	blocking := newBlockingService()
	v1.RegisterGreeterServiceServer(grpcServer, blocking)

	srv, err := server.NewServer(grpcServer, &config.Config{Address: ":0"})
	require.NoError(t, err)

	return srv, blocking
}

func serve(srv *server.Server) <-chan error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Serve()
	}()

	return errCh
}

func startBlockingSayHello(t *testing.T, srv *server.Server, blocking *blockingService) <-chan error {
	t.Helper()

	conn, err := grpc.NewClient(srv.String(), grpc.WithTransportCredentials(grpc.NewInsecureCredentials()))
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, conn.Close())
	})

	rpcCtx, cancelRPC := context.WithCancel(t.Context())
	t.Cleanup(cancelRPC)

	errCh := make(chan error, 1)
	go func() {
		_, err := v1.NewGreeterServiceClient(conn).SayHello(rpcCtx, &v1.SayHelloRequest{Name: "test"})
		errCh <- err
	}()

	select {
	case <-blocking.started:
	case <-time.After(time.Second):
		cancelRPC()
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_ = srv.Shutdown(ctx)
		require.FailNow(t, "timeout waiting for blocking gRPC method to start")
	}

	return errCh
}

type blockingService struct {
	v1.UnimplementedGreeterServiceServer

	started chan struct{}
}

func newBlockingService() *blockingService {
	return &blockingService{started: make(chan struct{})}
}

func (s *blockingService) SayHello(ctx context.Context, _ *v1.SayHelloRequest) (*v1.SayHelloResponse, error) {
	close(s.started)
	<-ctx.Done()

	return nil, ctx.Err()
}
