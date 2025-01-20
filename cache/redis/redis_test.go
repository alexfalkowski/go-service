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

		ca, err := world.Cache.NewRedisCache()
		So(err, ShouldBeNil)

		ctx, cancel := test.Timeout()
		defer cancel()

		ctx = meta.WithAttribute(ctx, "test", meta.String("test"))

		world.Start()

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

		world.Stop()
	})
}

func TestSetXXCache(t *testing.T) {
	Convey("Given I have a cache", t, func() {
		world := test.NewWorld(t)

		ca, err := world.Cache.NewRedisCache()
		So(err, ShouldBeNil)

		ctx, cancel := test.Timeout()
		defer cancel()

		ctx = meta.WithAttribute(ctx, "test", meta.String("test"))

		world.Start()

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

		world.Stop()
	})
}

func TestSetNXCache(t *testing.T) {
	Convey("Given I have a cache", t, func() {
		world := test.NewWorld(t)

		ca, err := world.Cache.NewRedisCache()
		So(err, ShouldBeNil)

		ctx, cancel := test.Timeout()
		defer cancel()

		ctx = meta.WithAttribute(ctx, "test", meta.String("test"))

		world.Start()

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

		world.Stop()
	})
}

func TestInvalidHostCache(t *testing.T) {
	Convey("Given I have a cache", t, func() {
		world := test.NewWorld(t, test.WithWorldRedisConfig(test.NewRedisConfig("redis_invalid", "snappy", "proto")))

		ca, err := world.Cache.NewRedisCache()
		So(err, ShouldBeNil)

		world.Start()

		Convey("When I try to cache an item", func() {
			value := &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}
			err := ca.Set(&cache.Item{Ctx: context.Background(), Key: "test", Value: value, TTL: time.Minute})

			Convey("Then I should have an error", func() {
				So(err, ShouldNotBeNil)
			})
		})

		world.Stop()
	})
}

func TestInvalidMarshallerCache(t *testing.T) {
	Convey("Given I have a cache", t, func() {
		world := test.NewWorld(t, test.WithWorldRedisConfig(test.NewRedisConfig("redis", "snappy", "error")))

		ca, err := world.Cache.NewRedisCache()
		So(err, ShouldBeNil)

		world.Start()

		Convey("When I try to cache an item", func() {
			value := &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}
			err := ca.Set(&cache.Item{Ctx: context.Background(), Key: "test", Value: value, TTL: time.Minute})

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, "failed")
			})
		})

		world.Stop()
	})
}

func TestMissingMarshallerCache(t *testing.T) {
	Convey("When I try to create a cache", t, func() {
		world := test.NewWorld(t, test.WithWorldRedisConfig(test.NewRedisConfig("redis", "snappy", "test")))
		world.Start()

		_, err := world.Cache.NewRedisCache()

		Convey("Then I should have no error", func() {
			So(err, ShouldBeNil)
		})

		world.Stop()
	})
}

func TestInvalidCompressorCache(t *testing.T) {
	Convey("Given I have a cache", t, func() {
		world := test.NewWorld(t, test.WithWorldRedisConfig(test.NewRedisConfig("redis", "error", "proto")))

		ca, err := world.Cache.NewRedisCache()
		So(err, ShouldBeNil)

		ctx := context.Background()

		world.Start()

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

		world.Stop()
	})
}

func TestMissingCompressorCache(t *testing.T) {
	Convey("When I try to create a cache", t, func() {
		world := test.NewWorld(t, test.WithWorldRedisConfig(test.NewRedisConfig("redis", "test", "proto")))
		world.Start()

		_, err := world.Cache.NewRedisCache()

		Convey("Then I should have no error", func() {
			So(err, ShouldBeNil)
		})

		world.Stop()
	})
}
