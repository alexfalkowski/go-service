package grpc_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/test"
	v1 "github.com/alexfalkowski/go-service/test/greet/v1"
	"github.com/alexfalkowski/go-service/transport/grpc/security/token"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestInsecureUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)

		s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}
		s.Register()

		cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}

		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := meta.WithAttribute(context.Background(), "test", meta.Redacted("test"))

			conn := cl.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)
			req := &v1.SayHelloRequest{Name: "test"}

			resp, err := client.SayHello(ctx, req)
			So(err, ShouldBeNil)

			Convey("Then I should have a valid reply", func() {
				So(resp.GetMessage(), ShouldEqual, "Hello test")
			})

			lc.RequireStop()
		})
	})
}

func TestSecureUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		cfg := test.NewSecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)

		s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}
		s.Register()

		cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Security: test.NewSecureClientConfig()}

		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := context.Background()

			conn := cl.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)
			req := &v1.SayHelloRequest{Name: "test"}

			resp, err := client.SayHello(ctx, req)
			So(err, ShouldBeNil)

			Convey("Then I should have a valid reply", func() {
				So(resp.GetMessage(), ShouldEqual, "Hello test")
			})

			lc.RequireStop()
		})
	})
}

func TestValidAuthUnary(t *testing.T) {
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

		Convey("When I query for an authenticated greet", func() {
			ctx := context.Background()

			conn := cl.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)
			req := &v1.SayHelloRequest{Name: "test"}

			resp, err := client.SayHello(ctx, req)
			So(err, ShouldBeNil)

			Convey("Then I should have a valid reply", func() {
				So(resp.GetMessage(), ShouldEqual, "Hello test")
			})

			lc.RequireStop()
		})
	})
}

//nolint:dupl
func TestInvalidAuthUnary(t *testing.T) {
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

		Convey("When I query for a unauthenticated greet", func() {
			ctx := context.Background()

			conn := cl.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)
			req := &v1.SayHelloRequest{Name: "test"}

			_, err := client.SayHello(ctx, req)

			Convey("Then I should have a unauthenticated reply", func() {
				So(status.Code(err), ShouldEqual, codes.Unauthenticated)
			})

			lc.RequireStop()
		})
	})
}

//nolint:dupl
func TestEmptyAuthUnary(t *testing.T) {
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

		Convey("When I query for a unauthenticated greet", func() {
			ctx := context.Background()

			conn := cl.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)
			req := &v1.SayHelloRequest{Name: "test"}

			_, err := client.SayHello(ctx, req)

			Convey("Then I should have a unauthenticated reply", func() {
				So(status.Code(err), ShouldEqual, codes.Unauthenticated)
			})

			lc.RequireStop()
		})
	})
}

func TestMissingClientAuthUnary(t *testing.T) {
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

		Convey("When I query for a unauthenticated greet", func() {
			ctx := context.Background()

			conn := cl.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)
			req := &v1.SayHelloRequest{Name: "test"}

			_, err := client.SayHello(ctx, req)

			Convey("Then I should have a unauthenticated reply", func() {
				So(status.Code(err), ShouldEqual, codes.Unauthenticated)
			})

			lc.RequireStop()
		})
	})
}

func TestTokenErrorAuthUnary(t *testing.T) {
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
			Credentials: token.NewPerRPCCredentials(test.NewGenerator("bob", errors.New("token error"))),
		}

		lc.RequireStart()

		Convey("When I query for a unauthenticated greet", func() {
			ctx := context.Background()

			conn := cl.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)
			req := &v1.SayHelloRequest{Name: "test"}

			_, err := client.SayHello(ctx, req)

			Convey("Then I should have a unauthenticated reply", func() {
				So(status.Code(err), ShouldEqual, codes.Unauthenticated)
			})

			lc.RequireStop()
		})
	})
}

func TestBreakerUnary(t *testing.T) {
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

		Convey("When I query for a unauthenticated greet multiple times", func() {
			ctx := context.Background()

			conn := cl.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)
			req := &v1.SayHelloRequest{Name: "test"}

			var err error

			for i := 0; i < 10; i++ {
				_, err = client.SayHello(ctx, req)
			}

			Convey("Then I should have a unauthenticated reply", func() {
				So(status.Code(err), ShouldEqual, codes.Unavailable)
			})
		})

		lc.RequireStop()
	})
}

func TestStream(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)

		s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}
		s.Register()

		cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}

		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := meta.WithAttribute(context.Background(), "test", meta.Redacted("test"))

			conn := cl.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)

			ctx, cancel := context.WithDeadline(ctx, time.Now().Add(10*time.Minute))
			defer cancel()

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
