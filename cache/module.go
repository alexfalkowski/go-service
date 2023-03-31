package cache

import (
	"github.com/alexfalkowski/go-service/cache/compressor"
	"github.com/alexfalkowski/go-service/cache/marshaller"
	"github.com/alexfalkowski/go-service/cache/redis"
	"github.com/alexfalkowski/go-service/cache/redis/otel"
	"github.com/alexfalkowski/go-service/cache/ristretto"
	"go.uber.org/fx"
)

var (
	// RedisModule for fx.
	RedisModule = fx.Options(fx.Provide(redis.NewClient), fx.Provide(redis.NewOptions), fx.Provide(redis.NewCache), fx.Provide(redis.NewRingOptions))

	// RedisOTELModule for fx.
	RedisOTELModule = fx.Provide(otel.NewTracer)

	// RistrettoModule for fx.
	RistrettoModule = fx.Provide(ristretto.NewCache)

	// SnappyCompressorModule for fx.
	SnappyCompressorModule = fx.Provide(compressor.NewSnappy)

	// ProtoMarshallerModule for fx.
	ProtoMarshallerModule = fx.Provide(marshaller.NewProto)
)
