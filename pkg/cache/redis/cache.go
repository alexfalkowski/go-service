package redis

import (
	"context"

	"github.com/alexfalkowski/go-service/pkg/config"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/golang/snappy"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
)

// NewCache from config.
// The cache is based on https://github.com/go-redis/cache
func NewCache(lc fx.Lifecycle, cfg *config.Config) *cache.Cache {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server": cfg.RedisCacheHost,
		},
	})

	cache := cache.New(&cache.Options{
		Redis: ring,
		Marshal: func(v interface{}) ([]byte, error) {
			m, err := proto.Marshal(v.(proto.Message))
			if err != nil {
				return nil, err
			}

			return snappy.Encode(nil, m), nil
		},
		Unmarshal: func(b []byte, v interface{}) error {
			m, err := snappy.Decode(nil, b)
			if err != nil {
				return err
			}

			return proto.Unmarshal(m, v.(proto.Message))
		},
	})

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return ring.Close()
		},
	})

	return cache
}
