package redis

import (
	"context"

	"github.com/alexfalkowski/go-service/pkg/config"
	"github.com/go-redis/redis/v8"
	"go.uber.org/fx"
)

// NewRing for redis.
func NewRing(lc fx.Lifecycle, cfg *config.Config) *redis.Ring {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server": cfg.RedisCacheHost,
		},
	})

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return ring.Close()
		},
	})

	return ring
}
