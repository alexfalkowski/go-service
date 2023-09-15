package cache

import (
	"github.com/alexfalkowski/go-service/cache/compressor"
	"github.com/alexfalkowski/go-service/cache/marshaller"
	"github.com/alexfalkowski/go-service/cache/redis"
	rtel "github.com/alexfalkowski/go-service/cache/redis/telemetry"
	"github.com/alexfalkowski/go-service/cache/ristretto"
	ritel "github.com/alexfalkowski/go-service/cache/ristretto/telemetry"
	"go.uber.org/fx"
)

var (
	// RedisModule for fx.
	RedisModule = fx.Options(
		fx.Provide(redis.NewClient),
		fx.Provide(redis.NewOptions),
		fx.Provide(redis.NewCache),
		fx.Provide(redis.NewRingOptions),
		fx.Provide(rtel.NewTracer),
		fx.Provide(rtel.NewMetrics),
		fx.Invoke(rtel.RegisterMetrics),
	)

	// RistrettoModule for fx.
	RistrettoModule = fx.Options(
		fx.Provide(ristretto.NewCache),
		fx.Provide(ritel.NewMetrics),
		fx.Invoke(ritel.RegisterMetrics),
	)

	// SnappyCompressorModule for fx.
	SnappyCompressorModule = fx.Provide(compressor.NewSnappy)

	// ProtoMarshallerModule for fx.
	ProtoMarshallerModule = fx.Provide(marshaller.NewProto)
)
