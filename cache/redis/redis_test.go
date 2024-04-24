package redis_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/test"
	"github.com/go-redis/cache/v8"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func init() {
	tracer.Register()
	test.Marshaller.Register("error", test.NewMarshaller(errors.New("failed")))
	test.Compressor.Register("error", test.NewCompressor(errors.New("failed")))
}

func TestSetCache(t *testing.T) {
	Convey("Given I have a cache", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		m := test.NewMeter(lc)
		c, _ := test.NewRedisCache(lc, test.NewRedisConfig("localhost:6379", "snappy", "proto"), logger, m)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		ctx = meta.WithAttribute(ctx, "test", meta.String("test"))

		lc.RequireStart()

		Convey("When I try to cache an item", func() {
			value := &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}
			err := c.Set(&cache.Item{Ctx: ctx, Key: "test", Value: value, TTL: time.Minute})
			So(err, ShouldBeNil)

			Convey("Then I should have a cached item", func() {
				var v grpc_health_v1.HealthCheckResponse

				err := c.Get(ctx, "test", &v)
				So(err, ShouldBeNil)

				So(v.GetStatus(), ShouldEqual, grpc_health_v1.HealthCheckResponse_SERVING)

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
		m := test.NewMeter(lc)
		c, _ := test.NewRedisCache(lc, test.NewRedisConfig("localhost:6379", "snappy", "proto"), logger, m)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		ctx = meta.WithAttribute(ctx, "test", meta.String("test"))

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
		m := test.NewMeter(lc)
		c, _ := test.NewRedisCache(lc, test.NewRedisConfig("localhost:6379", "snappy", "proto"), logger, m)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		ctx = meta.WithAttribute(ctx, "test", meta.String("test"))

		lc.RequireStart()

		Convey("When I try to cache an item", func() {
			value := &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}
			err := c.Set(&cache.Item{Ctx: ctx, Key: "test", Value: value, TTL: time.Minute, SetNX: true})
			So(err, ShouldBeNil)

			Convey("Then I should have a cached item", func() {
				var v grpc_health_v1.HealthCheckResponse

				err := c.Get(ctx, "test", &v)
				So(err, ShouldBeNil)

				So(v.GetStatus(), ShouldEqual, grpc_health_v1.HealthCheckResponse_SERVING)

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
		m := test.NewMeter(lc)
		c, _ := test.NewRedisCache(lc, test.NewRedisConfig("invalid_host", "snappy", "proto"), logger, m)
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
		m := test.NewMeter(lc)

		c, _ := test.NewRedisCache(lc, test.NewRedisConfig("localhost:6379", "snappy", "error"), logger, m)
		ctx := context.Background()

		lc.RequireStart()

		Convey("When I try to cache an item", func() {
			value := &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}
			err := c.Set(&cache.Item{Ctx: ctx, Key: "test", Value: value, TTL: time.Minute})

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, "failed")
			})
		})

		lc.RequireStop()
	})
}

func TestMissingMarshallerCache(t *testing.T) {
	Convey("When I try to create a cache", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		m := test.NewMeter(lc)

		_, err := test.NewRedisCache(lc, test.NewRedisConfig("localhost:6379", "snappy", "test"), logger, m)

		lc.RequireStart()

		Convey("Then I should have an error", func() {
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "invalid kind")
		})

		lc.RequireStop()
	})
}

func TestInvalidCompressorCache(t *testing.T) {
	Convey("Given I have a cache", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		m := test.NewMeter(lc)
		c, _ := test.NewRedisCache(lc, test.NewRedisConfig("localhost:6379", "error", "proto"), logger, m)
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

func TestMissingCompressorCache(t *testing.T) {
	Convey("When I try to create a cache", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		m := test.NewMeter(lc)

		_, err := test.NewRedisCache(lc, test.NewRedisConfig("localhost:6379", "test", "proto"), logger, m)

		lc.RequireStart()

		Convey("Then I should have an error", func() {
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "invalid kind")
		})

		lc.RequireStop()
	})
}
