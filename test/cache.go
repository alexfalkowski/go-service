package test

import (
	"github.com/alexfalkowski/go-service/cache/redis"
	"github.com/alexfalkowski/go-service/cache/redis/client"
	"github.com/alexfalkowski/go-service/cache/redis/otel"
	cristretto "github.com/alexfalkowski/go-service/cache/ristretto"
	"github.com/alexfalkowski/go-service/compressor"
	"github.com/alexfalkowski/go-service/marshaller"
	"github.com/dgraph-io/ristretto"
	"github.com/go-redis/cache/v8"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// NewRedisCache for test.
func NewRedisCache(lc fx.Lifecycle, host string, logger *zap.Logger, compressor compressor.Compressor, marshaller marshaller.Marshaller) *cache.Cache {
	params := redis.OptionsParams{Client: NewRedisClient(lc, host, logger), Compressor: compressor, Marshaller: marshaller}
	opts := redis.NewOptions(params)

	return redis.NewCache(redis.CacheParams{Lifecycle: lc, Config: NewRedisConfig(host), Options: opts, Version: Version})
}

// NewRedisClient for test.
func NewRedisClient(lc fx.Lifecycle, host string, logger *zap.Logger) client.Client {
	tracer, _ := otel.NewTracer(otel.TracerParams{Lifecycle: lc, Config: NewOTELConfig(), Version: Version})
	client := redis.NewClient(redis.ClientParams{Lifecycle: lc, RingOptions: redis.NewRingOptions(NewRedisConfig(host)), Tracer: tracer, Logger: logger})

	return client
}

// NewRedisConfig for test.
func NewRedisConfig(host string) *redis.Config {
	return &redis.Config{Addresses: map[string]string{"server": host}}
}

// NewRistrettoCache for test.
func NewRistrettoCache(lc fx.Lifecycle) *ristretto.Cache {
	cfg := &cristretto.Config{NumCounters: 1e7, MaxCost: 1 << 30, BufferItems: 64}
	c, _ := cristretto.NewCache(cristretto.CacheParams{Lifecycle: lc, Config: cfg, Version: Version})

	return c
}
