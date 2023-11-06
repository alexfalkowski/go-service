package checker

import (
	"context"
	"time"

	gr "github.com/alexfalkowski/go-service/redis"
)

// NewRedisChecker for health.
func NewRedisChecker(client gr.Client, timeout time.Duration) *RedisChecker {
	return &RedisChecker{client: client, timeout: timeout}
}

// RedisChecker for health.
type RedisChecker struct {
	client  gr.Client
	timeout time.Duration
}

// Check redis health.
func (c *RedisChecker) Check(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	return c.client.Ping(ctx).Err()
}
