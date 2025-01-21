package redis_test

import (
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestClientIncr(t *testing.T) {
	Convey("Given I have a cache", t, func() {
		world := test.NewWorld(t)
		world.Register()

		client, err := world.NewRedisClient()
		So(err, ShouldBeNil)

		ctx, cancel := test.Timeout()
		defer cancel()

		ctx = meta.WithAttribute(ctx, "password", meta.Redacted("test-1234"))

		world.RequireStart()

		Convey("When I try to cache an item", func() {
			cmd := client.Set(ctx, "test-incr", 1, time.Hour)
			So(cmd.Err(), ShouldBeNil)

			Convey("Then I should have a cached item", func() {
				cmd := client.Incr(ctx, "test-incr")
				So(cmd.Err(), ShouldBeNil)

				r, err := cmd.Result()
				So(err, ShouldBeNil)

				So(r, ShouldEqual, 2)

				d := client.Del(ctx, "test-incr")
				So(d.Err(), ShouldBeNil)
			})
		})

		world.RequireStop()
	})
}
