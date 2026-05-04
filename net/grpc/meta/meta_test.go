package meta_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/net"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/meta"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/peer"
)

func TestUnaryClientInterceptorReplacesOutgoingMetadata(t *testing.T) {
	ctx := meta.WithAttributes(context.Background(),
		meta.WithUserAgent(meta.String("current-agent")),
		meta.WithRequestID(meta.String("current-id")),
	)
	ctx = meta.NewOutgoingContext(ctx, meta.Pairs(
		"user-agent", "stale-agent",
		"request-id", "stale-id",
	))
	interceptor := meta.UnaryClientInterceptor(env.UserAgent("fallback-agent"), staticGenerator("generated-id"))

	err := interceptor(ctx, "/greet.v1.Greeter/SayHello", nil, nil, nil, func(ctx context.Context, _ string, _, _ any, _ *grpc.ClientConn, _ ...grpc.CallOption) error {
		md, ok := meta.FromOutgoingContext(ctx)
		require.True(t, ok)
		require.Equal(t, []string{"current-agent"}, md.Get("user-agent"))
		require.Equal(t, []string{"current-id"}, md.Get("request-id"))

		return nil
	})
	require.NoError(t, err)
}

func TestStreamClientInterceptorReplacesOutgoingMetadata(t *testing.T) {
	ctx := meta.WithAttributes(context.Background(),
		meta.WithUserAgent(meta.String("current-agent")),
		meta.WithRequestID(meta.String("current-id")),
	)
	ctx = meta.NewOutgoingContext(ctx, meta.Pairs(
		"user-agent", "stale-agent",
		"request-id", "stale-id",
	))
	interceptor := meta.StreamClientInterceptor(env.UserAgent("fallback-agent"), staticGenerator("generated-id"))
	streamer := func(ctx context.Context, _ *grpc.StreamDesc, _ *grpc.ClientConn, _ string, _ ...grpc.CallOption) (grpc.ClientStream, error) {
		md, ok := meta.FromOutgoingContext(ctx)
		require.True(t, ok)
		require.Equal(t, []string{"current-agent"}, md.Get("user-agent"))
		require.Equal(t, []string{"current-id"}, md.Get("request-id"))

		return nil, nil
	}

	stream, err := interceptor(
		ctx, &grpc.StreamDesc{ServerStreams: true}, nil, "/greet.v1.Greeter/SayStreamHello", streamer,
	)
	require.NoError(t, err)
	require.Nil(t, stream)
}

func TestUnaryServerInterceptorHandlesMissingPeer(t *testing.T) {
	interceptor := meta.UnaryServerInterceptor(env.UserAgent("fallback-agent"), env.Version("v1"), staticGenerator("generated-id"))
	ctx := meta.NewIncomingContext(context.Background(), meta.Map{})

	resp, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/greet.v1.Greeter/SayHello"}, func(ctx context.Context, _ any) (any, error) {
		require.Equal(t, meta.String("peer"), meta.Attribute(ctx, meta.IPAddrKindKey))
		require.True(t, meta.IPAddr(ctx).IsEmpty())

		return "ok", nil
	})
	require.NoError(t, err)
	require.Equal(t, "ok", resp)
}

func TestUnaryServerInterceptorHandlesPeerWithoutAddr(t *testing.T) {
	interceptor := meta.UnaryServerInterceptor(env.UserAgent("fallback-agent"), env.Version("v1"), staticGenerator("generated-id"))
	ctx := meta.NewIncomingContext(context.Background(), meta.Map{})
	ctx = peer.NewContext(ctx, &peer.Peer{})

	resp, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/greet.v1.Greeter/SayHello"}, func(ctx context.Context, _ any) (any, error) {
		require.Equal(t, meta.String("peer"), meta.Attribute(ctx, meta.IPAddrKindKey))
		require.True(t, meta.IPAddr(ctx).IsEmpty())

		return "ok", nil
	})
	require.NoError(t, err)
	require.Equal(t, "ok", resp)
}

func TestUnaryServerInterceptorStoresPeerIPAddr(t *testing.T) {
	interceptor := meta.UnaryServerInterceptor(env.UserAgent("fallback-agent"), env.Version("v1"), staticGenerator("generated-id"))
	ctx := meta.NewIncomingContext(context.Background(), meta.Map{})
	ctx = peer.NewContext(ctx, &peer.Peer{Addr: &net.TCPAddr{IP: net.IP{127, 0, 0, 1}, Port: 8080}})

	resp, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/greet.v1.Greeter/SayHello"}, func(ctx context.Context, _ any) (any, error) {
		require.Equal(t, meta.String("peer"), meta.Attribute(ctx, meta.IPAddrKindKey))
		require.Equal(t, meta.String("127.0.0.1"), meta.IPAddr(ctx))

		return "ok", nil
	})
	require.NoError(t, err)
	require.Equal(t, "ok", resp)
}

func TestExtractIncomingReturnsMutableCopy(t *testing.T) {
	ctx := meta.NewIncomingContext(context.Background(), meta.Pairs("request-id", "original"))

	md := meta.ExtractIncoming(ctx)
	md.Set("request-id", "changed")

	original, ok := meta.FromIncomingContext(ctx)
	require.True(t, ok)
	require.Equal(t, []string{"original"}, original.Get("request-id"))
}

func TestExtractOutgoingReturnsMutableCopy(t *testing.T) {
	ctx := meta.NewOutgoingContext(context.Background(), meta.Pairs("request-id", "original"))

	md := meta.ExtractOutgoing(ctx)
	md.Set("request-id", "changed")

	original, ok := meta.FromOutgoingContext(ctx)
	require.True(t, ok)
	require.Equal(t, []string{"original"}, original.Get("request-id"))
}

type staticGenerator string

func (g staticGenerator) Generate() string {
	return string(g)
}
