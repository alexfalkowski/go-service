package cache

import (
	"io"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/cache/config"
	"github.com/alexfalkowski/go-service/v2/cache/driver"
	"github.com/alexfalkowski/go-service/v2/compress"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/encoding/base64"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/sync"
	"github.com/alexfalkowski/go-service/v2/time"
	"google.golang.org/protobuf/proto"
)

// CacheParams defines dependencies for constructing a Cache.
//
// It is intended for dependency injection (Fx/Dig). The constructor will typically be wired via `Module`.
type CacheParams struct {
	di.In
	Lifecycle  di.Lifecycle
	Config     *config.Config
	Encoder    *encoding.Map
	Pool       *sync.BufferPool
	Compressor *compress.Map
	Driver     driver.Driver
}

// NewCache constructs a Cache from configuration and registers a shutdown hook.
//
// If caching is disabled (i.e. params.Config is nil), NewCache returns nil. Callers are expected to
// tolerate a nil cache instance.
//
// When enabled, NewCache registers an OnStop lifecycle hook that calls (*Cache).Close to flush the
// underlying driver on shutdown.
func NewCache(params CacheParams) *Cache {
	if !params.Config.IsEnabled() {
		return nil
	}

	cache := &Cache{
		cm:     params.Compressor,
		em:     params.Encoder,
		cfg:    params.Config,
		pool:   params.Pool,
		driver: params.Driver,
	}

	params.Lifecycle.Append(di.Hook{
		OnStop: func(ctx context.Context) error {
			return cache.Close(ctx)
		},
	})

	return cache
}

// Cache provides a typed-ish cache facade on top of a cache driver.
//
// It serializes values using an encoder, optionally compresses the serialized bytes, base64-encodes the
// final bytes, and stores the resulting string via the configured driver.
//
// Encoding selection is value-dependent (see encoder):
//   - io.ReaderFrom / io.WriterTo use "plain"
//   - proto.Message uses "proto"
//   - otherwise the configured encoder is used, falling back to "json"
//
// Compression is selected from configuration, falling back to "none" when unknown/unavailable.
type Cache struct {
	cm     *compress.Map
	em     *encoding.Map
	cfg    *config.Config
	pool   *sync.BufferPool
	driver driver.Driver
}

// Close flushes the underlying driver.
//
// This is typically invoked automatically via the lifecycle hook registered by NewCache.
func (c *Cache) Close(_ context.Context) error {
	return c.driver.Flush()
}

// Remove deletes a cached key.
//
// If the key does not exist, driver behavior is implementation-specific.
func (c *Cache) Remove(_ context.Context, key string) error {
	return c.driver.Delete(key)
}

// Get loads a cached value for key into value.
//
// Cache misses are not treated as errors: if the entry is missing or expired, Get returns nil and
// leaves value unchanged.
//
// The value parameter should be a pointer to the destination value (for example *MyStruct).
func (c *Cache) Get(_ context.Context, key string, value any) error {
	val, err := c.driver.Fetch(key)
	if err != nil {
		if driver.IsExpiredError(err) {
			return nil
		}

		return err
	}

	return c.decode(val, value)
}

// Persist stores value under key with the provided TTL.
//
// The value is encoded, compressed, and base64-encoded before being saved via the driver.
// A TTL <= 0 is passed through to the driver; semantics are driver-specific (for example, it may mean
// "no expiration" or "immediate expiration").
func (c *Cache) Persist(_ context.Context, key string, value any, ttl time.Duration) error {
	enc, err := c.encode(value)
	if err != nil {
		return err
	}

	return c.driver.Save(key, enc, ttl)
}

func (c *Cache) encode(value any) (string, error) {
	buf := c.pool.Get()
	defer c.pool.Put(buf)

	if err := c.encoder(value).Encode(buf, value); err != nil {
		return strings.Empty, err
	}

	data := buf.Bytes()
	compressed := c.compressor().Compress(data)
	encoded := base64.Encode(compressed)

	return encoded, nil
}

func (c *Cache) decode(value string, field any) error {
	decoded, err := base64.Decode(value)
	if err != nil {
		return err
	}

	decompressed, err := c.compressor().Decompress(decoded)
	if err != nil {
		return err
	}

	return c.encoder(field).Decode(bytes.NewReader(decompressed), field)
}

func (c *Cache) compressor() compress.Compressor {
	if cmp := c.cm.Get(c.cfg.Compressor); cmp != nil {
		return cmp
	}

	return c.cm.Get("none")
}

func (c *Cache) encoder(value any) encoding.Encoder {
	switch value.(type) {
	case io.ReaderFrom, io.WriterTo:
		return c.em.Get("plain")
	case proto.Message:
		return c.em.Get("proto")
	default:
		if enc := c.em.Get(c.cfg.Encoder); enc != nil {
			return enc
		}

		return c.em.Get("json")
	}
}
