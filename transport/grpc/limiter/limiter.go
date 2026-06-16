package limiter

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/transport/limiter"
)

// KeyMap is an alias for [limiter.KeyMap].
//
// It maps limiter key kinds (for example, "user-agent" or "ip") to functions that derive a rate-limit key
// from the request context.
type KeyMap = limiter.KeyMap

func take(ctx context.Context, rateLimiter *limiter.Limiter) (limiter.Decision, error) {
	decision, err := rateLimiter.TakeDecision(ctx)
	if err != nil {
		return decision, status.SafeError(codes.Internal, err)
	}

	return decision, nil
}

func limitError() error {
	return status.Error(codes.ResourceExhausted, grpc.StatusText(codes.ResourceExhausted))
}
