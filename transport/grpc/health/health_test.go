package health_test

import (
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/health"
	. "github.com/smartystreets/goconvey/convey"
	v1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
)

func TestCheck(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldGRPC())
		world.Register()

		so := world.HealthServer(test.Name.String(), test.StatusURL("200"))

		err := so.Observe(test.Name.String(), "grpc", "http")
		So(err, ShouldBeNil)

		server := health.NewServer(health.ServerParams{Server: so})
		health.Register(health.RegisterParams{Registrar: world.GRPCServer.ServiceRegistrar(), Server: server})

		world.RequireStart()
		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			ctx := t.Context()
			ctx = meta.WithRequestID(ctx, meta.String("test-id"))
			ctx = meta.WithUserAgent(ctx, meta.String("test-user-agent"))

			conn := world.NewGRPC()
			defer conn.Close()

			client := v1.NewHealthClient(conn)
			req := &v1.HealthCheckRequest{Service: test.Name.String()}

			resp, err := client.Check(ctx, req)
			So(err, ShouldBeNil)

			Convey("Then I should have a healthy response", func() {
				So(resp.GetStatus(), ShouldEqual, v1.HealthCheckResponse_SERVING)
			})

			world.RequireStop()
		})
	})
}

func TestInvalidCheck(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldGRPC())
		world.Register()

		so := world.HealthServer(test.Name.String(), test.StatusURL("500"))

		err := so.Observe(test.Name.String(), "grpc", "http")
		So(err, ShouldBeNil)

		server := health.NewServer(health.ServerParams{Server: so})
		health.Register(health.RegisterParams{Registrar: world.GRPCServer.ServiceRegistrar(), Server: server})

		world.RequireStart()
		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			conn := world.NewGRPC()
			defer conn.Close()

			client := v1.NewHealthClient(conn)
			req := &v1.HealthCheckRequest{Service: test.Name.String()}

			md := metadata.New(map[string]string{"request-id": "test-id", "user-agent": "test-user-agent"})
			ctx := metadata.NewOutgoingContext(t.Context(), md)

			resp, err := client.Check(ctx, req)
			So(err, ShouldBeNil)

			Convey("Then I should have an unhealthy response", func() {
				So(resp.GetStatus(), ShouldEqual, v1.HealthCheckResponse_NOT_SERVING)
			})

			world.RequireStop()
		})
	})
}

func TestNotFoundCheck(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldGRPC())
		world.Register()

		so := world.HealthServer(test.Name.String(), test.StatusURL("500"))

		err := so.Observe(test.Name.String(), "grpc", "http")
		So(err, ShouldBeNil)

		server := health.NewServer(health.ServerParams{Server: so})
		health.Register(health.RegisterParams{Registrar: world.GRPCServer.ServiceRegistrar(), Server: server})

		world.RequireStart()
		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			conn := world.NewGRPC()
			defer conn.Close()

			client := v1.NewHealthClient(conn)
			req := &v1.HealthCheckRequest{Service: "bob"}

			md := metadata.New(map[string]string{"request-id": "test-id", "user-agent": "test-user-agent"})
			ctx := metadata.NewOutgoingContext(t.Context(), md)

			_, err := client.Check(ctx, req)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
				So(status.Code(err), ShouldEqual, codes.NotFound)
			})

			world.RequireStop()
		})
	})
}

func TestIgnoreAuthCheck(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(nil, test.NewVerifier("test")), test.WithWorldGRPC())
		world.Register()

		so := world.HealthServer(test.Name.String(), test.StatusURL("200"))

		err := so.Observe(test.Name.String(), "grpc", "http")
		So(err, ShouldBeNil)

		server := health.NewServer(health.ServerParams{Server: so})
		health.Register(health.RegisterParams{Registrar: world.GRPCServer.ServiceRegistrar(), Server: server})

		world.RequireStart()
		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			conn := world.NewGRPC()
			defer conn.Close()

			client := v1.NewHealthClient(conn)
			req := &v1.HealthCheckRequest{Service: test.Name.String()}

			resp, err := client.Check(t.Context(), req)
			So(err, ShouldBeNil)

			Convey("Then I should have a healthy response", func() {
				So(resp.GetStatus(), ShouldEqual, v1.HealthCheckResponse_SERVING)
			})

			world.RequireStop()
		})
	})
}

func TestList(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldGRPC())
		world.Register()

		so := world.HealthServer(test.Name.String(), test.StatusURL("200"))

		err := so.Observe(test.Name.String(), "grpc", "http")
		So(err, ShouldBeNil)

		server := health.NewServer(health.ServerParams{Server: so})
		health.Register(health.RegisterParams{Registrar: world.GRPCServer.ServiceRegistrar(), Server: server})

		world.RequireStart()
		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			ctx := t.Context()
			ctx = meta.WithRequestID(ctx, meta.String("test-id"))
			ctx = meta.WithUserAgent(ctx, meta.String("test-user-agent"))

			conn := world.NewGRPC()
			defer conn.Close()

			client := v1.NewHealthClient(conn)
			req := &v1.HealthListRequest{}

			resp, err := client.List(ctx, req)
			So(err, ShouldBeNil)

			Convey("Then I should have a healthy response", func() {
				expected := map[string]*v1.HealthCheckResponse{
					test.Name.String(): {Status: v1.HealthCheckResponse_SERVING},
				}
				So(resp.GetStatuses(), ShouldEqual, expected)
			})

			world.RequireStop()
		})
	})
}

func TestWatch(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 10)), test.WithWorldGRPC())
		world.Register()

		so := world.HealthServer(test.Name.String(), test.StatusURL("200"))

		err := so.Observe(test.Name.String(), "grpc", "http")
		So(err, ShouldBeNil)

		server := health.NewServer(health.ServerParams{Server: so})
		health.Register(health.RegisterParams{Registrar: world.GRPCServer.ServiceRegistrar(), Server: server})

		world.RequireStart()
		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			conn := world.NewGRPC()
			defer conn.Close()

			client := v1.NewHealthClient(conn)
			req := &v1.HealthCheckRequest{Service: test.Name.String()}

			wc, err := client.Watch(t.Context(), req)
			So(err, ShouldBeNil)

			resp, err := wc.Recv()
			So(err, ShouldBeNil)

			Convey("Then I should have a healthy response", func() {
				So(resp.GetStatus(), ShouldEqual, v1.HealthCheckResponse_SERVING)
			})

			world.RequireStop()
		})
	})
}

func TestInvalidWatch(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldGRPC())
		world.Register()

		so := world.HealthServer(test.Name.String(), test.StatusURL("500"))

		err := so.Observe(test.Name.String(), "grpc", "http")
		So(err, ShouldBeNil)

		server := health.NewServer(health.ServerParams{Server: so})
		health.Register(health.RegisterParams{Registrar: world.GRPCServer.ServiceRegistrar(), Server: server})

		world.RequireStart()
		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			conn := world.NewGRPC()
			defer conn.Close()

			client := v1.NewHealthClient(conn)
			req := &v1.HealthCheckRequest{Service: test.Name.String()}

			wc, err := client.Watch(t.Context(), req)
			So(err, ShouldBeNil)

			resp, err := wc.Recv()
			So(err, ShouldBeNil)

			Convey("Then I should have a healthy response", func() {
				So(resp.GetStatus(), ShouldEqual, v1.HealthCheckResponse_NOT_SERVING)
			})

			world.RequireStop()
		})
	})
}

func TestNotFoundWatch(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldGRPC())
		world.Register()

		so := world.HealthServer(test.Name.String(), test.StatusURL("500"))

		err := so.Observe(test.Name.String(), "grpc", "http")
		So(err, ShouldBeNil)

		server := health.NewServer(health.ServerParams{Server: so})
		health.Register(health.RegisterParams{Registrar: world.GRPCServer.ServiceRegistrar(), Server: server})

		world.RequireStart()
		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			conn := world.NewGRPC()
			defer conn.Close()

			client := v1.NewHealthClient(conn)
			req := &v1.HealthCheckRequest{Service: "bob"}

			wc, err := client.Watch(t.Context(), req)
			So(err, ShouldBeNil)

			_, err = wc.Recv()

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
				So(status.Code(err), ShouldEqual, codes.NotFound)
			})

			world.RequireStop()
		})
	})
}

func TestIgnoreAuthWatch(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(nil, test.NewVerifier("test")), test.WithWorldGRPC())
		world.Register()

		so := world.HealthServer(test.Name.String(), test.StatusURL("200"))

		err := so.Observe(test.Name.String(), "grpc", "http")
		So(err, ShouldBeNil)

		server := health.NewServer(health.ServerParams{Server: so})
		health.Register(health.RegisterParams{Registrar: world.GRPCServer.ServiceRegistrar(), Server: server})

		world.RequireStart()
		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			conn := world.NewGRPC()
			defer conn.Close()

			client := v1.NewHealthClient(conn)
			req := &v1.HealthCheckRequest{Service: test.Name.String()}

			wc, err := client.Watch(t.Context(), req)
			So(err, ShouldBeNil)

			resp, err := wc.Recv()
			So(err, ShouldBeNil)

			Convey("Then I should have a healthy response", func() {
				So(resp.GetStatus(), ShouldEqual, v1.HealthCheckResponse_SERVING)
			})

			world.RequireStop()
		})
	})
}
