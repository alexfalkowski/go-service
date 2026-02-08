package server_test

import (
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/server"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

type testServer struct {
	err  error
	done chan struct{}
}

func newTestServer(err error) *testServer {
	return &testServer{err: err, done: make(chan struct{})}
}

func (s *testServer) Serve() error {
	close(s.done)

	return s.err
}

func (*testServer) Shutdown(_ context.Context) error {
	return nil
}

func (*testServer) String() string {
	return "test"
}

func TestInvalidServer(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	l := test.NewLogger(lc, test.NewJSONLoggerConfig())
	sh := test.NewShutdowner()
	srv := newTestServer(test.ErrFailed)
	svc := server.NewService("test", srv, l, sh)

	svc.Start()

	select {
	case <-srv.done:
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

func TestValidServer(t *testing.T) {
	sh := test.NewShutdowner()
	srv := newTestServer(nil)
	svc := server.NewService("test", srv, nil, sh)

	svc.Start()

	select {
	case <-srv.done:
	case <-time.After(2 * time.Second):
		require.FailNow(t, "timeout waiting for server to serve")
	}

	require.False(t, sh.Called())

	svc.Stop(t.Context())
}
