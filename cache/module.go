package cache

import (
	"github.com/alexfalkowski/go-service/cache/compressor"
	"github.com/alexfalkowski/go-service/cache/marshaller"
	"github.com/alexfalkowski/go-service/cache/redis"
	"github.com/alexfalkowski/go-service/cache/redis/trace/opentracing"
	"github.com/alexfalkowski/go-service/cache/ristretto"
	"go.uber.org/fx"
)

var (
	// RedisModule for fx.
	RedisModule = fx.Options(fx.Provide(redis.NewRing), fx.Provide(redis.NewOptions), fx.Provide(redis.NewCache))

	// RedisOpentracingModule for fx.
	RedisOpentracingModule = fx.Provide(opentracing.NewTracer)

	// RistrettoModule for fx.
	RistrettoModule = fx.Provide(ristretto.NewCache)

	// SnappyCompressorModule for fx.
	SnappyCompressorModule = fx.Provide(compressor.NewSnappy)

	// ProtoMarshallerModule for fx.
	ProtoMarshallerModule = fx.Provide(marshaller.NewProto)
)
