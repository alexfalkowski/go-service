package breaker_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	basebreaker "github.com/alexfalkowski/go-service/v2/transport/breaker"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/breaker"
	"github.com/stretchr/testify/require"
)

func TestUnaryClientInterceptorUsesConfigFailureCodes(t *testing.T) {
	interceptor := breaker.UnaryClientInterceptor(
		breaker.NewConfig(test.NewBreaker(1), codes.InvalidArgument).Options()...,
	)

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
	require.Equal(t, codes.ResourceExhausted, status.Code(err))
	require.Equal(t, 1, calls)
}

func TestUnaryClientInterceptorDoesNotOpenOnNonFailureCode(t *testing.T) {
	interceptor := breaker.UnaryClientInterceptor(
		breaker.WithSettings(settings()),
		breaker.WithFailureCodes(codes.Unavailable),
	)

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

func TestUnaryClientInterceptorOpensOnClassifiedFailures(t *testing.T) {
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
			interceptor := breaker.UnaryClientInterceptor(
				breaker.WithSettings(settings(test.isSuccessful)),
				breaker.WithFailureCodes(codes.Unavailable),
			)

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
			require.Equal(t, codes.ResourceExhausted, status.Code(err))
			require.Equal(t, 1, calls)
		})
	}
}

func TestUnaryClientInterceptorIsolatesBreakersByFullMethod(t *testing.T) {
	interceptor := breaker.UnaryClientInterceptor(
		breaker.WithSettings(settings()),
		breaker.WithFailureCodes(codes.Unavailable),
	)

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
		require.Equal(t, codes.ResourceExhausted, status.Code(err))
		require.Equal(t, 1, calls["/test.Service/GetBook"])
		require.Equal(t, 1, calls["/test.Service/ListBooks"])
	})
}

func settings(isSuccessful ...func(error) bool) breaker.Settings {
	settings := breaker.Settings{
		ReadyToTrip: func(counts basebreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 1
		},
	}
	if len(isSuccessful) > 0 {
		settings.IsSuccessful = isSuccessful[0]
	}

	return settings
}
