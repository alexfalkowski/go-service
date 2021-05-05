package grpc_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/alexfalkowski/go-service/pkg/config"
	pkgGRPC "github.com/alexfalkowski/go-service/pkg/grpc"
	tokenGRPC "github.com/alexfalkowski/go-service/pkg/grpc/security/token"
	"github.com/alexfalkowski/go-service/pkg/logger/zap"
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

		logger, err := zap.NewLogger(lc)
		So(err, ShouldBeNil)

		cfg := &config.Config{GRPCPort: "10007"}
		gs := pkgGRPC.NewServer(lc, test.NewShutdowner(), cfg, logger)

		test.RegisterGreeterServer(gs, test.NewServer())

		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := context.Background()
			opts := []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()}

			conn, err := pkgGRPC.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", cfg.GRPCPort), logger, opts...)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := test.NewGreeterClient(conn)
			req := &test.HelloRequest{Name: "test"}

			resp, err := client.SayHello(ctx, req)
			So(err, ShouldBeNil)

			lc.RequireStop()

			Convey("Then I should have a valid reply", func() {
				So(resp.GetMessage(), ShouldEqual, "Hello test")
			})
		})
	})
}

func TestValidAuthUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc)
		So(err, ShouldBeNil)

		cfg := &config.Config{GRPCPort: "10007"}
		verifier := test.NewVerifier("test")
		serverUnaryOpt := pkgGRPC.UnaryServerOption(logger, tokenGRPC.UnaryServerInterceptor(verifier))
		serverStreamOpt := pkgGRPC.StreamServerOption(logger, tokenGRPC.StreamServerInterceptor(verifier))
		gs := pkgGRPC.NewServer(lc, test.NewShutdowner(), cfg, logger, serverUnaryOpt, serverStreamOpt)

		test.RegisterGreeterServer(gs, test.NewServer())

		lc.RequireStart()

		Convey("When I query for an authenticated greet", func() {
			ctx := context.Background()
			clientOpts := []grpc.DialOption{
				grpc.WithBlock(),
				grpc.WithInsecure(),
				pkgGRPC.UnaryDialOption(logger),
				pkgGRPC.StreamDialOption(logger),
				grpc.WithPerRPCCredentials(tokenGRPC.NewPerRPCCredentials(test.NewGenerator("test"))),
			}

			conn, err := pkgGRPC.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", cfg.GRPCPort), logger, clientOpts...)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := test.NewGreeterClient(conn)
			req := &test.HelloRequest{Name: "test"}

			resp, err := client.SayHello(ctx, req)
			So(err, ShouldBeNil)

			lc.RequireStop()

			Convey("Then I should have a valid reply", func() {
				So(resp.GetMessage(), ShouldEqual, "Hello test")
			})
		})
	})
}

// nolint:dupl
func TestInvalidAuthUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc)
		So(err, ShouldBeNil)

		cfg := &config.Config{GRPCPort: "10007"}
		verifier := test.NewVerifier("test")
		serverUnaryOpt := pkgGRPC.UnaryServerOption(logger, tokenGRPC.UnaryServerInterceptor(verifier))
		serverStreamOpt := pkgGRPC.StreamServerOption(logger, tokenGRPC.StreamServerInterceptor(verifier))
		gs := pkgGRPC.NewServer(lc, test.NewShutdowner(), cfg, logger, serverUnaryOpt, serverStreamOpt)

		test.RegisterGreeterServer(gs, test.NewServer())

		lc.RequireStart()

		Convey("When I query for a unauthenticated greet", func() {
			ctx := context.Background()
			clientOpts := []grpc.DialOption{
				grpc.WithBlock(),
				grpc.WithInsecure(),
				pkgGRPC.UnaryDialOption(logger),
				pkgGRPC.StreamDialOption(logger),
				grpc.WithPerRPCCredentials(tokenGRPC.NewPerRPCCredentials(test.NewGenerator("bob"))),
			}

			conn, err := pkgGRPC.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", cfg.GRPCPort), logger, clientOpts...)
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

		logger, err := zap.NewLogger(lc)
		So(err, ShouldBeNil)

		cfg := &config.Config{GRPCPort: "10007"}
		verifier := test.NewVerifier("test")
		serverUnaryOpt := pkgGRPC.UnaryServerOption(logger, tokenGRPC.UnaryServerInterceptor(verifier))
		serverStreamOpt := pkgGRPC.StreamServerOption(logger, tokenGRPC.StreamServerInterceptor(verifier))
		gs := pkgGRPC.NewServer(lc, test.NewShutdowner(), cfg, logger, serverUnaryOpt, serverStreamOpt)

		test.RegisterGreeterServer(gs, test.NewServer())

		lc.RequireStart()

		Convey("When I query for a unauthenticated greet", func() {
			ctx := context.Background()
			clientOpts := []grpc.DialOption{
				grpc.WithBlock(),
				grpc.WithInsecure(),
				pkgGRPC.UnaryDialOption(logger),
				pkgGRPC.StreamDialOption(logger),
				grpc.WithPerRPCCredentials(tokenGRPC.NewPerRPCCredentials(test.NewGenerator(""))),
			}

			conn, err := pkgGRPC.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", cfg.GRPCPort), logger, clientOpts...)
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

		logger, err := zap.NewLogger(lc)
		So(err, ShouldBeNil)

		cfg := &config.Config{GRPCPort: "10007"}
		verifier := test.NewVerifier("test")
		serverUnaryOpt := pkgGRPC.UnaryServerOption(logger, tokenGRPC.UnaryServerInterceptor(verifier))
		serverStreamOpt := pkgGRPC.StreamServerOption(logger, tokenGRPC.StreamServerInterceptor(verifier))
		gs := pkgGRPC.NewServer(lc, test.NewShutdowner(), cfg, logger, serverUnaryOpt, serverStreamOpt)

		test.RegisterGreeterServer(gs, test.NewServer())

		lc.RequireStart()

		Convey("When I query for a unauthenticated greet", func() {
			ctx := context.Background()
			clientOpts := []grpc.DialOption{
				grpc.WithBlock(),
				grpc.WithInsecure(),
				pkgGRPC.UnaryDialOption(logger),
				pkgGRPC.StreamDialOption(logger),
			}

			conn, err := pkgGRPC.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", cfg.GRPCPort), logger, clientOpts...)
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

		logger, err := zap.NewLogger(lc)
		So(err, ShouldBeNil)

		cfg := &config.Config{GRPCPort: "10008"}
		gs := pkgGRPC.NewServer(lc, test.NewShutdowner(), cfg, logger)

		test.RegisterGreeterServer(gs, test.NewServer())

		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := context.Background()
			opts := []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()}

			conn, err := pkgGRPC.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", cfg.GRPCPort), logger, opts...)
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

func TestValidAuthStream(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc)
		So(err, ShouldBeNil)

		cfg := &config.Config{GRPCPort: "10008"}
		verifier := test.NewVerifier("test")
		serverUnaryOpt := pkgGRPC.UnaryServerOption(logger, tokenGRPC.UnaryServerInterceptor(verifier))
		serverStreamOpt := pkgGRPC.StreamServerOption(logger, tokenGRPC.StreamServerInterceptor(verifier))
		gs := pkgGRPC.NewServer(lc, test.NewShutdowner(), cfg, logger, serverUnaryOpt, serverStreamOpt)

		test.RegisterGreeterServer(gs, test.NewServer())

		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := context.Background()
			opts := []grpc.DialOption{
				grpc.WithBlock(),
				grpc.WithInsecure(),
				pkgGRPC.UnaryDialOption(logger),
				pkgGRPC.StreamDialOption(logger),
				grpc.WithPerRPCCredentials(tokenGRPC.NewPerRPCCredentials(test.NewGenerator("test"))),
			}

			conn, err := pkgGRPC.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", cfg.GRPCPort), logger, opts...)
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

// nolint:dupl
func TestInvalidAuthStream(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc)
		So(err, ShouldBeNil)

		cfg := &config.Config{GRPCPort: "10008"}
		verifier := test.NewVerifier("test")
		serverUnaryOpt := pkgGRPC.UnaryServerOption(logger, tokenGRPC.UnaryServerInterceptor(verifier))
		serverStreamOpt := pkgGRPC.StreamServerOption(logger, tokenGRPC.StreamServerInterceptor(verifier))
		gs := pkgGRPC.NewServer(lc, test.NewShutdowner(), cfg, logger, serverUnaryOpt, serverStreamOpt)

		test.RegisterGreeterServer(gs, test.NewServer())

		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := context.Background()
			opts := []grpc.DialOption{
				grpc.WithBlock(),
				grpc.WithInsecure(),
				pkgGRPC.UnaryDialOption(logger),
				pkgGRPC.StreamDialOption(logger),
				grpc.WithPerRPCCredentials(tokenGRPC.NewPerRPCCredentials(test.NewGenerator("bob"))),
			}

			conn, err := pkgGRPC.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", cfg.GRPCPort), logger, opts...)
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

// nolint:dupl
func TestEmptyAuthStream(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc)
		So(err, ShouldBeNil)

		cfg := &config.Config{GRPCPort: "10008"}
		verifier := test.NewVerifier("test")
		serverUnaryOpt := pkgGRPC.UnaryServerOption(logger, tokenGRPC.UnaryServerInterceptor(verifier))
		serverStreamOpt := pkgGRPC.StreamServerOption(logger, tokenGRPC.StreamServerInterceptor(verifier))
		gs := pkgGRPC.NewServer(lc, test.NewShutdowner(), cfg, logger, serverUnaryOpt, serverStreamOpt)

		test.RegisterGreeterServer(gs, test.NewServer())

		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := context.Background()
			opts := []grpc.DialOption{
				grpc.WithBlock(),
				grpc.WithInsecure(),
				pkgGRPC.UnaryDialOption(logger),
				pkgGRPC.StreamDialOption(logger),
				grpc.WithPerRPCCredentials(tokenGRPC.NewPerRPCCredentials(test.NewGenerator(""))),
			}

			conn, err := pkgGRPC.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", cfg.GRPCPort), logger, opts...)
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

func TestMissingClientAuthStream(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc)
		So(err, ShouldBeNil)

		cfg := &config.Config{GRPCPort: "10008"}
		verifier := test.NewVerifier("test")
		serverUnaryOpt := pkgGRPC.UnaryServerOption(logger, tokenGRPC.UnaryServerInterceptor(verifier))
		serverStreamOpt := pkgGRPC.StreamServerOption(logger, tokenGRPC.StreamServerInterceptor(verifier))
		gs := pkgGRPC.NewServer(lc, test.NewShutdowner(), cfg, logger, serverUnaryOpt, serverStreamOpt)

		test.RegisterGreeterServer(gs, test.NewServer())

		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := context.Background()
			opts := []grpc.DialOption{
				grpc.WithBlock(),
				grpc.WithInsecure(),
				pkgGRPC.UnaryDialOption(logger),
				pkgGRPC.StreamDialOption(logger),
			}

			conn, err := pkgGRPC.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", cfg.GRPCPort), logger, opts...)
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
