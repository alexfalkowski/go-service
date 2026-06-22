package cache_test

import (
	"math"
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/cache"
	"github.com/alexfalkowski/go-service/v2/cache/config"
	"github.com/alexfalkowski/go-service/v2/cache/driver"
	drivererrors "github.com/alexfalkowski/go-service/v2/cache/driver/errors"
	compresserrors "github.com/alexfalkowski/go-service/v2/compress/errors"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/encoding/base64"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/ptr"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestValidCache(t *testing.T) {
	for _, tt := range cacheRoundTripCases() {
		t.Run(tt.name, func(t *testing.T) {
			test.RequireCacheRoundTrip(t, tt.config, tt.persist(), tt.get())
		})
	}
}

func TestGenericValidCache(t *testing.T) {
	cfg := test.NewCacheConfig("ttlcache", "snappy", "json", "redis")
	world := test.NewStartedWorld(t, test.WithWorldCacheConfig(cfg), test.WithWorldRegisterCache())

	require.NoError(t, cache.Persist(t.Context(), "test", new("hello?"), time.Minute))

	value, err := cache.Get[string](t.Context(), "test")
	require.NoError(t, err)
	require.Equal(t, "hello?", *value)

	require.NoError(t, world.Remove(t.Context(), "test"))
}

func TestModuleRegistersGenericCache(t *testing.T) {
	cache.Register(nil)
	t.Cleanup(func() {
		cache.Register(nil)
	})

	app := di.New(
		di.NoLogger,
		cache.Module,
		di.Constructor(func() *logger.Logger { return nil }),
		fx.Supply(
			test.NewCacheConfig("ttlcache", "none", "json", "redis"),
			test.FS,
			test.Encoder,
			test.Pool,
			test.Compressor,
		),
	)
	require.NoError(t, app.Err())

	require.NoError(t, app.Start(t.Context()))
	t.Cleanup(func() {
		require.NoError(t, app.Stop(context.Background()))
	})

	require.NoError(t, cache.Persist(t.Context(), "test", new("hello?"), time.Minute))

	value, err := cache.Get[string](t.Context(), "test")
	require.NoError(t, err)
	require.Equal(t, "hello?", *value)
}

func TestGenericDisabledCache(t *testing.T) {
	cache.Register(nil)

	require.NoError(t, cache.Persist(t.Context(), "test", new("hello?"), time.Minute))

	value, err := cache.Get[string](t.Context(), "test")
	require.NoError(t, err)
	require.Equal(t, strings.Empty, *value)
}

func TestMaxSizeOnPersist(t *testing.T) {
	cfg := test.NewCacheConfig("ttlcache", "snappy", "json", "redis")
	cfg.MaxSize = 4
	world := test.NewStartedWorld(t, test.WithWorldCacheConfig(cfg), test.WithWorldRegisterCache())

	err := world.Persist(t.Context(), "test", new("hello?"), time.Minute)
	require.ErrorIs(t, err, compresserrors.ErrTooLarge)
}

func TestMaxSizeOnPersistCompressedValue(t *testing.T) {
	cfg := test.NewCacheConfig("ttlcache", "snappy", "json", "redis")
	cfg.MaxSize = 4
	world := test.NewStartedWorld(t, test.WithWorldCacheConfig(cfg), test.WithWorldRegisterCache())

	err := world.Persist(t.Context(), "test", new("a"), time.Minute)
	require.ErrorIs(t, err, compresserrors.ErrTooLarge)
}

func TestMaxSizeOnPersistStopsOversizedEncodedValue(t *testing.T) {
	cfg := test.NewCacheConfig("ttlcache", "none", "json", "redis")
	cfg.MaxSize = 4
	world := test.NewStartedWorld(t, test.WithWorldCacheConfig(cfg), test.WithWorldRegisterCache())
	value := &oversizedWriterTo{size: cfg.MaxSize.Bytes() + 1}

	err := world.Persist(t.Context(), "test", value, time.Minute)
	require.ErrorIs(t, err, compresserrors.ErrTooLarge)
	require.Equal(t, cfg.MaxSize.Bytes(), value.written)
}

func TestMaxSizeOnGet(t *testing.T) {
	cfg := test.NewCacheConfig("ttlcache", "snappy", "json", "redis")
	world := test.NewStartedWorld(t, test.WithWorldCacheConfig(cfg), test.WithWorldRegisterCache())

	require.NoError(t, world.Persist(t.Context(), "test", new("hello?"), time.Minute))

	cfg.MaxSize = 4

	err := world.Get(t.Context(), "test", ptr.Zero[string]())
	require.ErrorIs(t, err, compresserrors.ErrTooLarge)
}

func TestMaxSizeOnGetAllowsEncodedValueLargerThanLimit(t *testing.T) {
	cfg := test.NewCacheConfig("ttlcache", "none", "json", "redis")
	cfg.MaxSize = 4
	world := test.NewStartedWorld(t,
		test.WithWorldCacheConfig(cfg),
		test.WithWorldCacheDriver(&test.Cache{Value: base64.Encode(strings.Bytes(`"ok"`))}),
	)

	value := ptr.Zero[string]()
	require.NoError(t, world.Get(t.Context(), "test", value))
	require.Equal(t, "ok", *value)
}

func TestMaxSizeOnGetRejectsEncodedValueTooLarge(t *testing.T) {
	cfg := test.NewCacheConfig("ttlcache", "none", "json", "redis")
	cfg.MaxSize = 4
	world := test.NewStartedWorld(t,
		test.WithWorldCacheConfig(cfg),
		test.WithWorldCacheDriver(&test.Cache{Value: base64.Encode(make([]byte, 7))}),
	)

	err := world.Get(t.Context(), "test", ptr.Zero[string]())
	require.ErrorIs(t, err, compresserrors.ErrTooLarge)
}

func TestMaxSizeOnGetAllowsHugeConfiguredLimit(t *testing.T) {
	cfg := test.NewCacheConfig("ttlcache", "none", "json", "redis")
	cfg.MaxSize = bytes.Size(math.MaxInt64)
	world := test.NewStartedWorld(t,
		test.WithWorldCacheConfig(cfg),
		test.WithWorldCacheDriver(&test.Cache{Value: base64.Encode(strings.Bytes(`"ok"`))}),
	)

	value := ptr.Zero[string]()
	require.NoError(t, world.Get(t.Context(), "test", value))
	require.Equal(t, "ok", *value)
}

func TestGetMissesAfterCacheFormatChange(t *testing.T) {
	tests := []struct {
		persist   *config.Config
		retrieval *config.Config
		name      string
	}{
		{
			persist:   test.NewCacheConfig("ttlcache", "none", "json", "redis"),
			retrieval: test.NewCacheConfig("ttlcache", "snappy", "json", "redis"),
			name:      "compressor",
		},
		{
			persist:   test.NewCacheConfig("ttlcache", "none", "json", "redis"),
			retrieval: test.NewCacheConfig("ttlcache", "none", "gob", "redis"),
			name:      "encoder",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			drv, err := driver.NewDriver(driver.DriverParams{Config: tt.persist})
			require.NoError(t, err)

			writer := test.NewStartedWorld(t,
				test.WithWorldCacheConfig(tt.persist),
				test.WithWorldCacheDriver(drv),
			)
			require.NoError(t, writer.Persist(t.Context(), "test", new("hello?"), time.Minute))

			reader := test.NewStartedWorld(t,
				test.WithWorldCacheConfig(tt.retrieval),
				test.WithWorldCacheDriver(drv),
			)
			value := "unchanged"
			require.NoError(t, reader.Get(t.Context(), "test", &value))
			require.Equal(t, "unchanged", value)
		})
	}
}

func TestExpiredCache(t *testing.T) {
	cfg := test.NewCacheConfig("ttlcache", "snappy", "json", "redis")
	world := test.NewStartedWorld(t, test.WithWorldCacheConfig(cfg), test.WithWorldRegisterCache())
	require.NoError(t, cache.Persist(t.Context(), "test", new("hello?"), time.Nanosecond))

	// Simulate expiry.
	time.Sleep(time.Second)

	_, err := cache.Get[string](t.Context(), "test")
	require.NoError(t, err)

	require.NoError(t, world.Remove(t.Context(), "test"))
}

func TestErroneousCache(t *testing.T) {
	tests := []struct {
		config *config.Config
		name   string
	}{
		{name: "missing driver secret", config: test.NewCacheConfig("redis", "snappy", "json", "none")},
		{name: "invalid driver secret", config: test.NewCacheConfig("redis", "snappy", "json", "hooks")},
		{name: "unsupported cache kind", config: test.NewCacheConfig("test", "snappy", "json", "hooks")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := driver.NewDriver(driver.DriverParams{
				Lifecycle: fxtest.NewLifecycle(t),
				FS:        test.FS,
				Config:    tt.config,
			})
			require.Error(t, err)
		})
	}
}

func TestDisabledCache(t *testing.T) {
	t.Run("driver", func(t *testing.T) {
		_, err := driver.NewDriver(driver.DriverParams{
			Lifecycle: fxtest.NewLifecycle(t),
			FS:        test.FS,
		})
		require.NoError(t, err)
	})

	t.Run("cache", func(t *testing.T) {
		world := test.NewStartedWorld(t, test.WithWorldCacheConfig(nil), test.WithWorldRegisterCache())

		require.Nil(t, world.Cache)
	})
}

func TestErroneousSave(t *testing.T) {
	t.Run("invalid encoder", func(t *testing.T) {
		cfg := test.NewCacheConfig("ttlcache", "snappy", "error", "redis")
		test.NewStartedWorld(t, test.WithWorldCacheConfig(cfg), test.WithWorldRegisterCache())
		require.Error(t, cache.Persist(t.Context(), "test", new("test"), time.Minute))
	})

	t.Run("read from only falls back to configured encoder", func(t *testing.T) {
		cfg := test.NewCacheConfig("ttlcache", "none", "json", "redis")
		world := test.NewStartedWorld(t, test.WithWorldCacheConfig(cfg))
		require.NoError(t, world.Persist(t.Context(), "test", &test.ReadFromOnly{Name: "hello?"}, time.Minute))

		var get test.Request
		require.NoError(t, world.Get(t.Context(), "test", &get))
		require.Equal(t, "hello?", get.Name)
	})
}

func TestErroneousGet(t *testing.T) {
	tests := []struct {
		driver driver.Driver
		config *config.Config
		name   string
	}{
		{name: "decode error", config: test.NewCacheConfig("ttlcache", "snappy", "error", "redis"), driver: &test.Cache{Value: "d2hhdD8="}},
		{name: "decompress error", config: test.NewCacheConfig("ttlcache", "error", "json", "redis"), driver: &test.Cache{Value: "d2hhdD8="}},
		{name: "unmarshal error", config: test.NewCacheConfig("ttlcache", "snappy", "json", "redis"), driver: &test.Cache{Value: "what?"}},
		{name: "driver error", config: test.NewCacheConfig("ttlcache", "snappy", "json", "redis"), driver: &test.ErrCache{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			world := test.NewStartedWorld(t,
				test.WithWorldCacheConfig(tt.config),
				test.WithWorldCacheDriver(tt.driver),
			)

			require.Error(t, world.Get(t.Context(), "test", ptr.Zero[string]()))
		})
	}
}

func TestWriteToOnlyUsesConfiguredEncoderOnGet(t *testing.T) {
	cfg := test.NewCacheConfig("ttlcache", "none", "json", "redis")
	world := test.NewStartedWorld(t, test.WithWorldCacheConfig(cfg))
	require.NoError(t, world.Persist(t.Context(), "test", &test.Request{Name: "hello?"}, time.Minute))

	get := &test.WriteToOnly{}
	require.NoError(t, world.Get(t.Context(), "test", get))
	require.Equal(t, "hello?", get.Name)
}

func TestMissingCache(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	cfg := &config.Config{Kind: "ttlcache", MaxEntries: config.DefaultMaxEntries}

	d, err := driver.NewDriver(driver.DriverParams{
		Lifecycle: lc,
		Config:    cfg,
	})
	require.NoError(t, err)

	params := cache.CacheParams{
		Config:     cfg,
		Compressor: test.Compressor,
		Encoder:    test.Encoder,
		Pool:       test.Pool,
		Driver:     d,
	}

	c := cache.NewCache(params)
	t.Cleanup(func() {
		require.NoError(t, c.Flush(context.Background()))
	})
	value := "existing"

	require.NoError(t, c.Get(t.Context(), "missing", &value))
	require.Equal(t, "existing", value)
}

func TestExpiredCacheLeavesDestinationUnchanged(t *testing.T) {
	cfg := &config.Config{Kind: "ttlcache", MaxEntries: config.DefaultMaxEntries}
	c := cache.NewCache(cache.CacheParams{
		Config:     cfg,
		Compressor: test.Compressor,
		Encoder:    test.Encoder,
		Pool:       test.Pool,
		Driver:     expiredCacheDriver{},
	})
	value := "existing"

	require.NoError(t, c.Get(t.Context(), "expired", &value))
	require.Equal(t, "existing", value)
}

func cacheRoundTripCases() []cacheRoundTripCase {
	compressors := []string{"none", "snappy", "s2", "zstd"}
	encoders := []string{"json", "hjson", "yaml", "yml", "toml", "gob", "msgpack"}
	cases := []cacheRoundTripCase{
		{
			name:    "redis/default/request",
			config:  test.NewCacheConfig("redis", "snappy", strings.Empty, "redis"),
			persist: func() any { return &test.Request{Name: "hello?"} },
			get:     func() any { return &test.Request{} },
		},
		{
			name:    "ttlcache/default/string",
			config:  test.NewCacheConfig("ttlcache", strings.Empty, strings.Empty, "redis"),
			persist: func() any { return new("hello?") },
			get:     func() any { return ptr.Zero[string]() },
		},
		{
			name:    "ttlcache/default/buffer",
			config:  test.NewCacheConfig("ttlcache", strings.Empty, strings.Empty, "redis"),
			persist: func() any { return bytes.NewBufferString("hello?") },
			get:     func() any { return &bytes.Buffer{} },
		},
		{
			name:    "ttlcache/default/protobuf",
			config:  test.NewCacheConfig("ttlcache", strings.Empty, strings.Empty, "redis"),
			persist: func() any { return &v1.SayHelloRequest{Name: "hello?"} },
			get:     func() any { return &v1.SayHelloRequest{} },
		},
		{
			name:    "ttlcache/unknown/unknown/request",
			config:  test.NewCacheConfig("ttlcache", "unknown", "unknown", "redis"),
			persist: func() any { return &test.Request{Name: "hello?"} },
			get:     func() any { return &test.Request{} },
		},
	}

	tests := make([]cacheRoundTripCase, 0, len(cases)+len(compressors)*len(encoders))
	tests = append(tests, cases...)

	for _, compressor := range compressors {
		for _, encoder := range encoders {
			tests = append(tests, cacheRoundTripCase{
				name:    "ttlcache/" + compressor + "/" + encoder + "/request",
				config:  test.NewCacheConfig("ttlcache", compressor, encoder, "redis"),
				persist: func() any { return &test.Request{Name: "hello?"} },
				get:     func() any { return &test.Request{} },
			})
		}
	}

	return tests
}

type cacheRoundTripCase struct {
	persist func() any
	get     func() any
	config  *config.Config
	name    string
}

type oversizedWriterTo struct {
	size    int64
	written int64
}

func (w *oversizedWriterTo) WriteTo(dst io.Writer) (int64, error) {
	n, err := dst.Write(make([]byte, w.size))
	w.written = int64(n)

	return w.written, err
}

type expiredCacheDriver struct{}

func (expiredCacheDriver) Delete(context.Context, string) error {
	return nil
}

func (expiredCacheDriver) Fetch(context.Context, string) (string, error) {
	return strings.Empty, drivererrors.ErrExpired
}

func (expiredCacheDriver) Flush(context.Context) error {
	return nil
}

func (expiredCacheDriver) Save(context.Context, string, string, time.Duration) error {
	return nil
}
