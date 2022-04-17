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
	"github.com/alexfalkowski/go-service/trace/opentracing"
	shttp "github.com/alexfalkowski/go-service/transport/http"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestHealth(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		tracer, err := opentracing.NewJaegerTransportTracer(lc, logger, test.NewJaegerConfig())
		So(err, ShouldBeNil)

		rcfg := &redis.Config{Host: "localhost:6379"}
		r := redis.NewRing(lc, rcfg, logger)

		cc := checker.NewHTTPChecker("https://httpstat.us/200", test.NewHTTPClient(logger, tracer))
		hr := server.NewRegistration("http", 10*time.Millisecond, cc)
		rc := hchecker.NewRedisChecker(r, 1*time.Second)
		rr := server.NewRegistration("redis", 10*time.Millisecond, rc)
		regs := health.Registrations{hr, rr}
		server := health.NewServer(lc, regs)
		o := server.Observe("http")
		cfg := &shttp.Config{Port: test.GenerateRandomPort()}
		params := shttp.ServerParams{Lifecycle: lc, Shutdowner: test.NewShutdowner(), Config: cfg, Logger: logger, Tracer: tracer}
		httpServer := shttp.NewServer(params)

		err = hhttp.Register(httpServer, &hhttp.HealthObserver{Observer: o}, &hhttp.LivenessObserver{Observer: o}, &hhttp.ReadinessObserver{Observer: o})
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I query health", func() {
			client := test.NewHTTPClient(logger, tracer)

			req, err := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("http://localhost:%s/health", cfg.Port), nil)
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

// nolint:dupl
func TestLiveness(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		tracer, err := opentracing.NewJaegerTransportTracer(lc, logger, test.NewJaegerConfig())
		So(err, ShouldBeNil)

		cc := checker.NewHTTPChecker("https://httpstat.us/200", test.NewHTTPClient(logger, tracer))
		hr := server.NewRegistration("http", 10*time.Millisecond, cc)
		regs := health.Registrations{hr}
		server := health.NewServer(lc, regs)
		o := server.Observe("http")
		cfg := &shttp.Config{Port: test.GenerateRandomPort()}
		params := shttp.ServerParams{Lifecycle: lc, Shutdowner: test.NewShutdowner(), Config: cfg, Logger: logger, Tracer: tracer}
		httpServer := shttp.NewServer(params)

		err = hhttp.Register(httpServer, &hhttp.HealthObserver{Observer: o}, &hhttp.LivenessObserver{Observer: o}, &hhttp.ReadinessObserver{Observer: o})
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I query health", func() {
			client := test.NewHTTPClient(logger, tracer)

			req, err := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("http://localhost:%s/liveness", cfg.Port), nil)
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

// nolint:dupl
func TestReadiness(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		tracer, err := opentracing.NewJaegerTransportTracer(lc, logger, test.NewJaegerConfig())
		So(err, ShouldBeNil)

		cc := checker.NewHTTPChecker("https://httpstat.us/200", test.NewHTTPClient(logger, tracer))
		hr := server.NewRegistration("http", 10*time.Millisecond, cc)
		regs := health.Registrations{hr}
		server := health.NewServer(lc, regs)
		o := server.Observe("http")
		cfg := &shttp.Config{Port: test.GenerateRandomPort()}
		params := shttp.ServerParams{Lifecycle: lc, Shutdowner: test.NewShutdowner(), Config: cfg, Logger: logger, Tracer: tracer}
		httpServer := shttp.NewServer(params)

		err = hhttp.Register(httpServer, &hhttp.HealthObserver{Observer: o}, &hhttp.LivenessObserver{Observer: o}, &hhttp.ReadinessObserver{Observer: o})
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I query health", func() {
			client := test.NewHTTPClient(logger, tracer)

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

// nolint:dupl
func TestInvalidHealth(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		tracer, err := opentracing.NewJaegerTransportTracer(lc, logger, test.NewJaegerConfig())
		So(err, ShouldBeNil)

		cc := checker.NewHTTPChecker("https://httpstat.us/500", test.NewHTTPClient(logger, tracer))
		hr := server.NewRegistration("http", 10*time.Millisecond, cc)
		regs := health.Registrations{hr}
		server := health.NewServer(lc, regs)
		o := server.Observe("http")
		cfg := &shttp.Config{Port: test.GenerateRandomPort()}
		params := shttp.ServerParams{Lifecycle: lc, Shutdowner: test.NewShutdowner(), Config: cfg, Logger: logger, Tracer: tracer}
		httpServer := shttp.NewServer(params)

		err = hhttp.Register(httpServer, &hhttp.HealthObserver{Observer: o}, &hhttp.LivenessObserver{Observer: o}, &hhttp.ReadinessObserver{Observer: o})
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I query health", func() {
			client := test.NewHTTPClient(logger, tracer)

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
