//nolint:varnamelen
package redis_test

import (
	"context"
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/test"
	"github.com/go-redis/cache/v9"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"google.golang.org/grpc/health/grpc_health_v1"
)

func TestSetCache(t *testing.T) {
	Convey("Given I have a cache", t, func() {
		world := test.NewWorld(t)
		world.Register()

		ca, err := world.NewRedisCache()
		So(err, ShouldBeNil)

		ctx, cancel := test.Timeout()
		defer cancel()

		ctx = meta.WithAttribute(ctx, "test", meta.String("test"))

		world.RequireStart()

		Convey("When I try to cache an item", func() {
			value := &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}
			err := ca.Set(&cache.Item{Ctx: ctx, Key: "test", Value: value, TTL: time.Minute})
			So(err, ShouldBeNil)

			Convey("Then I should have a cached item", func() {
				var v grpc_health_v1.HealthCheckResponse

				err := ca.Get(ctx, "test", &v)
				So(err, ShouldBeNil)

				So(v.GetStatus(), ShouldEqual, grpc_health_v1.HealthCheckResponse_SERVING)

				err = ca.Delete(ctx, "test")
				So(err, ShouldBeNil)
			})
		})

		world.RequireStop()
	})
}

func TestSetXXCache(t *testing.T) {
	Convey("Given I have a cache", t, func() {
		world := test.NewWorld(t)
		world.Register()

		ca, err := world.NewRedisCache()
		So(err, ShouldBeNil)

		ctx, cancel := test.Timeout()
		defer cancel()

		ctx = meta.WithAttribute(ctx, "test", meta.String("test"))

		world.RequireStart()

		Convey("When I try to cache an item", func() {
			value := &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}
			err := ca.Set(&cache.Item{Ctx: ctx, Key: "test", Value: value, TTL: time.Minute, SetXX: true})
			So(err, ShouldBeNil)

			Convey("Then I should have a cached item", func() {
				var v grpc_health_v1.HealthCheckResponse

				err := ca.Get(ctx, "test", &v)
				So(err, ShouldBeError)
			})
		})

		world.RequireStop()
	})
}

func TestSetNXCache(t *testing.T) {
	Convey("Given I have a cache", t, func() {
		world := test.NewWorld(t)
		world.Register()

		ca, err := world.NewRedisCache()
		So(err, ShouldBeNil)

		ctx, cancel := test.Timeout()
		defer cancel()

		ctx = meta.WithAttribute(ctx, "test", meta.String("test"))

		world.RequireStart()

		Convey("When I try to cache an item", func() {
			value := &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}
			err := ca.Set(&cache.Item{Ctx: ctx, Key: "test", Value: value, TTL: time.Minute, SetNX: true})
			So(err, ShouldBeNil)

			Convey("Then I should have a cached item", func() {
				var v grpc_health_v1.HealthCheckResponse

				err := ca.Get(ctx, "test", &v)
				So(err, ShouldBeNil)

				So(v.GetStatus(), ShouldEqual, grpc_health_v1.HealthCheckResponse_SERVING)

				err = ca.Delete(ctx, "test")
				So(err, ShouldBeNil)
			})
		})

		world.RequireStop()
	})
}

func TestInvalidHostCache(t *testing.T) {
	Convey("Given I have a cache", t, func() {
		world := test.NewWorld(t, test.WithWorldRedisConfig(test.NewRedisConfig("redis_invalid", "snappy", "proto")))
		world.Register()

		ca, err := world.NewRedisCache()
		So(err, ShouldBeNil)

		world.RequireStart()

		Convey("When I try to cache an item", func() {
			value := &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}
			err := ca.Set(&cache.Item{Ctx: context.Background(), Key: "test", Value: value, TTL: time.Minute})

			Convey("Then I should have an error", func() {
				So(err, ShouldNotBeNil)
			})
		})

		world.RequireStop()
	})
}

func TestInvalidMarshallerCache(t *testing.T) {
	Convey("Given I have a cache", t, func() {
		world := test.NewWorld(t, test.WithWorldRedisConfig(test.NewRedisConfig("redis", "snappy", "error")))
		world.Register()

		ca, err := world.NewRedisCache()
		So(err, ShouldBeNil)

		world.RequireStart()

		Convey("When I try to cache an item", func() {
			value := &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}
			err := ca.Set(&cache.Item{Ctx: context.Background(), Key: "test", Value: value, TTL: time.Minute})

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, "failed")
			})
		})

		world.RequireStop()
	})
}

func TestMissingMarshallerCache(t *testing.T) {
	Convey("When I try to create a cache", t, func() {
		world := test.NewWorld(t, test.WithWorldRedisConfig(test.NewRedisConfig("redis", "snappy", "test")))
		world.Register()
		world.RequireStart()

		_, err := world.NewRedisCache()

		Convey("Then I should have no error", func() {
			So(err, ShouldBeNil)
		})

		world.RequireStop()
	})
}

func TestInvalidCompressorCache(t *testing.T) {
	Convey("Given I have a cache", t, func() {
		world := test.NewWorld(t, test.WithWorldRedisConfig(test.NewRedisConfig("redis", "error", "proto")))
		world.Register()

		ca, err := world.NewRedisCache()
		So(err, ShouldBeNil)

		ctx := context.Background()

		world.RequireStart()

		Convey("When I try to cache an item", func() {
			value := &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}
			err := ca.Set(&cache.Item{Ctx: context.Background(), Key: "test", Value: value, TTL: time.Minute})
			So(err, ShouldBeNil)

			Convey("Then I should have an error", func() {
				var v grpc_health_v1.HealthCheckResponse

				err := ca.Get(ctx, "test", &v)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "failed")
			})
		})

		world.RequireStop()
	})
}

func TestMissingCompressorCache(t *testing.T) {
	Convey("When I try to create a cache", t, func() {
		world := test.NewWorld(t, test.WithWorldRedisConfig(test.NewRedisConfig("redis", "test", "proto")))
		world.Register()
		world.RequireStart()

		_, err := world.NewRedisCache()

		Convey("Then I should have no error", func() {
			So(err, ShouldBeNil)
		})

		world.RequireStop()
	})
}
