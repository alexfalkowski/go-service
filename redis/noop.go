package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// NewNoopClient for redis.
func NewNoopClient() *NoopClient {
	return &NoopClient{}
}

// NoopClient for redis.
type NoopClient struct{}

func (*NoopClient) Set(ctx context.Context, _ string, _ any, _ time.Duration) *redis.StatusCmd {
	return redis.NewStatusCmd(ctx)
}

func (*NoopClient) SetXX(ctx context.Context, _ string, _ any, _ time.Duration) *redis.BoolCmd {
	return redis.NewBoolCmd(ctx)
}

func (*NoopClient) SetNX(ctx context.Context, _ string, _ any, _ time.Duration) *redis.BoolCmd {
	return redis.NewBoolCmd(ctx)
}

func (*NoopClient) Get(ctx context.Context, _ string) *redis.StringCmd {
	return redis.NewStringCmd(ctx)
}

func (*NoopClient) Del(ctx context.Context, _ ...string) *redis.IntCmd {
	return redis.NewIntCmd(ctx)
}

func (*NoopClient) Incr(ctx context.Context, _ string) *redis.IntCmd {
	return redis.NewIntCmd(ctx)
}

func (*NoopClient) Ping(ctx context.Context) *redis.StatusCmd {
	return redis.NewStatusCmd(ctx)
}

func (*NoopClient) Close() error {
	return nil
}
