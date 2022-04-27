package cache

import (
	"github.com/alexfalkowski/go-service/cache/compressor"
	"github.com/alexfalkowski/go-service/cache/marshaller"
	"github.com/alexfalkowski/go-service/cache/redis"
	"github.com/alexfalkowski/go-service/cache/redis/trace/opentracing/datadog"
	"github.com/alexfalkowski/go-service/cache/redis/trace/opentracing/jaeger"
	"github.com/alexfalkowski/go-service/cache/ristretto"
	"go.uber.org/fx"
)

var (
	// RedisModule for fx.
	// nolint:gochecknoglobals
	RedisModule = fx.Options(fx.Provide(redis.NewRing), fx.Provide(redis.NewOptions), fx.Provide(redis.NewCache))

	// RedisJaegerModule for fx.
	// nolint:gochecknoglobals
	RedisJaegerModule = fx.Provide(jaeger.NewTracer)

	// RedisDataDogModule for fx.
	// nolint:gochecknoglobals
	RedisDataDogModule = fx.Provide(datadog.NewTracer)

	// RistrettoModule for fx.
	// nolint:gochecknoglobals
	RistrettoModule = fx.Provide(ristretto.NewCache)

	// SnappyCompressorModule for fx.
	// nolint:gochecknoglobals
	SnappyCompressorModule = fx.Provide(compressor.NewSnappy)

	// ProtoMarshallerModule for fx.
	// nolint:gochecknoglobals
	ProtoMarshallerModule = fx.Provide(marshaller.NewProto)
)
