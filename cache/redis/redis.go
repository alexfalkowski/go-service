package redis

import (
	"github.com/alexfalkowski/go-service/cache/redis/metrics/prometheus"
	"github.com/alexfalkowski/go-service/os"
	"github.com/go-redis/cache/v8"
	"go.uber.org/fx"
)

// NewCache from config.
// The cache is based on https://github.com/go-redis/cache
func NewCache(lc fx.Lifecycle, cfg *Config, opts *cache.Options) (*cache.Cache, error) {
	cache := cache.New(opts)

	name, err := os.ExecutableName()
	if err != nil {
		return nil, err
	}

	prometheus.Register(lc, name, cache)

	return cache, nil
}
