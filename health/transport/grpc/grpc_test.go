package grpc_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/alexfalkowski/go-health/checker"
	"github.com/alexfalkowski/go-health/server"
	"github.com/alexfalkowski/go-service/health"
	hgrpc "github.com/alexfalkowski/go-service/health/transport/grpc"
	"github.com/alexfalkowski/go-service/logger/zap"
	"github.com/alexfalkowski/go-service/test"
	tgrpc "github.com/alexfalkowski/go-service/transport/grpc"
	gprometheus "github.com/alexfalkowski/go-service/transport/grpc/metrics/prometheus"
	"github.com/alexfalkowski/go-service/transport/grpc/security/jwt"
	"github.com/alexfalkowski/go-service/transport/grpc/trace/opentracing"
	hprometheus "github.com/alexfalkowski/go-service/transport/http/metrics/prometheus"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func TestUnary(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		tracer, err := opentracing.NewTracer(opentracing.TracerParams{Lifecycle: lc, Config: test.NewJaegerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		cc := checker.NewHTTPChecker("https://httpstat.us/200", test.NewHTTPClient(logger, tracer, test.Version, hprometheus.NewClientMetrics(lc, test.Version)))
		hr := server.NewRegistration("http", 10*time.Millisecond, cc)
		regs := health.Registrations{hr}
		hs := health.NewServer(lc, regs)
		o := hs.Observe("http")
		cfg := test.NewGRPCConfig()
		params := tgrpc.ServerParams{
			Lifecycle: lc, Shutdowner: test.NewShutdowner(),
			Config: cfg, Logger: logger, Tracer: tracer,
			Metrics: gprometheus.NewServerMetrics(lc, test.Version),
		}
		gs := tgrpc.NewServer(params)
		metrics := gprometheus.NewClientMetrics(lc, test.Version)

		hgrpc.Register(gs, &hgrpc.Observer{Observer: o})
		lc.RequireStart()
		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			ctx := context.Background()

			conn, err := tgrpc.NewClient(
				tgrpc.ClientParams{Context: ctx, Host: fmt.Sprintf("127.0.0.1:%s", cfg.Port), Version: test.Version, Config: cfg},
				tgrpc.WithClientLogger(logger), tgrpc.WithClientTracer(tracer),
				tgrpc.WithClientBreaker(), tgrpc.WithClientRetry(),
				tgrpc.WithClientDialOption(grpc.WithBlock()),
				tgrpc.WithClientMetrics(metrics),
			)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := grpc_health_v1.NewHealthClient(conn)
			req := &grpc_health_v1.HealthCheckRequest{}

			resp, err := client.Check(ctx, req)
			So(err, ShouldBeNil)

			lc.RequireStop()

			Convey("Then I should have a healthy response", func() {
				So(resp.Status, ShouldEqual, grpc_health_v1.HealthCheckResponse_SERVING)
			})
		})
	})
}

func TestInvalidUnary(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		tracer, err := opentracing.NewTracer(opentracing.TracerParams{Lifecycle: lc, Config: test.NewJaegerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		cc := checker.NewHTTPChecker("https://httpstat.us/500", test.NewHTTPClient(logger, tracer, test.Version, hprometheus.NewClientMetrics(lc, test.Version)))
		hr := server.NewRegistration("http", 10*time.Millisecond, cc)
		regs := health.Registrations{hr}
		hs := health.NewServer(lc, regs)
		o := hs.Observe("http")
		cfg := test.NewGRPCConfig()
		params := tgrpc.ServerParams{
			Lifecycle: lc, Shutdowner: test.NewShutdowner(),
			Config: cfg, Logger: logger, Tracer: tracer,
			Metrics: gprometheus.NewServerMetrics(lc, test.Version),
		}
		gs := tgrpc.NewServer(params)

		hgrpc.Register(gs, &hgrpc.Observer{Observer: o})

		lc.RequireStart()

		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			ctx := context.Background()

			conn, err := tgrpc.NewClient(
				tgrpc.ClientParams{Context: ctx, Host: fmt.Sprintf("127.0.0.1:%s", cfg.Port), Version: test.Version, Config: cfg},
				tgrpc.WithClientLogger(logger), tgrpc.WithClientTracer(tracer),
				tgrpc.WithClientBreaker(), tgrpc.WithClientRetry(),
				tgrpc.WithClientDialOption(grpc.WithBlock()),
			)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := grpc_health_v1.NewHealthClient(conn)
			req := &grpc_health_v1.HealthCheckRequest{}

			resp, err := client.Check(ctx, req)
			So(err, ShouldBeNil)

			lc.RequireStop()

			Convey("Then I should have an unhealthy response", func() {
				So(resp.Status, ShouldEqual, grpc_health_v1.HealthCheckResponse_NOT_SERVING)
			})
		})
	})
}

func TestIgnoreAuthUnary(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		tracer, err := opentracing.NewTracer(opentracing.TracerParams{Lifecycle: lc, Config: test.NewJaegerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		cc := checker.NewHTTPChecker("https://httpstat.us/200", test.NewHTTPClient(logger, tracer, test.Version, hprometheus.NewClientMetrics(lc, test.Version)))
		hr := server.NewRegistration("http", 10*time.Millisecond, cc)
		regs := health.Registrations{hr}
		hs := health.NewServer(lc, regs)
		o := hs.Observe("http")
		cfg := test.NewGRPCConfig()
		verifier := test.NewVerifier("test")
		params := tgrpc.ServerParams{
			Lifecycle:  lc,
			Shutdowner: test.NewShutdowner(),
			Config:     cfg,
			Logger:     logger,
			Tracer:     tracer,
			Unary:      []grpc.UnaryServerInterceptor{jwt.UnaryServerInterceptor(verifier)},
			Stream:     []grpc.StreamServerInterceptor{jwt.StreamServerInterceptor(verifier)},
			Metrics:    gprometheus.NewServerMetrics(lc, test.Version),
		}
		gs := tgrpc.NewServer(params)

		hgrpc.Register(gs, &hgrpc.Observer{Observer: o})

		lc.RequireStart()

		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			ctx := context.Background()

			conn, err := tgrpc.NewClient(
				tgrpc.ClientParams{Context: ctx, Host: fmt.Sprintf("127.0.0.1:%s", cfg.Port), Version: test.Version, Config: cfg},
				tgrpc.WithClientLogger(logger), tgrpc.WithClientTracer(tracer),
				tgrpc.WithClientBreaker(), tgrpc.WithClientRetry(),
				tgrpc.WithClientDialOption(grpc.WithBlock()),
			)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := grpc_health_v1.NewHealthClient(conn)
			req := &grpc_health_v1.HealthCheckRequest{}

			resp, err := client.Check(ctx, req)
			So(err, ShouldBeNil)

			lc.RequireStop()

			Convey("Then I should have a healthy response", func() {
				So(resp.Status, ShouldEqual, grpc_health_v1.HealthCheckResponse_SERVING)
			})
		})
	})
}

func TestStream(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		tracer, err := opentracing.NewTracer(opentracing.TracerParams{Lifecycle: lc, Config: test.NewJaegerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		cc := checker.NewHTTPChecker("https://httpstat.us/200", test.NewHTTPClient(logger, tracer, test.Version, hprometheus.NewClientMetrics(lc, test.Version)))
		hr := server.NewRegistration("http", 10*time.Millisecond, cc)
		regs := health.Registrations{hr}
		hs := health.NewServer(lc, regs)
		o := hs.Observe("http")
		cfg := test.NewGRPCConfig()
		params := tgrpc.ServerParams{
			Lifecycle: lc, Shutdowner: test.NewShutdowner(),
			Config: cfg, Logger: logger, Tracer: tracer,
			Metrics: gprometheus.NewServerMetrics(lc, test.Version),
		}
		gs := tgrpc.NewServer(params)
		metrics := gprometheus.NewClientMetrics(lc, test.Version)

		hgrpc.Register(gs, &hgrpc.Observer{Observer: o})
		lc.RequireStart()
		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			ctx := context.Background()

			conn, err := tgrpc.NewClient(
				tgrpc.ClientParams{Context: ctx, Host: fmt.Sprintf("127.0.0.1:%s", cfg.Port), Version: test.Version, Config: cfg},
				tgrpc.WithClientLogger(logger), tgrpc.WithClientTracer(tracer),
				tgrpc.WithClientBreaker(), tgrpc.WithClientRetry(),
				tgrpc.WithClientDialOption(grpc.WithBlock()),
				tgrpc.WithClientMetrics(metrics),
			)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := grpc_health_v1.NewHealthClient(conn)
			req := &grpc_health_v1.HealthCheckRequest{}

			wc, err := client.Watch(ctx, req)
			So(err, ShouldBeNil)

			resp, err := wc.Recv()
			So(err, ShouldBeNil)

			lc.RequireStop()

			Convey("Then I should have a healthy response", func() {
				So(resp.Status, ShouldEqual, grpc_health_v1.HealthCheckResponse_SERVING)
			})
		})
	})
}

func TestInvalidStream(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		tracer, err := opentracing.NewTracer(opentracing.TracerParams{Lifecycle: lc, Config: test.NewJaegerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		cc := checker.NewHTTPChecker("https://httpstat.us/500", test.NewHTTPClient(logger, tracer, test.Version, hprometheus.NewClientMetrics(lc, test.Version)))
		hr := server.NewRegistration("http", 10*time.Millisecond, cc)
		regs := health.Registrations{hr}
		hs := health.NewServer(lc, regs)
		o := hs.Observe("http")
		cfg := test.NewGRPCConfig()
		params := tgrpc.ServerParams{
			Lifecycle: lc, Shutdowner: test.NewShutdowner(),
			Config: cfg, Logger: logger, Tracer: tracer,
			Metrics: gprometheus.NewServerMetrics(lc, test.Version),
		}
		gs := tgrpc.NewServer(params)

		hgrpc.Register(gs, &hgrpc.Observer{Observer: o})

		lc.RequireStart()

		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			ctx := context.Background()

			conn, err := tgrpc.NewClient(
				tgrpc.ClientParams{Context: ctx, Host: fmt.Sprintf("127.0.0.1:%s", cfg.Port), Version: test.Version, Config: cfg},
				tgrpc.WithClientLogger(logger), tgrpc.WithClientTracer(tracer),
				tgrpc.WithClientBreaker(), tgrpc.WithClientRetry(),
				tgrpc.WithClientDialOption(grpc.WithBlock()),
			)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := grpc_health_v1.NewHealthClient(conn)
			req := &grpc_health_v1.HealthCheckRequest{}

			wc, err := client.Watch(ctx, req)
			So(err, ShouldBeNil)

			resp, err := wc.Recv()
			So(err, ShouldBeNil)

			lc.RequireStop()

			Convey("Then I should have a healthy response", func() {
				So(resp.Status, ShouldEqual, grpc_health_v1.HealthCheckResponse_NOT_SERVING)
			})
		})
	})
}

// nolint:funlen
func TestIgnoreAuthStream(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		tracer, err := opentracing.NewTracer(opentracing.TracerParams{Lifecycle: lc, Config: test.NewJaegerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		cc := checker.NewHTTPChecker("https://httpstat.us/200", test.NewHTTPClient(logger, tracer, test.Version, hprometheus.NewClientMetrics(lc, test.Version)))
		hr := server.NewRegistration("http", 10*time.Millisecond, cc)
		regs := health.Registrations{hr}
		hs := health.NewServer(lc, regs)
		o := hs.Observe("http")
		cfg := test.NewGRPCConfig()
		verifier := test.NewVerifier("test")
		params := tgrpc.ServerParams{
			Lifecycle:  lc,
			Shutdowner: test.NewShutdowner(),
			Config:     cfg,
			Logger:     logger,
			Tracer:     tracer,
			Unary:      []grpc.UnaryServerInterceptor{jwt.UnaryServerInterceptor(verifier)},
			Stream:     []grpc.StreamServerInterceptor{jwt.StreamServerInterceptor(verifier)},
			Metrics:    gprometheus.NewServerMetrics(lc, test.Version),
		}
		gs := tgrpc.NewServer(params)

		hgrpc.Register(gs, &hgrpc.Observer{Observer: o})

		lc.RequireStart()

		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			ctx := context.Background()

			conn, err := tgrpc.NewClient(
				tgrpc.ClientParams{Context: ctx, Host: fmt.Sprintf("127.0.0.1:%s", cfg.Port), Version: test.Version, Config: cfg},
				tgrpc.WithClientLogger(logger), tgrpc.WithClientTracer(tracer),
				tgrpc.WithClientBreaker(), tgrpc.WithClientRetry(),
				tgrpc.WithClientDialOption(grpc.WithBlock()),
			)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := grpc_health_v1.NewHealthClient(conn)
			req := &grpc_health_v1.HealthCheckRequest{}

			wc, err := client.Watch(ctx, req)
			So(err, ShouldBeNil)

			resp, err := wc.Recv()
			So(err, ShouldBeNil)

			lc.RequireStop()

			Convey("Then I should have a healthy response", func() {
				So(resp.Status, ShouldEqual, grpc_health_v1.HealthCheckResponse_SERVING)
			})
		})
	})
}
