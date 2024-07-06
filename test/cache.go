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
	"github.com/go-redis/cache/v9"
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
	cl, err := c.NewRedisClient()
	if err != nil {
		return nil, err
	}

	params := redis.OptionsParams{Client: cl, Config: c.Redis, Marshaller: Marshaller, Map: Compressor}

	opts, err := redis.NewOptions(params)
	if err != nil {
		return nil, err
	}

	cache := redis.NewCache(opts)
	rem.Register(cache, c.Meter)

	return cache, nil
}

// NewRedisClient for test.
func (c *Cache) NewRedisClient() (gr.Client, error) {
	tracer, err := tracer.NewTracer(c.Lifecycle, Environment, Version, Name, c.Tracer, c.Logger)
	if err != nil {
		return nil, err
	}

	opts, err := redis.NewRingOptions(c.Redis)
	runtime.Must(err)

	client := redis.NewClient(redis.ClientParams{Lifecycle: c.Lifecycle, RingOptions: opts, Tracer: tracer, Logger: c.Logger})

	return client, nil
}

// NewRistrettoCache for test.
func (c *Cache) NewRistrettoCache() sr.Cache {
	cfg := &ristretto.Config{NumCounters: 1e7, MaxCost: 1 << 30, BufferItems: 64}

	ca, err := ristretto.NewCache(ristretto.CacheParams{Lifecycle: c.Lifecycle, Config: cfg})
	runtime.Must(err)

	rim.Register(ca, c.Meter)

	return ca
}
