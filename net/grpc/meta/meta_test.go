package meta_test

import (
	"context"
	"testing"

	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	grpcmeta "github.com/alexfalkowski/go-service/v2/net/grpc/meta"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func TestUnaryClientInterceptorReplacesOutgoingMetadata(t *testing.T) {
	ctx := context.Background()
	ctx = meta.WithUserAgent(ctx, meta.String("current-agent"))
	ctx = meta.WithRequestID(ctx, meta.String("current-id"))
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(
		"user-agent", "stale-agent",
		"request-id", "stale-id",
	))
	interceptor := grpcmeta.UnaryClientInterceptor(env.UserAgent("fallback-agent"), staticGenerator("generated-id"))

	err := interceptor(ctx, "/greet.v1.Greeter/SayHello", nil, nil, nil, func(ctx context.Context, _ string, _, _ any, _ *grpc.ClientConn, _ ...grpc.CallOption) error {
		md, ok := metadata.FromOutgoingContext(ctx)
		require.True(t, ok)
		require.Equal(t, []string{"current-agent"}, md.Get("user-agent"))
		require.Equal(t, []string{"current-id"}, md.Get("request-id"))

		return nil
	})
	require.NoError(t, err)
}

func TestStreamClientInterceptorReplacesOutgoingMetadata(t *testing.T) {
	ctx := context.Background()
	ctx = meta.WithUserAgent(ctx, meta.String("current-agent"))
	ctx = meta.WithRequestID(ctx, meta.String("current-id"))
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(
		"user-agent", "stale-agent",
		"request-id", "stale-id",
	))
	interceptor := grpcmeta.StreamClientInterceptor(env.UserAgent("fallback-agent"), staticGenerator("generated-id"))
	streamer := func(ctx context.Context, _ *grpc.StreamDesc, _ *grpc.ClientConn, _ string, _ ...grpc.CallOption) (grpc.ClientStream, error) {
		md, ok := metadata.FromOutgoingContext(ctx)
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
	interceptor := grpcmeta.UnaryServerInterceptor(env.UserAgent("fallback-agent"), env.Version("v1"), staticGenerator("generated-id"))
	ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{})

	resp, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/greet.v1.Greeter/SayHello"}, func(ctx context.Context, _ any) (any, error) {
		require.Equal(t, meta.String("peer"), meta.Attribute(ctx, meta.IPAddrKindKey))
		require.True(t, meta.IPAddr(ctx).IsEmpty())

		return "ok", nil
	})
	require.NoError(t, err)
	require.Equal(t, "ok", resp)
}

func TestUnaryServerInterceptorHandlesPeerWithoutAddr(t *testing.T) {
	interceptor := grpcmeta.UnaryServerInterceptor(env.UserAgent("fallback-agent"), env.Version("v1"), staticGenerator("generated-id"))
	ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{})
	ctx = peer.NewContext(ctx, &peer.Peer{})

	resp, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/greet.v1.Greeter/SayHello"}, func(ctx context.Context, _ any) (any, error) {
		require.Equal(t, meta.String("peer"), meta.Attribute(ctx, meta.IPAddrKindKey))
		require.True(t, meta.IPAddr(ctx).IsEmpty())

		return "ok", nil
	})
	require.NoError(t, err)
	require.Equal(t, "ok", resp)
}

type staticGenerator string

func (g staticGenerator) Generate() string {
	return string(g)
}
