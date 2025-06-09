package cache

import (
	"context"
	"io"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/cache/cacher"
	"github.com/alexfalkowski/go-service/v2/cache/config"
	"github.com/alexfalkowski/go-service/v2/cache/driver"
	cl "github.com/alexfalkowski/go-service/v2/cache/telemetry/logger"
	cm "github.com/alexfalkowski/go-service/v2/cache/telemetry/metrics"
	ct "github.com/alexfalkowski/go-service/v2/cache/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/compress"
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/encoding/base64"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/sync"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/time"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
)

// CacheParams for cache.
type CacheParams struct {
	fx.In

	Lifecycle  fx.Lifecycle
	Config     *config.Config
	Encoder    *encoding.Map
	Pool       *sync.BufferPool
	Compressor *compress.Map
	Driver     driver.Driver
	Tracer     *tracer.Tracer
	Logger     *logger.Logger
	Meter      *metrics.Meter
	Name       env.Name
}

// NewCache from config.
func NewCache(params CacheParams) cacher.Cache {
	if !config.IsEnabled(params.Config) {
		return nil
	}

	var cache cacher.Cache = &Cache{
		cm:     params.Compressor,
		em:     params.Encoder,
		cfg:    params.Config,
		pool:   params.Pool,
		driver: params.Driver,
	}

	if params.Tracer != nil {
		cache = ct.NewCache(params.Config.Kind, params.Tracer, cache)
	}

	if params.Logger != nil {
		cache = cl.NewCache(params.Config.Kind, params.Logger, cache)
	}

	if params.Meter != nil {
		cache = cm.NewCache(params.Config.Kind, params.Name, params.Meter, cache)
	}

	params.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return cache.Close(ctx)
		},
	})

	return cache
}

// Cache allows marshaling and compressing items to the cache.
type Cache struct {
	cm     *compress.Map
	em     *encoding.Map
	cfg    *config.Config
	pool   *sync.BufferPool
	driver driver.Driver
	rw     sync.RWMutex
}

// Close the cache.
func (c *Cache) Close(_ context.Context) error {
	c.rw.Lock()
	defer c.rw.Unlock()

	return c.driver.Flush()
}

// Remove a cached key.
func (c *Cache) Remove(_ context.Context, key string) (bool, error) {
	c.rw.Lock()
	defer c.rw.Unlock()

	if !c.driver.Contains(key) {
		return false, nil
	}

	return true, c.driver.Delete(key)
}

// Get a cached value.
func (c *Cache) Get(_ context.Context, key string, value any) (bool, error) {
	c.rw.RLock()
	defer c.rw.RUnlock()

	if !c.driver.Contains(key) {
		return false, nil
	}

	val, err := c.driver.Fetch(key)
	if err != nil {
		if errors.Is(err, driver.ErrExpired) {
			return false, nil
		}

		return false, err
	}

	return true, c.decode(val, value)
}

// Persist a value with key and TTL.
func (c *Cache) Persist(_ context.Context, key string, value any, ttl time.Duration) error {
	c.rw.Lock()
	defer c.rw.Unlock()

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
		return "", err
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
