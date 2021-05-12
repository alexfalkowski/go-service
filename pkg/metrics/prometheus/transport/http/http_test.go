package http_test

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/alexfalkowski/go-service/pkg/cache/redis"
	"github.com/alexfalkowski/go-service/pkg/cache/ristretto"
	"github.com/alexfalkowski/go-service/pkg/config"
	"github.com/alexfalkowski/go-service/pkg/logger/zap"
	prometheusHTTP "github.com/alexfalkowski/go-service/pkg/metrics/prometheus/transport/http"
	"github.com/alexfalkowski/go-service/pkg/sql"
	pkgHTTP "github.com/alexfalkowski/go-service/pkg/transport/http"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestHTTP(t *testing.T) {
	Convey("Given I register the metrics handler", t, func() {
		lc := fxtest.NewLifecycle(t)
		mux := pkgHTTP.NewMux()

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cfg := &config.Config{AppName: "test", HTTPPort: "10002"}
		rcfg := &redis.Config{Host: "localhost:6379"}

		_, err = sql.NewDB(lc, &sql.Config{PostgresURL: "postgres://test:test@localhost:5432/test?sslmode=disable"})
		So(err, ShouldBeNil)

		r := redis.NewRing(lc, rcfg)
		opts := redis.NewOptions(r)
		_ = redis.NewCache(lc, rcfg, opts)

		_, err = ristretto.NewCache(lc, cfg, ristretto.NewConfig())
		So(err, ShouldBeNil)

		pkgHTTP.Register(lc, test.NewShutdowner(), mux, cfg, logger)

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
	})
}
