package grpc_test

import (
	"context"
	"testing"

	"github.com/alexfalkowski/go-service/limiter"
	"github.com/alexfalkowski/go-service/test"
	v1 "github.com/alexfalkowski/go-service/test/greet/v1"
	gl "github.com/alexfalkowski/go-service/transport/grpc/limiter"
	tm "github.com/alexfalkowski/go-service/transport/meta"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestLimiterLimitedUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)

		l, err := limiter.New(test.NewLimiterConfig("0-S"))
		So(err, ShouldBeNil)

		cfg := test.NewInsecureTransportConfig()
		m := test.NewMeter(lc)
		hs := test.NewHTTPServer(lc, logger, test.NewOTLPTracerConfig(), cfg, m, nil)
		gs := test.NewGRPCServer(lc, logger, test.NewOTLPTracerConfig(), cfg, false, m,
			[]grpc.UnaryServerInterceptor{gl.UnaryServerInterceptor(l, tm.UserAgent)},
			[]grpc.StreamServerInterceptor{gl.StreamServerInterceptor(l, tm.UserAgent)},
		)

		test.RegisterTransport(lc, gs, hs)
		lc.RequireStart()

		Convey("When I query repeatedly", func() {
			ctx := context.Background()
			conn := test.NewGRPCClient(lc, logger, cfg, test.NewOTLPTracerConfig(), nil, m)

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

		l, err := limiter.New(test.NewLimiterConfig("10-S"))
		So(err, ShouldBeNil)

		cfg := test.NewInsecureTransportConfig()
		m := test.NewMeter(lc)
		hs := test.NewHTTPServer(lc, logger, test.NewOTLPTracerConfig(), cfg, m, nil)
		gs := test.NewGRPCServer(lc, logger, test.NewOTLPTracerConfig(), cfg, false, m,
			[]grpc.UnaryServerInterceptor{gl.UnaryServerInterceptor(l, tm.UserAgent)},
			[]grpc.StreamServerInterceptor{gl.StreamServerInterceptor(l, tm.UserAgent)},
		)

		test.RegisterTransport(lc, gs, hs)
		lc.RequireStart()

		Convey("When I query repeatedly", func() {
			ctx := context.Background()
			conn := test.NewGRPCClient(lc, logger, cfg, test.NewOTLPTracerConfig(), nil, m)

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

		l, err := limiter.New(test.NewLimiterConfig("0-S"))
		So(err, ShouldBeNil)

		cfg := test.NewInsecureTransportConfig()
		m := test.NewMeter(lc)
		hs := test.NewHTTPServer(lc, logger, test.NewOTLPTracerConfig(), cfg, m, nil)
		gs := test.NewGRPCServer(lc, logger, test.NewOTLPTracerConfig(), cfg, false, m,
			[]grpc.UnaryServerInterceptor{gl.UnaryServerInterceptor(l, tm.UserAgent)},
			[]grpc.StreamServerInterceptor{gl.StreamServerInterceptor(l, tm.UserAgent)},
		)

		test.RegisterTransport(lc, gs, hs)
		lc.RequireStart()

		Convey("When I stream repeatedly", func() {
			ctx := context.Background()
			conn := test.NewGRPCClient(lc, logger, cfg, test.NewOTLPTracerConfig(), nil, m)

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

		l, err := limiter.New(test.NewLimiterConfig("10-S"))
		So(err, ShouldBeNil)

		cfg := test.NewInsecureTransportConfig()
		m := test.NewMeter(lc)
		hs := test.NewHTTPServer(lc, logger, test.NewOTLPTracerConfig(), cfg, m, nil)
		gs := test.NewGRPCServer(lc, logger, test.NewOTLPTracerConfig(), cfg, false, m,
			[]grpc.UnaryServerInterceptor{gl.UnaryServerInterceptor(l, tm.UserAgent)},
			[]grpc.StreamServerInterceptor{gl.StreamServerInterceptor(l, tm.UserAgent)},
		)

		test.RegisterTransport(lc, gs, hs)
		lc.RequireStart()

		Convey("When I stream repeatedly", func() {
			ctx := context.Background()
			conn := test.NewGRPCClient(lc, logger, cfg, test.NewOTLPTracerConfig(), nil, m)

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
