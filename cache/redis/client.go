package redis

import (
	"context"

	"github.com/alexfalkowski/go-service/cache/redis/client"
	"github.com/alexfalkowski/go-service/cache/redis/logger"
	"github.com/go-redis/redis/v8"
	"go.uber.org/fx"
)

// RingParams for redis.
type RingParams struct {
	fx.In

	Lifecycle   fx.Lifecycle
	RingOptions *redis.RingOptions
}

// NewClient for redis.
func NewClient(params RingParams) client.Client {
	redis.SetLogger(logger.NewLogger())

	ring := redis.NewRing(params.RingOptions)

	params.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return ring.Close()
		},
	})

	return ring
}
