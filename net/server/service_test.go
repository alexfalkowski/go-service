package server_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/server"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-sync"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func TestStartRequestsShutdownOnServeError(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	l, err := test.NewLogger(lc, test.NewJSONLoggerConfig())
	require.NoError(t, err)
	sh := test.NewShutdowner()
	srv := test.NewObservableServer(test.ErrFailed, nil)
	svc := server.NewService("test", srv, l, sh)

	svc.Start()

	select {
	case <-srv.Done:
	case <-time.After(2 * time.Second):
		require.FailNow(t, "timeout waiting for server to serve")
	}

	select {
	case <-sh.Done():
	case <-time.After(2 * time.Second):
		require.FailNow(t, "timeout waiting for shutdown")
	}

	require.True(t, sh.Called())
}

func TestStartDoesNotShutdownOnNilServeReturn(t *testing.T) {
	sh := test.NewShutdowner()
	srv := test.NewObservableServer(nil, nil)
	svc := server.NewService("test", srv, nil, sh)

	svc.Start()

	select {
	case <-srv.Done:
	case <-time.After(2 * time.Second):
		require.FailNow(t, "timeout waiting for server to serve")
	}

	require.False(t, sh.Called())

	require.NoError(t, svc.Stop(t.Context()))
}

func TestStartReturnsWhileServeIsRunning(t *testing.T) {
	sh := test.NewShutdowner()
	srv := newBlockingServer()
	svc := server.NewService("test", srv, nil, sh)

	startDone := startService(svc)

	requireStartReturned(t, startDone, srv)
	requireServeStarted(t, srv)

	require.NoError(t, svc.Stop(t.Context()))
	require.False(t, sh.Called())
}

func TestInvalidStop(t *testing.T) {
	sh := test.NewShutdowner()
	srv := test.NewObservableServer(nil, test.ErrFailed)
	svc := server.NewService("test", srv, nil, sh)

	err := svc.Stop(t.Context())
	require.EqualError(t, err, "test: failed")
	require.ErrorIs(t, err, test.ErrFailed)
}

func TestStopPassesContextToServer(t *testing.T) {
	sh := test.NewShutdowner()
	srv := &contextServer{}
	svc := server.NewService("test", srv, nil, sh)
	ctx := context.WithValue(t.Context(), contextServerKey{}, "caller-value")

	require.NoError(t, svc.Stop(ctx))

	require.Equal(t, "caller-value", srv.ctx.Value(contextServerKey{}))
}

type blockingServer struct {
	started chan struct{}
	stop    chan struct{}
	once    sync.Once
}

func newBlockingServer() *blockingServer {
	return &blockingServer{
		started: make(chan struct{}),
		stop:    make(chan struct{}),
	}
}

func (s *blockingServer) Serve() error {
	close(s.started)
	<-s.stop

	return nil
}

func (s *blockingServer) Shutdown(context.Context) error {
	s.release()

	return nil
}

func (*blockingServer) String() string {
	return "test"
}

func (s *blockingServer) release() {
	s.once.Do(func() { close(s.stop) })
}

func startService(svc *server.Service) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		svc.Start()
		close(done)
	}()

	return done
}

func requireStartReturned(t *testing.T, done <-chan struct{}, srv *blockingServer) {
	t.Helper()

	select {
	case <-done:
	case <-time.After(time.Second):
		srv.release()
		require.FailNow(t, "timeout waiting for Start to return")
	}
}

func requireServeStarted(t *testing.T, srv *blockingServer) {
	t.Helper()

	select {
	case <-srv.started:
	case <-time.After(time.Second):
		srv.release()
		require.FailNow(t, "timeout waiting for server to serve")
	}
}

type contextServerKey struct{}

type contextServer struct {
	ctx context.Context
}

func (*contextServer) Serve() error {
	return nil
}

func (s *contextServer) Shutdown(ctx context.Context) error {
	s.ctx = ctx

	return nil
}

func (*contextServer) String() string {
	return "test"
}
