package breaker

import (
	"context"
	"errors"

	breaker "github.com/sony/gobreaker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryClientInterceptor for breaker.
func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	cb := breaker.NewCircuitBreaker(breaker.Settings{})

	return func(ctx context.Context, fullMethod string, req, resp any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		operation := func() (any, error) {
			return nil, invoker(ctx, fullMethod, req, resp, cc, opts...)
		}

		_, err := cb.Execute(operation)
		if err != nil {
			if errors.Is(err, breaker.ErrOpenState) || errors.Is(err, breaker.ErrTooManyRequests) {
				return status.Error(codes.Unavailable, err.Error())
			}

			return err
		}

		return nil
	}
}
