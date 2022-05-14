package grpc_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/test"
	v1 "github.com/alexfalkowski/go-service/test/greet/v1"
	"github.com/alexfalkowski/go-service/transport/grpc/security/jwt"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		_, gconfig := test.NewGRPCServer(lc, logger, test.NewJaegerConfig(), false, nil, nil)

		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := context.Background()
			conn := test.NewGRPCClient(ctx, lc, gconfig, logger, test.NewJaegerConfig(), nil)
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)
			req := &v1.SayHelloRequest{Name: "test"}

			ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Minute))
			defer cancel()

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
		_, gconfig := test.NewGRPCServer(lc, logger, test.NewJaegerConfig(), true,
			[]grpc.UnaryServerInterceptor{jwt.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{jwt.StreamServerInterceptor(verifier)},
		)

		lc.RequireStart()

		Convey("When I query for an authenticated greet", func() {
			ctx := context.Background()
			conn := test.NewGRPCClient(ctx, lc, gconfig, logger, test.NewJaegerConfig(), jwt.NewPerRPCCredentials(test.NewGenerator("test", nil)))
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

// nolint:dupl
func TestInvalidAuthUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		verifier := test.NewVerifier("test")
		_, gconfig := test.NewGRPCServer(lc, logger, test.NewJaegerConfig(), true,
			[]grpc.UnaryServerInterceptor{jwt.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{jwt.StreamServerInterceptor(verifier)},
		)

		lc.RequireStart()

		Convey("When I query for a unauthenticated greet", func() {
			ctx := context.Background()
			conn := test.NewGRPCClient(ctx, lc, gconfig, logger, test.NewJaegerConfig(), jwt.NewPerRPCCredentials(test.NewGenerator("bob", nil)))
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

// nolint:dupl
func TestEmptyAuthUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		verifier := test.NewVerifier("test")
		_, gconfig := test.NewGRPCServer(lc, logger, test.NewJaegerConfig(), true,
			[]grpc.UnaryServerInterceptor{jwt.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{jwt.StreamServerInterceptor(verifier)},
		)

		lc.RequireStart()

		Convey("When I query for a unauthenticated greet", func() {
			ctx := context.Background()
			conn := test.NewGRPCClient(ctx, lc, gconfig, logger, test.NewJaegerConfig(), jwt.NewPerRPCCredentials(test.NewGenerator("", nil)))
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
		_, gconfig := test.NewGRPCServer(lc, logger, test.NewJaegerConfig(), true,
			[]grpc.UnaryServerInterceptor{jwt.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{jwt.StreamServerInterceptor(verifier)},
		)

		lc.RequireStart()

		Convey("When I query for a unauthenticated greet", func() {
			ctx := context.Background()
			conn := test.NewGRPCClient(ctx, lc, gconfig, logger, test.NewJaegerConfig(), nil)
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
		_, gconfig := test.NewGRPCServer(lc, logger, test.NewJaegerConfig(), true,
			[]grpc.UnaryServerInterceptor{jwt.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{jwt.StreamServerInterceptor(verifier)},
		)

		lc.RequireStart()

		Convey("When I query for a unauthenticated greet", func() {
			ctx := context.Background()
			conn := test.NewGRPCClient(ctx, lc, gconfig, logger, test.NewJaegerConfig(), jwt.NewPerRPCCredentials(test.NewGenerator("bob", errors.New("token error"))))
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

func TestStream(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		_, gconfig := test.NewGRPCServer(lc, logger, test.NewJaegerConfig(), false, nil, nil)

		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := context.Background()
			conn := test.NewGRPCClient(ctx, lc, gconfig, logger, test.NewJaegerConfig(), nil)
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)

			ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Minute))
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
		_, gconfig := test.NewGRPCServer(lc, logger, test.NewJaegerConfig(), true,
			[]grpc.UnaryServerInterceptor{jwt.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{jwt.StreamServerInterceptor(verifier)},
		)

		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := context.Background()
			conn := test.NewGRPCClient(ctx, lc, gconfig, logger, test.NewJaegerConfig(), jwt.NewPerRPCCredentials(test.NewGenerator("test", nil)))
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
		_, gconfig := test.NewGRPCServer(lc, logger, test.NewJaegerConfig(), true,
			[]grpc.UnaryServerInterceptor{jwt.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{jwt.StreamServerInterceptor(verifier)},
		)

		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := context.Background()
			conn := test.NewGRPCClient(ctx, lc, gconfig, logger, test.NewJaegerConfig(), jwt.NewPerRPCCredentials(test.NewGenerator("bob", nil)))
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
		_, gconfig := test.NewGRPCServer(lc, logger, test.NewJaegerConfig(), true,
			[]grpc.UnaryServerInterceptor{jwt.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{jwt.StreamServerInterceptor(verifier)},
		)

		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := context.Background()
			conn := test.NewGRPCClient(ctx, lc, gconfig, logger, test.NewJaegerConfig(), jwt.NewPerRPCCredentials(test.NewGenerator("", nil)))
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
		_, gconfig := test.NewGRPCServer(lc, logger, test.NewJaegerConfig(), true,
			[]grpc.UnaryServerInterceptor{jwt.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{jwt.StreamServerInterceptor(verifier)},
		)

		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := context.Background()
			conn := test.NewGRPCClient(ctx, lc, gconfig, logger, test.NewJaegerConfig(), nil)
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
		_, gconfig := test.NewGRPCServer(lc, logger, test.NewJaegerConfig(), true,
			[]grpc.UnaryServerInterceptor{jwt.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{jwt.StreamServerInterceptor(verifier)},
		)

		lc.RequireStart()

		Convey("When I query for a greet that will generate a token error", func() {
			ctx := context.Background()
			conn := test.NewGRPCClient(ctx, lc, gconfig, logger, test.NewJaegerConfig(), jwt.NewPerRPCCredentials(test.NewGenerator("", errors.New("token error"))))
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

func TestBreakerUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		verifier := test.NewVerifier("test")
		_, gconfig := test.NewGRPCServer(lc, logger, test.NewJaegerConfig(), true,
			[]grpc.UnaryServerInterceptor{jwt.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{jwt.StreamServerInterceptor(verifier)},
		)

		lc.RequireStart()

		Convey("When I query for a unauthenticated greet multiple times", func() {
			ctx := context.Background()
			conn := test.NewGRPCClient(ctx, lc, gconfig, logger, test.NewJaegerConfig(), jwt.NewPerRPCCredentials(test.NewGenerator("bob", nil)))

			client := v1.NewGreeterServiceClient(conn)
			req := &v1.SayHelloRequest{Name: "test"}

			var err error

			for i := 0; i < 10; i++ {
				_, err = client.SayHello(ctx, req)
			}

			Convey("Then I should have a unauthenticated reply", func() {
				So(status.Code(err), ShouldEqual, codes.Unavailable)
			})

			lc.RequireStop()
		})
	})
}
