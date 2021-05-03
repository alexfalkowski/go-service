package grpc_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/alexfalkowski/go-service/pkg/config"
	pkgGRPC "github.com/alexfalkowski/go-service/pkg/grpc"
	pkgHTTP "github.com/alexfalkowski/go-service/pkg/http"
	"github.com/alexfalkowski/go-service/pkg/logger"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc"
)

func TestUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := logger.NewLogger(lc)
		So(err, ShouldBeNil)

		cfg := &config.Config{GRPCPort: "10007"}
		gs := pkgGRPC.NewServer(lc, test.NewShutdowner(), cfg, logger, pkgGRPC.NewServerOptions())

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

		logger, err := logger.NewLogger(lc)
		So(err, ShouldBeNil)

		cfg := &config.Config{GRPCPort: "10008"}
		gs := pkgGRPC.NewServer(lc, test.NewShutdowner(), cfg, logger, pkgGRPC.NewServerOptions())

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

func TestUnaryGateway(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		sh := test.NewShutdowner()
		lc := fxtest.NewLifecycle(t)

		logger, err := logger.NewLogger(lc)
		So(err, ShouldBeNil)

		cfg := &config.Config{GRPCPort: "10009", HTTPPort: "10010"}

		mux := pkgHTTP.NewMux()
		pkgHTTP.Register(lc, sh, mux, cfg, logger)

		gs := pkgGRPC.NewServer(lc, sh, cfg, logger, pkgGRPC.NewServerOptions())

		test.RegisterGreeterServer(gs, test.NewServer())

		lc.RequireStart()

		ctx := context.Background()
		opts := []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()}

		conn, err := pkgGRPC.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", cfg.GRPCPort), logger, opts...)
		So(err, ShouldBeNil)

		defer conn.Close()

		err = test.RegisterGreeterHandler(ctx, mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a greet", func() {
			client := &http.Client{Transport: pkgHTTP.NewRoundTripper(logger)}

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:10010/v1/greet/hello", bytes.NewBuffer(message))
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Request-ID", "test")

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			actual := strings.TrimSpace(string(body))

			lc.RequireStop()

			Convey("Then I should have a valid reply", func() {
				So(actual, ShouldEqual, `{"message":"Hello test"}`)
			})

			lc.RequireStop()
		})
	})
}
