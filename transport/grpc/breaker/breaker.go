package breaker

import (
	"context"

	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	breaker "github.com/sony/gobreaker"
)

// UnaryClientInterceptor for breaker.
func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	cbkr := breaker.NewCircuitBreaker(breaker.Settings{})

	return func(ctx context.Context, fullMethod string, req, resp any, conn *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		operation := func() (any, error) {
			return nil, invoker(ctx, fullMethod, req, resp, conn, opts...)
		}

		_, err := cbkr.Execute(operation)
		if err != nil {
			if errors.Is(err, breaker.ErrOpenState) || errors.Is(err, breaker.ErrTooManyRequests) {
				return status.Error(codes.Unavailable, err.Error())
			}

			return err
		}

		return nil
	}
}
