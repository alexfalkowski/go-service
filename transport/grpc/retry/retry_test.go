package retry_test

import (
	"testing"
	"testing/synctest"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/retry"
	config "github.com/alexfalkowski/go-service/v2/transport/retry"
	"github.com/stretchr/testify/require"
)

func TestConfigRejectsInvalidCodes(t *testing.T) {
	require.NoError(t, test.Validator.Struct(test.NewGRPCRetryConfig(0, 0, codes.Unavailable)))
	require.Error(t, test.Validator.Struct(test.NewGRPCRetryConfig(0, 0, codes.OK)))
	require.Error(t, test.Validator.Struct(test.NewGRPCRetryConfig(0, 0, codes.Code(17))))
}

func TestNewConfigReturnsGRPCConfig(t *testing.T) {
	cfg := &config.Config{Attempts: 2, Backoff: time.Millisecond}
	retrying := retry.NewConfig(cfg, codes.ResourceExhausted)

	require.Same(t, cfg, retrying.Config)
	require.Equal(t, []codes.Code{codes.ResourceExhausted}, retrying.Codes)
	require.Nil(t, retry.NewConfig(nil))
}

func TestClientUnaryInterceptorDoesNotRetryWhenAttemptsIsOne(t *testing.T) {
	interceptor := retry.NewClient(test.NewGRPCRetryConfig(1, time.Millisecond)).UnaryInterceptor()

	calls := 0
	err := interceptor(t.Context(), "/test.Service/SayHello", nil, nil, nil, func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
		calls++
		return status.Error(codes.Unavailable, "unavailable")
	})

	require.Error(t, err)
	require.Equal(t, 1, calls)
}

func TestClientUnaryInterceptorDoesNotRetryWhenAttemptsIsZero(t *testing.T) {
	interceptor := retry.NewClient(test.NewGRPCRetryConfig(0, time.Millisecond)).UnaryInterceptor()

	calls := 0
	err := interceptor(t.Context(), "/test.Service/GetHello", nil, nil, nil, func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
		calls++
		return status.Error(codes.Unavailable, "unavailable")
	})

	require.Error(t, err)
	require.Equal(t, 1, calls)
}

func TestClientUnaryInterceptorClampsAttemptsAboveMax(t *testing.T) {
	interceptor := retry.NewClient(test.NewGRPCRetryConfig(config.MaxAttempts+1, time.Nanosecond)).UnaryInterceptor()

	calls := 0
	err := interceptor(t.Context(), "/test.Service/GetHello", nil, nil, nil, func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
		calls++
		return status.Error(codes.Unavailable, "unavailable")
	})

	require.Error(t, err)
	require.Equal(t, int(config.MaxAttempts), calls)
}

func TestClientUnaryInterceptorRetriesSafeMethodWhenAttemptsIsTwo(t *testing.T) {
	interceptor := retry.NewClient(test.NewGRPCRetryConfig(2, time.Millisecond)).UnaryInterceptor()

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

func TestClientUnaryInterceptorRetriesConfiguredCode(t *testing.T) {
	interceptor := retry.NewClient(test.NewGRPCRetryConfig(2, time.Millisecond, codes.ResourceExhausted)).UnaryInterceptor()

	calls := 0
	err := interceptor(t.Context(), "/test.Service/GetHello", nil, nil, nil, func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
		calls++
		if calls == 1 {
			return status.Error(codes.ResourceExhausted, "resource exhausted")
		}

		return nil
	})

	require.NoError(t, err)
	require.Equal(t, 2, calls)
}

func TestClientUnaryInterceptorDoesNotRetryLocalStatusError(t *testing.T) {
	interceptor := retry.NewClient(test.NewGRPCRetryConfig(2, time.Millisecond, codes.ResourceExhausted)).UnaryInterceptor()

	calls := 0
	err := interceptor(t.Context(), "/test.Service/GetHello", nil, nil, nil, func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
		calls++
		return status.LocalError(status.SafeError(codes.ResourceExhausted, test.ErrInvalid))
	})

	require.Error(t, err)
	require.ErrorIs(t, err, test.ErrInvalid)
	require.True(t, status.IsLocalError(err))
	require.Equal(t, codes.ResourceExhausted, status.Code(err))
	require.Equal(t, 1, calls)
}

func TestClientUnaryInterceptorConfiguredCodesReplaceDefaults(t *testing.T) {
	interceptor := retry.NewClient(test.NewGRPCRetryConfig(2, time.Millisecond, codes.ResourceExhausted)).UnaryInterceptor()

	calls := 0
	err := interceptor(t.Context(), "/test.Service/GetHello", nil, nil, nil, func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
		calls++
		return status.Error(codes.Unavailable, "unavailable")
	})

	require.Error(t, err)
	require.Equal(t, 1, calls)
}

func TestClientUnaryInterceptorCapsGrowingBackoffWithMaxBackoff(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		cfg := config.Config{Strategy: "exponential", Backoff: 2 * time.Millisecond, MaxBackoff: 10 * time.Millisecond, Timeout: time.Second, Attempts: config.MaxAttempts}
		interceptor := retry.NewClient(retry.NewConfig(&cfg)).UnaryInterceptor()

		calls := 0
		start := time.Now()
		err := interceptor(t.Context(), "/test.Service/GetHello", nil, nil, nil, func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
			calls++
			return status.Error(codes.Unavailable, "unavailable")
		})
		elapsed := time.Since(start)

		require.Error(t, err)
		require.Equal(t, int(config.MaxAttempts), calls)
		// Uncapped exponential growth from a 2ms base across 9 retries sums to ~1022ms; capping
		// at 10ms bounds the same schedule to well under that, so a generous bound below the
		// uncapped total still proves the cap is applied.
		require.Less(t, elapsed, 400*time.Millisecond)
	})
}

func TestClientUnaryInterceptorRetriesWithDefaultBackoff(t *testing.T) {
	interceptor := retry.NewClient(test.NewGRPCRetryConfig(2, 0)).UnaryInterceptor()

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

func TestClientUnaryInterceptorDoesNotRetryWhenRetryInfoDelayExceedsMinimumBackoff(t *testing.T) {
	interceptor := retry.NewClient(test.NewGRPCRetryConfig(2, 10*time.Millisecond)).UnaryInterceptor()

	calls := 0
	err := interceptor(t.Context(), "/test.Service/GetHello", nil, nil, nil, func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
		calls++
		return retryInfoError(t, 20*time.Millisecond)
	})

	require.Error(t, err)
	require.Equal(t, 1, calls)
}

func TestClientUnaryInterceptorHonorsGrownBackoffForRetryInfo(t *testing.T) {
	tests := []struct {
		name  string
		delay time.Duration
		calls int
		err   bool
	}{
		{name: "below the grown backoff retries", delay: 15 * time.Millisecond, calls: 3, err: false},
		{name: "above the grown backoff still suppresses", delay: 25 * time.Millisecond, calls: 2, err: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Exponential growth from a 10ms base gives ~10ms before attempt 2 and ~20ms before
			// attempt 3; the second attempt's gate must compare against the latter, not the base.
			cfg := config.Config{Strategy: "exponential", Backoff: 10 * time.Millisecond, Attempts: 3}
			interceptor := retry.NewClient(retry.NewConfig(&cfg)).UnaryInterceptor()

			calls := 0
			err := interceptor(t.Context(), "/test.Service/GetHello", nil, nil, nil, func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
				calls++
				switch calls {
				case 1:
					return status.Error(codes.Unavailable, "unavailable")
				case 2:
					return retryInfoError(t, tt.delay)
				default:
					return nil
				}
			})

			if tt.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.calls, calls)
		})
	}
}

func TestClientUnaryInterceptorDoesNotRetryWhenRetryInfoDelayExceedsCappedBackoff(t *testing.T) {
	cfg := config.Config{Strategy: "exponential", Backoff: 100 * time.Millisecond, MaxBackoff: 10 * time.Millisecond, Timeout: time.Second, Attempts: 2}
	interceptor := retry.NewClient(retry.NewConfig(&cfg)).UnaryInterceptor()

	calls := 0
	err := interceptor(t.Context(), "/test.Service/GetHello", nil, nil, nil, func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
		calls++
		return retryInfoError(t, 20*time.Millisecond)
	})

	require.Error(t, err)
	require.Equal(t, 1, calls)
}

func TestClientUnaryInterceptorRetriesWhenRetryInfoDelayDoesNotExceedBackoff(t *testing.T) {
	tests := map[string]time.Duration{
		"minimum": 8 * time.Millisecond,
		"zero":    0,
		"smaller": time.Millisecond,
	}

	for name, delay := range tests {
		t.Run(name, func(t *testing.T) {
			interceptor := retry.NewClient(test.NewGRPCRetryConfig(2, 10*time.Millisecond)).UnaryInterceptor()

			calls := 0
			err := interceptor(t.Context(), "/test.Service/GetHello", nil, nil, nil, func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
				calls++
				if calls == 1 {
					return retryInfoError(t, delay)
				}

				return nil
			})

			require.NoError(t, err)
			require.Equal(t, 2, calls)
		})
	}
}

func TestClientUnaryInterceptorRetriesWhenRetryInfoDelayIsMissing(t *testing.T) {
	interceptor := retry.NewClient(test.NewGRPCRetryConfig(2, 10*time.Millisecond)).UnaryInterceptor()

	calls := 0
	err := interceptor(t.Context(), "/test.Service/GetHello", nil, nil, nil, func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
		calls++
		if calls == 1 {
			grpcStatus, err := status.New(codes.Unavailable, "retry later").WithDetails(&status.RetryInfo{})
			require.NoError(t, err)

			return grpcStatus.Err()
		}

		return nil
	})

	require.NoError(t, err)
	require.Equal(t, 2, calls)
}

func TestClientUnaryInterceptorSetsDeadlineForEachAttempt(t *testing.T) {
	interceptor := retry.NewClient(test.NewGRPCRetryConfig(2, time.Millisecond)).UnaryInterceptor()

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

func TestClientUnaryInterceptorDoesNotRetryUnsafeMethodByDefault(t *testing.T) {
	interceptor := retry.NewClient(test.NewGRPCRetryConfig(2, time.Millisecond)).UnaryInterceptor()

	calls := 0
	err := interceptor(t.Context(), "/test.Service/CreateBook", nil, nil, nil, func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
		calls++
		return status.Error(codes.Unavailable, "unavailable")
	})

	require.Error(t, err)
	require.Equal(t, 1, calls)
}

func TestClientUnaryInterceptorDoesNotRetryDataLossByDefault(t *testing.T) {
	interceptor := retry.NewClient(test.NewGRPCRetryConfig(2, time.Millisecond)).UnaryInterceptor()

	calls := 0
	err := interceptor(t.Context(), "/test.Service/GetBook", nil, nil, nil, func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
		calls++
		return status.Error(codes.DataLoss, "data loss")
	})

	require.Error(t, err)
	require.Equal(t, 1, calls)
}

func TestClientUnaryInterceptorDoesNotRetryContextStatusCodesByDefault(t *testing.T) {
	tests := []struct {
		name string
		code codes.Code
	}{
		{name: "deadline exceeded", code: codes.DeadlineExceeded},
		{name: "canceled", code: codes.Canceled},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interceptor := retry.NewClient(test.NewGRPCRetryConfig(2, time.Millisecond)).UnaryInterceptor()

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

func TestClientUnaryInterceptorDoesNotRetryWhenPolicyDeniesMethod(t *testing.T) {
	interceptor := retry.NewClient(test.NewGRPCRetryConfig(2, time.Millisecond), retry.StandardReadMethods).UnaryInterceptor()

	calls := 0
	err := interceptor(t.Context(), "/test.Service/CreateBook", nil, nil, nil, func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
		calls++
		return status.Error(codes.Unavailable, "unavailable")
	})

	require.Error(t, err)
	require.Equal(t, 1, calls)
}

func TestClientUnaryInterceptorRetriesWhenPolicyAllowsReadMethod(t *testing.T) {
	tests := []struct {
		name       string
		fullMethod string
	}{
		{name: "get", fullMethod: "/test.Service/GetBook"},
		{name: "list", fullMethod: "/test.Service/ListBooks"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interceptor := retry.NewClient(test.NewGRPCRetryConfig(2, time.Millisecond), retry.StandardReadMethods).UnaryInterceptor()

			calls := 0
			err := interceptor(t.Context(), tt.fullMethod, nil, nil, nil, func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
				calls++
				if calls == 1 {
					return status.Error(codes.Unavailable, "unavailable")
				}

				return nil
			})

			require.NoError(t, err)
			require.Equal(t, 2, calls)
		})
	}
}

func TestClientUnaryInterceptorCombinesMultiplePolicies(t *testing.T) {
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
			interceptor := retry.NewClient(test.NewGRPCRetryConfig(2, time.Millisecond), tt.policies...).UnaryInterceptor()

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

func TestClientUnaryInterceptorRetriesWhenPolicyAllowsRequestID(t *testing.T) {
	interceptor := retry.NewClient(test.NewGRPCRetryConfig(2, time.Millisecond), retry.IdempotentMethods).UnaryInterceptor()

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

func retryInfoError(t *testing.T, delay time.Duration) error {
	t.Helper()

	grpcStatus, err := status.New(codes.Unavailable, "retry later").WithDetails(&status.RetryInfo{
		RetryDelay: status.NewDuration(delay),
	})
	require.NoError(t, err)

	return grpcStatus.Err()
}
