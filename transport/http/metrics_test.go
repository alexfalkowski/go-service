package http_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/database/sql/pg"
	sm "github.com/alexfalkowski/go-service/database/sql/telemetry/metrics"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/test"
	ht "github.com/alexfalkowski/go-service/transport/http"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
)

//nolint:funlen
func TestPrometheusInsecureHTTP(t *testing.T) {
	Convey("Given I register the metrics handler", t, func() {
		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		tc := test.NewOTLPTracerConfig()
		tracer, err := tracer.NewTracer(lc, test.Environment, test.Version, test.Name, tc, logger)
		So(err, ShouldBeNil)

		pg.Register(tracer, logger)

		mc := test.NewPrometheusMetricsConfig()
		m := test.NewMeter(lc, mc)

		dbs, err := pg.Open(pg.OpenParams{Lifecycle: lc, Config: test.NewPGConfig()})
		So(err, ShouldBeNil)

		sm.Register(dbs, m)

		c := &test.Cache{Lifecycle: lc, Redis: test.NewRedisConfig("redis", "snappy", "proto"), Logger: logger, Meter: m}
		_, _ = c.NewRedisCache()
		_ = c.NewRistrettoCache()
		cfg := test.NewInsecureTransportConfig()

		s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Mux: mux}
		s.Register()

		cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}

		ht.RegisterMetrics(mc, mux)
		lc.RequireStart()

		Convey("When I query metrics", func() {
			client := cl.NewHTTP()

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
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
				So(response, ShouldContainSubstring, "system")
				So(response, ShouldContainSubstring, "process")
				So(response, ShouldContainSubstring, "runtime")
			})
		})

		lc.RequireStop()
	})
}

//nolint:funlen
func TestPrometheusSecureHTTP(t *testing.T) {
	Convey("Given I register the metrics handler", t, func() {
		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		tc := test.NewOTLPTracerConfig()
		tracer, err := tracer.NewTracer(lc, test.Environment, test.Version, test.Name, tc, logger)
		So(err, ShouldBeNil)

		pg.Register(tracer, logger)

		mc := test.NewPrometheusMetricsConfig()
		m := test.NewMeter(lc, mc)

		dbs, err := pg.Open(pg.OpenParams{Lifecycle: lc, Config: test.NewPGConfig()})
		So(err, ShouldBeNil)

		sm.Register(dbs, m)

		c := &test.Cache{Lifecycle: lc, Redis: test.NewRedisConfig("redis", "snappy", "proto"), Logger: logger, Meter: m}
		_, _ = c.NewRedisCache()
		_ = c.NewRistrettoCache()
		cfg := test.NewSecureTransportConfig()

		s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Mux: mux}
		s.Register()

		cl := &test.Client{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m,
			TLS: test.NewTLSClientConfig(),
		}

		ht.RegisterMetrics(mc, mux)
		lc.RequireStart()

		Convey("When I query metrics", func() {
			client := cl.NewHTTP()

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
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
