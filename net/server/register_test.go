package server_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/server"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func TestRegisterStopErrors(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	first := test.NewObservableServer(nil, test.ErrFailed)
	secondErr := errors.New("other failed")
	second := test.NewObservableServer(nil, secondErr)
	sh := test.NewShutdowner()

	server.Register(lc, []*server.Service{
		server.NewService("first", first, nil, sh),
		server.NewService("second", second, nil, sh),
	})

	require.NoError(t, lc.Start(t.Context()))

	err := lc.Stop(t.Context())
	require.Error(t, err)
	require.ErrorIs(t, err, test.ErrFailed)
	require.ErrorIs(t, err, secondErr)
	require.ErrorContains(t, err, "first: failed")
	require.ErrorContains(t, err, "second: other failed")
	require.Equal(t, 1, first.Shutdowns, "first shutdowns")
	require.Equal(t, 1, second.Shutdowns, "second shutdowns")
}

func TestRegisterStartsAllServices(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	first := test.NewObservableServer(nil, nil)
	second := test.NewObservableServer(nil, nil)
	sh := test.NewShutdowner()

	server.Register(lc, []*server.Service{
		server.NewService("first", first, nil, sh),
		server.NewService("second", second, nil, sh),
	})

	require.NoError(t, lc.Start(t.Context()))
	waitForRegisteredService(t, first.Done)
	waitForRegisteredService(t, second.Done)

	require.NoError(t, lc.Stop(t.Context()))
	require.Equal(t, 1, first.Shutdowns, "first shutdowns")
	require.Equal(t, 1, second.Shutdowns, "second shutdowns")
}

func TestRegisterStopsServicesConcurrently(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	first := newBlockingShutdownServer()
	second := newShutdownContextServer()
	sh := test.NewShutdowner()

	server.Register(lc, []*server.Service{
		server.NewService("first", first, nil, sh),
		server.NewService("second", second, nil, sh),
	})

	require.NoError(t, lc.Start(t.Context()))
	waitForRegisteredService(t, first.Done)
	waitForRegisteredService(t, second.Done)

	ctx, cancel := context.WithTimeout(t.Context(), 25*time.Millisecond)
	defer cancel()

	err := lc.Stop(ctx)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.NoError(t, <-second.ShutdownErrs)
}

func TestRegisterSnapshotsServices(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	original := test.NewObservableServer(nil, nil)
	replacement := test.NewObservableServer(nil, nil)
	sh := test.NewShutdowner()
	services := []*server.Service{
		server.NewService("original", original, nil, sh),
	}

	server.Register(lc, services)
	services[0] = server.NewService("replacement", replacement, nil, sh)

	require.NoError(t, lc.Start(t.Context()))
	waitForRegisteredService(t, original.Done)

	require.NoError(t, lc.Stop(t.Context()))
	require.Equal(t, 1, original.Shutdowns, "original shutdowns")
	require.Equal(t, 0, replacement.Shutdowns, "replacement shutdowns")
}

func waitForRegisteredService(t *testing.T, done <-chan struct{}) {
	t.Helper()

	select {
	case <-done:
	case <-time.After(time.Second):
		require.Fail(t, "registered service was not started")
	}
}

func newBlockingShutdownServer() *blockingShutdownServer {
	return &blockingShutdownServer{Done: make(chan struct{})}
}

type blockingShutdownServer struct {
	Done chan struct{}
}

func (s *blockingShutdownServer) Serve() error {
	close(s.Done)

	return nil
}

func (*blockingShutdownServer) Shutdown(ctx context.Context) error {
	<-ctx.Done()

	return ctx.Err()
}

func (*blockingShutdownServer) String() string {
	return "test"
}

func newShutdownContextServer() *shutdownContextServer {
	return &shutdownContextServer{
		Done:         make(chan struct{}),
		ShutdownErrs: make(chan error, 1),
	}
}

type shutdownContextServer struct {
	Done         chan struct{}
	ShutdownErrs chan error
}

func (s *shutdownContextServer) Serve() error {
	close(s.Done)

	return nil
}

func (s *shutdownContextServer) Shutdown(ctx context.Context) error {
	s.ShutdownErrs <- ctx.Err()

	return nil
}

func (*shutdownContextServer) String() string {
	return "test"
}
