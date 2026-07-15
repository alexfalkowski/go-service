package grpc_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	tls "github.com/alexfalkowski/go-service/v2/crypto/tls/config"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/v2/net"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/transport/breaker"
	transportgrpc "github.com/alexfalkowski/go-service/v2/transport/grpc"
	grpcbreaker "github.com/alexfalkowski/go-service/v2/transport/grpc/breaker"
	grpclimiter "github.com/alexfalkowski/go-service/v2/transport/grpc/limiter"
	grpcretry "github.com/alexfalkowski/go-service/v2/transport/grpc/retry"
	"github.com/alexfalkowski/go-service/v2/transport/limiter"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	transportgrpc.Register(test.FS)

	_, err := transportgrpc.NewClient("none", transportgrpc.WithClientTLS(&tls.Config{Cert: "bob", Key: "bob"}))
	require.Error(t, err)

	_, err = transportgrpc.NewClient("none", transportgrpc.WithClientTLS(&tls.Config{}))
	require.NoError(t, err)
}

func TestClientEmptyTLSConfigUsesTLS(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldGRPC())
	_, target := net.ListenNetworkAddress(world.TransportConfig.GRPC.Address)

	conn, err := transportgrpc.NewClient(target, transportgrpc.WithClientTLS(&tls.Config{}))
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, conn.Close())
	})

	client := v1.NewGreeterServiceClient(conn)
	_, err = client.SayHello(t.Context(), &v1.SayHelloRequest{Name: "test"})
	require.Error(t, err)
}

func TestRetryDoesNotReenterClientLimiter(t *testing.T) {
	clientLimiter, err := grpclimiter.NewClientLimiter(
		test.NoopLifecycle{},
		limiter.NewKeyMap(),
		test.NewLimiterConfig("user-agent", "1m", 0),
	)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, clientLimiter.Close(t.Context())) })

	limiting := grpclimiter.UnaryClientInterceptor(clientLimiter)
	err = limiting(t.Context(), "/test.Service/GetBook", nil, nil, nil, func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
		return nil
	})
	require.NoError(t, err)

	retrying := grpcretry.UnaryClientInterceptor(test.NewGRPCRetryConfig(3, time.Nanosecond, codes.ResourceExhausted))
	limiterCalls := 0
	downstreamCalls := 0
	invoker := func(ctx context.Context, fullMethod string, req, resp any, conn *grpc.ClientConn, opts ...grpc.CallOption) error {
		limiterCalls++
		return limiting(ctx, fullMethod, req, resp, conn, func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
			downstreamCalls++
			return nil
		}, opts...)
	}
	err = retrying(t.Context(), "/test.Service/GetBook", nil, nil, nil, invoker)

	require.Error(t, err)
	require.True(t, status.IsLocalError(err))
	require.Equal(t, codes.ResourceExhausted, status.Code(err))
	require.Equal(t, 1, limiterCalls)
	require.Equal(t, 0, downstreamCalls)
}

func TestRetryDoesNotReenterClientBreaker(t *testing.T) {
	breaking := grpcbreaker.UnaryClientInterceptor(
		grpcbreaker.NewConfig(test.NewBreaker(1), codes.Unavailable).Options()...,
	)
	err := breaking(t.Context(), "/test.Service/GetBook", nil, nil, nil, func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
		return status.Error(codes.Unavailable, "unavailable")
	})
	require.Error(t, err)
	require.Equal(t, codes.Unavailable, status.Code(err))

	retrying := grpcretry.UnaryClientInterceptor(test.NewGRPCRetryConfig(3, time.Nanosecond, codes.ResourceExhausted))
	breakerCalls := 0
	downstreamCalls := 0
	invoker := func(ctx context.Context, fullMethod string, req, resp any, conn *grpc.ClientConn, opts ...grpc.CallOption) error {
		breakerCalls++
		return breaking(ctx, fullMethod, req, resp, conn, func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
			downstreamCalls++
			return nil
		}, opts...)
	}
	err = retrying(t.Context(), "/test.Service/GetBook", nil, nil, nil, invoker)

	require.Error(t, err)
	require.ErrorIs(t, err, breaker.ErrOpenState)
	require.True(t, status.IsLocalError(err))
	require.Equal(t, codes.ResourceExhausted, status.Code(err))
	require.Equal(t, 1, breakerCalls)
	require.Equal(t, 0, downstreamCalls)
}

func TestRetryPreservesRemoteResourceExhausted(t *testing.T) {
	clientLimiter, err := grpclimiter.NewClientLimiter(
		test.NoopLifecycle{},
		limiter.NewKeyMap(),
		test.NewLimiterConfig("user-agent", "1m", 2),
	)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, clientLimiter.Close(t.Context())) })

	limiting := grpclimiter.UnaryClientInterceptor(clientLimiter)
	retrying := grpcretry.UnaryClientInterceptor(test.NewGRPCRetryConfig(3, time.Nanosecond, codes.ResourceExhausted))
	limiterCalls := 0
	downstreamCalls := 0
	invoker := func(ctx context.Context, fullMethod string, req, resp any, conn *grpc.ClientConn, opts ...grpc.CallOption) error {
		limiterCalls++
		return limiting(ctx, fullMethod, req, resp, conn, func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
			downstreamCalls++
			if downstreamCalls == 1 {
				return status.Error(codes.ResourceExhausted, "remote resource exhausted")
			}

			return nil
		}, opts...)
	}
	err = retrying(t.Context(), "/test.Service/GetBook", nil, nil, nil, invoker)

	require.NoError(t, err)
	require.Equal(t, 2, limiterCalls)
	require.Equal(t, 2, downstreamCalls)
}
