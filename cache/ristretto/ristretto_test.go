package ristretto_test

import (
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/cache/ristretto"
	"github.com/alexfalkowski/go-service/version"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// nolint:forcetypeassert
func TestCache(t *testing.T) {
	Convey("Given I have a cache", t, func() {
		cfg := &ristretto.Config{
			NumCounters: 1e7,
			MaxCost:     1 << 30,
			BufferItems: 64,
		}
		lc := fxtest.NewLifecycle(t)

		c, err := ristretto.NewCache(ristretto.CacheParams{Lifecycle: lc, Config: cfg, Version: version.Version("1.0.0")})
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I try to cache an item", func() {
			value := &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}

			ok := c.SetWithTTL("test", value, 0, time.Minute)
			So(ok, ShouldBeTrue)

			time.Sleep(1 * time.Second)

			Convey("Then I should have a cached item", func() {
				v, ok := c.Get("test")
				So(ok, ShouldBeTrue)

				r := v.(*grpc_health_v1.HealthCheckResponse)

				So(r.Status, ShouldEqual, grpc_health_v1.HealthCheckResponse_SERVING)
			})
		})

		lc.RequireStop()
	})
}

func TestInvalidCache(t *testing.T) {
	Convey("Given I have an invalid config", t, func() {
		lc := fxtest.NewLifecycle(t)
		cfg := &ristretto.Config{}

		lc.RequireStart()

		Convey("When I try to create a cache", func() {
			_, err := ristretto.NewCache(ristretto.CacheParams{Lifecycle: lc, Config: cfg, Version: version.Version("1.0.0")})

			Convey("Then I should have an error", func() {
				So(err, ShouldNotBeNil)
			})
		})

		lc.RequireStop()
	})
}
