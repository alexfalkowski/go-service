package grpc_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/alexfalkowski/go-health/checker"
	"github.com/alexfalkowski/go-health/server"
	"github.com/alexfalkowski/go-health/subscriber"
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
	"go.uber.org/fx"
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

		o := observer(lc, "https://httpstat.us/200", test.NewHTTPClient(logger, tracer, test.Version, hprometheus.NewClientMetrics(lc, test.Version)))
		gs, gconfig := test.NewGRPCServer(lc, logger, test.NewJaegerConfig(), false, nil, nil)

		hgrpc.Register(gs, &hgrpc.Observer{Observer: o})
		lc.RequireStart()
		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			ctx := context.Background()

			conn, err := tgrpc.NewClient(
				tgrpc.ClientParams{Context: ctx, Host: fmt.Sprintf("127.0.0.1:%s", gconfig.Port), Version: test.Version, Config: gconfig},
				tgrpc.WithClientLogger(logger), tgrpc.WithClientTracer(tracer),
				tgrpc.WithClientBreaker(), tgrpc.WithClientRetry(),
				tgrpc.WithClientDialOption(grpc.WithBlock()),
				tgrpc.WithClientMetrics(gprometheus.NewClientMetrics(lc, test.Version)),
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

		o := observer(lc, "https://httpstat.us/500", test.NewHTTPClient(logger, tracer, test.Version, hprometheus.NewClientMetrics(lc, test.Version)))
		gs, gconfig := test.NewGRPCServer(lc, logger, test.NewJaegerConfig(), false, nil, nil)

		hgrpc.Register(gs, &hgrpc.Observer{Observer: o})
		lc.RequireStart()
		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			ctx := context.Background()

			conn, err := tgrpc.NewClient(
				tgrpc.ClientParams{Context: ctx, Host: fmt.Sprintf("127.0.0.1:%s", gconfig.Port), Version: test.Version, Config: gconfig},
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

		o := observer(lc, "https://httpstat.us/200", test.NewHTTPClient(logger, tracer, test.Version, hprometheus.NewClientMetrics(lc, test.Version)))
		verifier := test.NewVerifier("test")
		gs, gconfig := test.NewGRPCServer(lc, logger, test.NewJaegerConfig(), false,
			[]grpc.UnaryServerInterceptor{jwt.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{jwt.StreamServerInterceptor(verifier)},
		)

		hgrpc.Register(gs, &hgrpc.Observer{Observer: o})
		lc.RequireStart()
		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			ctx := context.Background()

			conn, err := tgrpc.NewClient(
				tgrpc.ClientParams{Context: ctx, Host: fmt.Sprintf("127.0.0.1:%s", gconfig.Port), Version: test.Version, Config: gconfig},
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

		o := observer(lc, "https://httpstat.us/200", test.NewHTTPClient(logger, tracer, test.Version, hprometheus.NewClientMetrics(lc, test.Version)))
		gs, gconfig := test.NewGRPCServer(lc, logger, test.NewJaegerConfig(), false, nil, nil)

		hgrpc.Register(gs, &hgrpc.Observer{Observer: o})
		lc.RequireStart()
		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			ctx := context.Background()

			conn, err := tgrpc.NewClient(
				tgrpc.ClientParams{Context: ctx, Host: fmt.Sprintf("127.0.0.1:%s", gconfig.Port), Version: test.Version, Config: gconfig},
				tgrpc.WithClientLogger(logger), tgrpc.WithClientTracer(tracer),
				tgrpc.WithClientBreaker(), tgrpc.WithClientRetry(),
				tgrpc.WithClientDialOption(grpc.WithBlock()),
				tgrpc.WithClientMetrics(gprometheus.NewClientMetrics(lc, test.Version)),
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

		o := observer(lc, "https://httpstat.us/500", test.NewHTTPClient(logger, tracer, test.Version, hprometheus.NewClientMetrics(lc, test.Version)))
		gs, gconfig := test.NewGRPCServer(lc, logger, test.NewJaegerConfig(), false, nil, nil)

		hgrpc.Register(gs, &hgrpc.Observer{Observer: o})
		lc.RequireStart()
		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			ctx := context.Background()

			conn, err := tgrpc.NewClient(
				tgrpc.ClientParams{Context: ctx, Host: fmt.Sprintf("127.0.0.1:%s", gconfig.Port), Version: test.Version, Config: gconfig},
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

func TestIgnoreAuthStream(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		tracer, err := opentracing.NewTracer(opentracing.TracerParams{Lifecycle: lc, Config: test.NewJaegerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		o := observer(lc, "https://httpstat.us/200", test.NewHTTPClient(logger, tracer, test.Version, hprometheus.NewClientMetrics(lc, test.Version)))
		verifier := test.NewVerifier("test")
		gs, gconfig := test.NewGRPCServer(lc, logger, test.NewJaegerConfig(), false,
			[]grpc.UnaryServerInterceptor{jwt.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{jwt.StreamServerInterceptor(verifier)},
		)

		hgrpc.Register(gs, &hgrpc.Observer{Observer: o})
		lc.RequireStart()
		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			ctx := context.Background()

			conn, err := tgrpc.NewClient(
				tgrpc.ClientParams{Context: ctx, Host: fmt.Sprintf("127.0.0.1:%s", gconfig.Port), Version: test.Version, Config: gconfig},
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

func observer(lc fx.Lifecycle, url string, client *http.Client) *subscriber.Observer {
	cc := checker.NewHTTPChecker(url, client)
	hr := server.NewRegistration("http", 10*time.Millisecond, cc)
	regs := health.Registrations{hr}
	hs := health.NewServer(lc, regs)

	return hs.Observe("http")
}
