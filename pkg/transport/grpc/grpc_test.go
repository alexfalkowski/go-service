package grpc_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/pkg/logger/zap"
	pkgGRPC "github.com/alexfalkowski/go-service/pkg/transport/grpc"
	tokenGRPC "github.com/alexfalkowski/go-service/pkg/transport/grpc/security/token"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cfg := &pkgGRPC.Config{GRPCPort: "10007"}
		serverParams := pkgGRPC.ServerParams{Config: cfg, Logger: logger}
		gs := pkgGRPC.NewServer(lc, test.NewShutdowner(), serverParams)

		test.RegisterGreeterServer(gs, test.NewServer())

		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := context.Background()
			clientParams := &pkgGRPC.ClientParams{Logger: logger}
			clientOpts := []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()}

			conn, err := pkgGRPC.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", cfg.GRPCPort), clientParams, clientOpts...)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := test.NewGreeterClient(conn)
			req := &test.HelloRequest{Name: "test"}

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

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cfg := &pkgGRPC.Config{GRPCPort: "10007"}
		verifier := test.NewVerifier("test")
		serverParams := pkgGRPC.ServerParams{
			Config: cfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{tokenGRPC.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{tokenGRPC.StreamServerInterceptor(verifier)},
		}
		gs := pkgGRPC.NewServer(lc, test.NewShutdowner(), serverParams)

		test.RegisterGreeterServer(gs, test.NewServer())

		lc.RequireStart()

		Convey("When I query for an authenticated greet", func() {
			ctx := context.Background()
			clientParams := &pkgGRPC.ClientParams{Logger: logger}
			clientOpts := []grpc.DialOption{
				grpc.WithBlock(),
				grpc.WithInsecure(),
				grpc.WithPerRPCCredentials(tokenGRPC.NewPerRPCCredentials(test.NewGenerator("test", nil))),
			}

			conn, err := pkgGRPC.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", cfg.GRPCPort), clientParams, clientOpts...)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := test.NewGreeterClient(conn)
			req := &test.HelloRequest{Name: "test"}

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

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cfg := &pkgGRPC.Config{GRPCPort: "10007"}
		verifier := test.NewVerifier("test")
		serverParams := pkgGRPC.ServerParams{
			Config: cfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{tokenGRPC.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{tokenGRPC.StreamServerInterceptor(verifier)},
		}
		gs := pkgGRPC.NewServer(lc, test.NewShutdowner(), serverParams)

		test.RegisterGreeterServer(gs, test.NewServer())

		lc.RequireStart()

		Convey("When I query for a unauthenticated greet", func() {
			ctx := context.Background()
			clientParams := &pkgGRPC.ClientParams{Logger: logger}
			clientOpts := []grpc.DialOption{
				grpc.WithBlock(),
				grpc.WithInsecure(),
				grpc.WithPerRPCCredentials(tokenGRPC.NewPerRPCCredentials(test.NewGenerator("bob", nil))),
			}

			conn, err := pkgGRPC.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", cfg.GRPCPort), clientParams, clientOpts...)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := test.NewGreeterClient(conn)
			req := &test.HelloRequest{Name: "test"}

			Convey("Then I should have a unauthenticated reply", func() {
				_, err := client.SayHello(ctx, req)
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

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cfg := &pkgGRPC.Config{GRPCPort: "10007"}
		verifier := test.NewVerifier("test")
		serverParams := pkgGRPC.ServerParams{
			Config: cfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{tokenGRPC.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{tokenGRPC.StreamServerInterceptor(verifier)},
		}
		gs := pkgGRPC.NewServer(lc, test.NewShutdowner(), serverParams)

		test.RegisterGreeterServer(gs, test.NewServer())

		lc.RequireStart()

		Convey("When I query for a unauthenticated greet", func() {
			ctx := context.Background()
			clientParams := &pkgGRPC.ClientParams{Logger: logger}
			clientOpts := []grpc.DialOption{
				grpc.WithBlock(),
				grpc.WithInsecure(),
				grpc.WithPerRPCCredentials(tokenGRPC.NewPerRPCCredentials(test.NewGenerator("", nil))),
			}

			conn, err := pkgGRPC.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", cfg.GRPCPort), clientParams, clientOpts...)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := test.NewGreeterClient(conn)
			req := &test.HelloRequest{Name: "test"}

			Convey("Then I should have a unauthenticated reply", func() {
				_, err := client.SayHello(ctx, req)
				So(status.Code(err), ShouldEqual, codes.Unauthenticated)
			})

			lc.RequireStop()
		})
	})
}

func TestMissingClientAuthUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cfg := &pkgGRPC.Config{GRPCPort: "10007"}
		verifier := test.NewVerifier("test")
		serverParams := pkgGRPC.ServerParams{
			Config: cfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{tokenGRPC.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{tokenGRPC.StreamServerInterceptor(verifier)},
		}
		gs := pkgGRPC.NewServer(lc, test.NewShutdowner(), serverParams)

		test.RegisterGreeterServer(gs, test.NewServer())

		lc.RequireStart()

		Convey("When I query for a unauthenticated greet", func() {
			ctx := context.Background()
			clientParams := &pkgGRPC.ClientParams{Logger: logger}
			clientOpts := []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()}

			conn, err := pkgGRPC.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", cfg.GRPCPort), clientParams, clientOpts...)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := test.NewGreeterClient(conn)
			req := &test.HelloRequest{Name: "test"}

			Convey("Then I should have a unauthenticated reply", func() {
				_, err := client.SayHello(ctx, req)
				So(status.Code(err), ShouldEqual, codes.Unauthenticated)
			})

			lc.RequireStop()
		})
	})
}

func TestTokenErrorAuthUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cfg := &pkgGRPC.Config{GRPCPort: "10007"}
		verifier := test.NewVerifier("test")
		serverParams := pkgGRPC.ServerParams{
			Config: cfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{tokenGRPC.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{tokenGRPC.StreamServerInterceptor(verifier)},
		}
		gs := pkgGRPC.NewServer(lc, test.NewShutdowner(), serverParams)

		test.RegisterGreeterServer(gs, test.NewServer())

		lc.RequireStart()

		Convey("When I query for a unauthenticated greet", func() {
			ctx := context.Background()
			clientParams := &pkgGRPC.ClientParams{Logger: logger}
			clientOpts := []grpc.DialOption{
				grpc.WithBlock(),
				grpc.WithInsecure(),
				grpc.WithPerRPCCredentials(tokenGRPC.NewPerRPCCredentials(test.NewGenerator("bob", errors.New("token error")))),
			}

			conn, err := pkgGRPC.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", cfg.GRPCPort), clientParams, clientOpts...)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := test.NewGreeterClient(conn)
			req := &test.HelloRequest{Name: "test"}

			Convey("Then I should have a unauthenticated reply", func() {
				_, err := client.SayHello(ctx, req)
				So(status.Code(err), ShouldEqual, codes.Unauthenticated)
			})

			lc.RequireStop()
		})
	})
}

func TestStream(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cfg := &pkgGRPC.Config{GRPCPort: "10008"}
		serverParams := pkgGRPC.ServerParams{Config: cfg, Logger: logger}
		gs := pkgGRPC.NewServer(lc, test.NewShutdowner(), serverParams)

		test.RegisterGreeterServer(gs, test.NewServer())

		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := context.Background()
			clientParams := &pkgGRPC.ClientParams{Logger: logger}
			clientOpts := []grpc.DialOption{
				grpc.WithBlock(),
				grpc.WithInsecure(),
			}

			conn, err := pkgGRPC.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", cfg.GRPCPort), clientParams, clientOpts...)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := test.NewGreeterClient(conn)

			ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Minute))
			defer cancel()

			stream, err := client.SayStreamHello(ctx)
			So(err, ShouldBeNil)

			err = stream.Send(&test.HelloRequest{Name: "test"})
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

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cfg := &pkgGRPC.Config{GRPCPort: "10008"}
		verifier := test.NewVerifier("test")
		serverParams := pkgGRPC.ServerParams{
			Config: cfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{tokenGRPC.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{tokenGRPC.StreamServerInterceptor(verifier)},
		}
		gs := pkgGRPC.NewServer(lc, test.NewShutdowner(), serverParams)

		test.RegisterGreeterServer(gs, test.NewServer())

		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := context.Background()
			clientParams := &pkgGRPC.ClientParams{Logger: logger}
			clientOpts := []grpc.DialOption{
				grpc.WithBlock(),
				grpc.WithInsecure(),
				grpc.WithPerRPCCredentials(tokenGRPC.NewPerRPCCredentials(test.NewGenerator("test", nil))),
			}

			conn, err := pkgGRPC.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", cfg.GRPCPort), clientParams, clientOpts...)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := test.NewGreeterClient(conn)

			stream, err := client.SayStreamHello(ctx)
			So(err, ShouldBeNil)

			err = stream.Send(&test.HelloRequest{Name: "test"})
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

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cfg := &pkgGRPC.Config{GRPCPort: "10008"}
		verifier := test.NewVerifier("test")
		serverParams := pkgGRPC.ServerParams{
			Config: cfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{tokenGRPC.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{tokenGRPC.StreamServerInterceptor(verifier)},
		}
		gs := pkgGRPC.NewServer(lc, test.NewShutdowner(), serverParams)

		test.RegisterGreeterServer(gs, test.NewServer())

		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := context.Background()
			clientParams := &pkgGRPC.ClientParams{Logger: logger}
			clientOpts := []grpc.DialOption{
				grpc.WithBlock(),
				grpc.WithInsecure(),
				grpc.WithPerRPCCredentials(tokenGRPC.NewPerRPCCredentials(test.NewGenerator("bob", nil))),
			}

			conn, err := pkgGRPC.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", cfg.GRPCPort), clientParams, clientOpts...)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := test.NewGreeterClient(conn)

			stream, err := client.SayStreamHello(ctx)
			So(err, ShouldBeNil)

			err = stream.Send(&test.HelloRequest{Name: "test"})
			So(err, ShouldBeNil)

			Convey("Then I should have a unauthenticated reply", func() {
				_, err := stream.Recv()
				So(status.Code(err), ShouldEqual, codes.Unauthenticated)
			})

			lc.RequireStop()
		})
	})
}

func TestEmptyAuthStream(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cfg := &pkgGRPC.Config{GRPCPort: "10008"}
		verifier := test.NewVerifier("test")
		serverParams := pkgGRPC.ServerParams{
			Config: cfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{tokenGRPC.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{tokenGRPC.StreamServerInterceptor(verifier)},
		}
		gs := pkgGRPC.NewServer(lc, test.NewShutdowner(), serverParams)

		test.RegisterGreeterServer(gs, test.NewServer())

		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := context.Background()
			clientParams := &pkgGRPC.ClientParams{Logger: logger}
			clientOpts := []grpc.DialOption{
				grpc.WithBlock(),
				grpc.WithInsecure(),
				grpc.WithPerRPCCredentials(tokenGRPC.NewPerRPCCredentials(test.NewGenerator("", nil))),
			}

			conn, err := pkgGRPC.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", cfg.GRPCPort), clientParams, clientOpts...)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := test.NewGreeterClient(conn)

			Convey("Then I should have an auth error", func() {
				_, err := client.SayStreamHello(ctx)
				So(status.Code(err), ShouldEqual, codes.Unauthenticated)
			})

			lc.RequireStop()
		})
	})
}

func TestMissingClientAuthStream(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cfg := &pkgGRPC.Config{GRPCPort: "10008"}
		verifier := test.NewVerifier("test")
		serverParams := pkgGRPC.ServerParams{
			Config: cfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{tokenGRPC.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{tokenGRPC.StreamServerInterceptor(verifier)},
		}
		gs := pkgGRPC.NewServer(lc, test.NewShutdowner(), serverParams)

		test.RegisterGreeterServer(gs, test.NewServer())

		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := context.Background()
			clientParams := &pkgGRPC.ClientParams{Logger: logger}
			clientOpts := []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()}

			conn, err := pkgGRPC.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", cfg.GRPCPort), clientParams, clientOpts...)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := test.NewGreeterClient(conn)

			stream, err := client.SayStreamHello(ctx)
			So(err, ShouldBeNil)

			err = stream.Send(&test.HelloRequest{Name: "test"})
			So(err, ShouldBeNil)

			Convey("Then I should have a unauthenticated reply", func() {
				_, err := stream.Recv()
				So(status.Code(err), ShouldEqual, codes.Unauthenticated)
			})

			lc.RequireStop()
		})
	})
}

func TestTokenErrorAuthStream(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cfg := &pkgGRPC.Config{GRPCPort: "10008"}
		verifier := test.NewVerifier("test")
		serverParams := pkgGRPC.ServerParams{
			Config: cfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{tokenGRPC.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{tokenGRPC.StreamServerInterceptor(verifier)},
		}
		gs := pkgGRPC.NewServer(lc, test.NewShutdowner(), serverParams)

		test.RegisterGreeterServer(gs, test.NewServer())

		lc.RequireStart()

		Convey("When I query for a greet that will generate a token error", func() {
			ctx := context.Background()
			clientParams := &pkgGRPC.ClientParams{Logger: logger}
			clientOpts := []grpc.DialOption{
				grpc.WithBlock(),
				grpc.WithInsecure(),
				grpc.WithPerRPCCredentials(tokenGRPC.NewPerRPCCredentials(test.NewGenerator("", errors.New("token error")))),
			}

			conn, err := pkgGRPC.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", cfg.GRPCPort), clientParams, clientOpts...)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := test.NewGreeterClient(conn)

			Convey("Then I should have an error", func() {
				_, err := client.SayStreamHello(ctx)
				So(status.Code(err), ShouldEqual, codes.Unauthenticated)
			})

			lc.RequireStop()
		})
	})
}
