package test

import (
	"github.com/alexfalkowski/go-service/cache/redis"
	"github.com/alexfalkowski/go-service/cache/redis/client"
	rem "github.com/alexfalkowski/go-service/cache/redis/telemetry/metrics"
	"github.com/alexfalkowski/go-service/cache/redis/telemetry/tracer"
	cristretto "github.com/alexfalkowski/go-service/cache/ristretto"
	rim "github.com/alexfalkowski/go-service/cache/ristretto/telemetry/metrics"
	"github.com/alexfalkowski/go-service/compressor"
	"github.com/alexfalkowski/go-service/marshaller"
	"github.com/dgraph-io/ristretto"
	"github.com/go-redis/cache/v8"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// NewRedisCache for test.
func NewRedisCache(lc fx.Lifecycle, host string, logger *zap.Logger, compressor compressor.Compressor, marshaller marshaller.Marshaller, meter metric.Meter) *cache.Cache {
	params := redis.OptionsParams{Client: NewRedisClient(lc, host, logger), Compressor: compressor, Marshaller: marshaller}
	opts := redis.NewOptions(params)
	cache := redis.NewCache(redis.CacheParams{Lifecycle: lc, Config: NewRedisConfig(host), Options: opts, Version: Version})

	rem.Register(cache, Version, meter)

	return redis.NewCache(redis.CacheParams{Lifecycle: lc, Config: NewRedisConfig(host), Options: opts, Version: Version})
}

// NewRedisClient for test.
func NewRedisClient(lc fx.Lifecycle, host string, logger *zap.Logger) client.Client {
	tracer, _ := tracer.NewTracer(tracer.Params{Lifecycle: lc, Config: NewTracerConfig(), Version: Version})
	client := redis.NewClient(redis.ClientParams{Lifecycle: lc, RingOptions: redis.NewRingOptions(NewRedisConfig(host)), Tracer: tracer, Logger: logger})

	return client
}

// NewRedisConfig for test.
func NewRedisConfig(host string) *redis.Config {
	return &redis.Config{Addresses: map[string]string{"server": host}}
}

// NewRistrettoCache for test.
func NewRistrettoCache(lc fx.Lifecycle, meter metric.Meter) *ristretto.Cache {
	cfg := &cristretto.Config{NumCounters: 1e7, MaxCost: 1 << 30, BufferItems: 64}
	c, _ := cristretto.NewCache(cristretto.CacheParams{Lifecycle: lc, Config: cfg, Version: Version})

	rim.Register(c, Version, meter)

	return c
}
