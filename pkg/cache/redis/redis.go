package redis

import (
	"github.com/alexfalkowski/go-service/pkg/cache/redis/metrics/prometheus"
	"github.com/go-redis/cache/v8"
	"go.uber.org/fx"
)

// NewCache from config.
// The cache is based on https://github.com/go-redis/cache
func NewCache(lc fx.Lifecycle, cfg *Config, opts *cache.Options) *cache.Cache {
	cache := cache.New(opts)

	prometheus.Register(lc, cfg.AppName, cache)

	return cache
}
