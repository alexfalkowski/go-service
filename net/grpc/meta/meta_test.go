package meta_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	grpcmeta "github.com/alexfalkowski/go-service/v2/net/grpc/meta"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/peer"
)

func TestUnaryClientInterceptorReplacesOutgoingMetadata(t *testing.T) {
	ctx := grpcmeta.WithAttributes(t.Context(),
		grpcmeta.WithUserAgent(grpcmeta.String("current-agent")),
		grpcmeta.WithRequestID(grpcmeta.String("current-id")),
	)
	ctx = grpcmeta.NewOutgoingContext(ctx, grpcmeta.Pairs(
		"user-agent", "stale-agent",
		"request-id", "stale-id",
	))
	interceptor := grpcmeta.UnaryClientInterceptor(env.UserAgent("fallback-agent"), test.StaticIDGenerator("generated-id"))

	err := interceptor(ctx, "/greet.v1.Greeter/SayHello", nil, nil, nil, func(ctx context.Context, _ string, _, _ any, _ *grpc.ClientConn, _ ...grpc.CallOption) error {
		md, ok := grpcmeta.FromOutgoingContext(ctx)
		require.True(t, ok)
		require.Equal(t, []string{"current-agent"}, md.Get("user-agent"))
		require.Equal(t, []string{"current-id"}, md.Get("request-id"))

		return nil
	})
	require.NoError(t, err)
}

func TestUnaryClientInterceptorIgnoresBlankOutgoingMetadata(t *testing.T) {
	ctx := grpcmeta.NewOutgoingContext(t.Context(), grpcmeta.Pairs(
		"user-agent", "",
		"request-id", "",
	))
	interceptor := grpcmeta.UnaryClientInterceptor(env.UserAgent("fallback-agent"), test.StaticIDGenerator("generated-id"))

	err := interceptor(ctx, "/greet.v1.Greeter/SayHello", nil, nil, nil, func(ctx context.Context, _ string, _, _ any, _ *grpc.ClientConn, _ ...grpc.CallOption) error {
		md, ok := grpcmeta.FromOutgoingContext(ctx)
		require.True(t, ok)
		require.Equal(t, []string{"fallback-agent"}, md.Get("user-agent"))
		require.Equal(t, []string{"generated-id"}, md.Get("request-id"))
		require.Equal(t, grpcmeta.String("fallback-agent"), grpcmeta.UserAgent(ctx))
		require.Equal(t, grpcmeta.String("generated-id"), meta.RequestID(ctx))

		return nil
	})
	require.NoError(t, err)
}

func TestStreamClientInterceptorReplacesOutgoingMetadata(t *testing.T) {
	ctx := grpcmeta.WithAttributes(t.Context(),
		grpcmeta.WithUserAgent(grpcmeta.String("current-agent")),
		grpcmeta.WithRequestID(grpcmeta.String("current-id")),
	)
	ctx = grpcmeta.NewOutgoingContext(ctx, grpcmeta.Pairs(
		"user-agent", "stale-agent",
		"request-id", "stale-id",
	))
	interceptor := grpcmeta.StreamClientInterceptor(env.UserAgent("fallback-agent"), test.StaticIDGenerator("generated-id"))
	streamer := func(ctx context.Context, _ *grpc.StreamDesc, _ *grpc.ClientConn, _ string, _ ...grpc.CallOption) (grpc.ClientStream, error) {
		md, ok := grpcmeta.FromOutgoingContext(ctx)
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

func TestStreamClientInterceptorIgnoresBlankOutgoingMetadata(t *testing.T) {
	ctx := grpcmeta.NewOutgoingContext(t.Context(), grpcmeta.Pairs(
		"user-agent", "",
		"request-id", "",
	))
	interceptor := grpcmeta.StreamClientInterceptor(env.UserAgent("fallback-agent"), test.StaticIDGenerator("generated-id"))
	streamer := func(ctx context.Context, _ *grpc.StreamDesc, _ *grpc.ClientConn, _ string, _ ...grpc.CallOption) (grpc.ClientStream, error) {
		md, ok := grpcmeta.FromOutgoingContext(ctx)
		require.True(t, ok)
		require.Equal(t, []string{"fallback-agent"}, md.Get("user-agent"))
		require.Equal(t, []string{"generated-id"}, md.Get("request-id"))
		require.Equal(t, grpcmeta.String("fallback-agent"), grpcmeta.UserAgent(ctx))
		require.Equal(t, grpcmeta.String("generated-id"), meta.RequestID(ctx))

		return nil, nil
	}

	stream, err := interceptor(
		ctx, &grpc.StreamDesc{ServerStreams: true}, nil, "/greet.v1.Greeter/SayStreamHello", streamer,
	)
	require.NoError(t, err)
	require.Nil(t, stream)
}

func TestUnaryServerInterceptorHandlesMissingPeer(t *testing.T) {
	interceptor := grpcmeta.UnaryServerInterceptor(env.UserAgent("fallback-agent"), env.Version("v1"), test.StaticIDGenerator("generated-id"))
	ctx := grpcmeta.NewIncomingContext(t.Context(), grpcmeta.Map{})

	resp, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/greet.v1.Greeter/SayHello"}, func(ctx context.Context, _ any) (any, error) {
		require.Equal(t, grpcmeta.String("peer"), grpcmeta.Attribute(ctx, grpcmeta.IPAddrKindKey))
		require.True(t, grpcmeta.IPAddr(ctx).IsEmpty())

		return "ok", nil
	})
	require.NoError(t, err)
	require.Equal(t, "ok", resp)
}

func TestUnaryServerInterceptorHandlesPeerWithoutAddr(t *testing.T) {
	interceptor := grpcmeta.UnaryServerInterceptor(env.UserAgent("fallback-agent"), env.Version("v1"), test.StaticIDGenerator("generated-id"))
	ctx := grpcmeta.NewIncomingContext(t.Context(), grpcmeta.Map{})
	ctx = peer.NewContext(ctx, &peer.Peer{})

	resp, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/greet.v1.Greeter/SayHello"}, func(ctx context.Context, _ any) (any, error) {
		require.Equal(t, grpcmeta.String("peer"), grpcmeta.Attribute(ctx, grpcmeta.IPAddrKindKey))
		require.True(t, grpcmeta.IPAddr(ctx).IsEmpty())

		return "ok", nil
	})
	require.NoError(t, err)
	require.Equal(t, "ok", resp)
}

func TestUnaryServerInterceptorStoresPeerIPAddr(t *testing.T) {
	interceptor := grpcmeta.UnaryServerInterceptor(env.UserAgent("fallback-agent"), env.Version("v1"), test.StaticIDGenerator("generated-id"))
	ctx := grpcmeta.NewIncomingContext(t.Context(), grpcmeta.Map{})
	ctx = peer.NewContext(ctx, &peer.Peer{Addr: &net.TCPAddr{IP: net.IP{127, 0, 0, 1}, Port: 8080}})

	resp, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/greet.v1.Greeter/SayHello"}, func(ctx context.Context, _ any) (any, error) {
		require.Equal(t, grpcmeta.String("peer"), grpcmeta.Attribute(ctx, grpcmeta.IPAddrKindKey))
		require.Equal(t, grpcmeta.String("127.0.0.1"), grpcmeta.IPAddr(ctx))

		return "ok", nil
	})
	require.NoError(t, err)
	require.Equal(t, "ok", resp)
}

func TestUnaryServerInterceptorStoresGeolocationAsIgnored(t *testing.T) {
	interceptor := grpcmeta.UnaryServerInterceptor(env.UserAgent("fallback-agent"), env.Version("v1"), test.StaticIDGenerator("generated-id"))
	ctx := grpcmeta.NewIncomingContext(t.Context(), grpcmeta.Pairs("geolocation", "geo:47,11"))

	resp, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/greet.v1.Greeter/SayHello"}, func(ctx context.Context, _ any) (any, error) {
		geolocation := meta.Geolocation(ctx)

		require.Equal(t, "geo:47,11", geolocation.Value())
		require.Empty(t, geolocation.String())
		require.NotContains(t, meta.CamelStrings(ctx, meta.NoPrefix), meta.GeolocationKey)

		return "ok", nil
	})
	require.NoError(t, err)
	require.Equal(t, "ok", resp)
}

func TestStreamServerInterceptorAppendDoesNotOverwriteRequestID(t *testing.T) {
	interceptor := grpcmeta.StreamServerInterceptor(env.UserAgent("fallback-agent"), env.Version("v1"), test.StaticIDGenerator("generated-id"))
	ctx := grpcmeta.NewIncomingContext(t.Context(), grpcmeta.Map{})
	stream := &test.MetaServerStream{Ctx: ctx}

	err := interceptor(nil, stream, &grpc.StreamServerInfo{FullMethod: "/greet.v1.Greeter/SayStreamHello"}, func(any, grpc.ServerStream) error {
		return nil
	})
	require.NoError(t, err)
	require.Equal(t, []string{"1", "v2"}, stream.Header.Get("service-version"))
	require.Equal(t, []string{"generated-id"}, stream.Header.Get("request-id"))
}

func TestStreamServerInterceptorExtractsOperationMetadata(t *testing.T) {
	interceptor := grpcmeta.StreamServerInterceptor(env.UserAgent("fallback-agent"), env.Version("v1"), test.StaticIDGenerator("generated-id"))
	ctx := grpcmeta.NewIncomingContext(t.Context(), grpcmeta.Pairs(
		"authorization", "invalid",
		"user-agent", "watch-agent",
	))
	stream := &test.MetaServerStream{Ctx: ctx}

	err := interceptor(nil, stream, &grpc.StreamServerInfo{FullMethod: "/grpc.health.v1.Health/Watch"}, func(_ any, stream grpc.ServerStream) error {
		require.Equal(t, grpcmeta.String("watch-agent"), grpcmeta.UserAgent(stream.Context()))

		return nil
	})
	require.NoError(t, err)
}

func TestExtractIncomingReturnsMutableCopy(t *testing.T) {
	ctx := grpcmeta.NewIncomingContext(t.Context(), grpcmeta.Pairs("request-id", "original"))

	md := grpcmeta.ExtractIncoming(ctx)
	md.Set("request-id", "changed")

	original, ok := grpcmeta.FromIncomingContext(ctx)
	require.True(t, ok)
	require.Equal(t, []string{"original"}, original.Get("request-id"))
}

func TestExtractIncomingReturnsEmptyMapWithoutMetadata(t *testing.T) {
	md := grpcmeta.ExtractIncoming(t.Context())

	require.NotNil(t, md)
	require.Empty(t, md)
}

func TestExtractOutgoingReturnsMutableCopy(t *testing.T) {
	ctx := grpcmeta.NewOutgoingContext(t.Context(), grpcmeta.Pairs("request-id", "original"))

	md := grpcmeta.ExtractOutgoing(ctx)
	md.Set("request-id", "changed")

	original, ok := grpcmeta.FromOutgoingContext(ctx)
	require.True(t, ok)
	require.Equal(t, []string{"original"}, original.Get("request-id"))
}

func TestExtractOutgoingReturnsEmptyMapWithoutMetadata(t *testing.T) {
	md := grpcmeta.ExtractOutgoing(t.Context())

	require.NotNil(t, md)
	require.Empty(t, md)
}
