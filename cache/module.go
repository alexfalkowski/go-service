package cache

import (
	"github.com/alexfalkowski/go-service/cache/redis"
	"github.com/alexfalkowski/go-service/cache/ristretto"
	"go.uber.org/fx"
)

var (
	// RedisModule for fx.
	// nolint:gochecknoglobals
	RedisModule = fx.Options(fx.Provide(redis.NewRing), fx.Provide(redis.NewOptions), fx.Provide(redis.NewCache))

	// RistrettoModule for fx.
	// nolint:gochecknoglobals
	RistrettoModule = fx.Options(fx.Provide(ristretto.NewCache))
)
