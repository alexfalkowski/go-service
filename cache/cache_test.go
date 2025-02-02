package cache_test

import (
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/cache"
	"github.com/alexfalkowski/go-service/test"
	"github.com/alexfalkowski/go-service/types"
	"github.com/faabiosr/cachego"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestValidCache(t *testing.T) {
	configs := []*cache.Config{
		{Kind: "redis", Encoder: "json", Compressor: "snappy", Options: map[string]any{"url": test.Path("secrets/redis")}},
		{Encoder: "json", Compressor: "snappy"},
	}

	for _, config := range configs {
		Convey("Given I have a cache", t, func() {
			world := test.NewWorld(t)
			world.Register()

			params := cache.Params{
				Lifecycle:  world.Lifecycle,
				Config:     config,
				Compressor: test.Compressor,
				Encoder:    test.Encoder,
				Pool:       test.Pool,
			}

			ca, err := cache.New(params)
			So(err, ShouldBeNil)

			cache.Register(ca)
			world.RequireStart()

			Convey("When I save an item", func() {
				value := "hello?"
				err := cache.Persist("test", &value, time.Minute)

				Convey("Then I should have no error", func() {
					So(err, ShouldBeNil)
				})
			})

			Convey("When I get an item", func() {
				value := "wassup?"

				err := cache.Persist("test", &value, time.Minute)
				So(err, ShouldBeNil)

				ptr := types.Pointer[string]()

				err = cache.Get("test", ptr)
				So(err, ShouldBeNil)

				Convey("Then I should have a value", func() {
					So(*ptr, ShouldEqual, "wassup?")
				})
			})

			err = cache.Remove("test")
			So(err, ShouldBeNil)

			world.RequireStop()
		})
	}
}

func TestErroneousCache(t *testing.T) {
	configs := []*cache.Config{
		{Kind: "redis", Encoder: "json", Compressor: "snappy", Options: map[string]any{"url": test.Path("secrets/none")}},
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
			}

			_, err := cache.New(params)

			world.RequireStart()

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})

			world.RequireStop()
		})
	}
}

func TestErroneousSave(t *testing.T) {
	configs := []*cache.Config{
		{Encoder: "error", Compressor: "snappy"},
	}

	for _, config := range configs {
		Convey("Given I have a cache", t, func() {
			world := test.NewWorld(t)
			world.Register()

			params := cache.Params{
				Lifecycle:  world.Lifecycle,
				Config:     config,
				Compressor: test.Compressor,
				Encoder:    test.Encoder,
				Pool:       test.Pool,
			}

			ca, err := cache.New(params)
			So(err, ShouldBeNil)

			cache.Register(ca)
			world.RequireStart()

			Convey("When I try to save a value", func() {
				value := "what?"
				err := cache.Persist("test", &value, time.Minute)

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
		Config *cache.Config
		cachego.Cache
	}

	tuples := []*Tuple{
		{Config: &cache.Config{Encoder: "error", Compressor: "snappy"}, Cache: &test.Cache{Value: "d2hhdD8="}},
		{Config: &cache.Config{Encoder: "json", Compressor: "error"}, Cache: &test.Cache{Value: "d2hhdD8="}},
		{Config: &cache.Config{Encoder: "json", Compressor: "snappy"}, Cache: &test.Cache{Value: "what?"}},
		{Config: &cache.Config{Encoder: "json", Compressor: "snappy"}, Cache: &test.ErrCache{}},
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
			}

			ca, err := cache.New(params)
			So(err, ShouldBeNil)

			ca.Cache = tuple.Cache

			cache.Register(ca)

			Convey("When I try to encode a value", func() {
				ptr := types.Pointer[string]()
				err := cache.Get("test", ptr)

				Convey("Then I should have an error", func() {
					So(err, ShouldBeError)
				})

				world.RequireStop()
			})
		})
	}
}
