package http_test

import (
	"net/http"
	"testing"

	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestPrometheusInsecureHTTP(t *testing.T) {
	Convey("Given I register the metrics handler", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("prometheus"), test.WithWorldLimiter(test.NewLimiterConfig("user-agent", "1s", 100)))
		world.Register()

		_, err := world.OpenDatabase()
		So(err, ShouldBeNil)

		_, err = world.NewRedisCache()
		So(err, ShouldBeNil)

		world.RequireStart()

		Convey("When I query metrics", func() {
			ctx, cancel := test.Timeout()
			defer cancel()

			header := http.Header{}

			res, body, err := world.ResponseWithBody(ctx, "http", world.ServerHost(), http.MethodGet, "metrics", header, http.NoBody)
			So(err, ShouldBeNil)

			Convey("Then I should have valid metrics", func() {
				So(res.StatusCode, ShouldEqual, http.StatusOK)

				So(body, ShouldContainSubstring, "go_info")
				So(body, ShouldContainSubstring, "redis_hits_total")
				So(body, ShouldContainSubstring, "sql_max_open_total")
				So(body, ShouldContainSubstring, "system")
				So(body, ShouldContainSubstring, "process")
				So(body, ShouldContainSubstring, "runtime")
			})
		})

		world.RequireStop()
	})
}

func TestPrometheusSecureHTTP(t *testing.T) {
	Convey("Given I register the metrics handler", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("prometheus"), test.WithWorldSecure())
		world.Register()

		_, err := world.OpenDatabase()
		So(err, ShouldBeNil)

		_, err = world.NewRedisCache()
		So(err, ShouldBeNil)

		world.RequireStart()

		Convey("When I query metrics", func() {
			ctx, cancel := test.Timeout()
			defer cancel()

			header := http.Header{}

			res, body, err := world.ResponseWithBody(ctx, "https", world.ServerHost(), http.MethodGet, "metrics", header, http.NoBody)
			So(err, ShouldBeNil)

			Convey("Then I should have valid metrics", func() {
				So(res.StatusCode, ShouldEqual, http.StatusOK)

				So(body, ShouldContainSubstring, "go_info")
				So(body, ShouldContainSubstring, "redis_hits_total")
				So(body, ShouldContainSubstring, "sql_max_open_total")
			})
		})

		world.RequireStop()
	})
}
