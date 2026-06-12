package checker_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/health/checker"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

func TestCacheCheckerWithoutPingableDriver(t *testing.T) {
	check := checker.NewCacheChecker(nil, time.Second)

	require.ErrorIs(t, check.Check(t.Context()), checker.ErrNoCachePinger)
}

func TestCacheCheckerPingsDriver(t *testing.T) {
	check := checker.NewCacheChecker(cachePinger{}, time.Second)

	require.NoError(t, check.Check(t.Context()))
}

func TestCacheCheckerReturnsPingError(t *testing.T) {
	check := checker.NewCacheChecker(cachePinger{err: errCachePing}, time.Second)

	require.ErrorIs(t, check.Check(t.Context()), errCachePing)
}

func TestCacheCheckerReturnsTimeoutCause(t *testing.T) {
	check := checker.NewCacheChecker(cachePinger{wait: true}, time.Millisecond)

	require.ErrorIs(t, check.Check(t.Context()), checker.ErrCachePingTimeout)
}

var errCachePing = errors.New("cache ping")

type cachePinger struct {
	err  error
	wait bool
}

func (c cachePinger) Ping(ctx context.Context) error {
	if c.wait {
		<-ctx.Done()

		return context.Cause(ctx)
	}

	return c.err
}
