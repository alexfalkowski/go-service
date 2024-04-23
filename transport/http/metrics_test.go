package http_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/alexfalkowski/go-service/database/sql/pg"
	sm "github.com/alexfalkowski/go-service/database/sql/telemetry/metrics"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/test"
	"github.com/alexfalkowski/go-service/time"
	ht "github.com/alexfalkowski/go-service/transport/http"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
)

func init() {
	tracer.Register()
}

//nolint:dupl
func TestPrometheusInsecureHTTP(t *testing.T) {
	Convey("Given I register the metrics handler", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		tracer := test.NewTracer(lc)

		pg.Register(tracer, logger)

		m := test.NewMeter(lc)

		dbs, err := pg.Open(pg.OpenParams{Lifecycle: lc, Config: test.NewPGConfig()})
		So(err, ShouldBeNil)

		sm.Register(dbs, test.Version, m)

		_, _ = test.NewRedisCache(lc, test.NewRedisConfig("localhost:6379", "snappy", "proto"), logger, m)
		_ = test.NewRistrettoCache(lc, m)
		cfg := test.NewInsecureTransportConfig()
		hs := test.NewHTTPServer(lc, logger, test.NewOTLPTracerConfig(), cfg, m, nil)
		gs := test.NewGRPCServer(lc, logger, test.NewOTLPTracerConfig(), cfg, false, m, nil, nil)

		test.RegisterTransport(lc, gs, hs)

		err = ht.RegisterMetrics(test.Mux)
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I query metrics", func() {
			client := test.NewHTTPClient(lc, logger, test.NewOTLPTracerConfig(), cfg, m)

			ctx, cancel := context.WithTimeout(context.Background(), time.Timeout)
			defer cancel()

			req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://localhost:%s/metrics", cfg.HTTP.Port), http.NoBody)
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

//nolint:dupl
func TestPrometheusSecureHTTP(t *testing.T) {
	Convey("Given I register the metrics handler", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		tracer := test.NewTracer(lc)

		pg.Register(tracer, logger)

		m := test.NewMeter(lc)

		dbs, err := pg.Open(pg.OpenParams{Lifecycle: lc, Config: test.NewPGConfig()})
		So(err, ShouldBeNil)

		sm.Register(dbs, test.Version, m)

		_, _ = test.NewRedisCache(lc, test.NewRedisConfig("localhost:6379", "snappy", "proto"), logger, m)
		_ = test.NewRistrettoCache(lc, m)
		cfg := test.NewSecureTransportConfig()
		hs := test.NewHTTPServer(lc, logger, test.NewOTLPTracerConfig(), cfg, m, nil)
		gs := test.NewGRPCServer(lc, logger, test.NewOTLPTracerConfig(), cfg, false, m, nil, nil)

		test.RegisterTransport(lc, gs, hs)

		err = ht.RegisterMetrics(test.Mux)
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I query metrics", func() {
			client := test.NewHTTPClient(lc, logger, test.NewOTLPTracerConfig(), cfg, m)

			ctx, cancel := context.WithTimeout(context.Background(), time.Timeout)
			defer cancel()

			req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://localhost:%s/metrics", cfg.HTTP.Port), http.NoBody)
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
