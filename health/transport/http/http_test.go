package http_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/alexfalkowski/go-health/checker"
	"github.com/alexfalkowski/go-health/server"
	"github.com/alexfalkowski/go-service/cache/redis"
	"github.com/alexfalkowski/go-service/health"
	hchecker "github.com/alexfalkowski/go-service/health/checker"
	hhttp "github.com/alexfalkowski/go-service/health/transport/http"
	"github.com/alexfalkowski/go-service/logger/zap"
	"github.com/alexfalkowski/go-service/test"
	shttp "github.com/alexfalkowski/go-service/transport/http"
	"github.com/alexfalkowski/go-service/transport/http/metrics/prometheus"
	"github.com/alexfalkowski/go-service/transport/http/trace/opentracing"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestHealth(t *testing.T) {
	checks := []string{"health", "liveness", "readiness"}

	for _, check := range checks {
		Convey("Given I register the health handler", t, func() {
			lc := fxtest.NewLifecycle(t)

			logger, err := zap.NewLogger(lc, zap.NewConfig())
			So(err, ShouldBeNil)

			tracer, err := opentracing.NewTracer(opentracing.TracerParams{Lifecycle: lc, Config: test.NewJaegerConfig(), Version: test.Version})
			So(err, ShouldBeNil)

			o := observer(lc, "https://httpstat.us/200", test.NewHTTPClient(logger, tracer, test.Version, prometheus.NewClientMetrics(lc, test.Version))).Observe("http")
			cfg := &shttp.Config{Port: test.GenerateRandomPort()}
			params := shttp.ServerParams{
				Lifecycle: lc, Shutdowner: test.NewShutdowner(),
				Config: cfg, Logger: logger, Tracer: tracer,
				Metrics: prometheus.NewServerMetrics(lc, test.Version),
			}
			httpServer := shttp.NewServer(params)

			err = hhttp.Register(httpServer, &hhttp.HealthObserver{Observer: o}, &hhttp.LivenessObserver{Observer: o}, &hhttp.ReadinessObserver{Observer: o})
			So(err, ShouldBeNil)

			lc.RequireStart()

			Convey(fmt.Sprintf("When I query %s", check), func() {
				client := test.NewHTTPClient(logger, tracer, test.Version, prometheus.NewClientMetrics(lc, test.Version))

				req, err := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("http://localhost:%s/%s", cfg.Port, check), nil)
				So(err, ShouldBeNil)

				resp, err := client.Do(req)
				So(err, ShouldBeNil)

				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				So(err, ShouldBeNil)

				actual := strings.TrimSpace(string(body))

				lc.RequireStop()

				Convey("Then I should have a healthy response", func() {
					So(actual, ShouldEqual, `{"status":"SERVING"}`)
				})
			})
		})
	}
}

func TestReadinessNoop(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		tracer, err := opentracing.NewTracer(opentracing.TracerParams{Lifecycle: lc, Config: test.NewJaegerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		server := observer(lc, "https://httpstat.us/500", test.NewHTTPClient(logger, tracer, test.Version, prometheus.NewClientMetrics(lc, test.Version)))
		o := server.Observe("http")
		cfg := &shttp.Config{Port: test.GenerateRandomPort()}
		params := shttp.ServerParams{
			Lifecycle: lc, Shutdowner: test.NewShutdowner(),
			Config: cfg, Logger: logger, Tracer: tracer,
			Metrics: prometheus.NewServerMetrics(lc, test.Version),
		}
		httpServer := shttp.NewServer(params)

		err = hhttp.Register(httpServer, &hhttp.HealthObserver{Observer: o}, &hhttp.LivenessObserver{Observer: o}, &hhttp.ReadinessObserver{Observer: server.Observe("noop")})
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I query health", func() {
			client := test.NewHTTPClient(logger, tracer, test.Version, prometheus.NewClientMetrics(lc, test.Version))

			req, err := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("http://localhost:%s/readiness", cfg.Port), nil)
			So(err, ShouldBeNil)

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			actual := strings.TrimSpace(string(body))

			lc.RequireStop()

			Convey("Then I should have a healthy response", func() {
				So(actual, ShouldEqual, `{"status":"SERVING"}`)
			})
		})
	})
}

func TestInvalidHealth(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		tracer, err := opentracing.NewTracer(opentracing.TracerParams{Lifecycle: lc, Config: test.NewJaegerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		o := observer(lc, "https://httpstat.us/500", test.NewHTTPClient(logger, tracer, test.Version, prometheus.NewClientMetrics(lc, test.Version))).Observe("http")
		cfg := &shttp.Config{Port: test.GenerateRandomPort()}
		params := shttp.ServerParams{
			Lifecycle: lc, Shutdowner: test.NewShutdowner(),
			Config: cfg, Logger: logger, Tracer: tracer,
			Metrics: prometheus.NewServerMetrics(lc, test.Version),
		}
		httpServer := shttp.NewServer(params)

		err = hhttp.Register(httpServer, &hhttp.HealthObserver{Observer: o}, &hhttp.LivenessObserver{Observer: o}, &hhttp.ReadinessObserver{Observer: o})
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I query health", func() {
			client := test.NewHTTPClient(logger, tracer, test.Version, prometheus.NewClientMetrics(lc, test.Version))

			req, err := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("http://localhost:%s/health", cfg.Port), nil)
			So(err, ShouldBeNil)

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			actual := strings.TrimSpace(string(body))

			lc.RequireStop()

			Convey("Then I should have an unhealthy response", func() {
				So(actual, ShouldEqual, `{"errors":{"http":"invalid status code"},"status":"NOT_SERVING"}`)
			})
		})
	})
}

func observer(lc fx.Lifecycle, url string, client *http.Client) *server.Server {
	rcfg := &redis.Config{Host: "localhost:6379"}
	r := redis.NewRing(lc, rcfg)
	rc := hchecker.NewRedisChecker(r, 1*time.Second)
	rr := server.NewRegistration("redis", 10*time.Millisecond, rc)

	cc := checker.NewHTTPChecker(url, client)
	hr := server.NewRegistration("http", 10*time.Millisecond, cc)

	no := checker.NewNoopChecker()
	nr := server.NewRegistration("noop", 10*time.Millisecond, no)

	regs := health.Registrations{hr, nr, rr}

	return health.NewServer(lc, regs)
}
