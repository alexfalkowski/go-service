package cache

import (
	"github.com/alexfalkowski/go-service/pkg/cache/redis"
	"go.uber.org/fx"
)

var (
	// RedisRingModule for fx.
	RedisRingModule = fx.Provide(redis.NewRing)

	// RedisOptionsModule for fx.
	RedisOptionsModule = fx.Provide(redis.NewOptions)

	// RedisCacheModule for fx.
	RedisCacheModule = fx.Provide(redis.NewCache)

	// RedisModule for fx.
	RedisModule = fx.Options(RedisRingModule, RedisOptionsModule, RedisCacheModule)

	// Module for fx.
	Module = fx.Options(RedisModule)
)
