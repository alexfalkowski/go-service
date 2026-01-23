package driver

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	cache "github.com/alexfalkowski/go-service/v2/cache/config"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/faabiosr/cachego"
	"github.com/faabiosr/cachego/redis"
	"github.com/faabiosr/cachego/sync"
	otel "github.com/redis/go-redis/extra/redisotel/v9"
	client "github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
)

// ErrExpired is an alias for cachego.ErrCacheExpired.
const ErrExpired = cachego.ErrCacheExpired

// ErrNotFound for driver.
var ErrNotFound = errors.New("cache: driver not found")

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

		opts.MaintNotificationsConfig = &maintnotifications.Config{
			Mode: maintnotifications.ModeDisabled,
		}

		client := client.NewClient(opts)
		runtime.Must(otel.InstrumentTracing(client))
		runtime.Must(otel.InstrumentMetrics(client))

		return redis.New(client), nil
	case "sync":
		return sync.New(), nil
	default:
		return nil, ErrNotFound
	}
}

// IsExpiredError checks if the error is an expired error.
func IsExpiredError(err error) bool {
	return errors.Is(err, ErrExpired)
}

// Driver is a alias of cachego.
type Driver = cachego.Cache
