package cache_test

import (
	"fmt"
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
	. "github.com/smartystreets/goconvey/convey"
)

//nolint:funlen
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
			Convey("Given I have a cache of kind "+config.Kind, t, func() {
				world := test.NewWorld(t)
				world.Register()

				driver, err := driver.NewDriver(test.FS, config)
				So(err, ShouldBeNil)

				params := cache.CacheParams{
					Lifecycle:  world.Lifecycle,
					Config:     config,
					Compressor: test.Compressor,
					Encoder:    test.Encoder,
					Pool:       test.Pool,
					Driver:     driver,
					Tracer:     world.NewTracer(),
					Logger:     world.Logger,
					Meter:      world.Server.Meter,
				}

				cache := cache.NewCache(params)

				world.RequireStart()

				Convey(fmt.Sprintf("When I get an item of type %T", value), func() {
					persist, get := value[0], value[1]

					err := cache.Persist(t.Context(), "test", persist, time.Minute)
					So(err, ShouldBeNil)

					ok, err := cache.Get(t.Context(), "test", get)
					So(err, ShouldBeNil)
					So(ok, ShouldBeTrue)

					Convey(fmt.Sprintf("Then I should have a value %T", value), func() {
						switch kind := get.(type) {
						case *string:
							So(*kind, ShouldEqual, "hello?")
						case *bytes.Buffer:
							So(kind.Bytes(), ShouldEqual, strings.Bytes("hello?"))
						case *v1.SayHelloRequest:
							So(kind.GetName(), ShouldEqual, "hello?")
						case *test.Request:
							So(kind.Name, ShouldEqual, "hello?")
						default:
							So(true, ShouldBeFalse) // should never happen.
						}
					})

					_, err = cache.Remove(t.Context(), "test")
					So(err, ShouldBeNil)
				})

				world.RequireStop()
			})
		}
	}
}

func TestGenericValidCache(t *testing.T) {
	Convey("Given I have a cache", t, func() {
		world := test.NewWorld(t)
		world.Register()

		config := test.NewCacheConfig("sync", "snappy", "json", "redis")

		driver, err := driver.NewDriver(test.FS, config)
		So(err, ShouldBeNil)

		params := cache.CacheParams{
			Lifecycle:  world.Lifecycle,
			Config:     config,
			Compressor: test.Compressor,
			Encoder:    test.Encoder,
			Pool:       test.Pool,
			Driver:     driver,
			Tracer:     world.NewTracer(),
			Logger:     world.Logger,
			Meter:      world.Server.Meter,
		}

		kind := cache.NewCache(params)
		cache.Register(kind)

		world.RequireStart()

		Convey("When I get an item", func() {
			err := cache.Persist(t.Context(), "test", ptr.Value("hello?"), time.Minute)
			So(err, ShouldBeNil)

			value, ok, err := cache.Get[string](t.Context(), "test")
			So(err, ShouldBeNil)
			So(ok, ShouldBeTrue)

			Convey("Then I should have a value", func() {
				So(*value, ShouldEqual, "hello?")
			})

			_, err = kind.Remove(t.Context(), "test")
			So(err, ShouldBeNil)
		})

		world.RequireStop()
	})
}

func TestExpiredCache(t *testing.T) {
	Convey("Given I have a cache", t, func() {
		world := test.NewWorld(t)
		world.Register()

		config := test.NewCacheConfig("sync", "snappy", "json", "redis")

		driver, err := driver.NewDriver(test.FS, config)
		So(err, ShouldBeNil)

		params := cache.CacheParams{
			Lifecycle:  world.Lifecycle,
			Config:     config,
			Compressor: test.Compressor,
			Encoder:    test.Encoder,
			Pool:       test.Pool,
			Driver:     driver,
			Tracer:     world.NewTracer(),
			Logger:     world.Logger,
			Meter:      world.Server.Meter,
		}

		kind := cache.NewCache(params)
		cache.Register(kind)

		world.RequireStart()

		Convey("When I get an item", func() {
			err := cache.Persist(t.Context(), "test", ptr.Value("hello?"), time.Nanosecond)
			So(err, ShouldBeNil)

			// Simulate expiry.
			time.Sleep(time.Second)

			_, ok, err := cache.Get[string](t.Context(), "test")
			So(err, ShouldBeNil)

			Convey("Then I should not have an item in cache", func() {
				So(ok, ShouldBeFalse)
			})

			_, err = kind.Remove(t.Context(), "test")
			So(err, ShouldBeNil)
		})

		world.RequireStop()
	})
}

func TestErroneousCache(t *testing.T) {
	configs := []*config.Config{
		test.NewCacheConfig("redis", "snappy", "json", "none"),
		test.NewCacheConfig("redis", "snappy", "json", "hooks"),
		test.NewCacheConfig("test", "snappy", "json", "hooks"),
	}

	for _, config := range configs {
		Convey("When I create a cache", t, func() {
			world := test.NewWorld(t)
			world.Register()

			_, err := driver.NewDriver(test.FS, config)

			world.RequireStart()

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})

			world.RequireStop()
		})
	}
}

func TestDisabledCache(t *testing.T) {
	configs := []*config.Config{
		nil,
	}

	for _, config := range configs {
		Convey("When I create a cache", t, func() {
			world := test.NewWorld(t)
			world.Register()

			_, err := driver.NewDriver(test.FS, config)

			world.RequireStart()

			Convey("Then I should not have an error", func() {
				So(err, ShouldBeNil)
			})

			world.RequireStop()
		})
	}

	for _, config := range configs {
		Convey("When I create a cache", t, func() {
			world := test.NewWorld(t)
			world.Register()

			params := cache.CacheParams{
				Lifecycle:  world.Lifecycle,
				Config:     config,
				Compressor: test.Compressor,
				Encoder:    test.Encoder,
				Pool:       test.Pool,
				Tracer:     world.NewTracer(),
				Logger:     world.Logger,
				Meter:      world.Server.Meter,
				Driver:     &test.Cache{},
			}

			kind := cache.NewCache(params)
			cache.Register(kind)

			world.RequireStart()

			Convey("Then I should have no cache", func() {
				So(kind, ShouldBeNil)
			})

			world.RequireStop()
		})
	}
}

func TestErroneousSave(t *testing.T) {
	configs := []*config.Config{
		test.NewCacheConfig("sync", "snappy", "error", "redis"),
	}

	for _, config := range configs {
		Convey("Given I have a cache", t, func() {
			world := test.NewWorld(t)
			world.Register()

			driver, err := driver.NewDriver(test.FS, config)
			So(err, ShouldBeNil)

			params := cache.CacheParams{
				Lifecycle:  world.Lifecycle,
				Config:     config,
				Compressor: test.Compressor,
				Encoder:    test.Encoder,
				Pool:       test.Pool,
				Driver:     driver,
				Tracer:     world.NewTracer(),
				Logger:     world.Logger,
				Meter:      world.Server.Meter,
			}

			kind := cache.NewCache(params)
			cache.Register(kind)

			world.RequireStart()

			Convey("When I try to save a value", func() {
				err := cache.Persist(t.Context(), "test", ptr.Value("test"), time.Minute)

				Convey("Then I should have an error", func() {
					So(err, ShouldBeError)
				})

				world.RequireStop()
			})
		})
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
		Convey("Given I have a cache", t, func() {
			world := test.NewWorld(t)
			world.Register()

			params := cache.CacheParams{
				Lifecycle:  world.Lifecycle,
				Config:     value.Key,
				Compressor: test.Compressor,
				Encoder:    test.Encoder,
				Pool:       test.Pool,
				Driver:     value.Value,
				Tracer:     world.NewTracer(),
				Logger:     world.Logger,
				Meter:      world.Server.Meter,
			}

			kind := cache.NewCache(params)
			cache.Register(kind)

			world.RequireStart()

			Convey("When I try to encode a value", func() {
				ptr := ptr.Zero[string]()
				_, err := kind.Get(t.Context(), "test", ptr)

				Convey("Then I should have an error", func() {
					So(err, ShouldBeError)
				})
			})

			world.RequireStop()
		})
	}
}
