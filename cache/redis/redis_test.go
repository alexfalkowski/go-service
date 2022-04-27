package redis_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/cache/compressor"
	"github.com/alexfalkowski/go-service/cache/marshaller"
	"github.com/alexfalkowski/go-service/cache/redis"
	"github.com/alexfalkowski/go-service/cache/redis/trace/opentracing"
	"github.com/alexfalkowski/go-service/cache/redis/trace/opentracing/jaeger"
	"github.com/alexfalkowski/go-service/test"
	"github.com/go-redis/cache/v8"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func TestCache(t *testing.T) {
	Convey("Given I have a cache", t, func() {
		cfg := &redis.Config{Host: "localhost:6379"}
		lc := fxtest.NewLifecycle(t)

		r := redis.NewRing(lc, cfg)
		params := redis.OptionsParams{Ring: r, Compressor: compressor.NewSnappy(), Marshaller: marshaller.NewProto()}
		opts := redis.NewOptions(params)

		c := redis.NewCache(lc, cfg, opts)
		ctx := context.Background()

		tracer, err := jaeger.NewTracer(lc, test.NewJaegerConfig())
		So(err, ShouldBeNil)

		ctx, span := opentracing.StartSpanFromContext(ctx, tracer, "test", "test")
		defer span.Finish()

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
			})
		})

		lc.RequireStop()
	})
}

func TestInvalidHostCache(t *testing.T) {
	Convey("Given I have a cache", t, func() {
		cfg := &redis.Config{Host: "invalid_host"}
		lc := fxtest.NewLifecycle(t)

		r := redis.NewRing(lc, cfg)
		params := redis.OptionsParams{Ring: r, Compressor: compressor.NewSnappy(), Marshaller: marshaller.NewProto()}
		opts := redis.NewOptions(params)

		c := redis.NewCache(lc, cfg, opts)
		ctx := context.Background()

		tracer, err := jaeger.NewTracer(lc, test.NewJaegerConfig())
		So(err, ShouldBeNil)

		ctx, span := opentracing.StartSpanFromContext(ctx, tracer, "test", "test")
		defer span.Finish()

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

// nolint:goerr113
func TestInvalidMarshallerCache(t *testing.T) {
	Convey("Given I have a cache", t, func() {
		cfg := &redis.Config{Host: "localhost:6379"}
		lc := fxtest.NewLifecycle(t)

		r := redis.NewRing(lc, cfg)
		params := redis.OptionsParams{Ring: r, Compressor: compressor.NewSnappy(), Marshaller: test.NewMarshaller(errors.New("failed"))}
		opts := redis.NewOptions(params)

		c := redis.NewCache(lc, cfg, opts)
		ctx := context.Background()

		tracer, err := jaeger.NewTracer(lc, test.NewJaegerConfig())
		So(err, ShouldBeNil)

		ctx, span := opentracing.StartSpanFromContext(ctx, tracer, "test", "test")
		defer span.Finish()

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

// nolint:goerr113
func TestInvalidCompressorCache(t *testing.T) {
	Convey("Given I have a cache", t, func() {
		cfg := &redis.Config{Host: "localhost:6379"}
		lc := fxtest.NewLifecycle(t)

		r := redis.NewRing(lc, cfg)
		params := redis.OptionsParams{Ring: r, Compressor: test.NewCompressor(errors.New("failed")), Marshaller: marshaller.NewProto()}
		opts := redis.NewOptions(params)

		c := redis.NewCache(lc, cfg, opts)
		ctx := context.Background()

		tracer, err := jaeger.NewTracer(lc, test.NewJaegerConfig())
		So(err, ShouldBeNil)

		ctx, span := opentracing.StartSpanFromContext(ctx, tracer, "test", "test")
		defer span.Finish()

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
