package cachego

import (
	cache "github.com/alexfalkowski/go-service/cache/config"
	"github.com/alexfalkowski/go-service/cache/errors"
	"github.com/alexfalkowski/go-service/os"
	"github.com/faabiosr/cachego"
	"github.com/faabiosr/cachego/redis"
	"github.com/faabiosr/cachego/sync"
	client "github.com/redis/go-redis/v9"
)

// Cache is a alias of cachego.
type Cache = cachego.Cache

// New cache from config.
func New(cfg *cache.Config) (Cache, error) {
	if !cache.IsEnabled(cfg) {
		return nil, nil
	}

	switch cfg.Kind {
	case "redis":
		url, err := os.ReadFile(cfg.Options["url"].(string))
		if err != nil {
			return nil, err
		}

		opts, err := client.ParseURL(string(url))
		if err != nil {
			return nil, err
		}

		return redis.New(client.NewClient(opts)), nil
	case "sync":
		return sync.New(), nil
	default:
		return nil, errors.ErrNotFound
	}
}
