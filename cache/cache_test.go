//nolint:varnamelen
package cache_test

import (
	"context"
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/cache"
	"github.com/alexfalkowski/go-service/cache/cachego"
	"github.com/alexfalkowski/go-service/cache/config"
	"github.com/alexfalkowski/go-service/test"
	"github.com/alexfalkowski/go-service/types/ptr"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

//nolint:funlen
func TestValidCache(t *testing.T) {
	configs := []*config.Config{
		test.NewCacheConfig("redis", "snappy", "json", "redis"),
		test.NewCacheConfig("sync", "snappy", "json", "redis"),
	}

	for _, config := range configs {
		Convey("Given I have a cache", t, func() {
			world := test.NewWorld(t)
			world.Register()

			cachego, err := cachego.New(config)
			So(err, ShouldBeNil)

			params := cache.Params{
				Lifecycle:  world.Lifecycle,
				Config:     config,
				Compressor: test.Compressor,
				Encoder:    test.Encoder,
				Pool:       test.Pool,
				Cache:      cachego,
				Tracer:     world.NewTracer(),
				Logger:     world.Logger,
				Meter:      world.Server.Meter,
			}

			c, err := cache.New(params)
			So(err, ShouldBeNil)

			world.RequireStart()

			ctx := context.Background()

			Convey("When I save an item", func() {
				value := "hello?"
				err := c.Persist(ctx, "test", &value, time.Minute)

				Convey("Then I should have no error", func() {
					So(err, ShouldBeNil)
				})
			})

			Convey("When I get an item", func() {
				value := "wassup?"

				err := cache.Persist(ctx, c, "test", &value, time.Minute)
				So(err, ShouldBeNil)

				v, err := cache.Get[string](ctx, c, "test")
				So(err, ShouldBeNil)

				Convey("Then I should have a value", func() {
					So(*v, ShouldEqual, "wassup?")
				})
			})

			err = c.Remove(ctx, "test")
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

			_, err := cachego.New(config)

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

			_, err := cachego.New(config)

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
				Cache:      &test.Cache{},
			}

			_, err := cache.New(params)

			world.RequireStart()

			Convey("Then I should not have an error", func() {
				So(err, ShouldBeNil)
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

			cachego, err := cachego.New(config)
			So(err, ShouldBeNil)

			params := cache.Params{
				Lifecycle:  world.Lifecycle,
				Config:     config,
				Compressor: test.Compressor,
				Encoder:    test.Encoder,
				Pool:       test.Pool,
				Cache:      cachego,
				Tracer:     world.NewTracer(),
				Logger:     world.Logger,
				Meter:      world.Server.Meter,
			}

			c, err := cache.New(params)
			So(err, ShouldBeNil)

			world.RequireStart()

			ctx := context.Background()

			Convey("When I try to save a value", func() {
				err := cache.Persist(ctx, c, "test", ptr.Value("test"), time.Minute)

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
		cachego.Cache
	}

	tuples := []*Tuple{
		{Config: test.NewCacheConfig("sync", "snappy", "error", "redis"), Cache: &test.Cache{Value: "d2hhdD8="}},
		{Config: test.NewCacheConfig("sync", "error", "json", "redis"), Cache: &test.Cache{Value: "d2hhdD8="}},
		{Config: test.NewCacheConfig("sync", "snappy", "json", "redis"), Cache: &test.Cache{Value: "what?"}},
		{Config: test.NewCacheConfig("sync", "snappy", "json", "redis"), Cache: &test.ErrCache{}},
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
				Cache:      tuple.Cache,
				Tracer:     world.NewTracer(),
				Logger:     world.Logger,
				Meter:      world.Server.Meter,
			}

			cache, err := cache.New(params)
			So(err, ShouldBeNil)

			world.RequireStart()

			ctx := context.Background()

			Convey("When I try to encode a value", func() {
				ptr := ptr.Zero[string]()
				err := cache.Get(ctx, "test", ptr)

				Convey("Then I should have an error", func() {
					So(err, ShouldBeError)
				})
			})

			world.RequireStop()
		})
	}
}
