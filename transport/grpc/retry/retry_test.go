package retry_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/retry"
	config "github.com/alexfalkowski/go-service/v2/transport/retry"
	"github.com/stretchr/testify/require"
)

func TestUnaryClientInterceptorDoesNotRetryWhenAttemptsIsOne(t *testing.T) {
	interceptor := retry.UnaryClientInterceptor(&config.Config{
		Attempts: 1,
		Timeout:  time.Second,
		Backoff:  time.Millisecond,
	})

	calls := 0
	err := interceptor(t.Context(), "/test.Service/SayHello", nil, nil, nil, func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
		calls++
		return status.Error(codes.Unavailable, "unavailable")
	})

	require.Error(t, err)
	require.Equal(t, 1, calls)
}

func TestUnaryClientInterceptorDoesNotRetryWhenAttemptsIsZero(t *testing.T) {
	interceptor := retry.UnaryClientInterceptor(&config.Config{
		Attempts: 0,
		Timeout:  time.Second,
		Backoff:  time.Millisecond,
	})

	calls := 0
	err := interceptor(t.Context(), "/test.Service/GetHello", nil, nil, nil, func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
		calls++
		return status.Error(codes.Unavailable, "unavailable")
	})

	require.Error(t, err)
	require.Equal(t, 1, calls)
}

func TestUnaryClientInterceptorRetriesSafeMethodWhenAttemptsIsTwo(t *testing.T) {
	interceptor := retry.UnaryClientInterceptor(&config.Config{
		Attempts: 2,
		Timeout:  time.Second,
		Backoff:  time.Millisecond,
	})

	calls := 0
	err := interceptor(t.Context(), "/test.Service/GetHello", nil, nil, nil, func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
		calls++
		if calls == 1 {
			return status.Error(codes.Unavailable, "unavailable")
		}

		return nil
	})

	require.NoError(t, err)
	require.Equal(t, 2, calls)
}

func TestUnaryClientInterceptorRetriesWithDefaultBackoff(t *testing.T) {
	interceptor := retry.UnaryClientInterceptor(&config.Config{
		Attempts: 2,
		Timeout:  time.Second,
	})

	calls := 0
	err := interceptor(t.Context(), "/test.Service/GetHello", nil, nil, nil, func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
		calls++
		if calls == 1 {
			return status.Error(codes.Unavailable, "unavailable")
		}

		return nil
	})

	require.NoError(t, err)
	require.Equal(t, 2, calls)
}

func TestUnaryClientInterceptorSetsDeadlineForEachAttempt(t *testing.T) {
	interceptor := retry.UnaryClientInterceptor(&config.Config{
		Attempts: 2,
		Timeout:  time.Second,
		Backoff:  time.Millisecond,
	})

	calls := 0
	err := interceptor(t.Context(), "/test.Service/GetHello", nil, nil, nil, func(ctx context.Context, _ string, _, _ any, _ *grpc.ClientConn, _ ...grpc.CallOption) error {
		calls++
		_, ok := ctx.Deadline()
		require.True(t, ok, "attempt %d should have a deadline", calls)
		if calls == 1 {
			return status.Error(codes.Unavailable, "unavailable")
		}

		return nil
	})

	require.NoError(t, err)
	require.Equal(t, 2, calls)
}

func TestUnaryClientInterceptorDoesNotRetryUnsafeMethodByDefault(t *testing.T) {
	interceptor := retry.UnaryClientInterceptor(&config.Config{
		Attempts: 2,
		Timeout:  time.Second,
		Backoff:  time.Millisecond,
	})

	calls := 0
	err := interceptor(t.Context(), "/test.Service/CreateBook", nil, nil, nil, func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
		calls++
		return status.Error(codes.Unavailable, "unavailable")
	})

	require.Error(t, err)
	require.Equal(t, 1, calls)
}

func TestUnaryClientInterceptorDoesNotRetryDataLossByDefault(t *testing.T) {
	interceptor := retry.UnaryClientInterceptor(&config.Config{
		Attempts: 2,
		Timeout:  time.Second,
		Backoff:  time.Millisecond,
	})

	calls := 0
	err := interceptor(t.Context(), "/test.Service/GetBook", nil, nil, nil, func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
		calls++
		return status.Error(codes.DataLoss, "data loss")
	})

	require.Error(t, err)
	require.Equal(t, 1, calls)
}

func TestUnaryClientInterceptorDoesNotRetryContextStatusCodesByDefault(t *testing.T) {
	tests := []struct {
		name string
		code codes.Code
	}{
		{name: "deadline exceeded", code: codes.DeadlineExceeded},
		{name: "canceled", code: codes.Canceled},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interceptor := retry.UnaryClientInterceptor(&config.Config{
				Attempts: 2,
				Timeout:  time.Second,
				Backoff:  time.Millisecond,
			})

			calls := 0
			err := interceptor(t.Context(), "/test.Service/GetBook", nil, nil, nil, func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
				calls++
				return status.Error(tt.code, tt.name)
			})

			require.Error(t, err)
			require.Equal(t, 1, calls)
		})
	}
}

func TestUnaryClientInterceptorDoesNotRetryWhenPolicyDeniesMethod(t *testing.T) {
	interceptor := retry.UnaryClientInterceptor(&config.Config{
		Attempts: 2,
		Timeout:  time.Second,
		Backoff:  time.Millisecond,
	}, retry.StandardReadMethods)

	calls := 0
	err := interceptor(t.Context(), "/test.Service/CreateBook", nil, nil, nil, func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
		calls++
		return status.Error(codes.Unavailable, "unavailable")
	})

	require.Error(t, err)
	require.Equal(t, 1, calls)
}

func TestUnaryClientInterceptorRetriesWhenPolicyAllowsReadMethod(t *testing.T) {
	interceptor := retry.UnaryClientInterceptor(&config.Config{
		Attempts: 2,
		Timeout:  time.Second,
		Backoff:  time.Millisecond,
	}, retry.StandardReadMethods)

	calls := 0
	err := interceptor(t.Context(), "/test.Service/GetBook", nil, nil, nil, func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
		calls++
		if calls == 1 {
			return status.Error(codes.Unavailable, "unavailable")
		}

		return nil
	})

	require.NoError(t, err)
	require.Equal(t, 2, calls)
}

func TestUnaryClientInterceptorComposesMultiplePolicies(t *testing.T) {
	allow := retry.Policy(func(context.Context, string, any) bool { return true })
	deny := retry.Policy(func(context.Context, string, any) bool { return false })

	tests := []struct {
		name     string
		policies []retry.Policy
		calls    int
	}{
		{name: "deny wins", policies: []retry.Policy{allow, deny}, calls: 1},
		{name: "nil policies are ignored", policies: []retry.Policy{allow, nil, allow}, calls: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interceptor := retry.UnaryClientInterceptor(&config.Config{
				Attempts: 2,
				Timeout:  time.Second,
				Backoff:  time.Millisecond,
			}, tt.policies...)

			calls := 0
			err := interceptor(t.Context(), "/test.Service/CreateBook", nil, nil, nil, func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
				calls++
				if calls == 1 {
					return status.Error(codes.Unavailable, "unavailable")
				}

				return nil
			})
			if tt.calls == 1 {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.calls, calls)
		})
	}
}

func TestUnaryClientInterceptorRetriesWhenPolicyAllowsRequestID(t *testing.T) {
	interceptor := retry.UnaryClientInterceptor(&config.Config{
		Attempts: 2,
		Timeout:  time.Second,
		Backoff:  time.Millisecond,
	}, retry.IdempotentMethods)

	ctx := meta.WithAttributes(t.Context(), meta.WithRequestID(meta.String("request-id")))
	calls := 0
	err := interceptor(ctx, "/test.Service/CreateBook", nil, nil, nil, func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
		calls++
		if calls == 1 {
			return status.Error(codes.Unavailable, "unavailable")
		}

		return nil
	})

	require.NoError(t, err)
	require.Equal(t, 2, calls)
}
