package token_test

import (
	"context"
	"testing"

	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/token"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func TestUnaryClientInterceptorReplacesOutgoingAuthorization(t *testing.T) {
	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("authorization", "Bearer stale-token"))
	interceptor := token.UnaryClientInterceptor(env.UserID("service-user"), staticTokenGenerator("fresh-token"))

	err := interceptor(ctx, "/greet.v1.Greeter/SayHello", nil, nil, nil, func(ctx context.Context, _ string, _, _ any, _ *grpc.ClientConn, _ ...grpc.CallOption) error {
		md, ok := metadata.FromOutgoingContext(ctx)
		require.True(t, ok)
		require.Equal(t, []string{"Bearer fresh-token"}, md.Get("authorization"))

		return nil
	})
	require.NoError(t, err)
}

func TestStreamClientInterceptorReplacesOutgoingAuthorization(t *testing.T) {
	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("authorization", "Bearer stale-token"))
	interceptor := token.StreamClientInterceptor(env.UserID("service-user"), staticTokenGenerator("fresh-token"))
	streamer := func(ctx context.Context, _ *grpc.StreamDesc, _ *grpc.ClientConn, _ string, _ ...grpc.CallOption) (grpc.ClientStream, error) {
		md, ok := metadata.FromOutgoingContext(ctx)
		require.True(t, ok)
		require.Equal(t, []string{"Bearer fresh-token"}, md.Get("authorization"))

		return nil, nil
	}

	stream, err := interceptor(
		ctx, &grpc.StreamDesc{ServerStreams: true}, nil, "/greet.v1.Greeter/SayStreamHello", streamer,
	)
	require.NoError(t, err)
	require.Nil(t, stream)
}

type staticTokenGenerator string

func (g staticTokenGenerator) Generate(_, _ string) ([]byte, error) {
	return []byte(g), nil
}
