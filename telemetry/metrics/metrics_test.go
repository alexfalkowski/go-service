package metrics_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/alexfalkowski/go-service/compressor"
	"github.com/alexfalkowski/go-service/database/sql/pg"
	st "github.com/alexfalkowski/go-service/database/sql/pg/telemetry/tracer"
	sm "github.com/alexfalkowski/go-service/database/sql/telemetry/metrics"
	"github.com/alexfalkowski/go-service/marshaller"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/test"
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

		tracer, err := st.NewTracer(st.Params{Lifecycle: lc, Config: test.NewOTLPTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		pg.Register(tracer, logger)

		m := test.NewMeter(lc)

		dbs, err := pg.Open(pg.OpenParams{Lifecycle: lc, Config: test.NewPGConfig()})
		So(err, ShouldBeNil)

		sm.Register(dbs, test.Version, m)

		_ = test.NewRedisCache(lc, "localhost:6379", logger, compressor.NewSnappy(), marshaller.NewProto(), m)
		_ = test.NewRistrettoCache(lc, m)
		cfg := test.NewInsecureTransportConfig()
		hs := test.NewHTTPServer(lc, logger, test.NewOTLPTracerConfig(), cfg, m, nil)
		gs := test.NewGRPCServer(lc, logger, test.NewOTLPTracerConfig(), cfg, false, m, nil, nil)

		test.RegisterTransport(lc, gs, hs)

		err = metrics.Register(hs)
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I query metrics", func() {
			client := test.NewHTTPClient(lc, logger, test.NewOTLPTracerConfig(), cfg, m)

			req, err := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("http://localhost:%s/metrics", cfg.HTTP.Port), nil)
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

		tracer, err := st.NewTracer(st.Params{Lifecycle: lc, Config: test.NewOTLPTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		pg.Register(tracer, logger)

		m := test.NewMeter(lc)

		dbs, err := pg.Open(pg.OpenParams{Lifecycle: lc, Config: test.NewPGConfig()})
		So(err, ShouldBeNil)

		sm.Register(dbs, test.Version, m)

		_ = test.NewRedisCache(lc, "localhost:6379", logger, compressor.NewSnappy(), marshaller.NewProto(), m)
		_ = test.NewRistrettoCache(lc, m)
		cfg := test.NewSecureTransportConfig()
		hs := test.NewHTTPServer(lc, logger, test.NewOTLPTracerConfig(), cfg, m, nil)
		gs := test.NewGRPCServer(lc, logger, test.NewOTLPTracerConfig(), cfg, false, m, nil, nil)

		test.RegisterTransport(lc, gs, hs)

		err = metrics.Register(hs)
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I query metrics", func() {
			client := test.NewHTTPClient(lc, logger, test.NewOTLPTracerConfig(), cfg, m)

			req, err := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("https://localhost:%s/metrics", cfg.HTTP.Port), nil)
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

func TestOTLP(t *testing.T) {
	Convey("Given I register OTLP metrics", t, func() {
		lc := fxtest.NewLifecycle(t)

		m, err := metrics.NewMeter(lc, test.Environment, test.Version, test.NewOTLPMetricsConfig())
		So(err, ShouldBeNil)

		Convey("When I create a metric", func() {
			c, err := m.Int64Counter("test_otlp")
			So(err, ShouldBeNil)

			lc.RequireStart()
			c.Add(context.Background(), 1)

			Convey("Then I should have a metric", func() {
				lc.RequireStop()
				So(c, ShouldNotBeNil)
			})
		})
	})
}
