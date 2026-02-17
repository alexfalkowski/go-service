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
// It returns nil when caching is disabled.
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

// Cache marshals values, compresses them, and stores them via the configured driver.
type Cache struct {
	cm     *compress.Map
	em     *encoding.Map
	cfg    *config.Config
	pool   *sync.BufferPool
	driver driver.Driver
}

// Close flushes the underlying driver.
func (c *Cache) Close(_ context.Context) error {
	return c.driver.Flush()
}

// Remove deletes a cached key.
func (c *Cache) Remove(_ context.Context, key string) error {
	return c.driver.Delete(key)
}

// Get loads a cached value into value and returns nil on cache misses.
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
