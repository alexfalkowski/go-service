package redis_test

import (
	"context"
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/otel"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func init() {
	otel.Register()
}

func TestClientIncr(t *testing.T) {
	Convey("Given I have a cache", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		c := test.NewRedisClient(lc, "localhost:6379", logger)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		ctx = meta.WithAttribute(ctx, "test", "test")

		lc.RequireStart()

		Convey("When I try to cache an item", func() {
			cmd := c.Set(ctx, "test-incr", 1, time.Hour)
			So(cmd.Err(), ShouldBeNil)

			Convey("Then I should have a cached item", func() {
				cmd := c.Incr(ctx, "test-incr")
				So(cmd.Err(), ShouldBeNil)

				r, err := cmd.Result()
				So(err, ShouldBeNil)

				So(r, ShouldEqual, 2)

				d := c.Del(ctx, "test-incr")
				So(d.Err(), ShouldBeNil)
			})
		})

		lc.RequireStop()
	})
}
