package retry

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/strings"
	config "github.com/alexfalkowski/go-service/v2/transport/retry"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
)

// Config is an alias for `github.com/alexfalkowski/go-service/v2/transport/retry.Config`.
//
// It describes the retry policy for gRPC unary client calls, including:
//   - `Attempts`: maximum number of attempts (initial call + retries).
//   - `Timeout`: per-attempt timeout duration.
//   - `Backoff`: backoff duration used by the chosen backoff strategy.
type Config = config.Config

// Policy decides whether a unary RPC is eligible for retry.
//
// Policies describe operation safety, not transient failure classification. The retry interceptor still only retries
// configured gRPC status codes after the policy allows the logical RPC.
type Policy func(ctx context.Context, fullMethod string, req any) bool

// AllowAll allows retries for every unary RPC.
func AllowAll(context.Context, string, any) bool {
	return true
}

// StandardReadMethods allows retries for AIP-style read methods named Get* or List*.
func StandardReadMethods(_ context.Context, fullMethod string, _ any) bool {
	_, method, _ := strings.SplitServiceMethod(fullMethod)

	return strings.HasPrefix(method, "Get") || strings.HasPrefix(method, "List")
}

// HasRequestID allows retries when request metadata contains a request id.
//
// A request id only makes write retries safe when the server treats it as an idempotency key and deduplicates
// repeated attempts for the same logical operation.
func HasRequestID(ctx context.Context, _ string, _ any) bool {
	return !meta.RequestID(ctx).IsEmpty()
}

// IdempotentMethods allows AIP-style read methods, or requests with a request-id idempotency contract.
func IdempotentMethods(ctx context.Context, fullMethod string, req any) bool {
	return StandardReadMethods(ctx, fullMethod, req) || HasRequestID(ctx, fullMethod, req)
}

// UnaryClientInterceptor returns a gRPC unary client interceptor that retries failed calls.
//
// The interceptor is built using `go-grpc-middleware`'s retry interceptor and is intended to be used on the
// client side.
//
// Behavior:
//   - It checks whether the logical RPC is eligible for retry using policies.
//   - It uses the typed per-attempt timeout and backoff durations from cfg.
//   - It retries up to `cfg.MaxAttempts()` total attempts (including the initial attempt).
//   - It applies a per-attempt timeout (`retry.WithPerRetryTimeout`) so each attempt is bounded.
//   - It uses a linear backoff strategy with a step duration derived from `cfg.Backoff`.
//
// Failure classification:
// Retries are only attempted for selected gRPC status codes. This implementation currently retries on
// `codes.Unavailable` and `codes.DataLoss` (see `retry.WithCodes` in the implementation).
//
// Notes:
// This interceptor does not automatically retry on every error; application-level errors that map to other
// status codes will not be retried by default.
func UnaryClientInterceptor(cfg *Config, policies ...Policy) grpc.UnaryClientInterceptor {
	policy := composePolicy(policies)
	interceptor := retry.UnaryClientInterceptor(
		retry.WithCodes(codes.Unavailable, codes.DataLoss),
		retry.WithMax(uint(cfg.MaxAttempts())),
		retry.WithBackoff(retry.BackoffLinear(cfg.Backoff.Duration())),
		retry.WithPerRetryTimeout(cfg.Timeout.Duration()),
	)

	return func(ctx context.Context, fullMethod string, req, resp any, conn *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if !policy(ctx, fullMethod, req) {
			return invoker(ctx, fullMethod, req, resp, conn, opts...)
		}

		return interceptor(ctx, fullMethod, req, resp, conn, invoker, opts...)
	}
}

func composePolicy(policies []Policy) Policy {
	filtered := make([]Policy, 0, len(policies))
	for _, policy := range policies {
		if policy != nil {
			filtered = append(filtered, policy)
		}
	}

	if len(filtered) == 0 {
		return AllowAll
	}

	return func(ctx context.Context, fullMethod string, req any) bool {
		for _, policy := range filtered {
			if !policy(ctx, fullMethod, req) {
				return false
			}
		}

		return true
	}
}
