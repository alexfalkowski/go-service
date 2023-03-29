package redis_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/cache/compressor"
	"github.com/alexfalkowski/go-service/cache/marshaller"
	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/otel"
	"github.com/alexfalkowski/go-service/test"
	"github.com/go-redis/cache/v8"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func init() {
	otel.Register()
}

func TestSetCache(t *testing.T) {
	Convey("Given I have a cache", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		c := test.NewRedisCache(lc, "localhost:6379", logger, compressor.NewSnappy(), marshaller.NewProto())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		ctx = meta.WithAttribute(ctx, "test", "test")

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

				err = c.Delete(ctx, "test")
				So(err, ShouldBeNil)
			})
		})

		lc.RequireStop()
	})
}

func TestSetXXCache(t *testing.T) {
	Convey("Given I have a cache", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		c := test.NewRedisCache(lc, "localhost:6379", logger, compressor.NewSnappy(), marshaller.NewProto())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		ctx = meta.WithAttribute(ctx, "test", "test")

		lc.RequireStart()

		Convey("When I try to cache an item", func() {
			value := &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}
			err := c.Set(&cache.Item{Ctx: ctx, Key: "test", Value: value, TTL: time.Minute, SetXX: true})
			So(err, ShouldBeNil)

			Convey("Then I should have a cached item", func() {
				var v grpc_health_v1.HealthCheckResponse

				err := c.Get(ctx, "test", &v)
				So(err, ShouldBeError)
			})
		})

		lc.RequireStop()
	})
}

func TestSetNXCache(t *testing.T) {
	Convey("Given I have a cache", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		c := test.NewRedisCache(lc, "localhost:6379", logger, compressor.NewSnappy(), marshaller.NewProto())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		ctx = meta.WithAttribute(ctx, "test", "test")

		lc.RequireStart()

		Convey("When I try to cache an item", func() {
			value := &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}
			err := c.Set(&cache.Item{Ctx: ctx, Key: "test", Value: value, TTL: time.Minute, SetNX: true})
			So(err, ShouldBeNil)

			Convey("Then I should have a cached item", func() {
				var v grpc_health_v1.HealthCheckResponse

				err := c.Get(ctx, "test", &v)
				So(err, ShouldBeNil)

				So(v.Status, ShouldEqual, grpc_health_v1.HealthCheckResponse_SERVING)

				err = c.Delete(ctx, "test")
				So(err, ShouldBeNil)
			})
		})

		lc.RequireStop()
	})
}

func TestInvalidHostCache(t *testing.T) {
	Convey("Given I have a cache", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		c := test.NewRedisCache(lc, "invalid_host", logger, compressor.NewSnappy(), marshaller.NewProto())
		ctx := context.Background()

		lc.RequireStart()

		Convey("When I try to cache an item", func() {
			value := &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}
			err := c.Set(&cache.Item{Ctx: ctx, Key: "test", Value: value, TTL: time.Minute})

			Convey("Then I should have an error", func() {
				So(err, ShouldNotBeNil)
			})
		})

		lc.RequireStop()
	})
}

func TestInvalidMarshallerCache(t *testing.T) {
	Convey("Given I have a cache", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		c := test.NewRedisCache(lc, "localhost:6379", logger, compressor.NewSnappy(), test.NewMarshaller(errors.New("failed")))
		ctx := context.Background()

		lc.RequireStart()

		Convey("When I try to cache an item", func() {
			value := &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}
			err := c.Set(&cache.Item{Ctx: ctx, Key: "test", Value: value, TTL: time.Minute})

			Convey("Then I should have an error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "failed")
			})
		})

		lc.RequireStop()
	})
}

func TestInvalidCompressorCache(t *testing.T) {
	Convey("Given I have a cache", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		c := test.NewRedisCache(lc, "localhost:6379", logger, test.NewCompressor(errors.New("failed")), marshaller.NewProto())
		ctx := context.Background()

		lc.RequireStart()

		Convey("When I try to cache an item", func() {
			value := &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}
			err := c.Set(&cache.Item{Ctx: ctx, Key: "test", Value: value, TTL: time.Minute})
			So(err, ShouldBeNil)

			Convey("Then I should have an error", func() {
				var v grpc_health_v1.HealthCheckResponse

				err := c.Get(ctx, "test", &v)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "failed")
			})
		})

		lc.RequireStop()
	})
}
