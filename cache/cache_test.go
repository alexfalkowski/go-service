package cache_test

import (
	"fmt"
	"testing"

	"github.com/alexfalkowski/go-service/bytes"
	"github.com/alexfalkowski/go-service/cache"
	"github.com/alexfalkowski/go-service/cache/config"
	"github.com/alexfalkowski/go-service/cache/driver"
	"github.com/alexfalkowski/go-service/internal/test"
	"github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/types/ptr"
	. "github.com/smartystreets/goconvey/convey"
)

type tuple [2]any

//nolint:funlen
func TestValidCache(t *testing.T) {
	configs := []*config.Config{
		test.NewCacheConfig("redis", "snappy", "json", "redis"),
		test.NewCacheConfig("sync", "snappy", "json", "redis"),
	}

	for _, config := range configs {
		for _, value := range []tuple{
			{ptr.Value("hello?"), ptr.Zero[string]()},
			{bytes.NewBufferString("hello?"), &bytes.Buffer{}},
		} {
			Convey("Given I have a cache of kind "+config.Kind, t, func() {
				world := test.NewWorld(t)
				world.Register()

				driver, err := driver.New(config)
				So(err, ShouldBeNil)

				params := cache.Params{
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
						case *[]byte:
							So(*kind, ShouldEqual, []byte("hello?"))
						case *bytes.Buffer:
							So(kind.Bytes(), ShouldEqual, []byte("hello?"))
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

		driver, err := driver.New(config)
		So(err, ShouldBeNil)

		params := cache.Params{
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

		driver, err := driver.New(config)
		So(err, ShouldBeNil)

		params := cache.Params{
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

			_, err := driver.New(config)

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

			_, err := driver.New(config)

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

			params := cache.Params{
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

			driver, err := driver.New(config)
			So(err, ShouldBeNil)

			params := cache.Params{
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
	type Tuple struct {
		Config *config.Config
		driver.Driver
	}

	tuples := []*Tuple{
		{Config: test.NewCacheConfig("sync", "snappy", "error", "redis"), Driver: &test.Cache{Value: "d2hhdD8="}},
		{Config: test.NewCacheConfig("sync", "error", "json", "redis"), Driver: &test.Cache{Value: "d2hhdD8="}},
		{Config: test.NewCacheConfig("sync", "snappy", "json", "redis"), Driver: &test.Cache{Value: "what?"}},
		{Config: test.NewCacheConfig("sync", "snappy", "json", "redis"), Driver: &test.ErrCache{}},
	}

	for _, tuple := range tuples {
		Convey("Given I have a cache", t, func() {
			world := test.NewWorld(t)
			world.Register()

			params := cache.Params{
				Lifecycle:  world.Lifecycle,
				Config:     tuple.Config,
				Compressor: test.Compressor,
				Encoder:    test.Encoder,
				Pool:       test.Pool,
				Driver:     tuple.Driver,
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
