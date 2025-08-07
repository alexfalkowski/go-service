package driver

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	cache "github.com/alexfalkowski/go-service/v2/cache/config"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/faabiosr/cachego"
	"github.com/faabiosr/cachego/redis"
	"github.com/faabiosr/cachego/sync"
	client "github.com/redis/go-redis/v9"
)

var (
	// ErrNotFound for driver.
	ErrNotFound = errors.New("cache: driver not found")

	// ErrExpired is an alias of cachego.ErrCacheExpired.
	ErrExpired = cachego.ErrCacheExpired
)

// NewDriver creates a new cache driver with different backends.
func NewDriver(fs *os.FS, cfg *cache.Config) (Driver, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	switch cfg.Kind {
	case "redis":
		data, err := fs.ReadSource(cfg.Options["url"].(string))
		if err != nil {
			return nil, err
		}

		opts, err := client.ParseURL(bytes.String(data))
		if err != nil {
			return nil, err
		}

		return redis.New(client.NewClient(opts)), nil
	case "sync":
		return sync.New(), nil
	default:
		return nil, ErrNotFound
	}
}

// Driver is a alias of cachego.
type Driver = cachego.Cache
