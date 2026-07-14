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

// limitError returns the ResourceExhausted error used when a limiter rejects a
// request. Client interceptors wrap it with status.LocalError; server
// interceptors leave it unmarked so remote retry configuration still applies.
func limitError() error {
	return status.Error(codes.ResourceExhausted, grpc.StatusText(codes.ResourceExhausted))
}

// serverLimitError returns ResourceExhausted with a google.rpc.RetryInfo detail
// built from the server limiter decision's reset window, mirroring the HTTP
// server limiter's Retry-After. RetryInfo is a status detail, so it can only
// ride the rejection error; the proactive quota state stays in the ratelimit
// response metadata. It advertises a delay only when reset timing is known, so
// clients are not told to retry immediately, and falls back to the bare
// limitError when the detail cannot be attached.
func serverLimitError(decision limiter.Decision) error {
	if reset := decision.ResetAfter(); reset > 0 {
		s, err := status.New(codes.ResourceExhausted, grpc.StatusText(codes.ResourceExhausted)).WithDetails(&status.RetryInfo{
			RetryDelay: status.NewDuration(reset),
		})
		if err == nil {
			return s.Err()
		}
	}

	return limitError()
}
