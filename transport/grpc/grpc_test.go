package grpc_test

import (
	"maps"
	"net"
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/meta"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

func TestInsecureUnary(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldGRPC())

	ctx := meta.WithAttributes(t.Context(),
		test.WithTest(meta.Ignored("test")),
		meta.NewPair("real-ip", meta.ToString(net.ParseIP("192.168.8.0"))),
		meta.NewPair("redacted-ip", meta.ToRedacted(net.ParseIP("192.168.8.0"))),
	)

	conn := test.RequireGRPCConn(t, world)
	t.Cleanup(func() {
		require.NoError(t, conn.Close())
	})

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}
	var header meta.Map

	resp, err := client.SayHello(ctx, req, grpc.Header(&header))
	require.NoError(t, err)

	h := header.Get("service-version")
	require.NotEmpty(t, h)
	require.Equal(t, "1.0.0", h[0])
	require.Equal(t, "Hello test", resp.GetMessage())
}

func TestCompressionUnary(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldGRPC(), test.WithWorldCompression())

	conn := test.RequireGRPCConn(t, world)
	t.Cleanup(func() {
		require.NoError(t, conn.Close())
	})

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	resp, err := client.SayHello(t.Context(), req)
	require.NoError(t, err)
	require.Equal(t, "Hello test", resp.GetMessage())
}

func TestSecureUnary(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldGRPC(), test.WithWorldSecure())

	ctx := meta.WithAttributes(t.Context(), meta.NewPair("ip", meta.ToIgnored(net.ParseIP("192.168.8.0"))))

	conn := test.RequireGRPCConn(t, world)
	t.Cleanup(func() {
		require.NoError(t, conn.Close())
	})

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	resp, err := client.SayHello(ctx, req)
	require.NoError(t, err)
	require.Equal(t, "Hello test", resp.GetMessage())
}

func TestStream(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldGRPC())

	ctx := meta.WithAttributes(t.Context(), test.WithTest(meta.Redacted("test")))

	conn := test.RequireGRPCConn(t, world)
	t.Cleanup(func() {
		require.NoError(t, conn.Close())
	})

	client := v1.NewGreeterServiceClient(conn)

	ctx, cancel := test.Timeout(ctx)
	defer cancel()

	stream, err := client.SayStreamHello(ctx)
	require.NoError(t, err)

	resp, err := test.SendStreamHello(t, stream, "test")
	require.NoError(t, err)
	require.Equal(t, "Hello test", resp.GetMessage())
}

func TestServerRecoversUnaryPanic(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldGRPC())
	conn := test.RequireGRPCConn(t, world)
	t.Cleanup(func() {
		require.NoError(t, conn.Close())
	})

	client := v1.NewGreeterServiceClient(conn)

	_, err := client.SayHello(t.Context(), &v1.SayHelloRequest{Name: "panic"})
	require.Error(t, err)
	require.Equal(t, codes.Internal, status.Code(err))
	assertSafePanicStatus(t, err)

	resp, err := client.SayHello(t.Context(), &v1.SayHelloRequest{Name: "test"})
	require.NoError(t, err)
	require.Equal(t, "Hello test", resp.GetMessage())
}

func TestBreakerUnary(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldGRPC(), test.WithWorldBreaker(test.NewBreaker(1)))

	conn := test.RequireGRPCConn(t, world)
	t.Cleanup(func() {
		require.NoError(t, conn.Close())
	})

	client := v1.NewGreeterServiceClient(conn)

	_, err := client.SayHello(t.Context(), &v1.SayHelloRequest{Name: "panic"})
	require.Equal(t, codes.Internal, status.Code(err))

	_, err = client.SayHello(t.Context(), &v1.SayHelloRequest{Name: "test"})
	require.True(t, status.IsLocalError(err))
	require.Equal(t, codes.ResourceExhausted, status.Code(err))
}

func TestServerRecoversStreamPanic(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldGRPC())
	conn := test.RequireGRPCConn(t, world)
	t.Cleanup(func() {
		require.NoError(t, conn.Close())
	})

	client := v1.NewGreeterServiceClient(conn)

	stream, err := client.SayStreamHello(t.Context())
	require.NoError(t, err)

	_, err = test.SendStreamHello(t, stream, "panic")
	require.Error(t, err)
	require.Equal(t, codes.Internal, status.Code(err))
	assertSafePanicStatus(t, err)

	stream, err = client.SayStreamHello(t.Context())
	require.NoError(t, err)
	resp, err := test.SendStreamHello(t, stream, "test")
	require.NoError(t, err)
	require.Equal(t, "Hello test", resp.GetMessage())
}

func TestUnaryMaxReceiveSize(t *testing.T) {
	world := newStartedGRPCWorld(t, 64)
	conn := test.RequireGRPCConn(t, world)
	t.Cleanup(func() {
		require.NoError(t, conn.Close())
	})

	client := v1.NewGreeterServiceClient(conn)

	_, err := client.SayHello(t.Context(), &v1.SayHelloRequest{Name: strings.Repeat("a", 256)})
	require.Error(t, err)
	require.Equal(t, codes.ResourceExhausted, status.Code(err))
}

func TestStreamMaxReceiveSize(t *testing.T) {
	world := newStartedGRPCWorld(t, 64)
	conn := test.RequireGRPCConn(t, world)
	t.Cleanup(func() {
		require.NoError(t, conn.Close())
	})

	client := v1.NewGreeterServiceClient(conn)
	stream, err := client.SayStreamHello(t.Context())
	require.NoError(t, err)

	_, err = test.SendStreamHello(t, stream, strings.Repeat("a", 256))
	require.Error(t, err)
	require.Equal(t, codes.ResourceExhausted, status.Code(err))
}

func TestUnaryMaxSendSize(t *testing.T) {
	world := newStartedGRPCWorldWithOptions(t, 0, map[string]string{"max_send_msg_size": "64B"})
	conn := test.RequireGRPCConn(t, world)
	t.Cleanup(func() {
		require.NoError(t, conn.Close())
	})

	client := v1.NewGreeterServiceClient(conn)

	_, err := client.SayHello(t.Context(), &v1.SayHelloRequest{Name: strings.Repeat("a", 256)})
	require.Error(t, err)
	require.Equal(t, codes.ResourceExhausted, status.Code(err))
}

func TestStreamMaxSendSize(t *testing.T) {
	world := newStartedGRPCWorldWithOptions(t, 0, map[string]string{"max_send_msg_size": "64B"})
	conn := test.RequireGRPCConn(t, world)
	t.Cleanup(func() {
		require.NoError(t, conn.Close())
	})

	client := v1.NewGreeterServiceClient(conn)
	stream, err := client.SayStreamHello(t.Context())
	require.NoError(t, err)

	_, err = test.SendStreamHello(t, stream, strings.Repeat("a", 256))
	require.Error(t, err)
	require.Equal(t, codes.ResourceExhausted, status.Code(err))
}

func newStartedGRPCWorld(t *testing.T, maxReceiveSize bytes.Size) *test.World {
	t.Helper()

	return newStartedGRPCWorldWithOptions(t, maxReceiveSize, nil)
}

func newStartedGRPCWorldWithOptions(t *testing.T, maxReceiveSize bytes.Size, opts map[string]string) *test.World {
	t.Helper()

	cfg := test.NewInsecureTransportConfig()
	if maxReceiveSize > 0 {
		cfg.GRPC.MaxReceiveSize = maxReceiveSize
	}
	if len(opts) > 0 {
		if cfg.GRPC.Options == nil {
			cfg.GRPC.Options = map[string]string{}
		}
		maps.Copy(cfg.GRPC.Options, opts)
	}

	return test.NewStartedWorld(t, test.WithWorldTransportConfig(cfg), test.WithWorldGRPC())
}

func assertSafePanicStatus(t *testing.T, err error) {
	t.Helper()

	stat, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, "grpc: internal", stat.Message())
	require.NotContains(t, stat.Message(), "test panic")
	require.NotContains(t, stat.Message(), "recovered")
}
