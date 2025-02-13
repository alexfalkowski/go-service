package grpc_test

import (
	"context"
	"net"
	"testing"

	"github.com/alexfalkowski/go-service/internal/test"
	v1 "github.com/alexfalkowski/go-service/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/time"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestInsecureUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldGRPC())
		world.Register()
		world.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := meta.WithAttribute(t.Context(), "test", meta.Ignored("test"))
			ctx = meta.WithAttribute(ctx, "ip", meta.ToRedacted(net.ParseIP("192.168.8.0")))

			conn := world.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)
			req := &v1.SayHelloRequest{Name: "test"}

			var header metadata.MD

			resp, err := client.SayHello(ctx, req, grpc.Header(&header))
			So(err, ShouldBeNil)

			h := header.Get("service-version")

			Convey("Then I should have a valid reply", func() {
				So(h, ShouldNotBeEmpty)
				So(h[0], ShouldEqual, "1.0.0")
				So(resp.GetMessage(), ShouldEqual, "Hello test")
			})

			world.RequireStop()
		})
	})
}

func TestSecureUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldGRPC())
		world.Register()
		world.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := meta.WithAttribute(t.Context(), "ip", meta.ToIgnored(net.ParseIP("192.168.8.0")))

			conn := world.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)
			req := &v1.SayHelloRequest{Name: "test"}

			resp, err := client.SayHello(ctx, req)
			So(err, ShouldBeNil)

			Convey("Then I should have a valid reply", func() {
				So(resp.GetMessage(), ShouldEqual, "Hello test")
			})

			world.RequireStop()
		})
	})
}

func TestStream(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldGRPC())
		world.Register()
		world.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := meta.WithAttribute(t.Context(), "test", meta.Redacted("test"))

			conn := world.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)

			test.Timeout()

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

			world.RequireStop()
		})
	})
}
