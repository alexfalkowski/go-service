package cache

import (
	"github.com/alexfalkowski/go-service/cache/compressor"
	"github.com/alexfalkowski/go-service/cache/marshaller"
	"github.com/alexfalkowski/go-service/cache/redis"
	reprom "github.com/alexfalkowski/go-service/cache/redis/metrics/prometheus"
	"github.com/alexfalkowski/go-service/cache/redis/otel"
	"github.com/alexfalkowski/go-service/cache/ristretto"
	riprom "github.com/alexfalkowski/go-service/cache/ristretto/metrics/prometheus"
	"go.uber.org/fx"
)

var (
	// RedisModule for fx.
	RedisModule = fx.Options(
		fx.Provide(redis.NewClient),
		fx.Provide(redis.NewOptions),
		fx.Provide(redis.NewCache),
		fx.Provide(redis.NewRingOptions),
		fx.Provide(otel.NewTracer),
		fx.Provide(reprom.NewCollector),
		fx.Invoke(reprom.Register),
	)

	// RistrettoModule for fx.
	RistrettoModule = fx.Options(
		fx.Provide(ristretto.NewCache),
		fx.Provide(riprom.NewCollector),
		fx.Invoke(riprom.Register),
	)

	// SnappyCompressorModule for fx.
	SnappyCompressorModule = fx.Provide(compressor.NewSnappy)

	// ProtoMarshallerModule for fx.
	ProtoMarshallerModule = fx.Provide(marshaller.NewProto)
)
