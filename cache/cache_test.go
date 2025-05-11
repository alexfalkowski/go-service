package cache_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/cache"
	"github.com/alexfalkowski/go-service/cache/config"
	"github.com/alexfalkowski/go-service/cache/driver"
	"github.com/alexfalkowski/go-service/internal/test"
	"github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/types/ptr"
	. "github.com/smartystreets/goconvey/convey"
)

func TestValidCache(t *testing.T) {
	configs := []*config.Config{
		test.NewCacheConfig("redis", "snappy", "json", "redis"),
		test.NewCacheConfig("sync", "snappy", "json", "redis"),
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

			Convey("When I save an item", func() {
				value := "hello?"
				err := kind.Persist(t.Context(), "test", &value, time.Minute)

				Convey("Then I should have no error", func() {
					So(err, ShouldBeNil)
				})
			})

			Convey("When I get an item", func() {
				value := "wassup?"

				err := cache.Persist(t.Context(), "test", &value, time.Minute)
				So(err, ShouldBeNil)

				v, ok, err := cache.Get[string](t.Context(), "test")
				So(err, ShouldBeNil)
				So(ok, ShouldBeTrue)

				Convey("Then I should have a value", func() {
					So(*v, ShouldEqual, "wassup?")
				})
			})

			_, err = kind.Remove(t.Context(), "test")
			So(err, ShouldBeNil)

			world.RequireStop()
		})
	}
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
