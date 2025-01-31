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
		{Kind: "redis", Options: map[string]any{"url": test.Path("secrets/redis")}},
		{},
	}

	for _, config := range configs {
		Convey("Given I have a cache", t, func() {
			world := test.NewWorld(t)
			world.Register()

			cache, err := cache.NewCache(world.Lifecycle, config)
			So(err, ShouldBeNil)

			world.RequireStart()

			Convey("When I save an item", func() {
				err = cache.Save("test", "what?", time.Minute)

				Convey("Then I should have no error", func() {
					So(err, ShouldBeNil)
				})

				err = cache.Delete("test")
				So(err, ShouldBeNil)
			})

			world.RequireStop()
		})
	}
}

func TestErroneousCache(t *testing.T) {
	configs := []*cache.Config{
		{Kind: "redis", Options: map[string]any{"url": test.Path("secrets/none")}},
	}

	for _, config := range configs {
		Convey("When I create a cache", t, func() {
			world := test.NewWorld(t)
			world.Register()

			_, err := cache.NewCache(world.Lifecycle, config)

			world.RequireStart()

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})

			world.RequireStop()
		})
	}
}
