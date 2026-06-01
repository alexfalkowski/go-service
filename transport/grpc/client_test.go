package grpc_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	tls "github.com/alexfalkowski/go-service/v2/crypto/tls/config"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	grpcmeta "github.com/alexfalkowski/go-service/v2/net/grpc/meta"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/slices"
	"github.com/alexfalkowski/go-service/v2/time"
	transportgrpc "github.com/alexfalkowski/go-service/v2/transport/grpc"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/retry"
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
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)
	_, err = client.SayHello(t.Context(), &v1.SayHelloRequest{Name: "test"})
	require.Error(t, err)
}

func TestUnaryClientInterceptorsGenerateStableRequestIDBeforeRetry(t *testing.T) {
	interceptors := transportgrpc.UnaryClientInterceptors(
		transportgrpc.WithClientID(&test.IDSequenceGenerator{IDs: []string{"generated-id"}}),
		transportgrpc.WithClientRetry(&retry.Config{
			Attempts: 2,
			Timeout:  time.Second,
			Backoff:  time.Millisecond,
		}),
	)

	calls := 0
	requestIDs := []string{}
	err := invokeUnaryClientMethod(t.Context(), "/test.Service/CreateBook", interceptors,
		func(ctx context.Context, _ string, _, _ any, _ *grpc.ClientConn, _ ...grpc.CallOption) error {
			calls++
			requestIDs = append(requestIDs, meta.RequestID(ctx).Value())
			if calls == 1 {
				return status.Error(codes.Unavailable, "unavailable")
			}

			return nil
		},
	)

	require.NoError(t, err)
	require.Equal(t, 2, calls)
	require.Equal(t, []string{"generated-id", "generated-id"}, requestIDs)
}

func TestUnaryClientInterceptorsGenerateTokenPerRetryAttempt(t *testing.T) {
	interceptors := transportgrpc.UnaryClientInterceptors(
		transportgrpc.WithClientRetry(&retry.Config{
			Attempts: 2,
			Timeout:  time.Second,
			Backoff:  time.Millisecond,
		}),
		transportgrpc.WithClientTokenGenerator(env.UserID("service-user"), &test.SequenceGenerator{}),
	)

	calls := 0
	authorizations := []string{}
	err := invokeUnaryClientMethod(t.Context(), "/test.Service/GetBook", interceptors,
		func(ctx context.Context, _ string, _, _ any, _ *grpc.ClientConn, _ ...grpc.CallOption) error {
			calls++
			md, ok := grpcmeta.FromOutgoingContext(ctx)
			require.True(t, ok)
			authorizations = append(authorizations, md.Get("authorization")...)
			if calls == 1 {
				return status.Error(codes.Unavailable, "unavailable")
			}

			return nil
		},
	)

	require.NoError(t, err)
	require.Equal(t, 2, calls)
	require.Equal(t, []string{"Bearer token-1", "Bearer token-2"}, authorizations)
}

func invokeUnaryClientMethod(ctx context.Context, fullMethod string, interceptors []grpc.UnaryClientInterceptor, invoker grpc.UnaryInvoker) error {
	chained := invoker
	for _, interceptor := range slices.Backward(interceptors) {
		next := chained
		chained = func(ctx context.Context, fullMethod string, req, resp any, conn *grpc.ClientConn, opts ...grpc.CallOption) error {
			return interceptor(ctx, fullMethod, req, resp, conn, next, opts...)
		}
	}

	return chained(ctx, fullMethod, nil, nil, nil)
}
