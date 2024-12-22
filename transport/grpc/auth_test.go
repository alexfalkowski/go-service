package grpc_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/alexfalkowski/go-service/test"
	v1 "github.com/alexfalkowski/go-service/test/greet/v1"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestTokenErrorAuthUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		verifier := test.NewVerifier("test")
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)

		s := &test.Server{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, VerifyAuth: true,
			Verifier: verifier, Mux: mux,
		}
		s.Register()

		cl := &test.Client{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m,
			Generator: test.NewGenerator("bob", test.ErrGenerate),
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

func TestEmptyAuthUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		verifier := test.NewVerifier("test")
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)

		s := &test.Server{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, VerifyAuth: true,
			Verifier: verifier, Mux: mux,
		}
		s.Register()

		cl := &test.Client{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m,
			Generator: test.NewGenerator("", nil),
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
		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		verifier := test.NewVerifier("test")
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)

		s := &test.Server{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, VerifyAuth: true,
			Verifier: verifier, Mux: mux,
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

func TestInvalidAuthUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		verifier := test.NewVerifier("test")
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)

		s := &test.Server{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, VerifyAuth: true,
			Verifier: verifier, Mux: mux,
		}
		s.Register()

		cl := &test.Client{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m,
			Generator: test.NewGenerator("bob", nil),
		}

		lc.RequireStart()

		Convey("When I query for a unauthenticated greet", func() {
			ctx := context.Background()
			ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "What Invalid")
			ctx = metadata.AppendToOutgoingContext(ctx, "x-forwarded-for", "127.0.0.1")
			ctx = metadata.AppendToOutgoingContext(ctx, "geolocation", "geo:47,11")

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

func TestValidAuthUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		verifier := test.NewVerifier("test")
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)

		s := &test.Server{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, VerifyAuth: true,
			Verifier: verifier, Mux: mux,
		}
		s.Register()

		cl := &test.Client{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m,
			Generator: test.NewGenerator("test", nil),
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

func TestBreakerAuthUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		verifier := test.NewVerifier("test")
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)

		s := &test.Server{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, VerifyAuth: true,
			Verifier: verifier, Mux: mux,
		}
		s.Register()

		cl := &test.Client{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m,
			Generator: test.NewGenerator("bob", nil),
		}

		lc.RequireStart()

		Convey("When I query for a unauthenticated greet multiple times", func() {
			ctx := context.Background()

			conn := cl.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)
			req := &v1.SayHelloRequest{Name: "test"}

			var err error

			for range 10 {
				_, err = client.SayHello(ctx, req)
			}

			Convey("Then I should have a unauthenticated reply", func() {
				So(status.Code(err), ShouldEqual, codes.Unavailable)
			})
		})

		lc.RequireStop()
	})
}

func TestValidAuthStream(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		verifier := test.NewVerifier("test")
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)

		s := &test.Server{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, VerifyAuth: true,
			Verifier: verifier, Mux: mux,
		}
		s.Register()

		cl := &test.Client{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m,
			Generator: test.NewGenerator("test", nil),
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
		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		verifier := test.NewVerifier("test")
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)

		s := &test.Server{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, VerifyAuth: true,
			Verifier: verifier, Mux: mux,
		}
		s.Register()

		cl := &test.Client{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m,
			Generator: test.NewGenerator("bob", nil),
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
		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		verifier := test.NewVerifier("test")
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)

		s := &test.Server{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, VerifyAuth: true,
			Verifier: verifier, Mux: mux,
		}
		s.Register()

		cl := &test.Client{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m,
			Generator: test.NewGenerator("", nil),
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
		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		verifier := test.NewVerifier("test")
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)

		s := &test.Server{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, VerifyAuth: true,
			Verifier: verifier, Mux: mux,
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
		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		verifier := test.NewVerifier("test")
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)

		s := &test.Server{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, VerifyAuth: true,
			Verifier: verifier, Mux: mux,
		}
		s.Register()

		cl := &test.Client{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m,
			Generator: test.NewGenerator("", test.ErrGenerate),
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
