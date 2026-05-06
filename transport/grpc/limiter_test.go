package grpc_test

import (
	"strconv"
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/v2/limiter"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	grpclimiter "github.com/alexfalkowski/go-service/v2/transport/grpc/limiter"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func TestServerLimiterUnary(t *testing.T) {
	serverLimiter := requireServerLimiter(t, test.NewLimiterConfig("user-agent", "1s", 0), true)
	conn := test.NewBufconnGRPCConn(t, test.WithBufconnServerLimiter(serverLimiter))

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	_, _ = client.SayHello(t.Context(), req)
	_, err := client.SayHello(t.Context(), req)
	require.Error(t, err)
	require.Equal(t, codes.ResourceExhausted, status.Code(err))
}

func TestClientLimiterUnary(t *testing.T) {
	clientLimiter := requireClientLimiter(t, test.NewLimiterConfig("user-agent", "1s", 0), true)
	conn := test.NewBufconnGRPCConn(t, test.WithBufconnClientLimiter(clientLimiter))

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	_, _ = client.SayHello(t.Context(), req)
	_, err := client.SayHello(t.Context(), req)
	require.Error(t, err)
	require.Equal(t, codes.ResourceExhausted, status.Code(err))
}

func TestLimiterUnlimitedUnary(t *testing.T) {
	cfg := test.NewLimiterConfig("user-agent", "1s", 10)
	serverLimiter := requireServerLimiter(t, cfg, true)
	clientLimiter := requireClientLimiter(t, cfg, true)
	conn := test.NewBufconnGRPCConn(t, test.WithBufconnServerLimiter(serverLimiter), test.WithBufconnClientLimiter(clientLimiter))

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	_, err := client.SayHello(t.Context(), req)
	require.NoError(t, err)
}

func TestLimiterAuthUnary(t *testing.T) {
	serverLimiter := requireServerLimiter(t, test.NewLimiterConfig("user-agent", "1s", 10), true)
	conn := test.NewBufconnGRPCConn(t,
		test.WithBufconnServerLimiter(serverLimiter),
		test.WithBufconnGenerator(test.NewGenerator("bob", nil)),
		test.WithBufconnVerifier(test.NewVerifier("bob")),
	)

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	var err error
	for i := range 10 {
		t.Run("attempt-"+strconv.Itoa(i+1), func(t *testing.T) {
			_, err = client.SayHello(t.Context(), req)
		})
	}
	require.NoError(t, err)
}

func TestServerClosedLimiterUnary(t *testing.T) {
	serverLimiter := requireServerLimiter(t, test.NewLimiterConfig("user-agent", "1s", 10), false)
	require.NoError(t, serverLimiter.Close(t.Context()))

	conn := test.NewBufconnGRPCConn(t, test.WithBufconnServerLimiter(serverLimiter))

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	_, err := client.SayHello(t.Context(), req)
	require.Error(t, err)
	require.Equal(t, codes.Internal, status.Code(err))
}

func TestClientClosedLimiterUnary(t *testing.T) {
	clientLimiter := requireClientLimiter(t, test.NewLimiterConfig("user-agent", "1s", 10), false)
	conn := test.NewBufconnGRPCConn(t, test.WithBufconnClientLimiter(clientLimiter))

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	require.NoError(t, clientLimiter.Close(t.Context()))

	_, err := client.SayHello(t.Context(), req)
	require.Error(t, err)
	require.Equal(t, codes.Internal, status.Code(err))
}

func requireServerLimiter(t *testing.T, cfg *limiter.Config, cleanup bool) *grpclimiter.Server {
	t.Helper()

	lc := fxtest.NewLifecycle(t)
	limiter, err := test.NewGRPCServerLimiter(lc, test.LimiterKeyMap, cfg)
	require.NoError(t, err)
	if cleanup {
		t.Cleanup(func() {
			_ = limiter.Close(t.Context())
		})
	}

	return limiter
}

func requireClientLimiter(t *testing.T, cfg *limiter.Config, cleanup bool) *grpclimiter.Client {
	t.Helper()

	lc := fxtest.NewLifecycle(t)
	limiter, err := test.NewGRPCClientLimiter(lc, test.LimiterKeyMap, cfg)
	require.NoError(t, err)
	if cleanup {
		t.Cleanup(func() {
			_ = limiter.Close(t.Context())
		})
	}

	return limiter
}
