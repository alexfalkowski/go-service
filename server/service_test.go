package server_test

import (
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/server"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func TestInvalidServer(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	l := test.NewLogger(lc, test.NewJSONLoggerConfig())
	sh := test.NewShutdowner()
	srv := &test.ErrServer{}
	svc := server.NewService("test", srv, l, sh)

	svc.Start()
	time.Sleep(1 * time.Second)
	require.True(t, sh.Called())
}

func TestValidServer(t *testing.T) {
	sh := test.NewShutdowner()
	srv := &test.NoopServer{}
	svc := server.NewService("test", srv, nil, sh)

	svc.Start()
	time.Sleep(1 * time.Second)
	require.False(t, sh.Called())

	svc.Stop(t.Context())
}
