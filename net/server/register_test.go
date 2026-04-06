package server_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/server"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func TestRegisterStopErrors(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	first := newTestServer(nil, test.ErrFailed)
	secondErr := errors.New("other failed")
	second := newTestServer(nil, secondErr)
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
	require.Equal(t, 1, first.shutdowns)
	require.Equal(t, 1, second.shutdowns)
}
