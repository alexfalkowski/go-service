package health_test

import (
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/health"
	"github.com/stretchr/testify/require"
	v1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
)

func TestCheck(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldGRPC())
	world.Register()

	so := world.HealthServer(test.Name.String(), test.StatusURL("200"))

	err := so.Observe(test.Name.String(), "grpc", "http")
	require.NoError(t, err)

	server := health.NewServer(health.ServerParams{Server: so})
	health.Register(health.RegisterParams{Registrar: world.GRPCServer.ServiceRegistrar(), Server: server})

	world.RequireStart()
	time.Sleep(1 * time.Second)

	ctx := t.Context()
	ctx = meta.WithRequestID(ctx, meta.String("test-id"))
	ctx = meta.WithUserAgent(ctx, meta.String("test-user-agent"))

	conn := world.NewGRPC()
	defer conn.Close()

	client := v1.NewHealthClient(conn)
	req := &v1.HealthCheckRequest{Service: test.Name.String()}

	resp, err := client.Check(ctx, req)
	require.NoError(t, err)

	require.Equal(t, v1.HealthCheckResponse_SERVING, resp.GetStatus())

	world.RequireStop()
}

func TestInvalidCheck(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldGRPC())
	world.Register()

	so := world.HealthServer(test.Name.String(), test.StatusURL("500"))

	err := so.Observe(test.Name.String(), "grpc", "http")
	require.NoError(t, err)

	server := health.NewServer(health.ServerParams{Server: so})
	health.Register(health.RegisterParams{Registrar: world.GRPCServer.ServiceRegistrar(), Server: server})

	world.RequireStart()
	time.Sleep(1 * time.Second)

	conn := world.NewGRPC()
	defer conn.Close()

	client := v1.NewHealthClient(conn)
	req := &v1.HealthCheckRequest{Service: test.Name.String()}

	md := metadata.New(map[string]string{"request-id": "test-id", "user-agent": "test-user-agent"})
	ctx := metadata.NewOutgoingContext(t.Context(), md)

	resp, err := client.Check(ctx, req)
	require.NoError(t, err)

	require.Equal(t, v1.HealthCheckResponse_NOT_SERVING, resp.GetStatus())

	world.RequireStop()
}

func TestNotFoundCheck(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldGRPC())
	world.Register()

	so := world.HealthServer(test.Name.String(), test.StatusURL("500"))

	err := so.Observe(test.Name.String(), "grpc", "http")
	require.NoError(t, err)

	server := health.NewServer(health.ServerParams{Server: so})
	health.Register(health.RegisterParams{Registrar: world.GRPCServer.ServiceRegistrar(), Server: server})

	world.RequireStart()
	time.Sleep(1 * time.Second)

	conn := world.NewGRPC()
	defer conn.Close()

	client := v1.NewHealthClient(conn)
	req := &v1.HealthCheckRequest{Service: "bob"}

	md := metadata.New(map[string]string{"request-id": "test-id", "user-agent": "test-user-agent"})
	ctx := metadata.NewOutgoingContext(t.Context(), md)

	_, err = client.Check(ctx, req)
	require.Error(t, err)
	require.Equal(t, codes.NotFound, status.Code(err))

	world.RequireStop()
}

func TestIgnoreAuthCheck(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(nil, test.NewVerifier("test")), test.WithWorldGRPC())
	world.Register()

	so := world.HealthServer(test.Name.String(), test.StatusURL("200"))

	err := so.Observe(test.Name.String(), "grpc", "http")
	require.NoError(t, err)

	server := health.NewServer(health.ServerParams{Server: so})
	health.Register(health.RegisterParams{Registrar: world.GRPCServer.ServiceRegistrar(), Server: server})

	world.RequireStart()
	time.Sleep(1 * time.Second)

	conn := world.NewGRPC()
	defer conn.Close()

	client := v1.NewHealthClient(conn)
	req := &v1.HealthCheckRequest{Service: test.Name.String()}

	resp, err := client.Check(t.Context(), req)
	require.NoError(t, err)

	require.Equal(t, v1.HealthCheckResponse_SERVING, resp.GetStatus())

	world.RequireStop()
}

func TestList(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldGRPC())
	world.Register()

	so := world.HealthServer(test.Name.String(), test.StatusURL("200"))

	err := so.Observe(test.Name.String(), "grpc", "http")
	require.NoError(t, err)

	server := health.NewServer(health.ServerParams{Server: so})
	health.Register(health.RegisterParams{Registrar: world.GRPCServer.ServiceRegistrar(), Server: server})

	world.RequireStart()
	time.Sleep(1 * time.Second)

	ctx := t.Context()
	ctx = meta.WithRequestID(ctx, meta.String("test-id"))
	ctx = meta.WithUserAgent(ctx, meta.String("test-user-agent"))

	conn := world.NewGRPC()
	defer conn.Close()

	client := v1.NewHealthClient(conn)
	req := &v1.HealthListRequest{}

	resp, err := client.List(ctx, req)
	require.NoError(t, err)

	expected := map[string]*v1.HealthCheckResponse{
		test.Name.String(): {Status: v1.HealthCheckResponse_SERVING},
	}
	require.Equal(t, expected, resp.GetStatuses())

	world.RequireStop()
}

func TestWatch(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 10)), test.WithWorldGRPC())
	world.Register()

	so := world.HealthServer(test.Name.String(), test.StatusURL("200"))

	err := so.Observe(test.Name.String(), "grpc", "http")
	require.NoError(t, err)

	server := health.NewServer(health.ServerParams{Server: so})
	health.Register(health.RegisterParams{Registrar: world.GRPCServer.ServiceRegistrar(), Server: server})

	world.RequireStart()
	time.Sleep(1 * time.Second)

	conn := world.NewGRPC()
	defer conn.Close()

	client := v1.NewHealthClient(conn)
	req := &v1.HealthCheckRequest{Service: test.Name.String()}

	wc, err := client.Watch(t.Context(), req)
	require.NoError(t, err)

	resp, err := wc.Recv()
	require.NoError(t, err)

	require.Equal(t, v1.HealthCheckResponse_SERVING, resp.GetStatus())

	world.RequireStop()
}

func TestInvalidWatch(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldGRPC())
	world.Register()

	so := world.HealthServer(test.Name.String(), test.StatusURL("500"))

	err := so.Observe(test.Name.String(), "grpc", "http")
	require.NoError(t, err)

	server := health.NewServer(health.ServerParams{Server: so})
	health.Register(health.RegisterParams{Registrar: world.GRPCServer.ServiceRegistrar(), Server: server})

	world.RequireStart()
	time.Sleep(1 * time.Second)

	conn := world.NewGRPC()
	defer conn.Close()

	client := v1.NewHealthClient(conn)
	req := &v1.HealthCheckRequest{Service: test.Name.String()}

	wc, err := client.Watch(t.Context(), req)
	require.NoError(t, err)

	resp, err := wc.Recv()
	require.NoError(t, err)

	require.Equal(t, v1.HealthCheckResponse_NOT_SERVING, resp.GetStatus())

	world.RequireStop()
}

func TestNotFoundWatch(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldGRPC())
	world.Register()

	so := world.HealthServer(test.Name.String(), test.StatusURL("500"))

	err := so.Observe(test.Name.String(), "grpc", "http")
	require.NoError(t, err)

	server := health.NewServer(health.ServerParams{Server: so})
	health.Register(health.RegisterParams{Registrar: world.GRPCServer.ServiceRegistrar(), Server: server})

	world.RequireStart()
	time.Sleep(1 * time.Second)

	conn := world.NewGRPC()
	defer conn.Close()

	client := v1.NewHealthClient(conn)
	req := &v1.HealthCheckRequest{Service: "bob"}

	wc, err := client.Watch(t.Context(), req)
	require.NoError(t, err)

	_, err = wc.Recv()
	require.Error(t, err)
	require.Equal(t, codes.NotFound, status.Code(err))

	world.RequireStop()
}

func TestIgnoreAuthWatch(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(nil, test.NewVerifier("test")), test.WithWorldGRPC())
	world.Register()

	so := world.HealthServer(test.Name.String(), test.StatusURL("200"))

	err := so.Observe(test.Name.String(), "grpc", "http")
	require.NoError(t, err)

	server := health.NewServer(health.ServerParams{Server: so})
	health.Register(health.RegisterParams{Registrar: world.GRPCServer.ServiceRegistrar(), Server: server})

	world.RequireStart()
	time.Sleep(1 * time.Second)

	conn := world.NewGRPC()
	defer conn.Close()

	client := v1.NewHealthClient(conn)
	req := &v1.HealthCheckRequest{Service: test.Name.String()}

	wc, err := client.Watch(t.Context(), req)
	require.NoError(t, err)

	resp, err := wc.Recv()
	require.NoError(t, err)

	require.Equal(t, v1.HealthCheckResponse_SERVING, resp.GetStatus())

	world.RequireStop()
}
