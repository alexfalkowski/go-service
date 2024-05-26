package grpc_test

import (
	"context"
	"testing"

	"github.com/alexfalkowski/go-service/limiter"
	"github.com/alexfalkowski/go-service/test"
	v1 "github.com/alexfalkowski/go-service/test/greet/v1"
	"github.com/alexfalkowski/go-service/transport/meta"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func init() {
	meta.RegisterKeys()
}

func TestLimiterLimitedUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)

		l, k, err := limiter.New(test.NewLimiterConfig("user-agent", "0-S"))
		So(err, ShouldBeNil)

		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)

		s := &test.Server{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m,
			Limiter: l, Key: k, Mux: test.GatewayMux,
		}
		s.Register()

		cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}

		lc.RequireStart()

		Convey("When I query repeatedly", func() {
			ctx := context.Background()

			conn := cl.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)
			req := &v1.SayHelloRequest{Name: "test"}

			_, err := client.SayHello(ctx, req)

			Convey("Then I should have exhausted resources", func() {
				So(err, ShouldBeError)
				So(status.Code(err), ShouldEqual, codes.ResourceExhausted)
			})
		})

		lc.RequireStop()
	})
}

func TestLimiterUnlimitedUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)

		l, k, err := limiter.New(test.NewLimiterConfig("user-agent", "10-S"))
		So(err, ShouldBeNil)

		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)

		s := &test.Server{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m,
			Limiter: l, Key: k, Mux: test.GatewayMux,
		}
		s.Register()

		cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}

		lc.RequireStart()

		Convey("When I query repeatedly", func() {
			ctx := context.Background()

			conn := cl.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)
			req := &v1.SayHelloRequest{Name: "test"}

			_, err := client.SayHello(ctx, req)

			Convey("Then I should not have exhausted resources", func() {
				So(err, ShouldBeNil)
			})
		})

		lc.RequireStop()
	})
}

func TestLimiterLimitedStream(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)

		l, k, err := limiter.New(test.NewLimiterConfig("user-agent", "0-S"))
		So(err, ShouldBeNil)

		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)

		s := &test.Server{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m,
			Limiter: l, Key: k, Mux: test.GatewayMux,
		}
		s.Register()

		cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}

		lc.RequireStart()

		Convey("When I stream repeatedly", func() {
			ctx := context.Background()

			conn := cl.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)
			req := &v1.SayStreamHelloRequest{Name: "test"}

			stream, err := client.SayStreamHello(ctx)
			So(err, ShouldBeNil)

			err = stream.Send(req)
			So(err, ShouldBeNil)

			_, err = stream.Recv()

			Convey("Then I should have exhausted resources", func() {
				So(err, ShouldBeError)
				So(status.Code(err), ShouldEqual, codes.ResourceExhausted)
			})
		})

		lc.RequireStop()
	})
}

func TestLimiterUnlimitedStream(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)

		l, k, err := limiter.New(test.NewLimiterConfig("user-agent", "10-S"))
		So(err, ShouldBeNil)

		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)

		s := &test.Server{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m,
			Limiter: l, Key: k, Mux: test.GatewayMux,
		}
		s.Register()

		cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}

		lc.RequireStart()

		Convey("When I stream repeatedly", func() {
			ctx := context.Background()

			conn := cl.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)
			req := &v1.SayStreamHelloRequest{Name: "test"}

			stream, err := client.SayStreamHello(ctx)
			So(err, ShouldBeNil)

			err = stream.Send(req)
			So(err, ShouldBeNil)

			_, err = stream.Recv()

			Convey("Then I should not have exhausted resources", func() {
				So(err, ShouldBeNil)
			})
		})

		lc.RequireStop()
	})
}

func TestLimiterAuthUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		verifier := test.NewVerifier("bob")

		l, k, err := limiter.New(test.NewLimiterConfig("token", "10-S"))
		So(err, ShouldBeNil)

		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)

		s := &test.Server{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, VerifyAuth: true,
			Limiter: l, Key: k, Verifier: verifier, Mux: test.GatewayMux,
		}
		s.Register()

		cl := &test.Client{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m,
			Generator: test.NewGenerator("bob", nil),
		}

		lc.RequireStart()

		Convey("When I query for a authenticated greet multiple times", func() {
			ctx := context.Background()

			conn := cl.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)
			req := &v1.SayHelloRequest{Name: "test"}

			var err error

			for i := 0; i < 10; i++ {
				_, err = client.SayHello(ctx, req)
			}

			Convey("Then I should not have exhausted resources", func() {
				So(err, ShouldBeNil)
			})
		})

		lc.RequireStop()
	})
}
