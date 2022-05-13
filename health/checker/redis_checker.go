package checker

import (
	"context"
	"time"

	"github.com/alexfalkowski/go-service/cache/redis/client"
)

// NewRedisChecker for health.
func NewRedisChecker(client client.Client, timeout time.Duration) *RedisChecker {
	return &RedisChecker{client: client, timeout: timeout}
}

// RedisChecker for health.
type RedisChecker struct {
	client  client.Client
	timeout time.Duration
}

// Check redis health.
func (c *RedisChecker) Check(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	return c.client.Ping(ctx).Err()
}
