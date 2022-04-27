package redis

import (
	"context"

	"github.com/alexfalkowski/go-service/cache/redis/logger"
	"github.com/go-redis/redis/v8"
	"go.uber.org/fx"
)

// NewRing for redis.
func NewRing(lc fx.Lifecycle, cfg *Config) *redis.Ring {
	redis.SetLogger(logger.NewLogger())

	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server": cfg.Host,
		},
	})

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return ring.Close()
		},
	})

	return ring
}
