package checker

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// NewRedisChecker for health.
func NewRedisChecker(ring *redis.Ring, timeout time.Duration) *RedisChecker {
	return &RedisChecker{ring: ring, timeout: timeout}
}

// RedisChecker for health.
type RedisChecker struct {
	ring    *redis.Ring
	timeout time.Duration
}

// Check redis health.
func (c *RedisChecker) Check(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	return c.ring.Ping(ctx).Err()
}
