package checker

import (
	"fmt"

	"github.com/alexfalkowski/go-health/v2/checker"
	"github.com/alexfalkowski/go-service/v2/cache"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-sync"
)

var _ checker.Checker = (*CacheChecker)(nil)

// ErrNoCachePinger is returned when a CacheChecker has no pingable cache driver to verify.
var ErrNoCachePinger = errors.New("cache: no pingable driver")

// ErrCachePingTimeout is the cause recorded when a cache health-check ping times out.
var ErrCachePingTimeout = fmt.Errorf("cache: ping timeout: %w", sync.ErrTimeout)

// NewCacheChecker constructs a CacheChecker that verifies cache backend connectivity.
//
// The timeout is applied per Ping invocation. Nil pingers are treated as unavailable for this checker;
// callers should register the checker only for backends where connectivity should affect readiness.
func NewCacheChecker(pinger cache.Pinger, timeout time.Duration) *CacheChecker {
	return &CacheChecker{pinger: pinger, timeout: timeout}
}

// CacheChecker is a health checker that verifies cache backend connectivity.
type CacheChecker struct {
	pinger  cache.Pinger
	timeout time.Duration
}

// Check verifies cache health by pinging the configured cache backend.
//
// If no pingable driver is configured, Check returns ErrNoCachePinger. If the ping times out, Check
// returns ErrCachePingTimeout.
func (c *CacheChecker) Check(ctx context.Context) error {
	if c.pinger == nil {
		return ErrNoCachePinger
	}

	ctx, cancel := context.WithTimeoutCause(ctx, c.timeout, ErrCachePingTimeout)
	defer cancel()

	err := c.pinger.Ping(ctx)
	if err != nil && ctx.Err() != nil {
		return context.Cause(ctx)
	}

	return err
}
