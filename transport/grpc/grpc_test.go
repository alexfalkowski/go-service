package grpc_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/limiter"
	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/test"
	v1 "github.com/alexfalkowski/go-service/test/greet/v1"
	gl "github.com/alexfalkowski/go-service/transport/grpc/limiter"
	"github.com/alexfalkowski/go-service/transport/grpc/security/token"
	tm "github.com/alexfalkowski/go-service/transport/meta"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func init() {
	tracer.Register()
}

func TestInsecureUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		cfg := test.NewInsecureTransportConfig()

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		hs := test.NewHTTPServer(lc, logger, test.NewDefaultTracerConfig(), cfg, m, nil)
		gs := test.NewGRPCServer(lc, logger, test.NewDefaultTracerConfig(), cfg, false, m, nil, nil)

		test.RegisterTransport(lc, gs, hs)
		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := meta.WithAttribute(context.Background(), "test", meta.SafeValue("test"))
			conn := test.NewGRPCClient(ctx, lc, logger, cfg, test.NewDefaultTracerConfig(), nil, m)

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

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		hs := test.NewHTTPServer(lc, logger, test.NewDefaultTracerConfig(), cfg, m, nil)
		gs := test.NewGRPCServer(lc, logger, test.NewDefaultTracerConfig(), cfg, false, m, nil, nil)

		test.RegisterTransport(lc, gs, hs)
		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := context.Background()
			conn := test.NewSecureGRPCClient(ctx, lc, logger, cfg, test.NewDefaultTracerConfig(), m)

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

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		hs := test.NewHTTPServer(lc, logger, test.NewDefaultTracerConfig(), cfg, m, nil)
		gs := test.NewGRPCServer(lc, logger, test.NewDefaultTracerConfig(), cfg, true, m,
			[]grpc.UnaryServerInterceptor{token.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{token.StreamServerInterceptor(verifier)},
		)

		test.RegisterTransport(lc, gs, hs)
		lc.RequireStart()

		Convey("When I query for an authenticated greet", func() {
			ctx := context.Background()
			conn := test.NewGRPCClient(ctx, lc, logger, cfg, test.NewDefaultTracerConfig(), token.NewPerRPCCredentials(test.NewGenerator("test", nil)), m)

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

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		hs := test.NewHTTPServer(lc, logger, test.NewDefaultTracerConfig(), cfg, m, nil)
		gs := test.NewGRPCServer(lc, logger, test.NewDefaultTracerConfig(), cfg, true, m,
			[]grpc.UnaryServerInterceptor{token.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{token.StreamServerInterceptor(verifier)},
		)

		test.RegisterTransport(lc, gs, hs)
		lc.RequireStart()

		Convey("When I query for a unauthenticated greet", func() {
			ctx := context.Background()
			conn := test.NewGRPCClient(ctx, lc, logger, cfg, test.NewDefaultTracerConfig(), token.NewPerRPCCredentials(test.NewGenerator("bob", nil)), m)

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

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		hs := test.NewHTTPServer(lc, logger, test.NewDefaultTracerConfig(), cfg, m, nil)
		gs := test.NewGRPCServer(lc, logger, test.NewDefaultTracerConfig(), cfg, true, m,
			[]grpc.UnaryServerInterceptor{token.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{token.StreamServerInterceptor(verifier)},
		)

		test.RegisterTransport(lc, gs, hs)
		lc.RequireStart()

		Convey("When I query for a unauthenticated greet", func() {
			ctx := context.Background()
			conn := test.NewGRPCClient(ctx, lc, logger, cfg, test.NewDefaultTracerConfig(), token.NewPerRPCCredentials(test.NewGenerator("", nil)), m)

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

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		hs := test.NewHTTPServer(lc, logger, test.NewDefaultTracerConfig(), cfg, m, nil)
		gs := test.NewGRPCServer(lc, logger, test.NewDefaultTracerConfig(), cfg, true, m,
			[]grpc.UnaryServerInterceptor{token.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{token.StreamServerInterceptor(verifier)},
		)

		test.RegisterTransport(lc, gs, hs)
		lc.RequireStart()

		Convey("When I query for a unauthenticated greet", func() {
			ctx := context.Background()
			conn := test.NewGRPCClient(ctx, lc, logger, cfg, test.NewDefaultTracerConfig(), nil, m)

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

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		hs := test.NewHTTPServer(lc, logger, test.NewDefaultTracerConfig(), cfg, m, nil)
		gs := test.NewGRPCServer(lc, logger, test.NewDefaultTracerConfig(), cfg, true, m,
			[]grpc.UnaryServerInterceptor{token.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{token.StreamServerInterceptor(verifier)},
		)

		test.RegisterTransport(lc, gs, hs)
		lc.RequireStart()

		Convey("When I query for a unauthenticated greet", func() {
			ctx := context.Background()
			conn := test.NewGRPCClient(ctx, lc, logger, cfg, test.NewDefaultTracerConfig(), token.NewPerRPCCredentials(test.NewGenerator("bob", errors.New("token error"))), m)

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

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		hs := test.NewHTTPServer(lc, logger, test.NewDefaultTracerConfig(), cfg, m, nil)
		gs := test.NewGRPCServer(lc, logger, test.NewDefaultTracerConfig(), cfg, true, m,
			[]grpc.UnaryServerInterceptor{token.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{token.StreamServerInterceptor(verifier)},
		)

		test.RegisterTransport(lc, gs, hs)
		lc.RequireStart()

		Convey("When I query for a unauthenticated greet multiple times", func() {
			ctx := context.Background()
			conn := test.NewGRPCClient(ctx, lc, logger, cfg, test.NewDefaultTracerConfig(), token.NewPerRPCCredentials(test.NewGenerator("bob", nil)), m)

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

func TestLimiterUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)

		l, err := limiter.New("0-S")
		So(err, ShouldBeNil)

		cfg := test.NewInsecureTransportConfig()

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		hs := test.NewHTTPServer(lc, logger, test.NewDefaultTracerConfig(), cfg, m, nil)
		gs := test.NewGRPCServer(lc, logger, test.NewDefaultTracerConfig(), cfg, false, m,
			[]grpc.UnaryServerInterceptor{gl.UnaryServerInterceptor(l, tm.UserAgent)},
			[]grpc.StreamServerInterceptor{gl.StreamServerInterceptor(l, tm.UserAgent)},
		)

		test.RegisterTransport(lc, gs, hs)
		lc.RequireStart()

		Convey("When I query repeatedly", func() {
			ctx := context.Background()
			conn := test.NewGRPCClient(ctx, lc, logger, cfg, test.NewDefaultTracerConfig(), nil, m)

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

func TestStream(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		cfg := test.NewInsecureTransportConfig()

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		hs := test.NewHTTPServer(lc, logger, test.NewDefaultTracerConfig(), cfg, m, nil)
		gs := test.NewGRPCServer(lc, logger, test.NewDefaultTracerConfig(), cfg, false, m, nil, nil)

		test.RegisterTransport(lc, gs, hs)
		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := meta.WithAttribute(context.Background(), "test", meta.SafeValue("test"))
			conn := test.NewGRPCClient(ctx, lc, logger, cfg, test.NewDefaultTracerConfig(), nil, m)

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

func TestValidAuthStream(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		verifier := test.NewVerifier("test")
		cfg := test.NewInsecureTransportConfig()

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		hs := test.NewHTTPServer(lc, logger, test.NewDefaultTracerConfig(), cfg, m, nil)
		gs := test.NewGRPCServer(lc, logger, test.NewDefaultTracerConfig(), cfg, true, m,
			[]grpc.UnaryServerInterceptor{token.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{token.StreamServerInterceptor(verifier)},
		)

		test.RegisterTransport(lc, gs, hs)
		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := context.Background()
			conn := test.NewGRPCClient(ctx, lc, logger, cfg, test.NewDefaultTracerConfig(), token.NewPerRPCCredentials(test.NewGenerator("test", nil)), m)

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

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		hs := test.NewHTTPServer(lc, logger, test.NewDefaultTracerConfig(), cfg, m, nil)
		gs := test.NewGRPCServer(lc, logger, test.NewDefaultTracerConfig(), cfg, true, m,
			[]grpc.UnaryServerInterceptor{token.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{token.StreamServerInterceptor(verifier)},
		)

		test.RegisterTransport(lc, gs, hs)
		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := context.Background()
			conn := test.NewGRPCClient(ctx, lc, logger, cfg, test.NewDefaultTracerConfig(), token.NewPerRPCCredentials(test.NewGenerator("bob", nil)), m)

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

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		hs := test.NewHTTPServer(lc, logger, test.NewDefaultTracerConfig(), cfg, m, nil)
		gs := test.NewGRPCServer(lc, logger, test.NewDefaultTracerConfig(), cfg, true, m,
			[]grpc.UnaryServerInterceptor{token.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{token.StreamServerInterceptor(verifier)},
		)

		test.RegisterTransport(lc, gs, hs)
		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := context.Background()
			conn := test.NewGRPCClient(ctx, lc, logger, cfg, test.NewDefaultTracerConfig(), token.NewPerRPCCredentials(test.NewGenerator("", nil)), m)

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

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		hs := test.NewHTTPServer(lc, logger, test.NewDefaultTracerConfig(), cfg, m, nil)
		gs := test.NewGRPCServer(lc, logger, test.NewDefaultTracerConfig(), cfg, true, m,
			[]grpc.UnaryServerInterceptor{token.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{token.StreamServerInterceptor(verifier)},
		)

		test.RegisterTransport(lc, gs, hs)
		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := context.Background()
			conn := test.NewGRPCClient(ctx, lc, logger, cfg, test.NewDefaultTracerConfig(), nil, m)

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

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		hs := test.NewHTTPServer(lc, logger, test.NewDefaultTracerConfig(), cfg, m, nil)
		gs := test.NewGRPCServer(lc, logger, test.NewDefaultTracerConfig(), cfg, true, m,
			[]grpc.UnaryServerInterceptor{token.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{token.StreamServerInterceptor(verifier)},
		)

		test.RegisterTransport(lc, gs, hs)
		lc.RequireStart()

		Convey("When I query for a greet that will generate a token error", func() {
			ctx := context.Background()
			conn := test.NewGRPCClient(ctx, lc, logger, cfg, test.NewDefaultTracerConfig(), token.NewPerRPCCredentials(test.NewGenerator("", errors.New("token error"))), m)

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

func TestLimiterStream(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)

		l, err := limiter.New("0-S")
		So(err, ShouldBeNil)

		cfg := test.NewInsecureTransportConfig()

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		hs := test.NewHTTPServer(lc, logger, test.NewDefaultTracerConfig(), cfg, m, nil)
		gs := test.NewGRPCServer(lc, logger, test.NewDefaultTracerConfig(), cfg, false, m,
			[]grpc.UnaryServerInterceptor{gl.UnaryServerInterceptor(l, tm.UserAgent)},
			[]grpc.StreamServerInterceptor{gl.StreamServerInterceptor(l, tm.UserAgent)},
		)

		test.RegisterTransport(lc, gs, hs)
		lc.RequireStart()

		Convey("When I stream repeatedly", func() {
			ctx := context.Background()
			conn := test.NewGRPCClient(ctx, lc, logger, cfg, test.NewDefaultTracerConfig(), nil, m)

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
