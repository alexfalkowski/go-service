package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Client for redis.
type Client interface {
	Set(ctx context.Context, key string, value any, ttl time.Duration) *redis.StatusCmd
	SetXX(ctx context.Context, key string, value any, ttl time.Duration) *redis.BoolCmd
	SetNX(ctx context.Context, key string, value any, ttl time.Duration) *redis.BoolCmd

	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd

	Incr(ctx context.Context, key string) *redis.IntCmd

	Ping(ctx context.Context) *redis.StatusCmd

	Close() error
}
