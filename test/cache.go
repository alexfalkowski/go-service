package test

import (
	"github.com/alexfalkowski/go-service/cache/redis"
	"github.com/alexfalkowski/go-service/cache/redis/trace/opentracing"
	"github.com/alexfalkowski/go-service/compressor"
	"github.com/alexfalkowski/go-service/marshaller"
	"github.com/go-redis/cache/v8"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// NewRedisCache for test.
func NewRedisCache(lc fx.Lifecycle, host string, logger *zap.Logger, compressor compressor.Compressor, marshaller marshaller.Marshaller) *cache.Cache {
	tracer, _ := opentracing.NewTracer(opentracing.TracerParams{Lifecycle: lc, Config: NewJaegerConfig(), Version: Version})
	cfg := &redis.Config{Host: host}
	client := redis.NewClient(redis.ClientParams{Lifecycle: lc, RingOptions: redis.NewRingOptions(cfg), Tracer: tracer, Logger: logger})
	params := redis.OptionsParams{Client: client, Compressor: compressor, Marshaller: marshaller}
	opts := redis.NewOptions(params)

	return redis.NewCache(redis.CacheParams{Lifecycle: lc, Config: cfg, Options: opts, Version: Version})
}
