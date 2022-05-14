package http_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/alexfalkowski/go-service/cache/ristretto"
	"github.com/alexfalkowski/go-service/compressor"
	"github.com/alexfalkowski/go-service/database/sql/pg"
	"github.com/alexfalkowski/go-service/marshaller"
	phttp "github.com/alexfalkowski/go-service/metrics/prometheus/transport/http"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestHTTP(t *testing.T) {
	Convey("Given I register the metrics handler", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)

		_, err := pg.NewDB(pg.DBParams{Lifecycle: lc, Config: &pg.Config{URL: "postgres://test:test@localhost:5432/test?sslmode=disable"}, Version: test.Version})
		So(err, ShouldBeNil)

		_ = test.NewRedisCache(lc, "localhost:6379", logger, compressor.NewSnappy(), marshaller.NewProto())

		ricfg := &ristretto.Config{NumCounters: 1e7, MaxCost: 1 << 30, BufferItems: 64}

		_, err = ristretto.NewCache(ristretto.CacheParams{Lifecycle: lc, Config: ricfg, Version: test.Version})
		So(err, ShouldBeNil)

		hs, hport := test.NewHTTPServer(lc, logger, test.NewJaegerConfig())

		err = phttp.Register(hs)
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I query metrics", func() {
			client := test.NewHTTPClient(lc, logger, test.NewJaegerConfig())

			req, err := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("http://localhost:%s/metrics", hport), nil)
			So(err, ShouldBeNil)

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			Convey("Then I should have valid metrics", func() {
				response := string(body)

				So(response, ShouldContainSubstring, "go_info")
				So(response, ShouldContainSubstring, "redis_hits_total")
				So(response, ShouldContainSubstring, "ristretto_hits_total")
				So(response, ShouldContainSubstring, "sql_max_open_total")
			})
		})

		lc.RequireStop()
	})
}
