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
	"github.com/alicebob/miniredis/v2"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
)

func init() {
	tracer.Register()
}

//nolint:dupl,funlen
func TestPrometheusInsecureHTTP(t *testing.T) {
	Convey("Given I register the metrics handler", t, func() {
		r := miniredis.RunT(t)
		defer r.Close()

		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		tc := test.NewOTLPTracerConfig()
		tracer := tracer.NewTracer(lc, test.Environment, test.Version, tc, logger)

		pg.Register(tracer, logger)

		m := test.NewPrometheusMeter(lc)

		dbs, err := pg.Open(pg.OpenParams{Lifecycle: lc, Config: test.NewPGConfig()})
		So(err, ShouldBeNil)

		sm.Register(dbs, test.Version, m)

		c := &test.Cache{Lifecycle: lc, Redis: test.NewRedisConfig(r.Addr(), "snappy", "proto"), Logger: logger, Meter: m}
		_, _ = c.NewRedisCache()
		_ = c.NewRistrettoCache()
		cfg := test.NewInsecureTransportConfig()

		s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}
		s.Register()

		cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}

		err = ht.RegisterMetrics(test.Mux)
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I query metrics", func() {
			client := cl.NewHTTP()

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

//nolint:dupl,funlen
func TestPrometheusSecureHTTP(t *testing.T) {
	Convey("Given I register the metrics handler", t, func() {
		r := miniredis.RunT(t)
		defer r.Close()

		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		tc := test.NewOTLPTracerConfig()
		tracer := tracer.NewTracer(lc, test.Environment, test.Version, tc, logger)

		pg.Register(tracer, logger)

		m := test.NewPrometheusMeter(lc)

		dbs, err := pg.Open(pg.OpenParams{Lifecycle: lc, Config: test.NewPGConfig()})
		So(err, ShouldBeNil)

		sm.Register(dbs, test.Version, m)

		c := &test.Cache{Lifecycle: lc, Redis: test.NewRedisConfig(r.Addr(), "snappy", "proto"), Logger: logger, Meter: m}
		_, _ = c.NewRedisCache()
		_ = c.NewRistrettoCache()
		cfg := test.NewSecureTransportConfig()

		s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}
		s.Register()

		cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}

		err = ht.RegisterMetrics(test.Mux)
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I query metrics", func() {
			client := cl.NewHTTP()

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
