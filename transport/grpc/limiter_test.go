package grpc_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestServerLimiterUnary(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 0)), test.WithWorldGRPC())
	world.Register()
	world.RequireStart()

	conn := world.NewGRPC()
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	_, _ = client.SayHello(t.Context(), req)
	_, err := client.SayHello(t.Context(), req)
	require.Error(t, err)
	require.Equal(t, codes.ResourceExhausted, status.Code(err))

	world.RequireStop()
}

func TestClientLimiterUnary(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldClientLimiter(test.NewLimiterConfig("user-agent", "1s", 0)), test.WithWorldGRPC())
	world.Register()
	world.RequireStart()

	conn := world.NewGRPC()
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	_, _ = client.SayHello(t.Context(), req)
	_, err := client.SayHello(t.Context(), req)
	require.Error(t, err)
	require.Equal(t, codes.ResourceExhausted, status.Code(err))

	world.RequireStop()
}

func TestLimiterUnlimitedUnary(t *testing.T) {
	cfg := test.NewLimiterConfig("user-agent", "1s", 10)
	world := test.NewWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldClientLimiter(cfg),
		test.WithWorldServerLimiter(cfg),
		test.WithWorldGRPC(),
	)
	world.Register()
	world.RequireStart()

	conn := world.NewGRPC()
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	_, err := client.SayHello(t.Context(), req)
	require.NoError(t, err)

	world.RequireStop()
}

func TestLimiterAuthUnary(t *testing.T) {
	world := test.NewWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 10)),
		test.WithWorldToken(test.NewGenerator("bob", nil), test.NewVerifier("bob")),
		test.WithWorldGRPC(),
	)
	world.Register()
	world.RequireStart()

	conn := world.NewGRPC()
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	var err error
	for range 10 {
		_, err = client.SayHello(t.Context(), req)
	}
	require.NoError(t, err)

	world.RequireStop()
}

func TestServerClosedLimiterUnary(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 10)), test.WithWorldGRPC())
	world.Register()
	world.RequireStart()

	require.NoError(t, world.Server.GRPCLimiter.Close(t.Context()))

	conn := world.NewGRPC()
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	_, err := client.SayHello(t.Context(), req)
	require.Error(t, err)
	require.Equal(t, codes.Internal, status.Code(err))

	world.RequireStop()
}

func TestClientClosedLimiterUnary(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldClientLimiter(test.NewLimiterConfig("user-agent", "1s", 10)), test.WithWorldGRPC())
	world.Register()
	world.RequireStart()

	conn := world.NewGRPC()
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	require.NoError(t, world.Client.GRPCLimiter.Close(t.Context()))

	_, err := client.SayHello(t.Context(), req)
	require.Error(t, err)
	require.Equal(t, codes.Internal, status.Code(err))

	world.RequireStop()
}
