package redis

import (
	"context"

	lzap "github.com/alexfalkowski/go-service/cache/redis/logger/zap"
	"github.com/go-redis/redis/v8"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// NewRing for redis.
func NewRing(lc fx.Lifecycle, cfg *Config, logger *zap.Logger) *redis.Ring {
	redis.SetLogger(lzap.NewLogger(logger))

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
