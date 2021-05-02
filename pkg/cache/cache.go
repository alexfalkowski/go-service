package cache

import (
	"context"
	"time"

	"github.com/alexfalkowski/go-service/pkg/config"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
)

// Item that is used for caching.
type Item = cache.Item

// NewCache from config.
// The cache is based on https://github.com/go-redis/cache
func NewCache(lc fx.Lifecycle, cfg *config.Config) *cache.Cache {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server": cfg.CacheHost,
		},
	})

	cache := cache.New(&cache.Options{
		Redis:      ring,
		LocalCache: cache.NewTinyLFU(1000, time.Minute), // nolint:gomnd
		Marshal: func(v interface{}) ([]byte, error) {
			return proto.Marshal(v.(proto.Message))
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return proto.Unmarshal(b, v.(proto.Message))
		},
	})

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return ring.Close()
		},
	})

	return cache
}
