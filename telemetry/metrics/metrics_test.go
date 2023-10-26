package metrics_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/alexfalkowski/go-service/compressor"
	"github.com/alexfalkowski/go-service/database/sql/pg"
	ptracer "github.com/alexfalkowski/go-service/database/sql/pg/telemetry/tracer"
	smetrics "github.com/alexfalkowski/go-service/database/sql/telemetry/metrics"
	"github.com/alexfalkowski/go-service/marshaller"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func init() {
	tracer.Register()
}

func TestHTTP(t *testing.T) {
	Convey("Given I register the metrics handler", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)

		tracer, err := ptracer.NewTracer(ptracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		pg.Register(tracer, logger)

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		dbs, err := pg.Open(pg.OpenParams{Lifecycle: lc, Config: test.NewPGConfig()})
		So(err, ShouldBeNil)

		smetrics.Register(dbs, test.Version, m)

		_ = test.NewRedisCache(lc, "localhost:6379", logger, compressor.NewSnappy(), marshaller.NewProto(), m)
		_ = test.NewRistrettoCache(lc, m)
		cfg := test.NewTransportConfig()
		hs := test.NewHTTPServer(lc, logger, test.NewTracerConfig(), cfg, m)
		gs := test.NewGRPCServer(lc, logger, test.NewTracerConfig(), cfg, false, nil, nil, m)

		test.RegisterTransport(lc, cfg, gs, hs)

		err = metrics.Register(hs)
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I query metrics", func() {
			client := test.NewHTTPClient(lc, logger, test.NewTracerConfig(), cfg, m)

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
