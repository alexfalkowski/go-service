package breaker_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/time"
	transportbreaker "github.com/alexfalkowski/go-service/v2/transport/breaker"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/breaker"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/retry"
	"github.com/stretchr/testify/require"
)

func TestClientUnaryInterceptorUsesConfigFailureCodes(t *testing.T) {
	interceptor := breaker.NewClient(
		breaker.NewConfig(test.NewBreaker(1), codes.InvalidArgument).Options()...,
	).UnaryInterceptor()

	calls := 0
	invoker := func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
		calls++
		return status.Error(codes.InvalidArgument, "invalid")
	}

	err := interceptor(t.Context(), "/test.Service/GetBook", nil, nil, nil, invoker)
	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))

	err = interceptor(t.Context(), "/test.Service/GetBook", nil, nil, nil, invoker)
	require.Error(t, err)
	require.True(t, status.IsLocalError(err))
	require.Equal(t, codes.ResourceExhausted, status.Code(err))
	require.Equal(t, 1, calls)
}

func TestClientUnaryInterceptorDoesNotOpenOnNonFailureCode(t *testing.T) {
	interceptor := breaker.NewClient(
		breaker.WithSettings(settings()),
		breaker.WithFailureCodes(codes.Unavailable),
	).UnaryInterceptor()

	calls := 0
	invoker := func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
		calls++
		return status.Error(codes.InvalidArgument, "invalid")
	}

	err := interceptor(t.Context(), "/test.Service/GetBook", nil, nil, nil, invoker)
	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))

	err = interceptor(t.Context(), "/test.Service/GetBook", nil, nil, nil, invoker)
	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
	require.Equal(t, 2, calls)
}

func TestClientUnaryInterceptorOpensOnClassifiedFailures(t *testing.T) {
	tests := map[string]struct {
		isSuccessful func(error) bool
		code         codes.Code
	}{
		"failure code wins over custom success": {
			isSuccessful: func(error) bool {
				return true
			},
			code: codes.Unavailable,
		},
		"custom failure handles non-failure code": {
			isSuccessful: func(error) bool {
				return false
			},
			code: codes.InvalidArgument,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			interceptor := breaker.NewClient(
				breaker.WithSettings(settings(test.isSuccessful)),
				breaker.WithFailureCodes(codes.Unavailable),
			).UnaryInterceptor()

			calls := 0
			invoker := func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
				calls++
				return status.Error(test.code, test.code.String())
			}

			err := interceptor(t.Context(), "/test.Service/GetBook", nil, nil, nil, invoker)
			require.Error(t, err)
			require.Equal(t, test.code, status.Code(err))

			err = interceptor(t.Context(), "/test.Service/GetBook", nil, nil, nil, invoker)
			require.Error(t, err)
			require.True(t, status.IsLocalError(err))
			require.Equal(t, codes.ResourceExhausted, status.Code(err))
			require.Equal(t, 1, calls)
		})
	}
}

func TestClientUnaryInterceptorDoesNotOpenOnCallerCancellation(t *testing.T) {
	tests := map[string][]breaker.Option{
		"canceled configured as a failure code": {
			breaker.WithSettings(settings()),
			breaker.WithFailureCodes(codes.Canceled),
		},
		"custom IsSuccessful treats canceled as a failure": {
			breaker.WithSettings(settings(func(error) bool { return false })),
		},
	}

	for name, opts := range tests {
		t.Run(name, func(t *testing.T) {
			interceptor := breaker.NewClient(opts...).UnaryInterceptor()

			invoker := func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
				return status.Error(codes.Canceled, "canceled")
			}

			err := interceptor(t.Context(), "/test.Service/GetBook", nil, nil, nil, invoker)
			require.Error(t, err)
			require.Equal(t, codes.Canceled, status.Code(err))
			require.False(t, status.IsLocalError(err))

			err = interceptor(t.Context(), "/test.Service/GetBook", nil, nil, nil, invoker)
			require.Error(t, err)
			require.Equal(t, codes.Canceled, status.Code(err))
			require.False(t, status.IsLocalError(err))
		})
	}
}

func TestClientUnaryInterceptorDoesNotOpenOnExpiredCallerDeadline(t *testing.T) {
	interceptor := breaker.NewClient(
		breaker.WithSettings(settings()),
	).UnaryInterceptor()

	expired, cancel := context.WithTimeout(t.Context(), 0)
	t.Cleanup(cancel)
	<-expired.Done()

	calls := 0
	invoker := func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
		calls++
		return status.Error(codes.DeadlineExceeded, "caller deadline exceeded")
	}

	err := interceptor(expired, "/test.Service/GetBook", nil, nil, nil, invoker)
	require.Error(t, err)
	require.Equal(t, codes.DeadlineExceeded, status.Code(err))

	err = interceptor(t.Context(), "/test.Service/GetBook", nil, nil, nil, func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
		calls++
		return nil
	})
	require.NoError(t, err)
	require.Equal(t, 2, calls)
}

func TestClientUnaryInterceptorOpensOnRemoteDeadline(t *testing.T) {
	interceptor := breaker.NewClient(
		breaker.WithSettings(settings()),
	).UnaryInterceptor()

	calls := 0
	invoker := func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
		calls++
		return status.Error(codes.DeadlineExceeded, "remote deadline exceeded")
	}

	err := interceptor(t.Context(), "/test.Service/GetBook", nil, nil, nil, invoker)
	require.Error(t, err)
	require.Equal(t, codes.DeadlineExceeded, status.Code(err))

	err = interceptor(t.Context(), "/test.Service/GetBook", nil, nil, nil, invoker)
	require.Error(t, err)
	require.True(t, status.IsLocalError(err))
	require.Equal(t, codes.ResourceExhausted, status.Code(err))
	require.Equal(t, 1, calls)
}

func TestClientUnaryInterceptorOpensOnAttemptDeadline(t *testing.T) {
	interceptor := breaker.NewClient(
		breaker.WithSettings(settings()),
	).UnaryInterceptor()
	retryConfig := test.NewGRPCRetryConfig(1, time.Nanosecond)
	retryConfig.Timeout = 10 * time.Millisecond
	retrying := retry.NewClient(retryConfig).UnaryInterceptor()

	calls := 0
	attempt := func(ctx context.Context, fullMethod string, req, resp any, conn *grpc.ClientConn, callOpts ...grpc.CallOption) error {
		return interceptor(ctx, fullMethod, req, resp, conn, func(ctx context.Context, _ string, _ any, _ any, _ *grpc.ClientConn, _ ...grpc.CallOption) error {
			calls++
			<-ctx.Done()
			return status.Error(codes.DeadlineExceeded, "attempt deadline exceeded")
		}, callOpts...)
	}
	err := retrying(t.Context(), "/test.Service/GetBook", nil, nil, nil, attempt)
	require.Error(t, err)
	require.Equal(t, codes.DeadlineExceeded, status.Code(err))

	err = interceptor(t.Context(), "/test.Service/GetBook", nil, nil, nil, func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
		calls++
		return nil
	})
	require.Error(t, err)
	require.True(t, status.IsLocalError(err))
	require.Equal(t, codes.ResourceExhausted, status.Code(err))
	require.Equal(t, 1, calls)
}

func TestClientUnaryInterceptorIsolatesBreakersByFullMethod(t *testing.T) {
	interceptor := breaker.NewClient(
		breaker.WithSettings(settings()),
		breaker.WithFailureCodes(codes.Unavailable),
	).UnaryInterceptor()

	calls := make(map[string]int)
	invoker := func(_ context.Context, fullMethod string, _, _ any, _ *grpc.ClientConn, _ ...grpc.CallOption) error {
		calls[fullMethod]++
		if fullMethod == "/test.Service/GetBook" {
			return status.Error(codes.Unavailable, "unavailable")
		}

		return nil
	}

	t.Run("opens breaker for failing method", func(t *testing.T) {
		err := interceptor(t.Context(), "/test.Service/GetBook", nil, nil, nil, invoker)
		require.Error(t, err)
		require.Equal(t, codes.Unavailable, status.Code(err))
	})

	t.Run("allows different method", func(t *testing.T) {
		err := interceptor(t.Context(), "/test.Service/ListBooks", nil, nil, nil, invoker)
		require.NoError(t, err)
	})

	t.Run("rejects failing method without invoker call", func(t *testing.T) {
		err := interceptor(t.Context(), "/test.Service/GetBook", nil, nil, nil, invoker)
		require.Error(t, err)
		require.True(t, status.IsLocalError(err))
		require.Equal(t, codes.ResourceExhausted, status.Code(err))
		require.Equal(t, 1, calls["/test.Service/GetBook"])
		require.Equal(t, 1, calls["/test.Service/ListBooks"])
	})
}

func settings(isSuccessful ...func(error) bool) transportbreaker.Settings {
	settings := transportbreaker.Settings{
		ReadyToTrip: func(counts transportbreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 1
		},
	}
	if len(isSuccessful) > 0 {
		settings.IsSuccessful = isSuccessful[0]
	}

	return settings
}
