package http_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/alexfalkowski/go-service/cache/redis"
	"github.com/alexfalkowski/go-service/cache/ristretto"
	"github.com/alexfalkowski/go-service/compressor"
	"github.com/alexfalkowski/go-service/database/sql/pg"
	"github.com/alexfalkowski/go-service/logger/zap"
	"github.com/alexfalkowski/go-service/marshaller"
	phttp "github.com/alexfalkowski/go-service/metrics/prometheus/transport/http"
	"github.com/alexfalkowski/go-service/test"
	shttp "github.com/alexfalkowski/go-service/transport/http"
	"github.com/alexfalkowski/go-service/transport/http/trace/opentracing"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

// nolint:funlen
func TestHTTP(t *testing.T) {
	Convey("Given I register the metrics handler", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		tracer, err := opentracing.NewTracer(lc, test.NewJaegerConfig())
		So(err, ShouldBeNil)

		rcfg := &redis.Config{Host: "localhost:6379"}

		_, err = pg.NewDB(lc, &pg.Config{URL: "postgres://test:test@localhost:5432/test?sslmode=disable"})
		So(err, ShouldBeNil)

		r := redis.NewRing(lc, rcfg)
		oparams := redis.OptionsParams{Ring: r, Compressor: compressor.NewSnappy(), Marshaller: marshaller.NewProto()}
		opts := redis.NewOptions(oparams)
		_ = redis.NewCache(lc, rcfg, opts)

		ricfg := &ristretto.Config{
			NumCounters: 1e7,
			MaxCost:     1 << 30,
			BufferItems: 64,
		}

		_, err = ristretto.NewCache(lc, ricfg)
		So(err, ShouldBeNil)

		cfg := &shttp.Config{Port: test.GenerateRandomPort()}
		params := shttp.ServerParams{Lifecycle: lc, Shutdowner: test.NewShutdowner(), Config: cfg, Logger: logger, Tracer: tracer}
		httpServer := shttp.NewServer(params)

		err = phttp.Register(httpServer)
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I query metrics", func() {
			client := test.NewHTTPClient(logger, tracer)

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
