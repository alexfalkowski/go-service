package test

import (
	"github.com/alexfalkowski/go-service/cache/redis"
	rem "github.com/alexfalkowski/go-service/cache/redis/telemetry/metrics"
	"github.com/alexfalkowski/go-service/cache/ristretto"
	rim "github.com/alexfalkowski/go-service/cache/ristretto/telemetry/metrics"
	gr "github.com/alexfalkowski/go-service/redis"
	sr "github.com/alexfalkowski/go-service/ristretto"
	"github.com/alexfalkowski/go-service/runtime"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/go-redis/cache/v8"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Cache for test.
type Cache struct {
	Lifecycle fx.Lifecycle
	Redis     *redis.Config
	Logger    *zap.Logger
	Tracer    *tracer.Config
	Meter     metric.Meter
}

// NewRedisCache for test.
func (c *Cache) NewRedisCache() (*cache.Cache, error) {
	params := redis.OptionsParams{Client: c.NewRedisClient(), Config: c.Redis, Marshaller: Marshaller, Compressor: Compressor}

	opts, err := redis.NewOptions(params)
	if err != nil {
		return nil, err
	}

	cache := redis.NewCache(opts)
	rem.Register(cache, c.Meter)

	return cache, nil
}

// NewRedisClient for test.
func (c *Cache) NewRedisClient() gr.Client {
	tracer := tracer.NewTracer(c.Lifecycle, Environment, Version, c.Tracer, c.Logger)
	client := redis.NewClient(redis.ClientParams{Lifecycle: c.Lifecycle, RingOptions: redis.NewRingOptions(c.Redis), Tracer: tracer, Logger: c.Logger})

	return client
}

// NewRistrettoCache for test.
func (c *Cache) NewRistrettoCache() sr.Cache {
	cfg := &ristretto.Config{NumCounters: 1e7, MaxCost: 1 << 30, BufferItems: 64}

	ca, err := ristretto.NewCache(ristretto.CacheParams{Lifecycle: c.Lifecycle, Config: cfg, Version: Version})
	runtime.Must(err)

	rim.Register(ca, c.Meter)

	return ca
}
