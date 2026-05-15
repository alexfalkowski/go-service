package cache_test

import (
	"io"
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/cache"
	"github.com/alexfalkowski/go-service/v2/cache/config"
	"github.com/alexfalkowski/go-service/v2/cache/driver"
	"github.com/alexfalkowski/go-service/v2/compress/errors"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/v2/ptr"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func TestValidCache(t *testing.T) {
	configs := []struct {
		config *config.Config
		name   string
	}{
		{name: "redis", config: test.NewCacheConfig("redis", "snappy", strings.Empty, "redis")},
		{name: "sync", config: test.NewCacheConfig("sync", strings.Empty, strings.Empty, "redis")},
	}
	values := []struct {
		persist func() any
		get     func() any
		name    string
	}{
		{name: "string", persist: func() any { return ptr.Value("hello?") }, get: func() any { return ptr.Zero[string]() }},
		{name: "buffer", persist: func() any { return bytes.NewBufferString("hello?") }, get: func() any { return &bytes.Buffer{} }},
		{name: "protobuf", persist: func() any { return &v1.SayHelloRequest{Name: "hello?"} }, get: func() any { return &v1.SayHelloRequest{} }},
		{name: "request", persist: func() any { return &test.Request{Name: "hello?"} }, get: func() any { return &test.Request{} }},
	}

	for _, cfg := range configs {
		t.Run(cfg.name, func(t *testing.T) {
			for _, value := range values {
				t.Run(value.name, func(t *testing.T) {
					testValidCacheCase(t, cfg.config, value.persist(), value.get())
				})
			}
		})
	}
}

func testValidCacheCase(t *testing.T, cfg *config.Config, persist, get any) {
	t.Helper()

	world := test.NewStartedWorld(t, test.WithWorldCacheConfig(cfg))

	require.NoError(t, world.Persist(t.Context(), "test", persist, time.Minute))
	require.NoError(t, world.Get(t.Context(), "test", get))
	assertValidCacheValue(t, get)
	require.NoError(t, world.Remove(t.Context(), "test"))
}

func assertValidCacheValue(t *testing.T, get any) {
	t.Helper()

	switch kind := get.(type) {
	case *string:
		require.Equal(t, "hello?", *kind)
	case *bytes.Buffer:
		require.Equal(t, strings.Bytes("hello?"), kind.Bytes())
	case *v1.SayHelloRequest:
		require.Equal(t, "hello?", kind.GetName())
	case *test.Request:
		require.Equal(t, "hello?", kind.Name)
	default:
		require.Fail(t, "invalid kind")
	}
}

func TestGenericValidCache(t *testing.T) {
	config := test.NewCacheConfig("sync", "snappy", "json", "redis")
	world := test.NewStartedWorld(t, test.WithWorldCacheConfig(config), test.WithWorldRegisterCache())

	require.NoError(t, cache.Persist(t.Context(), "test", ptr.Value("hello?"), time.Minute))

	value, err := cache.Get[string](t.Context(), "test")
	require.NoError(t, err)
	require.Equal(t, "hello?", *value)

	require.NoError(t, world.Remove(t.Context(), "test"))
}

func TestGenericDisabledCache(t *testing.T) {
	cache.Register(nil)

	require.NoError(t, cache.Persist(t.Context(), "test", ptr.Value("hello?"), time.Minute))

	value, err := cache.Get[string](t.Context(), "test")
	require.NoError(t, err)
	require.Equal(t, strings.Empty, *value)
}

func TestMaxSizeOnPersist(t *testing.T) {
	cfg := test.NewCacheConfig("sync", "snappy", "json", "redis")
	cfg.MaxSize = 4
	world := test.NewStartedWorld(t, test.WithWorldCacheConfig(cfg), test.WithWorldRegisterCache())

	err := world.Persist(t.Context(), "test", ptr.Value("hello?"), time.Minute)
	require.ErrorIs(t, err, errors.ErrTooLarge)
}

func TestMaxSizeOnPersistCompressedValue(t *testing.T) {
	cfg := test.NewCacheConfig("sync", "snappy", "json", "redis")
	cfg.MaxSize = 4
	world := test.NewStartedWorld(t, test.WithWorldCacheConfig(cfg), test.WithWorldRegisterCache())

	err := world.Persist(t.Context(), "test", ptr.Value("a"), time.Minute)
	require.ErrorIs(t, err, errors.ErrTooLarge)
}

func TestMaxSizeOnGet(t *testing.T) {
	cfg := test.NewCacheConfig("sync", "snappy", "json", "redis")
	world := test.NewStartedWorld(t, test.WithWorldCacheConfig(cfg), test.WithWorldRegisterCache())

	require.NoError(t, world.Persist(t.Context(), "test", ptr.Value("hello?"), time.Minute))

	cfg.MaxSize = 4

	err := world.Get(t.Context(), "test", ptr.Zero[string]())
	require.ErrorIs(t, err, errors.ErrTooLarge)
}

func TestExpiredCache(t *testing.T) {
	config := test.NewCacheConfig("sync", "snappy", "json", "redis")
	world := test.NewStartedWorld(t, test.WithWorldCacheConfig(config), test.WithWorldRegisterCache())
	require.NoError(t, cache.Persist(t.Context(), "test", ptr.Value("hello?"), time.Nanosecond))

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
		config := test.NewCacheConfig("sync", "snappy", "error", "redis")
		test.NewStartedWorld(t, test.WithWorldCacheConfig(config), test.WithWorldRegisterCache())
		require.Error(t, cache.Persist(t.Context(), "test", ptr.Value("test"), time.Minute))
	})

	t.Run("read from only falls back to configured encoder", func(t *testing.T) {
		config := test.NewCacheConfig("sync", "none", "json", "redis")
		world := test.NewStartedWorld(t, test.WithWorldCacheConfig(config))
		require.NoError(t, world.Persist(t.Context(), "test", &readFromOnly{Name: "hello?"}, time.Minute))

		var get test.Request
		require.NoError(t, world.Get(t.Context(), "test", &get))
		require.Equal(t, "hello?", get.Name)
	})
}

func TestErroneousGet(t *testing.T) {
	values := []struct {
		driver driver.Driver
		config *config.Config
		name   string
	}{
		{name: "decode error", config: test.NewCacheConfig("sync", "snappy", "error", "redis"), driver: &test.Cache{Value: "d2hhdD8="}},
		{name: "decompress error", config: test.NewCacheConfig("sync", "error", "json", "redis"), driver: &test.Cache{Value: "d2hhdD8="}},
		{name: "unmarshal error", config: test.NewCacheConfig("sync", "snappy", "json", "redis"), driver: &test.Cache{Value: "what?"}},
		{name: "driver error", config: test.NewCacheConfig("sync", "snappy", "json", "redis"), driver: &test.ErrCache{}},
	}

	for _, value := range values {
		t.Run(value.name, func(t *testing.T) {
			world := test.NewStartedWorld(t,
				test.WithWorldCacheConfig(value.config),
				test.WithWorldCacheDriver(value.driver),
			)

			require.Error(t, world.Get(t.Context(), "test", ptr.Zero[string]()))
		})
	}
}

func TestWriteToOnlyUsesConfiguredEncoderOnGet(t *testing.T) {
	config := test.NewCacheConfig("sync", "none", "json", "redis")
	world := test.NewStartedWorld(t, test.WithWorldCacheConfig(config))
	require.NoError(t, world.Persist(t.Context(), "test", &test.Request{Name: "hello?"}, time.Minute))

	get := &writeToOnly{}
	require.NoError(t, world.Get(t.Context(), "test", get))
	require.Equal(t, "hello?", get.Name)
}

func TestMissingCache(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	cfg := &config.Config{Kind: "sync"}

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

	kind := cache.NewCache(params)
	t.Cleanup(func() {
		require.NoError(t, kind.Flush(t.Context()))
	})
	value := "existing"

	require.NoError(t, kind.Get(t.Context(), "missing", &value))
	require.Equal(t, "existing", value)
}

type readFromOnly struct {
	Name string `json:"name"`
}

func (*readFromOnly) ReadFrom(io.Reader) (int64, error) {
	return 0, nil
}

type writeToOnly struct {
	Name string `json:"name"`
}

func (*writeToOnly) WriteTo(io.Writer) (int64, error) {
	return 0, nil
}
