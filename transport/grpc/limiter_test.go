package grpc_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestServerLimiterUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 0)), test.WithWorldGRPC())
		world.Register()
		world.RequireStart()

		Convey("When I query repeatedly", func() {
			conn := world.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)
			req := &v1.SayHelloRequest{Name: "test"}

			_, _ = client.SayHello(t.Context(), req)
			_, err := client.SayHello(t.Context(), req)

			Convey("Then I should have exhausted resources", func() {
				So(err, ShouldBeError)
				So(status.Code(err), ShouldEqual, codes.ResourceExhausted)
			})
		})

		world.RequireStop()
	})
}

func TestClientLimiterUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldClientLimiter(test.NewLimiterConfig("user-agent", "1s", 0)), test.WithWorldGRPC())
		world.Register()
		world.RequireStart()

		Convey("When I query repeatedly", func() {
			conn := world.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)
			req := &v1.SayHelloRequest{Name: "test"}

			_, _ = client.SayHello(t.Context(), req)
			_, err := client.SayHello(t.Context(), req)

			Convey("Then I should have exhausted resources", func() {
				So(err, ShouldBeError)
				So(status.Code(err), ShouldEqual, codes.ResourceExhausted)
			})
		})

		world.RequireStop()
	})
}

func TestLimiterUnlimitedUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		cfg := test.NewLimiterConfig("user-agent", "1s", 10)
		world := test.NewWorld(t,
			test.WithWorldTelemetry("otlp"),
			test.WithWorldClientLimiter(cfg),
			test.WithWorldServerLimiter(cfg),
			test.WithWorldGRPC(),
		)
		world.Register()
		world.RequireStart()

		Convey("When I query repeatedly", func() {
			conn := world.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)
			req := &v1.SayHelloRequest{Name: "test"}

			_, err := client.SayHello(t.Context(), req)

			Convey("Then I should not have exhausted resources", func() {
				So(err, ShouldBeNil)
			})
		})

		world.RequireStop()
	})
}

func TestLimiterAuthUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		world := test.NewWorld(t,
			test.WithWorldTelemetry("otlp"),
			test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 10)),
			test.WithWorldToken(test.NewGenerator("bob", nil), test.NewVerifier("bob")),
			test.WithWorldGRPC(),
		)
		world.Register()
		world.RequireStart()

		Convey("When I query for a authenticated greet multiple times", func() {
			conn := world.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)
			req := &v1.SayHelloRequest{Name: "test"}

			var err error
			for range 10 {
				_, err = client.SayHello(t.Context(), req)
			}

			Convey("Then I should not have exhausted resources", func() {
				So(err, ShouldBeNil)
			})
		})

		world.RequireStop()
	})
}

func TestServerClosedLimiterUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 10)), test.WithWorldGRPC())
		world.Register()
		world.RequireStart()

		err := world.Server.GRPCLimiter.Close(t.Context())
		So(err, ShouldBeNil)

		Convey("When  I query for a greet", func() {
			conn := world.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)
			req := &v1.SayHelloRequest{Name: "test"}

			_, err := client.SayHello(t.Context(), req)

			Convey("Then I should have an internal error", func() {
				So(err, ShouldBeError)
				So(status.Code(err), ShouldEqual, codes.Internal)
			})
		})

		world.RequireStop()
	})
}

func TestClientClosedLimiterUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldClientLimiter(test.NewLimiterConfig("user-agent", "1s", 10)), test.WithWorldGRPC())
		world.Register()
		world.RequireStart()

		Convey("When  I query for a greet", func() {
			conn := world.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)
			req := &v1.SayHelloRequest{Name: "test"}

			err := world.Client.GRPCLimiter.Close(t.Context())
			So(err, ShouldBeNil)

			_, err = client.SayHello(t.Context(), req)

			Convey("Then I should have an internal error", func() {
				So(err, ShouldBeError)
				So(status.Code(err), ShouldEqual, codes.Internal)
			})
		})

		world.RequireStop()
	})
}
