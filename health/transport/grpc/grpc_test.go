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
	shg "github.com/alexfalkowski/go-service/health/transport/grpc"
	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/test"
	"github.com/alexfalkowski/go-service/transport/grpc/security/token"
	tm "github.com/alexfalkowski/go-service/transport/meta"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
)

func init() {
	tracer.Register()
}

func TestUnary(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)
		cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}
		client := cl.NewHTTP()
		o := observer(lc, "http://localhost:6000/v1/status/200", client)

		s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}
		s.Register()

		shg.Register(s.GRPC, &shg.Observer{Observer: o})
		lc.RequireStart()
		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			ctx := context.Background()
			ctx = tm.WithRequestID(ctx, meta.String("test-id"))
			ctx = tm.WithUserAgent(ctx, meta.String("test-user-agent"))

			conn := cl.NewGRPC()
			defer conn.Close()

			client := grpc_health_v1.NewHealthClient(conn)
			req := &grpc_health_v1.HealthCheckRequest{}

			resp, err := client.Check(ctx, req)
			So(err, ShouldBeNil)

			lc.RequireStop()

			Convey("Then I should have a healthy response", func() {
				So(resp.GetStatus(), ShouldEqual, grpc_health_v1.HealthCheckResponse_SERVING)
			})
		})
	})
}

func TestInvalidUnary(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)
		cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}
		client := cl.NewHTTP()
		o := observer(lc, "http://localhost:6000/v1/status/500", client)

		s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}
		s.Register()

		shg.Register(s.GRPC, &shg.Observer{Observer: o})
		lc.RequireStart()
		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			ctx := context.Background()

			conn := cl.NewGRPC()
			defer conn.Close()

			client := grpc_health_v1.NewHealthClient(conn)
			req := &grpc_health_v1.HealthCheckRequest{}

			md := metadata.New(map[string]string{"request-id": "test-id", "user-agent": "test-user-agent"})
			ctx = metadata.NewOutgoingContext(ctx, md)

			resp, err := client.Check(ctx, req)
			So(err, ShouldBeNil)

			lc.RequireStop()

			Convey("Then I should have an unhealthy response", func() {
				So(resp.GetStatus(), ShouldEqual, grpc_health_v1.HealthCheckResponse_NOT_SERVING)
			})
		})
	})
}

func TestIgnoreAuthUnary(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)
		cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}
		client := cl.NewHTTP()
		o := observer(lc, "http://localhost:6000/v1/status/200", client)
		verifier := test.NewVerifier("test")

		s := &test.Server{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m,
			Unary:  []grpc.UnaryServerInterceptor{token.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{token.StreamServerInterceptor(verifier)},
		}
		s.Register()

		shg.Register(s.GRPC, &shg.Observer{Observer: o})
		lc.RequireStart()
		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			ctx := context.Background()

			conn := cl.NewGRPC()
			defer conn.Close()

			client := grpc_health_v1.NewHealthClient(conn)
			req := &grpc_health_v1.HealthCheckRequest{}

			resp, err := client.Check(ctx, req)
			So(err, ShouldBeNil)

			lc.RequireStop()

			Convey("Then I should have a healthy response", func() {
				So(resp.GetStatus(), ShouldEqual, grpc_health_v1.HealthCheckResponse_SERVING)
			})
		})
	})
}

//nolint:dupl
func TestStream(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)
		cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}
		client := cl.NewHTTP()
		o := observer(lc, "http://localhost:6000/v1/status/200", client)

		s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}
		s.Register()

		shg.Register(s.GRPC, &shg.Observer{Observer: o})
		lc.RequireStart()
		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			ctx := context.Background()

			conn := cl.NewGRPC()
			defer conn.Close()

			client := grpc_health_v1.NewHealthClient(conn)
			req := &grpc_health_v1.HealthCheckRequest{}

			wc, err := client.Watch(ctx, req)
			So(err, ShouldBeNil)

			resp, err := wc.Recv()
			So(err, ShouldBeNil)

			lc.RequireStop()

			Convey("Then I should have a healthy response", func() {
				So(resp.GetStatus(), ShouldEqual, grpc_health_v1.HealthCheckResponse_SERVING)
			})
		})
	})
}

//nolint:dupl
func TestInvalidStream(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)
		cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}
		client := cl.NewHTTP()
		o := observer(lc, "http://localhost:6000/v1/status/500", client)

		s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}
		s.Register()

		shg.Register(s.GRPC, &shg.Observer{Observer: o})
		lc.RequireStart()
		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			ctx := context.Background()

			conn := cl.NewGRPC()
			defer conn.Close()

			client := grpc_health_v1.NewHealthClient(conn)
			req := &grpc_health_v1.HealthCheckRequest{}

			wc, err := client.Watch(ctx, req)
			So(err, ShouldBeNil)

			resp, err := wc.Recv()
			So(err, ShouldBeNil)

			lc.RequireStop()

			Convey("Then I should have a healthy response", func() {
				So(resp.GetStatus(), ShouldEqual, grpc_health_v1.HealthCheckResponse_NOT_SERVING)
			})
		})
	})
}

func TestIgnoreAuthStream(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)
		cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}
		client := cl.NewHTTP()
		o := observer(lc, "http://localhost:6000/v1/status/200", client)
		verifier := test.NewVerifier("test")

		s := &test.Server{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m,
			Unary:  []grpc.UnaryServerInterceptor{token.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{token.StreamServerInterceptor(verifier)},
		}
		s.Register()

		shg.Register(s.GRPC, &shg.Observer{Observer: o})
		lc.RequireStart()
		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			ctx := context.Background()

			conn := cl.NewGRPC()
			defer conn.Close()

			client := grpc_health_v1.NewHealthClient(conn)
			req := &grpc_health_v1.HealthCheckRequest{}

			wc, err := client.Watch(ctx, req)
			So(err, ShouldBeNil)

			resp, err := wc.Recv()
			So(err, ShouldBeNil)

			lc.RequireStop()

			Convey("Then I should have a healthy response", func() {
				So(resp.GetStatus(), ShouldEqual, grpc_health_v1.HealthCheckResponse_SERVING)
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
