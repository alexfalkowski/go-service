package grpc_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/alexfalkowski/go-health/checker"
	"github.com/alexfalkowski/go-health/server"
	"github.com/alexfalkowski/go-health/subscriber"
	"github.com/alexfalkowski/go-service/health"
	hgrpc "github.com/alexfalkowski/go-service/health/transport/grpc"
	"github.com/alexfalkowski/go-service/otel"
	"github.com/alexfalkowski/go-service/test"
	"github.com/alexfalkowski/go-service/transport/grpc/security/jwt"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func init() {
	otel.Register()
}

//nolint:dupl
func TestUnary(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		cfg := test.NewTransportConfig()
		o := observer(lc, "https://httpstat.us/200", test.NewHTTPClient(lc, logger, test.NewOTELConfig(), cfg))
		hs := test.NewHTTPServer(lc, logger, test.NewOTELConfig(), cfg)
		gs := test.NewGRPCServer(lc, logger, test.NewOTELConfig(), cfg, false, nil, nil)

		test.RegisterTransport(lc, cfg, gs, hs)

		hgrpc.Register(gs, &hgrpc.Observer{Observer: o})
		lc.RequireStart()
		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			ctx := context.Background()
			conn := test.NewGRPCClient(ctx, lc, logger, cfg, test.NewOTELConfig(), nil)
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

//nolint:dupl
func TestInvalidUnary(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		cfg := test.NewTransportConfig()
		o := observer(lc, "https://httpstat.us/500", test.NewHTTPClient(lc, logger, test.NewOTELConfig(), cfg))
		hs := test.NewHTTPServer(lc, logger, test.NewOTELConfig(), cfg)
		gs := test.NewGRPCServer(lc, logger, test.NewOTELConfig(), cfg, false, nil, nil)

		test.RegisterTransport(lc, cfg, gs, hs)
		hgrpc.Register(gs, &hgrpc.Observer{Observer: o})
		lc.RequireStart()
		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			ctx := context.Background()
			conn := test.NewGRPCClient(ctx, lc, logger, cfg, test.NewOTELConfig(), nil)
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
		logger := test.NewLogger(lc)
		cfg := test.NewTransportConfig()
		o := observer(lc, "https://httpstat.us/200", test.NewHTTPClient(lc, logger, test.NewOTELConfig(), cfg))
		verifier := test.NewVerifier("test")
		hs := test.NewHTTPServer(lc, logger, test.NewOTELConfig(), cfg)
		gs := test.NewGRPCServer(lc, logger, test.NewOTELConfig(), cfg, false,
			[]grpc.UnaryServerInterceptor{jwt.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{jwt.StreamServerInterceptor(verifier)},
		)

		test.RegisterTransport(lc, cfg, gs, hs)
		hgrpc.Register(gs, &hgrpc.Observer{Observer: o})
		lc.RequireStart()
		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			ctx := context.Background()
			conn := test.NewGRPCClient(ctx, lc, logger, cfg, test.NewOTELConfig(), nil)
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

//nolint:dupl
func TestStream(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		cfg := test.NewTransportConfig()
		o := observer(lc, "https://httpstat.us/200", test.NewHTTPClient(lc, logger, test.NewOTELConfig(), cfg))
		hs := test.NewHTTPServer(lc, logger, test.NewOTELConfig(), cfg)
		gs := test.NewGRPCServer(lc, logger, test.NewOTELConfig(), cfg, false, nil, nil)

		test.RegisterTransport(lc, cfg, gs, hs)
		hgrpc.Register(gs, &hgrpc.Observer{Observer: o})
		lc.RequireStart()
		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			ctx := context.Background()
			conn := test.NewGRPCClient(ctx, lc, logger, cfg, test.NewOTELConfig(), nil)
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

//nolint:dupl
func TestInvalidStream(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		cfg := test.NewTransportConfig()
		o := observer(lc, "https://httpstat.us/500", test.NewHTTPClient(lc, logger, test.NewOTELConfig(), cfg))
		hs := test.NewHTTPServer(lc, logger, test.NewOTELConfig(), cfg)
		gs := test.NewGRPCServer(lc, logger, test.NewOTELConfig(), cfg, false, nil, nil)

		test.RegisterTransport(lc, cfg, gs, hs)
		hgrpc.Register(gs, &hgrpc.Observer{Observer: o})
		lc.RequireStart()
		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			ctx := context.Background()
			conn := test.NewGRPCClient(ctx, lc, logger, cfg, test.NewOTELConfig(), nil)
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
		logger := test.NewLogger(lc)
		cfg := test.NewTransportConfig()
		o := observer(lc, "https://httpstat.us/200", test.NewHTTPClient(lc, logger, test.NewOTELConfig(), cfg))
		verifier := test.NewVerifier("test")
		hs := test.NewHTTPServer(lc, logger, test.NewOTELConfig(), cfg)
		gs := test.NewGRPCServer(lc, logger, test.NewOTELConfig(), cfg, false,
			[]grpc.UnaryServerInterceptor{jwt.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{jwt.StreamServerInterceptor(verifier)},
		)

		test.RegisterTransport(lc, cfg, gs, hs)
		hgrpc.Register(gs, &hgrpc.Observer{Observer: o})
		lc.RequireStart()
		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			ctx := context.Background()
			conn := test.NewGRPCClient(ctx, lc, logger, cfg, test.NewOTELConfig(), nil)
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
