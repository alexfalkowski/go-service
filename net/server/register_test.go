package server_test

import (
	"testing"

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
	require.Equal(t, 1, first.Shutdowns)
	require.Equal(t, 1, second.Shutdowns)
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
	require.Equal(t, 1, first.Shutdowns)
	require.Equal(t, 1, second.Shutdowns)
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
	select {
	case <-original.Done:
	case <-time.After(time.Second):
		require.Fail(t, "registered service was not started")
	}

	require.NoError(t, lc.Stop(t.Context()))
	require.Equal(t, 1, original.Shutdowns)
	require.Equal(t, 0, replacement.Shutdowns)
}

func waitForRegisteredService(t *testing.T, done <-chan struct{}) {
	t.Helper()

	select {
	case <-done:
	case <-time.After(time.Second):
		require.Fail(t, "registered service was not started")
	}
}
