package test

import (
	"github.com/alexfalkowski/go-service/cache/redis"
	rem "github.com/alexfalkowski/go-service/cache/redis/telemetry/metrics"
	"github.com/alexfalkowski/go-service/cache/ristretto"
	rim "github.com/alexfalkowski/go-service/cache/ristretto/telemetry/metrics"
	gr "github.com/alexfalkowski/go-service/redis"
	"github.com/alexfalkowski/go-service/runtime"
	"github.com/go-redis/cache/v8"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// NewRedisCache for test.
func NewRedisCache(lc fx.Lifecycle, cfg *redis.Config, logger *zap.Logger, meter metric.Meter) (*cache.Cache, error) {
	params := redis.OptionsParams{Client: NewRedisClient(lc, cfg, logger), Config: cfg, Marshaller: Marshaller, Compressor: Compressor}

	opts, err := redis.NewOptions(params)
	if err != nil {
		return nil, err
	}

	cache := redis.NewCache(opts)
	rem.Register(cache, Version, meter)

	return cache, nil
}

// NewRedisClient for test.
func NewRedisClient(lc fx.Lifecycle, cfg *redis.Config, logger *zap.Logger) gr.Client {
	tracer := NewTracer(lc, logger)
	client := redis.NewClient(redis.ClientParams{Lifecycle: lc, RingOptions: redis.NewRingOptions(cfg), Tracer: tracer, Logger: logger})

	return client
}

// NewRistrettoCache for test.
func NewRistrettoCache(lc fx.Lifecycle, meter metric.Meter) ristretto.Cache {
	cfg := &ristretto.Config{NumCounters: 1e7, MaxCost: 1 << 30, BufferItems: 64}

	c, err := ristretto.NewCache(ristretto.CacheParams{Lifecycle: lc, Config: cfg, Version: Version})
	runtime.Must(err)

	rim.Register(c, Version, meter)

	return c
}
