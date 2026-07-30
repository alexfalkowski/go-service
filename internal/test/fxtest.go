package test

import (
	"testing"

	"go.uber.org/fx/fxtest"
)

// QuietLifecycle returns an [fxtest.Lifecycle] equivalent to [fxtest.NewLifecycle], except that its
// default per-hook logging is discarded, so a benchmark or test run isn't dominated by "[Fx] HOOK
// OnStart"/"OnStop" lines for every world or provider it starts. Errorf and FailNow are still
// forwarded to tb, so a lifecycle that fails to start or stop cleanly still fails the test exactly as
// fxtest.NewLifecycle would.
func QuietLifecycle(tb testing.TB) *fxtest.Lifecycle {
	tb.Helper()

	return fxtest.NewLifecycle(quietTB{TB: tb})
}

// quietTB adapts a testing.TB to fxtest.TB, discarding Logf calls.
type quietTB struct {
	testing.TB
}

// Logf discards the message.
func (quietTB) Logf(string, ...any) {}
