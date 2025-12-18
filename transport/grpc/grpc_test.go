package grpc_test

import (
	"net"
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestInsecureUnary(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldGRPC())
	world.Register()
	world.RequireStart()

	ctx := meta.WithAttribute(t.Context(), "test", meta.Ignored("test"))
	ctx = meta.WithAttribute(ctx, "real-ip", meta.ToString(net.ParseIP("192.168.8.0")))
	ctx = meta.WithAttribute(ctx, "redacted-ip", meta.ToRedacted(net.ParseIP("192.168.8.0")))

	conn := world.NewGRPC()
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}
	var header metadata.MD

	resp, err := client.SayHello(ctx, req, grpc.Header(&header))
	require.NoError(t, err)

	h := header.Get("service-version")
	require.NotEmpty(t, h)
	require.Equal(t, "1.0.0", h[0])
	require.Equal(t, "Hello test", resp.GetMessage())

	world.RequireStop()
}

func TestSecureUnary(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldGRPC(), test.WithWorldSecure())
	world.Register()
	world.RequireStart()

	ctx := meta.WithAttribute(t.Context(), "ip", meta.ToIgnored(net.ParseIP("192.168.8.0")))

	conn := world.NewGRPC()
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	resp, err := client.SayHello(ctx, req)
	require.NoError(t, err)
	require.Equal(t, "Hello test", resp.GetMessage())

	world.RequireStop()
}

func TestStream(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldGRPC())
	world.Register()
	world.RequireStart()

	ctx := meta.WithAttribute(t.Context(), "test", meta.Redacted("test"))

	conn := world.NewGRPC()
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)

	test.Timeout()

	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(10*time.Minute))
	defer cancel()

	stream, err := client.SayStreamHello(ctx)
	require.NoError(t, err)

	require.NoError(t, stream.Send(&v1.SayStreamHelloRequest{Name: "test"}))

	resp, err := stream.Recv()
	require.NoError(t, err)
	require.Equal(t, "Hello test", resp.GetMessage())

	world.RequireStop()
}
