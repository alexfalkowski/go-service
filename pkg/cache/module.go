package cache

import (
	"github.com/alexfalkowski/go-service/pkg/cache/redis"
	"github.com/alexfalkowski/go-service/pkg/cache/ristretto"
	"go.uber.org/fx"
)

var (
	// RedisModule for fx.
	RedisModule = fx.Options(fx.Provide(redis.NewRing), fx.Provide(redis.NewOptions), fx.Provide(redis.NewCache), fx.Provide(redis.NewConfig))

	// RistrettoModule for fx.
	RistrettoModule = fx.Options(fx.Provide(ristretto.NewConfig), fx.Provide(ristretto.NewCache))
)
