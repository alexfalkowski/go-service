package cache_test

import (
	"context"
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/pkg/cache"
	"github.com/alexfalkowski/go-service/pkg/config"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func TestCache(t *testing.T) {
	Convey("Given I have a cache", t, func() {
		cfg := &config.Config{CacheHost: "localhost:6379"}
		lc := fxtest.NewLifecycle(t)
		c := cache.NewCache(lc, cfg)
		ctx := context.TODO()

		lc.RequireStart()

		Convey("When I try to cache an item", func() {
			value := &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}
			err := c.Set(&cache.Item{Ctx: ctx, Key: "test", Value: value, TTL: time.Minute})
			So(err, ShouldBeNil)

			Convey("Then I should have a cached item", func() {
				var v grpc_health_v1.HealthCheckResponse

				err := c.Get(ctx, "test", &v)
				So(err, ShouldBeNil)

				So(v.Status, ShouldEqual, grpc_health_v1.HealthCheckResponse_SERVING)

				lc.RequireStop()
			})
		})
	})
}
