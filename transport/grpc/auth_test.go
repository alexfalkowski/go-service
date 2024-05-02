package grpc_test

import (
	"context"
	"errors"
	"testing"

	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/test"
	v1 "github.com/alexfalkowski/go-service/test/greet/v1"
	"github.com/alexfalkowski/go-service/transport/grpc/security/token"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func init() {
	tracer.Register()
}

func TestValidAuthStream(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		verifier := test.NewVerifier("test")
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)

		s := &test.Server{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, VerifyAuth: true,
			Unary:  []grpc.UnaryServerInterceptor{token.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{token.StreamServerInterceptor(verifier)},
		}
		s.Register()

		cl := &test.Client{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m,
			Credentials: token.NewPerRPCCredentials(test.NewGenerator("test", nil)),
		}

		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := context.Background()

			conn := cl.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)

			stream, err := client.SayStreamHello(ctx)
			So(err, ShouldBeNil)

			err = stream.Send(&v1.SayStreamHelloRequest{Name: "test"})
			So(err, ShouldBeNil)

			resp, err := stream.Recv()
			So(err, ShouldBeNil)

			Convey("Then I should have a valid reply", func() {
				So(resp.GetMessage(), ShouldEqual, "Hello test")
			})

			lc.RequireStop()
		})
	})
}

func TestInvalidAuthStream(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		verifier := test.NewVerifier("test")
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)

		s := &test.Server{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, VerifyAuth: true,
			Unary:  []grpc.UnaryServerInterceptor{token.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{token.StreamServerInterceptor(verifier)},
		}
		s.Register()

		cl := &test.Client{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m,
			Credentials: token.NewPerRPCCredentials(test.NewGenerator("bob", nil)),
		}

		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := context.Background()

			conn := cl.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)

			stream, err := client.SayStreamHello(ctx)
			So(err, ShouldBeNil)

			err = stream.Send(&v1.SayStreamHelloRequest{Name: "test"})
			So(err, ShouldBeNil)

			_, err = stream.Recv()

			Convey("Then I should have a unauthenticated reply", func() {
				So(status.Code(err), ShouldEqual, codes.Unauthenticated)
			})

			lc.RequireStop()
		})
	})
}

func TestEmptyAuthStream(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		verifier := test.NewVerifier("test")
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)

		s := &test.Server{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, VerifyAuth: true,
			Unary:  []grpc.UnaryServerInterceptor{token.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{token.StreamServerInterceptor(verifier)},
		}
		s.Register()

		cl := &test.Client{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m,
			Credentials: token.NewPerRPCCredentials(test.NewGenerator("", nil)),
		}

		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := context.Background()

			conn := cl.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)

			_, err := client.SayStreamHello(ctx)

			Convey("Then I should have an auth error", func() {
				So(status.Code(err), ShouldEqual, codes.Unauthenticated)
			})

			lc.RequireStop()
		})
	})
}

func TestMissingClientAuthStream(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		verifier := test.NewVerifier("test")
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)

		s := &test.Server{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, VerifyAuth: true,
			Unary:  []grpc.UnaryServerInterceptor{token.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{token.StreamServerInterceptor(verifier)},
		}
		s.Register()

		cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}

		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := context.Background()

			conn := cl.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)

			stream, err := client.SayStreamHello(ctx)
			So(err, ShouldBeNil)

			err = stream.Send(&v1.SayStreamHelloRequest{Name: "test"})
			So(err, ShouldBeNil)

			_, err = stream.Recv()

			Convey("Then I should have a unauthenticated reply", func() {
				So(status.Code(err), ShouldEqual, codes.Unauthenticated)
			})

			lc.RequireStop()
		})
	})
}

func TestTokenErrorAuthStream(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		verifier := test.NewVerifier("test")
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)

		s := &test.Server{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, VerifyAuth: true,
			Unary:  []grpc.UnaryServerInterceptor{token.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{token.StreamServerInterceptor(verifier)},
		}
		s.Register()

		cl := &test.Client{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m,
			Credentials: token.NewPerRPCCredentials(test.NewGenerator("", errors.New("token error"))),
		}

		lc.RequireStart()

		Convey("When I query for a greet that will generate a token error", func() {
			ctx := context.Background()

			conn := cl.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)

			_, err := client.SayStreamHello(ctx)

			Convey("Then I should have an error", func() {
				So(status.Code(err), ShouldEqual, codes.Unauthenticated)
			})

			lc.RequireStop()
		})
	})
}
