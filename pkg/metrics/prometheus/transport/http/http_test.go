package http_test

import (
	"context"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/alexfalkowski/go-service/pkg/cache/redis"
	"github.com/alexfalkowski/go-service/pkg/cache/ristretto"
	"github.com/alexfalkowski/go-service/pkg/logger/zap"
	prometheusHTTP "github.com/alexfalkowski/go-service/pkg/metrics/prometheus/transport/http"
	"github.com/alexfalkowski/go-service/pkg/sql/pg"
	pkgHTTP "github.com/alexfalkowski/go-service/pkg/transport/http"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestHTTP(t *testing.T) {
	Convey("Given I register the metrics handler", t, func() {
		os.Setenv("SERVICE_NAME", "test")

		lc := fxtest.NewLifecycle(t)
		mux := pkgHTTP.NewMux()

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		rcfg := &redis.Config{Host: "localhost:6379"}

		_, err = pg.NewDB(lc, &pg.Config{URL: "postgres://test:test@localhost:5432/test?sslmode=disable"})
		So(err, ShouldBeNil)

		r := redis.NewRing(lc, rcfg)
		opts := redis.NewOptions(r)
		_ = redis.NewCache(lc, rcfg, opts)

		ricfg, err := ristretto.NewConfig()
		So(err, ShouldBeNil)

		_, err = ristretto.NewCache(lc, ricfg)
		So(err, ShouldBeNil)

		pkgHTTP.Register(lc, test.NewShutdowner(), mux, &pkgHTTP.Config{Port: "10002"}, logger)

		err = prometheusHTTP.Register(mux)
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I query metrics", func() {
			client := pkgHTTP.NewClient(logger)

			req, err := http.NewRequestWithContext(context.Background(), "GET", "http://localhost:10002/metrics", nil)
			So(err, ShouldBeNil)

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			Convey("Then I should have valid metrics", func() {
				response := string(body)

				So(response, ShouldContainSubstring, "go_info")
				So(response, ShouldContainSubstring, "go_sql_stats")
				So(response, ShouldContainSubstring, "go_redis_stats")
				So(response, ShouldContainSubstring, "go_ristretto_stats")
			})
		})

		lc.RequireStop()
		So(os.Unsetenv("SERVICE_NAME"), ShouldBeNil)
	})
}
