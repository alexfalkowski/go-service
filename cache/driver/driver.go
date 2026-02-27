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
//
// Drivers may return this error (or wrap it) to indicate that a cache entry exists but is expired.
// Use IsExpiredError to classify this condition.
const ErrExpired = cachego.ErrCacheExpired

// ErrNotFound is returned when the configured cache driver kind is unknown.
//
// It is returned by NewDriver when Config.Kind does not match any backend compiled into this module.
var ErrNotFound = errors.New("cache: driver not found")

// NewDriver constructs a cache Driver for the configured backend.
//
// # Disabled behavior
//
// If cfg is nil (caching disabled), NewDriver returns (nil, nil). Callers are expected to tolerate a nil Driver.
//
// # Configuration expectations
//
// NewDriver dispatches on cfg.Kind. Some backends expect specific keys to be present in cfg.Options.
// For example, the "redis" backend expects:
//
//   - options["url"] to be a string "source string" (e.g. "env:REDIS_URL" or "file:/path/to/url" or a literal URL)
//
// The URL is read via fs.ReadSource, parsed using redis/go-redis ParseURL, and then the client is instrumented
// for tracing and metrics.
//
// # Backends
//
// Supported kinds include:
//   - "redis": Redis backend using github.com/redis/go-redis
//   - "sync": in-memory backend using github.com/faabiosr/cachego/sync
//
// If cfg.Kind is unknown, NewDriver returns ErrNotFound.
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

// IsExpiredError reports whether err represents an expired cache entry.
//
// This helper exists so higher-level code can treat expired entries as cache misses regardless of the
// underlying backend implementation.
func IsExpiredError(err error) bool {
	return errors.Is(err, ErrExpired)
}

// Driver is an alias for cachego.Cache.
//
// It is the minimal interface used by the cache facade (`cache.Cache`) for persistence operations:
// fetch/save/delete/flush.
type Driver = cachego.Cache
