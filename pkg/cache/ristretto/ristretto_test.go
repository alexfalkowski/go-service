package ristretto_test

import (
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/pkg/cache/ristretto"
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

		c, err := ristretto.NewCache(lc, cfg)
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
