package cache

import (
	"context"

	"github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/runtime"
	"github.com/faabiosr/cachego"
	"github.com/faabiosr/cachego/redis"
	"github.com/faabiosr/cachego/sync"
	client "github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

// NewCache from config.
func NewCache(lc fx.Lifecycle, cfg *Config) (cache cachego.Cache, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Prefix("cache", runtime.ConvertRecover(r))
		}
	}()

	switch cfg.Kind {
	case "redis":
		url, err := os.ReadFile(cfg.Options["url"].(string))
		runtime.Must(err)

		opts, err := client.ParseURL(url)
		runtime.Must(err)

		cache = redis.New(client.NewClient(opts))
	default:
		cache = sync.New()
	}

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			return cache.Flush()
		},
	})

	return
}
