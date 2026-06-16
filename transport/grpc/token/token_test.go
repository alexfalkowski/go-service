package token_test

import (
	"context"
	"testing"

	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/meta"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/token"
	"github.com/stretchr/testify/require"
)

func TestUnaryClientInterceptorReplacesOutgoingAuthorization(t *testing.T) {
	ctx := meta.NewOutgoingContext(context.Background(), meta.Pairs("authorization", "Bearer stale-token"))
	interceptor := token.UnaryClientInterceptor(env.UserID("service-user"), staticTokenGenerator("fresh-token"))

	err := interceptor(ctx, "/greet.v1.Greeter/SayHello", nil, nil, nil, func(ctx context.Context, _ string, _, _ any, _ *grpc.ClientConn, _ ...grpc.CallOption) error {
		md, ok := meta.FromOutgoingContext(ctx)
		require.True(t, ok)
		require.Equal(t, []string{"Bearer fresh-token"}, md.Get("authorization"))

		return nil
	})
	require.NoError(t, err)
}

func TestStreamClientInterceptorReplacesOutgoingAuthorization(t *testing.T) {
	ctx := meta.NewOutgoingContext(context.Background(), meta.Pairs("authorization", "Bearer stale-token"))
	interceptor := token.StreamClientInterceptor(env.UserID("service-user"), staticTokenGenerator("fresh-token"))
	streamer := func(ctx context.Context, _ *grpc.StreamDesc, _ *grpc.ClientConn, _ string, _ ...grpc.CallOption) (grpc.ClientStream, error) {
		md, ok := meta.FromOutgoingContext(ctx)
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

func TestUnaryAccessServerInterceptor(t *testing.T) {
	for _, tt := range accessServerTests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			if tt.user != "" {
				ctx = meta.WithAttributes(ctx, meta.WithUserID(meta.String(tt.user)))
			}
			interceptor := token.UnaryAccessServerInterceptor(tt.controller)
			called := false

			_, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: tt.method}, func(context.Context, any) (any, error) {
				called = true
				return nil, nil
			})

			require.Equal(t, tt.code, status.Code(err))
			require.Equal(t, tt.called, called)
		})
	}
}

func TestStreamAccessServerInterceptor(t *testing.T) {
	for _, tt := range accessServerTests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			if tt.user != "" {
				ctx = meta.WithAttributes(ctx, meta.WithUserID(meta.String(tt.user)))
			}
			interceptor := token.StreamAccessServerInterceptor(tt.controller)
			stream := &test.MetaServerStream{Ctx: ctx}
			called := false

			err := interceptor(nil, stream, &grpc.StreamServerInfo{FullMethod: tt.method}, func(any, grpc.ServerStream) error {
				called = true
				return nil
			})

			require.Equal(t, tt.code, status.Code(err))
			require.Equal(t, tt.called, called)
		})
	}
}

type staticTokenGenerator string

func (g staticTokenGenerator) Generate(_, _ string) ([]byte, error) {
	return []byte(g), nil
}

type accessServerTest struct {
	controller accessControllerFunc
	method     string
	name       string
	user       string
	code       codes.Code
	called     bool
}

var accessServerTests = []accessServerTest{
	{
		name:       "operation method bypasses access",
		method:     "/grpc.health.v1.Health/Check",
		code:       codes.OK,
		controller: accessControllerFunc(func(context.Context) (bool, error) { return false, test.ErrInvalid }),
		called:     true,
	},
	{
		name:       "missing user id is unauthenticated",
		method:     "/greet.v1.GreeterService/SayHello",
		code:       codes.Unauthenticated,
		controller: accessControllerFunc(func(context.Context) (bool, error) { return true, nil }),
	},
	{
		name:       "controller error is internal",
		method:     "/greet.v1.GreeterService/SayHello",
		code:       codes.Internal,
		user:       "frontend",
		controller: accessControllerFunc(func(context.Context) (bool, error) { return false, test.ErrInvalid }),
	},
	{
		name:       "access denial is permission denied",
		method:     "/greet.v1.GreeterService/SayHello",
		code:       codes.PermissionDenied,
		user:       "frontend",
		controller: accessControllerFunc(func(context.Context) (bool, error) { return false, nil }),
	},
	{
		name:       "access grant calls handler",
		method:     "/greet.v1.GreeterService/SayHello",
		code:       codes.OK,
		user:       "frontend",
		controller: accessControllerFunc(func(context.Context) (bool, error) { return true, nil }),
		called:     true,
	},
}

type accessControllerFunc func(context.Context) (bool, error)

func (f accessControllerFunc) HasAccess(ctx context.Context) (bool, error) {
	return f(ctx)
}
