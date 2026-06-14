package driver

import (
	"github.com/alexfalkowski/go-service/v2/cache/config"
	"github.com/alexfalkowski/go-service/v2/cache/driver/errors"
	"github.com/alexfalkowski/go-service/v2/cache/driver/internal/redis"
	"github.com/alexfalkowski/go-service/v2/cache/driver/internal/ttlcache"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/time"
)

// DriverParams defines dependencies for constructing a [Driver].
type DriverParams struct {
	di.In

	// Lifecycle registers backend shutdown hooks.
	Lifecycle di.Lifecycle

	// FS resolves backend source-string options.
	FS *os.FS

	// Config selects and configures the cache backend.
	Config *config.Config

	// Logger routes backend logs through the go-service logger when configured.
	Logger *logger.Logger
}

// NewDriver constructs a cache [Driver] for the configured backend.
//
// # Disabled behavior
//
// If cfg is nil (caching disabled), [NewDriver] returns (nil, nil). Callers are expected to tolerate a nil [Driver].
//
// # Configuration expectations
//
// NewDriver dispatches on [config.Config.Kind]. Some backends expect specific keys to be present in [config.Config.Options].
// For example, the "redis" backend expects:
//
//   - options["url"] to be a string "source string" (e.g. "env:REDIS_URL" or "file:/path/to/url" or a literal URL)
//
// The URL is read via [os.FS.ReadSource], parsed using
// [github.com/redis/go-redis/v9.ParseURL], and then the client is instrumented for tracing and metrics via
// [github.com/alexfalkowski/go-service/v2/cache/telemetry]
// when those telemetry providers are enabled.
//
// The Redis client is closed from the supplied lifecycle's [di.Hook.OnStop] hook.
//
// Instrumentation errors are treated as fatal configuration/runtime errors and
// are converted into panics via [github.com/alexfalkowski/go-service/v2/runtime.Must], matching the existing repository
// convention for mandatory telemetry wiring in internal constructors.
//
// # Backends
//
// Supported kinds include:
//   - "redis": Redis backend using [github.com/redis/go-redis/v9]
//   - "ttlcache": in-memory backend using [github.com/jellydator/ttlcache/v3]
//
// The built-in "ttlcache" backend stores values in process memory and expires entries lazily on access.
//
// If [config.Config.Kind] is unknown, [NewDriver] returns
// [github.com/alexfalkowski/go-service/v2/cache/driver/errors.ErrNotFound].
func NewDriver(params DriverParams) (Driver, error) {
	cfg := params.Config
	if !cfg.IsEnabled() {
		return nil, nil
	}

	switch cfg.Kind {
	case "redis":
		return redis.NewDriver(params.Lifecycle, params.FS, params.Config, params.Logger)
	case "ttlcache":
		return ttlcache.NewDriver(cfg.GetMaxEntries()), nil
	default:
		return nil, errors.ErrNotFound
	}
}

// Driver is the minimal cache backend interface used by the cache facade.
//
// Implementations must honor the provided context for blocking operations.
type Driver interface {
	// Delete removes the cached key.
	Delete(ctx context.Context, key string) error

	// Fetch retrieves the cached value for key.
	Fetch(ctx context.Context, key string) (string, error)

	// Flush removes cached data according to backend-specific semantics.
	//
	// Implementations may clear more than go-service cache namespace keys. For
	// example, the built-in Redis driver clears the entire selected Redis
	// database.
	Flush(ctx context.Context) error

	// Save stores value under key for the provided lifetime.
	Save(ctx context.Context, key, value string, lifetime time.Duration) error
}
