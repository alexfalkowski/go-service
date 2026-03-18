package cache_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/cache"
	"github.com/alexfalkowski/go-service/v2/cache/config"
	"github.com/alexfalkowski/go-service/v2/cache/driver"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/types/ptr"
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

	world := test.NewWorld(t)
	world.Register()

	driver, err := driver.NewDriver(test.FS, cfg)
	require.NoError(t, err)

	params := cache.CacheParams{
		Lifecycle:  world.Lifecycle,
		Config:     cfg,
		Compressor: test.Compressor,
		Encoder:    test.Encoder,
		Pool:       test.Pool,
		Driver:     driver,
	}

	cache := cache.NewCache(params)
	world.RequireStart()

	require.NoError(t, cache.Persist(t.Context(), "test", persist, time.Minute))
	require.NoError(t, cache.Get(t.Context(), "test", get))
	assertValidCacheValue(t, get)
	require.NoError(t, cache.Remove(t.Context(), "test"))

	world.RequireStop()
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
	world := test.NewWorld(t)
	world.Register()

	config := test.NewCacheConfig("sync", "snappy", "json", "redis")

	driver, err := driver.NewDriver(test.FS, config)
	require.NoError(t, err)

	params := cache.CacheParams{
		Lifecycle:  world.Lifecycle,
		Config:     config,
		Compressor: test.Compressor,
		Encoder:    test.Encoder,
		Pool:       test.Pool,
		Driver:     driver,
	}

	kind := cache.NewCache(params)
	cache.Register(kind)

	world.RequireStart()
	require.NoError(t, cache.Persist(t.Context(), "test", ptr.Value("hello?"), time.Minute))

	value, err := cache.Get[string](t.Context(), "test")
	require.NoError(t, err)
	require.Equal(t, "hello?", *value)

	require.NoError(t, kind.Remove(t.Context(), "test"))

	world.RequireStop()
}

func TestGenericDisabledCache(t *testing.T) {
	cache.Register(nil)

	require.NoError(t, cache.Persist(t.Context(), "test", ptr.Value("hello?"), time.Minute))

	value, err := cache.Get[string](t.Context(), "test")
	require.NoError(t, err)
	require.Equal(t, strings.Empty, *value)
}

func TestExpiredCache(t *testing.T) {
	world := test.NewWorld(t)
	world.Register()

	config := test.NewCacheConfig("sync", "snappy", "json", "redis")

	driver, err := driver.NewDriver(test.FS, config)
	require.NoError(t, err)

	params := cache.CacheParams{
		Lifecycle:  world.Lifecycle,
		Config:     config,
		Compressor: test.Compressor,
		Encoder:    test.Encoder,
		Pool:       test.Pool,
		Driver:     driver,
	}

	kind := cache.NewCache(params)
	cache.Register(kind)

	world.RequireStart()
	require.NoError(t, cache.Persist(t.Context(), "test", ptr.Value("hello?"), time.Nanosecond))

	// Simulate expiry.
	time.Sleep(time.Second)

	_, err = cache.Get[string](t.Context(), "test")
	require.NoError(t, err)

	require.NoError(t, kind.Remove(t.Context(), "test"))

	world.RequireStop()
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
			world := test.NewWorld(t)
			world.Register()

			_, err := driver.NewDriver(test.FS, tt.config)

			world.RequireStart()
			require.Error(t, err)
			world.RequireStop()
		})
	}
}

func TestDisabledCache(t *testing.T) {
	t.Run("driver", func(t *testing.T) {
		world := test.NewWorld(t)
		world.Register()

		_, err := driver.NewDriver(test.FS, nil)

		world.RequireStart()
		require.NoError(t, err)
		world.RequireStop()
	})

	t.Run("cache", func(t *testing.T) {
		world := test.NewWorld(t)
		world.Register()

		params := cache.CacheParams{
			Lifecycle:  world.Lifecycle,
			Config:     nil,
			Compressor: test.Compressor,
			Encoder:    test.Encoder,
			Pool:       test.Pool,
			Driver:     &test.Cache{},
		}

		kind := cache.NewCache(params)
		cache.Register(kind)

		world.RequireStart()
		require.Nil(t, kind)
		world.RequireStop()
	})
}

func TestErroneousSave(t *testing.T) {
	t.Run("invalid encoder", func(t *testing.T) {
		config := test.NewCacheConfig("sync", "snappy", "error", "redis")

		world := test.NewWorld(t)
		world.Register()

		driver, err := driver.NewDriver(test.FS, config)
		require.NoError(t, err)

		params := cache.CacheParams{
			Lifecycle:  world.Lifecycle,
			Config:     config,
			Compressor: test.Compressor,
			Encoder:    test.Encoder,
			Pool:       test.Pool,
			Driver:     driver,
		}

		kind := cache.NewCache(params)
		cache.Register(kind)

		world.RequireStart()
		require.Error(t, cache.Persist(t.Context(), "test", ptr.Value("test"), time.Minute))
		world.RequireStop()
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
			world := test.NewWorld(t)
			world.Register()

			params := cache.CacheParams{
				Lifecycle:  world.Lifecycle,
				Config:     value.config,
				Compressor: test.Compressor,
				Encoder:    test.Encoder,
				Pool:       test.Pool,
				Driver:     value.driver,
			}

			kind := cache.NewCache(params)
			cache.Register(kind)

			world.RequireStart()

			require.Error(t, kind.Get(t.Context(), "test", ptr.Zero[string]()))

			world.RequireStop()
		})
	}
}

func TestMissingCache(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	cfg := &config.Config{Kind: "sync"}

	d, err := driver.NewDriver(nil, cfg)
	require.NoError(t, err)

	params := cache.CacheParams{
		Lifecycle:  lc,
		Config:     cfg,
		Compressor: test.Compressor,
		Encoder:    test.Encoder,
		Pool:       test.Pool,
		Driver:     d,
	}

	kind := cache.NewCache(params)
	value := "existing"

	require.NoError(t, kind.Get(t.Context(), "missing", &value))
	require.Equal(t, "existing", value)
}
