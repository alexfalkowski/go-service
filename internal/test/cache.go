package test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/cache"
	"github.com/alexfalkowski/go-service/v2/cache/config"
	"github.com/alexfalkowski/go-service/v2/cache/driver"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

// Cache is a [cache.Cache] test double that returns a fixed value from Fetch.
type Cache struct {
	Value string
}

// Delete removes the given key from the cache.
func (c *Cache) Delete(context.Context, string) error {
	return nil
}

// Fetch returns the cached value.
func (c *Cache) Fetch(context.Context, string) (string, error) {
	return c.Value, nil
}

// Flush clears the cache.
func (c *Cache) Flush(context.Context) error {
	return nil
}

// Save stores the value in the cache for the given TTL.
func (c *Cache) Save(context.Context, string, string, time.Duration) error {
	return nil
}

// ErrCache is a [cache.Cache] test double that fails fetch, delete, and save operations with ErrFailed.
//
// Flush succeeds so tests can use ErrCache in started worlds without turning cleanup into the failure under test.
type ErrCache struct{}

// Delete returns ErrFailed.
func (*ErrCache) Delete(context.Context, string) error {
	return ErrFailed
}

// Fetch returns ErrFailed.
func (*ErrCache) Fetch(context.Context, string) (string, error) {
	return strings.Empty, ErrFailed
}

// Flush clears the cache.
func (*ErrCache) Flush(context.Context) error {
	return nil
}

// Save returns ErrFailed.
func (*ErrCache) Save(context.Context, string, string, time.Duration) error {
	return ErrFailed
}

// RequireCacheRoundTrip persists a value, reads it back, and asserts the shared hello payload.
func RequireCacheRoundTrip(tb testing.TB, cfg *config.Config, persist, get any) {
	tb.Helper()

	world := NewStartedWorld(tb, WithWorldCacheConfig(cfg))

	require.NoError(tb, world.Persist(tb.Context(), "test", persist, time.Minute))
	require.NoError(tb, world.Get(tb.Context(), "test", get))
	RequireCacheValue(tb, get)
	require.NoError(tb, world.Remove(tb.Context(), "test"))
}

// RequireCacheValue asserts that get contains the shared hello payload.
func RequireCacheValue(tb testing.TB, get any) {
	tb.Helper()

	switch kind := get.(type) {
	case *string:
		require.Equal(tb, "hello?", *kind)
	case *bytes.Buffer:
		require.Equal(tb, strings.Bytes("hello?"), kind.Bytes())
	case *v1.SayHelloRequest:
		require.Equal(tb, "hello?", kind.GetName())
	case *Request:
		require.Equal(tb, "hello?", kind.Name)
	default:
		require.Fail(tb, "invalid kind")
	}
}

// ReadFromOnly is a cache payload that implements [io.ReaderFrom] but relies on normal encoding.
type ReadFromOnly struct {
	Name string `json:"name"`
}

// ReadFrom implements [io.ReaderFrom] without consuming input.
func (*ReadFromOnly) ReadFrom(io.Reader) (int64, error) {
	return 0, nil
}

// WriteToOnly is a cache payload that implements [io.WriterTo] but relies on normal decoding.
type WriteToOnly struct {
	Name string `json:"name"`
}

// WriteTo implements [io.WriterTo] without writing output.
func (*WriteToOnly) WriteTo(io.Writer) (int64, error) {
	return 0, nil
}

func redisCache(lc di.Lifecycle) (*cache.Cache, cache.Pinger, error) {
	cfg := NewCacheConfig("redis", "snappy", "json", "redis")

	drv, err := driver.NewDriver(driver.DriverParams{
		Lifecycle: lc,
		FS:        FS,
		Config:    cfg,
	})
	if err != nil {
		return nil, nil, err
	}

	params := cache.CacheParams{
		Config:     cfg,
		Compressor: Compressor,
		Encoder:    Encoder,
		Pool:       Pool,
		Driver:     drv,
	}

	return cache.NewCache(params), cache.NewPinger(drv), nil
}

func newWorldCache(tb testing.TB, lc di.Lifecycle, opts *worldOpts) (*cache.Cache, cache.Pinger) {
	tb.Helper()

	var kind *cache.Cache
	var pinger cache.Pinger
	if opts.cache == nil {
		var err error
		kind, pinger, err = redisCache(lc)
		require.NoError(tb, err)
	} else {
		kind, pinger = createWorldCache(tb, lc, opts.cache)
	}

	if opts.registerCache {
		cache.Register(kind)
		tb.Cleanup(func() {
			cache.Register(nil)
		})
	}

	return kind, pinger
}

func createWorldCache(tb testing.TB, lc di.Lifecycle, opts *worldCacheOpts) (*cache.Cache, cache.Pinger) {
	tb.Helper()

	drv := opts.driver
	if drv == nil {
		var err error
		drv, err = driver.NewDriver(driver.DriverParams{
			Lifecycle: lc,
			FS:        FS,
			Config:    opts.config,
		})
		require.NoError(tb, err)
	}

	params := cache.CacheParams{
		Config:     opts.config,
		Compressor: Compressor,
		Encoder:    Encoder,
		Pool:       Pool,
		Driver:     drv,
	}

	return cache.NewCache(params), cache.NewPinger(drv)
}
