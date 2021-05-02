package grpc_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/alexfalkowski/go-service/pkg/config"
	pkgGRPC "github.com/alexfalkowski/go-service/pkg/grpc"
	"github.com/alexfalkowski/go-service/pkg/grpc/internal"
	"github.com/alexfalkowski/go-service/pkg/logger"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc"
)

type shutdowner struct{}

func (*shutdowner) Shutdown(...fx.ShutdownOption) error {
	return nil
}

type server struct {
	internal.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *internal.HelloRequest) (*internal.HelloReply, error) {
	return &internal.HelloReply{Message: fmt.Sprintf("Hello %s", req.GetName())}, nil
}
func (s *server) SayStreamHello(stream internal.Greeter_SayStreamHelloServer) error {
	req, err := stream.Recv()
	if err != nil {
		return err
	}

	return stream.Send(&internal.HelloReply{Message: fmt.Sprintf("Hello %s", req.GetName())})
}

func TestUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := logger.NewLogger(lc)
		So(err, ShouldBeNil)

		cfg := &config.Config{GRPCPort: "10007"}
		gs := pkgGRPC.NewServer(lc, &shutdowner{}, cfg, logger, pkgGRPC.NewServerOptions())

		internal.RegisterGreeterServer(gs, &server{})

		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := context.Background()
			opts := []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()}

			conn, err := pkgGRPC.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", cfg.GRPCPort), logger, opts...)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := internal.NewGreeterClient(conn)
			req := &internal.HelloRequest{Name: "test"}

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

		logger, err := logger.NewLogger(lc)
		So(err, ShouldBeNil)

		cfg := &config.Config{GRPCPort: "10006"}

		gs := pkgGRPC.NewServer(lc, &shutdowner{}, cfg, logger, pkgGRPC.NewServerOptions())
		defer gs.GracefulStop()

		internal.RegisterGreeterServer(gs, &server{})

		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := context.Background()
			opts := []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()}

			conn, err := pkgGRPC.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", cfg.GRPCPort), logger, opts...)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := internal.NewGreeterClient(conn)

			stream, err := client.SayStreamHello(ctx)
			So(err, ShouldBeNil)

			err = stream.Send(&internal.HelloRequest{Name: "test"})
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
