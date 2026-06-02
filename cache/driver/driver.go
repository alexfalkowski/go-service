package driver

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/cache/config"
	"github.com/alexfalkowski/go-service/v2/cache/telemetry"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/redis/go-redis/v9"
	notifications "github.com/redis/go-redis/v9/maintnotifications"
)

// ErrExpired is returned when a cache entry exists but is expired.
//
// Drivers may wrap this error; use [IsExpiredError] to classify this condition.
var ErrExpired = errors.New("cache: expired")

// ErrMissing is returned when a cache entry does not exist.
//
// Drivers may wrap this error; use [IsMissingError] to classify this condition.
var ErrMissing = errors.New("cache: missing")

// ErrNotFound is returned when the configured cache driver kind is unknown.
var ErrNotFound = errors.New("cache: driver not found")

// ErrInvalidURL is returned when a cache backend URL cannot be parsed.
var ErrInvalidURL = errors.New("cache: invalid driver url")

// DriverParams defines dependencies for constructing a [Driver].
type DriverParams struct {
	di.In
	Lifecycle di.Lifecycle
	FS        *os.FS
	Config    *config.Config

	// Logger routes Redis client logs through the go-service logger when configured.
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
// The URL is read via [os.FS.ReadSource], parsed using [redis.ParseURL], and
// then the client is instrumented for tracing and metrics via [telemetry]
// when those telemetry providers are enabled.
//
// The Redis client is closed from the supplied lifecycle's [di.Hook.OnStop] hook.
//
// Instrumentation errors are treated as fatal configuration/runtime errors and
// are converted into panics via [runtime.Must], matching the existing repository
// convention for mandatory telemetry wiring in internal constructors.
//
// # Backends
//
// Supported kinds include:
//   - "redis": Redis backend using [github.com/redis/go-redis/v9]
//   - "sync": in-memory backend using [github.com/alexfalkowski/go-sync] Map
//
// The built-in "sync" backend stores values in process memory and expires entries lazily on access.
//
// If [config.Config.Kind] is unknown, [NewDriver] returns [ErrNotFound].
func NewDriver(params DriverParams) (Driver, error) {
	cfg := params.Config
	if !cfg.IsEnabled() {
		return nil, nil
	}

	switch cfg.Kind {
	case "redis":
		return newRedisDriver(params)
	case "sync":
		return &syncDriver{}, nil
	default:
		return nil, ErrNotFound
	}
}

func newRedisDriver(params DriverParams) (Driver, error) {
	data, err := params.FS.ReadSource(params.Config.Options["url"].(string))
	if err != nil {
		return nil, err
	}

	opts, err := redis.ParseURL(bytes.String(data))
	if err != nil {
		return nil, ErrInvalidURL
	}

	opts.MaintNotificationsConfig = &notifications.Config{
		Mode: notifications.ModeDisabled,
	}
	if params.Logger != nil {
		redis.SetLogger(redisLogger{logger: params.Logger})
	}

	redisClient := redis.NewClient(opts)
	if tracer.IsEnabled() {
		runtime.Must(telemetry.InstrumentTracing(redisClient))
	}

	var metricsClose chan struct{}
	if metrics.IsEnabled() {
		metricsClose = make(chan struct{})
		runtime.Must(telemetry.InstrumentMetrics(redisClient, metricsClose))
	}

	params.Lifecycle.Append(di.Hook{
		OnStop: func(context.Context) error {
			if metricsClose != nil {
				close(metricsClose)
			}

			return redisClient.Close()
		},
	})

	return &redisDriver{client: redisClient}, nil
}

// Driver is the minimal cache backend interface used by the cache facade.
//
// Implementations must honor the provided context for blocking operations.
type Driver interface {
	// Delete removes the cached key.
	Delete(ctx context.Context, key string) error

	// Fetch retrieves the cached value for key.
	Fetch(ctx context.Context, key string) (string, error)

	// Flush removes all cached keys managed by the driver.
	Flush(ctx context.Context) error

	// Save stores value under key for the provided lifetime.
	Save(ctx context.Context, key, value string, lifetime time.Duration) error
}

// IsExpiredError reports whether err represents an expired cache entry.
//
// This helper exists so higher-level code can treat expired entries as cache misses regardless of the
// underlying backend implementation.
func IsExpiredError(err error) bool {
	return errors.Is(err, ErrExpired)
}

// IsMissingError reports whether err represents a missing cache entry.
//
// This helper normalizes the miss semantics of the backends currently supported by this package,
// including Redis nil replies ([redis.Nil]).
func IsMissingError(err error) bool {
	if errors.Is(err, redis.Nil) {
		return true
	}

	return errors.Is(err, ErrMissing)
}
