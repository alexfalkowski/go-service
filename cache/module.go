package cache

import (
	"github.com/alexfalkowski/go-service/cache/compressor"
	"github.com/alexfalkowski/go-service/cache/marshaller"
	"github.com/alexfalkowski/go-service/cache/redis"
	rem "github.com/alexfalkowski/go-service/cache/redis/telemetry/metrics"
	"github.com/alexfalkowski/go-service/cache/redis/telemetry/tracer"
	"github.com/alexfalkowski/go-service/cache/ristretto"
	rim "github.com/alexfalkowski/go-service/cache/ristretto/telemetry/metrics"
	"go.uber.org/fx"
)

var (
	// RedisModule for fx.
	RedisModule = fx.Options(
		fx.Provide(redis.NewClient),
		fx.Provide(redis.NewOptions),
		fx.Provide(redis.NewCache),
		fx.Provide(redis.NewRingOptions),
		fx.Provide(tracer.NewTracer),
		fx.Invoke(rem.Register),
	)

	// RistrettoModule for fx.
	RistrettoModule = fx.Options(
		fx.Provide(ristretto.NewCache),
		fx.Invoke(rim.Register),
	)

	// SnappyCompressorModule for fx.
	SnappyCompressorModule = fx.Provide(compressor.NewSnappy)

	// ProtoMarshallerModule for fx.
	ProtoMarshallerModule = fx.Provide(marshaller.NewProto)
)
