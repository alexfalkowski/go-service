package test

import (
	"github.com/alexfalkowski/go-service/cache/redis"
	"github.com/alexfalkowski/go-service/compressor"
	"github.com/alexfalkowski/go-service/marshaller"
	"github.com/go-redis/cache/v8"
	"go.uber.org/fx"
)

// NewRedisCache for test.
func NewRedisCache(lc fx.Lifecycle, host string, compressor compressor.Compressor, marshaller marshaller.Marshaller) *cache.Cache {
	cfg := &redis.Config{Host: host}
	r := redis.NewRing(lc, cfg)
	params := redis.OptionsParams{Ring: r, Compressor: compressor, Marshaller: marshaller}
	opts := redis.NewOptions(params)

	return redis.NewCache(redis.CacheParams{Lifecycle: lc, Config: cfg, Options: opts, Version: Version})
}
