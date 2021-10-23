package breaker

import (
	"context"

	breaker "github.com/sony/gobreaker"
	"google.golang.org/grpc"
)

// UnaryClientInterceptor for breaker.
func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	cb := breaker.NewCircuitBreaker(breaker.Settings{})

	return func(ctx context.Context, fullMethod string, req, resp interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		operation := func() (interface{}, error) {
			return nil, invoker(ctx, fullMethod, req, resp, cc, opts...)
		}

		_, err := cb.Execute(operation)

		return err
	}
}
