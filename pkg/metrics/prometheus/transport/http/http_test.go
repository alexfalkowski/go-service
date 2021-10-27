package http_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
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
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		rcfg := &redis.Config{
			Host: "localhost:6379",
		}

		_, err = pg.NewDB(lc, &pg.Config{URL: "postgres://test:test@localhost:5432/test?sslmode=disable"})
		So(err, ShouldBeNil)

		r := redis.NewRing(lc, rcfg)
		opts := redis.NewOptions(r)
		_, _ = redis.NewCache(lc, rcfg, opts)

		ricfg := &ristretto.Config{
			NumCounters: 1e7,
			MaxCost:     1 << 30,
			BufferItems: 64,
		}

		_, err = ristretto.NewCache(lc, ricfg)
		So(err, ShouldBeNil)

		cfg := &pkgHTTP.Config{Port: test.GenerateRandomPort()}
		httpServer := pkgHTTP.NewServer(lc, test.NewShutdowner(), cfg, logger)

		err = prometheusHTTP.Register(httpServer)
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I query metrics", func() {
			client := test.NewHTTPClient(logger)

			req, err := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("http://localhost:%s/metrics", cfg.Port), nil)
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
