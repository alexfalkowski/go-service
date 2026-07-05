package retry

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/retry"
	"github.com/alexfalkowski/go-service/v2/slices"
	"github.com/alexfalkowski/go-service/v2/strings"
	config "github.com/alexfalkowski/go-service/v2/transport/retry"
)

// Policy decides whether a unary RPC is eligible for retry.
//
// Policies describe operation safety, not transient failure classification. The retry interceptor still only retries
// configured gRPC status codes after the policy allows the logical RPC.
// When multiple non-nil policies are provided, all must allow the call. Nil policies are ignored; if every
// provided policy is nil, the default [IdempotentMethods] policy is used.
type Policy func(ctx context.Context, fullMethod string, req any) bool

// StandardReadMethods allows retries for AIP-style read methods named Get* or List*.
func StandardReadMethods(_ context.Context, fullMethod string, _ any) bool {
	_, method := grpc.ParseServiceMethod(fullMethod)

	return strings.HasPrefix(method, "Get") || strings.HasPrefix(method, "List")
}

// HasRequestID allows retries when request metadata contains a request id.
//
// In go-service, Request-Id identifies the logical request, not an individual wire attempt. The client metadata
// middleware installs it outside the retry middleware, so all retry attempts for one logical request share the same
// value. Services that retry writes should treat Request-Id as the idempotency key and deduplicate repeated attempts
// by it.
func HasRequestID(ctx context.Context, _ string, _ any) bool {
	return !meta.RequestID(ctx).IsEmpty()
}

// IdempotentMethods allows AIP-style read methods, or requests with a request-id idempotency contract.
func IdempotentMethods(ctx context.Context, fullMethod string, req any) bool {
	return StandardReadMethods(ctx, fullMethod, req) || HasRequestID(ctx, fullMethod, req)
}

// UnaryClientInterceptor returns a gRPC unary client interceptor that retries failed calls.
//
// Behavior:
//   - It checks whether the logical RPC is eligible for retry using policies.
//   - It uses the typed per-attempt timeout from cfg.GetTimeout() and the backoff duration from cfg.
//   - It retries up to `cfg.MaxAttempts()` total attempts (including the initial attempt).
//   - It applies a per-attempt timeout so each attempt is bounded.
//   - It uses a jittered backoff duration derived from `cfg.GetBackoff()`.
//
// Failure classification:
// Retries are only attempted for configured gRPC status codes. The default is [codes.Unavailable].
//
// Policy behavior:
// When no policy is provided, only side-effect-safe unary RPCs are eligible for retry: AIP-style read methods,
// or calls carrying a request-id idempotency contract. Callers that need different behavior can pass an
// explicit policy.
// Multiple non-nil policies are composed with logical AND, so any denying policy makes the RPC non-retryable.
// Nil policies are ignored, and the default policy is used only when no non-nil policy is supplied.
//
// Notes:
// This interceptor does not automatically retry on every error; application-level errors that map to other
// status codes will not be retried by default.
func UnaryClientInterceptor(cfg *Config, policies ...Policy) grpc.UnaryClientInterceptor {
	policy := composePolicy(policies)
	maxAttempts := cfg.MaxAttempts()
	timeout := cfg.GetTimeout()
	backoff := cfg.GetBackoff()
	codes := retryableCodes(cfg)

	invoke := func(ctx context.Context, attempt func(context.Context) error) error {
		if maxAttempts == 0 {
			return attempt(ctx)
		}

		retries := retry.WithJitterPercent(config.DefaultJitterPercent, retry.NewConstant(backoff))
		retries = retry.WithMaxRetries(maxAttempts-1, retries)
		return retry.Do(ctx, retries, func(ctx context.Context) error {
			err := attempt(ctx)
			if err == nil || !isRetryableCode(status.Code(err), codes) {
				return err
			}
			if retryInfoDelayExceedsBackoff(err, backoff) {
				return err
			}

			return retry.RetryableError(err)
		})
	}

	return func(ctx context.Context, fullMethod string, req, resp any, conn *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if !policy(ctx, fullMethod, req) {
			return invoker(ctx, fullMethod, req, resp, conn, opts...)
		}

		return invoke(ctx, func(ctx context.Context) error {
			attemptCtx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			return invoker(attemptCtx, fullMethod, req, resp, conn, opts...)
		})
	}
}

func retryableCodes(cfg *Config) []codes.Code {
	if cfg == nil || len(cfg.Codes) == 0 {
		return []codes.Code{codes.Unavailable}
	}

	codes := make([]codes.Code, 0, len(cfg.Codes))
	for _, code := range cfg.Codes {
		if isValidCode(code) {
			codes = append(codes, code)
		}
	}

	return codes
}

func isRetryableCode(code codes.Code, codes []codes.Code) bool {
	return slices.Contains(codes, code)
}

func isValidCode(code codes.Code) bool {
	return code > codes.OK && code <= codes.Unauthenticated
}

func composePolicy(policies []Policy) Policy {
	filtered := make([]Policy, 0, len(policies))
	for _, policy := range policies {
		if policy != nil {
			filtered = append(filtered, policy)
		}
	}

	if len(filtered) == 0 {
		return IdempotentMethods
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
