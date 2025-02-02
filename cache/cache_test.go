package cache_test

import (
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/cache"
	"github.com/alexfalkowski/go-service/test"
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

			cache, err := cache.New(params)
			So(err, ShouldBeNil)

			world.RequireStart()

			Convey("When I save an item", func() {
				value, err := cache.EncodeValue("what?")
				So(err, ShouldBeNil)

				err = cache.Save("test", value, time.Minute)

				Convey("Then I should have no error", func() {
					So(err, ShouldBeNil)
				})
			})

			Convey("When I get an item", func() {
				e, err := cache.EncodeValue("what?")
				So(err, ShouldBeNil)

				err = cache.Save("test", e, time.Minute)
				So(err, ShouldBeNil)

				v, err := cache.Fetch("test")
				So(err, ShouldBeNil)

				var value string

				err = cache.DecodeValue(v, &value)
				So(err, ShouldBeNil)

				Convey("Then I should have a value", func() {
					So(value, ShouldEqual, "what?")
				})
			})

			err = cache.Delete("test")
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

func TestErroneousEncode(t *testing.T) {
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

			cache, err := cache.New(params)
			So(err, ShouldBeNil)

			world.RequireStart()

			Convey("When I try to encode a value", func() {
				_, err := cache.EncodeValue("what?")

				Convey("Then I should have an error", func() {
					So(err, ShouldBeError)
				})

				world.RequireStop()
			})
		})
	}
}

func TestErroneousDecode(t *testing.T) {
	type Tuple struct {
		Config *cache.Config
		Value  string
	}

	tuples := []*Tuple{
		{Value: "d2hhdD8=", Config: &cache.Config{Encoder: "error", Compressor: "snappy"}},
		{Value: "d2hhdD8=", Config: &cache.Config{Encoder: "json", Compressor: "error"}},
		{Value: "what?", Config: &cache.Config{Encoder: "json", Compressor: "snappy"}},
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

			cache, err := cache.New(params)
			So(err, ShouldBeNil)

			world.RequireStart()

			Convey("When I try to encode a value", func() {
				var str string

				err := cache.DecodeValue(tuple.Value, &str)

				Convey("Then I should have an error", func() {
					So(err, ShouldBeError)
				})

				world.RequireStop()
			})
		})
	}
}
