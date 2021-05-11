package redis

import (
	"github.com/alexfalkowski/go-service/pkg/cache/redis/metrics/prometheus"
	"github.com/alexfalkowski/go-service/pkg/config"
	"github.com/go-redis/cache/v8"
)

// NewCache from config.
// The cache is based on https://github.com/go-redis/cache
func NewCache(cfg *config.Config, opts *cache.Options) (*cache.Cache, error) {
	cache := cache.New(opts)

	if err := prometheus.Register(cfg, cache); err != nil {
		return nil, err
	}

	return cache, nil
}
