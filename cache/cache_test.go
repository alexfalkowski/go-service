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
)

func TestValidCache(t *testing.T) {
	configs := []*config.Config{
		test.NewCacheConfig("redis", "snappy", strings.Empty, "redis"),
		test.NewCacheConfig("sync", strings.Empty, strings.Empty, "redis"),
	}

	for _, config := range configs {
		for _, value := range []test.AnyTuple{
			{ptr.Value("hello?"), ptr.Zero[string]()},
			{bytes.NewBufferString("hello?"), &bytes.Buffer{}},
			{&v1.SayHelloRequest{Name: "hello?"}, &v1.SayHelloRequest{}},
			{&test.Request{Name: "hello?"}, &test.Request{}},
		} {
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

			cache := cache.NewCache(params)
			world.RequireStart()

			persist, get := value[0], value[1]
			require.NoError(t, cache.Persist(t.Context(), "test", persist, time.Minute))

			require.NoError(t, cache.Get(t.Context(), "test", get))

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

			require.NoError(t, cache.Remove(t.Context(), "test"))

			world.RequireStop()
		}
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
	configs := []*config.Config{
		test.NewCacheConfig("redis", "snappy", "json", "none"),
		test.NewCacheConfig("redis", "snappy", "json", "hooks"),
		test.NewCacheConfig("test", "snappy", "json", "hooks"),
	}

	for _, config := range configs {
		world := test.NewWorld(t)
		world.Register()

		_, err := driver.NewDriver(test.FS, config)

		world.RequireStart()
		require.Error(t, err)
		world.RequireStop()
	}
}

func TestDisabledCache(t *testing.T) {
	configs := []*config.Config{
		nil,
	}

	for _, config := range configs {
		world := test.NewWorld(t)
		world.Register()

		_, err := driver.NewDriver(test.FS, config)

		world.RequireStart()
		require.NoError(t, err)
		world.RequireStop()
	}

	for _, config := range configs {
		world := test.NewWorld(t)
		world.Register()

		params := cache.CacheParams{
			Lifecycle:  world.Lifecycle,
			Config:     config,
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
	}
}

func TestErroneousSave(t *testing.T) {
	configs := []*config.Config{
		test.NewCacheConfig("sync", "snappy", "error", "redis"),
	}

	for _, config := range configs {
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
	}
}

func TestErroneousGet(t *testing.T) {
	values := []*test.KeyValue[*config.Config, driver.Driver]{
		{Key: test.NewCacheConfig("sync", "snappy", "error", "redis"), Value: &test.Cache{Value: "d2hhdD8="}},
		{Key: test.NewCacheConfig("sync", "error", "json", "redis"), Value: &test.Cache{Value: "d2hhdD8="}},
		{Key: test.NewCacheConfig("sync", "snappy", "json", "redis"), Value: &test.Cache{Value: "what?"}},
		{Key: test.NewCacheConfig("sync", "snappy", "json", "redis"), Value: &test.ErrCache{}},
	}

	for _, value := range values {
		world := test.NewWorld(t)
		world.Register()

		params := cache.CacheParams{
			Lifecycle:  world.Lifecycle,
			Config:     value.Key,
			Compressor: test.Compressor,
			Encoder:    test.Encoder,
			Pool:       test.Pool,
			Driver:     value.Value,
		}

		kind := cache.NewCache(params)
		cache.Register(kind)

		world.RequireStart()

		require.Error(t, kind.Get(t.Context(), "test", ptr.Zero[string]()))

		world.RequireStop()
	}
}
