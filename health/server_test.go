package health_test

import (
	"testing"

	"github.com/alexfalkowski/go-health/v2/server"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/health"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-sync"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func TestNewServerStopsWithLifecycle(t *testing.T) {
	checker := newLifecycleChecker()
	reg := server.NewRegistration("lifecycle", time.Millisecond.Duration(), checker)
	lc := fxtest.NewLifecycle(t)
	srv := health.NewServer(lc)

	srv.Register("test", reg)
	require.NoError(t, srv.Observe("test", "healthz", "lifecycle"))

	require.NoError(t, lc.Start(t.Context()))
	checker.waitForRunning(t)

	require.NoError(t, lc.Stop(t.Context()))
	checker.waitForStopped(t)
}

type lifecycleChecker struct {
	running     chan struct{}
	stopped     chan struct{}
	calls       sync.Int64
	runningOnce sync.Once
	stoppedOnce sync.Once
}

func newLifecycleChecker() *lifecycleChecker {
	return &lifecycleChecker{
		running: make(chan struct{}),
		stopped: make(chan struct{}),
	}
}

func (c *lifecycleChecker) Check(ctx context.Context) error {
	if c.calls.Add(1) == 1 {
		return nil
	}

	c.runningOnce.Do(func() { close(c.running) })
	<-ctx.Done()
	c.stoppedOnce.Do(func() { close(c.stopped) })
	return ctx.Err()
}

func (c *lifecycleChecker) waitForRunning(t *testing.T) {
	t.Helper()

	select {
	case <-c.running:
	case <-time.After(time.Second):
		require.Fail(t, "health checker did not run")
	}
}

func (c *lifecycleChecker) waitForStopped(t *testing.T) {
	t.Helper()

	select {
	case <-c.stopped:
	case <-time.After(time.Second):
		require.Fail(t, "health checker was not stopped")
	}
}
