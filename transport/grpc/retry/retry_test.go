package retry_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	gstatus "github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/retry"
	"github.com/stretchr/testify/require"
)

func TestUnaryClientInterceptorDoesNotRetryWhenAttemptsIsOne(t *testing.T) {
	interceptor := retry.UnaryClientInterceptor(&retry.Config{
		Attempts: 1,
		Timeout:  "1s",
		Backoff:  "1ms",
	})

	calls := 0
	err := interceptor(t.Context(), "/test.Service/SayHello", nil, nil, nil, func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
		calls++
		return gstatus.Error(codes.Unavailable, "unavailable")
	})

	require.Error(t, err)
	require.Equal(t, 1, calls)
}

func TestUnaryClientInterceptorRetriesWhenAttemptsIsTwo(t *testing.T) {
	interceptor := retry.UnaryClientInterceptor(&retry.Config{
		Attempts: 2,
		Timeout:  "1s",
		Backoff:  "1ms",
	})

	calls := 0
	err := interceptor(t.Context(), "/test.Service/SayHello", nil, nil, nil, func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
		calls++
		if calls == 1 {
			return gstatus.Error(codes.Unavailable, "unavailable")
		}

		return nil
	})

	require.NoError(t, err)
	require.Equal(t, 2, calls)
}
