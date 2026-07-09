package cache

import (
	"math"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/cache/config"
	"github.com/alexfalkowski/go-service/v2/cache/driver"
	drivererrors "github.com/alexfalkowski/go-service/v2/cache/driver/errors"
	"github.com/alexfalkowski/go-service/v2/compress"
	compresserrors "github.com/alexfalkowski/go-service/v2/compress/errors"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/encoding/base64"
	"github.com/alexfalkowski/go-service/v2/encoding/proto"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-sync"
)

const (
	// maxBase64EncodedLenInput is the largest decoded-size limit whose padded base64 length fits in int.
	maxBase64EncodedLenInput = math.MaxInt / 4 * 3

	cacheKeyNamespace = "cache:v1"
)

// CacheParams defines dependencies for constructing a [Cache].
//
// It is intended for dependency injection ([go.uber.org/fx]/[go.uber.org/dig]). The constructor will typically be wired via [Module].
type CacheParams struct {
	di.In

	// Config configures cache encoding, compression, and limits.
	Config *config.Config

	// Encoder provides value encoders by name.
	Encoder *encoding.Map

	// Pool provides reusable buffers for cache encoding.
	Pool *sync.BufferPool

	// Compressor provides compression implementations by name.
	Compressor *compress.Map

	// Driver stores encoded cache values.
	Driver driver.Driver
}

// NewCache constructs a [Cache] from configuration.
//
// If caching is disabled (i.e. [CacheParams.Config] is nil), [NewCache] returns nil. Callers are expected to
// tolerate a nil cache instance.
func NewCache(params CacheParams) *Cache {
	if !params.Config.IsEnabled() {
		return nil
	}

	return &Cache{
		compress: params.Compressor,
		encoding: params.Encoder,
		config:   params.Config,
		pool:     params.Pool,
		driver:   params.Driver,
	}
}

// Cache provides a typed cache facade on top of a cache driver.
//
// It serializes values using an encoder, optionally compresses the serialized bytes, base64-encodes the
// final bytes, and stores the resulting string via the configured driver.
//
// Encoding selection is operation-dependent:
//   - [Cache.Persist] uses "plain" only for [io.WriterTo] values
//   - [Cache.Get] uses "plain" only for [io.ReaderFrom] destinations
//   - [proto.Message] uses "proto"
//   - otherwise the configured encoder is used, falling back to "json"
//
// Compression is selected from configuration, falling back to "none" when unknown/unavailable.
type Cache struct {
	compress *compress.Map
	encoding *encoding.Map
	config   *config.Config
	pool     *sync.BufferPool
	driver   driver.Driver
	single   sync.AnySingleFlightGroup
}

// Flush removes cached data according to the underlying driver's flush semantics.
//
// For persistent backends such as Redis this can be a destructive operation:
// the built-in Redis driver uses FLUSHDB and clears the entire selected Redis
// database, including keys that were not created through this cache facade. It
// is intentionally not called during lifecycle shutdown.
func (c *Cache) Flush(ctx context.Context) error {
	return c.driver.Flush(ctx)
}

// Remove deletes a cached key.
//
// If the key does not exist, driver behavior is implementation-specific.
func (c *Cache) Remove(ctx context.Context, key string) error {
	return c.driver.Delete(ctx, c.driverKey(key))
}

// Get loads a cached value for key into value and reports whether a value was found.
//
// Cache misses are not treated as errors: if the entry is missing or expired, Get returns
// false, nil and leaves value unchanged.
//
// The value parameter should be a pointer to the destination value (for example *MyStruct).
func (c *Cache) Get(ctx context.Context, key string, value any) (bool, error) {
	val, err := c.driver.Get(ctx, c.driverKey(key))
	if err != nil {
		if drivererrors.IsMissingError(err) || drivererrors.IsExpiredError(err) {
			return false, nil
		}

		return false, err
	}

	if err := c.decode(val, value); err != nil {
		return false, err
	}

	return true, nil
}

// Persist stores value under key with the provided TTL.
//
// The value is encoded, compressed, and base64-encoded before being saved via the driver.
// A TTL <= 0 is passed through to the driver; semantics are driver-specific (for example, it may mean
// "no expiration" or "immediate expiration").
//
// TTL resolution is driver-specific.
func (c *Cache) Persist(ctx context.Context, key string, value any, ttl time.Duration) error {
	enc, err := c.encode(value)
	if err != nil {
		return err
	}

	return c.driver.Save(ctx, c.driverKey(key), enc, ttl)
}

// GetOrPersist returns the cached value for key, or produces and stores it via fn when the key is absent.
//
// On a cache hit fn is not called and no write occurs. On a miss fn populates value and the encoded value is
// published atomically, so concurrent callers converge on a single stored value. Concurrent in-process misses
// for the same key run fn once and share the produced value; separate processes may each run fn once, but the
// atomic publish still yields a single stored winner that every caller decodes.
//
// The value parameter should be a pointer to the destination value (for example *MyStruct). It is both
// populated by fn and used as the decode destination for the resolved value.
func (c *Cache) GetOrPersist(ctx context.Context, key string, value any, ttl time.Duration, fn func() error) error {
	ok, err := c.Get(ctx, key, value)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	driverKey := c.driverKey(key)

	result, err, _ := c.single.Do(driverKey, func() (any, error) {
		return c.loadOrSave(ctx, driverKey, value, ttl, fn)
	})
	if err != nil {
		return err
	}

	return c.decode(result.(string), value)
}

// loadOrSave re-checks the driver for driverKey, then produces and atomically publishes a value on a miss.
//
// It runs inside the single-flight group, so it executes once per key for concurrent in-process callers. It
// returns the encoded value that every caller decodes: the existing stored value when another writer won, or
// the value just produced and stored.
func (c *Cache) loadOrSave(ctx context.Context, driverKey string, value any, ttl time.Duration, fn func() error) (string, error) {
	if val, err := c.driver.Get(ctx, driverKey); err == nil {
		return val, nil
	} else if !drivererrors.IsMissingError(err) && !drivererrors.IsExpiredError(err) {
		return strings.Empty, err
	}

	if err := fn(); err != nil {
		return strings.Empty, err
	}

	enc, err := c.encode(value)
	if err != nil {
		return strings.Empty, err
	}

	existing, loaded, err := c.driver.GetOrSave(ctx, driverKey, enc, ttl)
	if err != nil {
		return strings.Empty, err
	}
	if loaded {
		return existing, nil
	}

	return enc, nil
}

func (c *Cache) driverKey(key string) string {
	return strings.Join(":", cacheKeyNamespace, c.compressorKind(), c.encoderKind(), key)
}

func (c *Cache) encode(value any) (string, error) {
	buf := c.pool.Get()
	defer c.pool.Put(buf)

	maxSize := c.config.GetMaxSize()
	writer := &maxSizeWriter{writer: buf, max: maxSize.Bytes()}
	if err := c.writeEncoder(value).Encode(writer, value); err != nil {
		return strings.Empty, err
	}

	data := buf.Bytes()
	compressed, err := c.compressor().Compress(data, maxSize)
	if err != nil {
		return strings.Empty, err
	}
	if int64(len(compressed)) > maxSize.Bytes() {
		return strings.Empty, compresserrors.ErrTooLarge
	}

	encoded := base64.Encode(compressed)

	return encoded, nil
}

func (c *Cache) decode(value string, field any) error {
	maxSize := c.config.GetMaxSize()
	if maxSize.Bytes() <= maxBase64EncodedLenInput && int64(len(value)) > base64.EncodedLen(maxSize) {
		return compresserrors.ErrTooLarge
	}

	decoded, err := base64.Decode(value)
	if err != nil {
		return err
	}

	// Enforce the read-side limit on decompressed payloads, not on the stored base64 wrapper.
	decompressed, err := c.compressor().Decompress(decoded, maxSize)
	if err != nil {
		return err
	}

	return c.readEncoder(field).Decode(bytes.NewReader(decompressed), field)
}

func (c *Cache) compressor() compress.Compressor {
	return c.compress.Get(c.compressorKind())
}

func (c *Cache) compressorKind() string {
	if cmp := c.compress.Get(c.config.Compressor); cmp != nil {
		return c.config.Compressor
	}

	return "none"
}

func (c *Cache) readEncoder(value any) encoding.Encoder {
	switch value.(type) {
	case io.ReaderFrom:
		return c.encoding.Get("plain")
	case proto.Message:
		return c.encoding.Get("proto")
	default:
		return c.configuredEncoder()
	}
}

func (c *Cache) writeEncoder(value any) encoding.Encoder {
	switch value.(type) {
	case io.WriterTo:
		return c.encoding.Get("plain")
	case proto.Message:
		return c.encoding.Get("proto")
	default:
		return c.configuredEncoder()
	}
}

func (c *Cache) configuredEncoder() encoding.Encoder {
	return c.encoding.Get(c.encoderKind())
}

func (c *Cache) encoderKind() string {
	if enc := c.encoding.Get(c.config.Encoder); enc != nil {
		return c.config.Encoder
	}

	return "json"
}
