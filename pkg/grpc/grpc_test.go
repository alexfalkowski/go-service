package grpc_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/alexfalkowski/go-service/pkg/config"
	pkgGRPC "github.com/alexfalkowski/go-service/pkg/grpc"
	"github.com/alexfalkowski/go-service/pkg/logger/zap"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc"
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
