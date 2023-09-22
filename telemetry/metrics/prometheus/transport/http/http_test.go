package http_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/alexfalkowski/go-service/compressor"
	"github.com/alexfalkowski/go-service/database/sql/pg"
	ptracer "github.com/alexfalkowski/go-service/database/sql/pg/telemetry/tracer"
	"github.com/alexfalkowski/go-service/marshaller"
	phttp "github.com/alexfalkowski/go-service/telemetry/metrics/prometheus/transport/http"
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

		_, _ = pg.Open(pg.DBParams{Lifecycle: lc, Config: test.NewPGConfig(), Version: test.Version})
		_ = test.NewRedisCache(lc, "localhost:6379", logger, compressor.NewSnappy(), marshaller.NewProto())
		_ = test.NewRistrettoCache(lc)
		cfg := test.NewTransportConfig()
		hs := test.NewHTTPServer(lc, logger, test.NewTracerConfig(), cfg)
		gs := test.NewGRPCServer(lc, logger, test.NewTracerConfig(), cfg, false, nil, nil)

		test.RegisterTransport(lc, cfg, gs, hs)

		err = phttp.Register(hs)
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I query metrics", func() {
			client := test.NewHTTPClient(lc, logger, test.NewTracerConfig(), cfg)

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
				So(response, ShouldContainSubstring, "redis_hits_total")
				So(response, ShouldContainSubstring, "ristretto_hits_total")
				So(response, ShouldContainSubstring, "pg_sql_max_open_total")
			})
		})

		lc.RequireStop()
	})
}
