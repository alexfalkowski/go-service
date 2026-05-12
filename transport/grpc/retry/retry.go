package retry

import (
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
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

// UnaryClientInterceptor returns a gRPC unary client interceptor that retries failed calls.
//
// The interceptor is built using `go-grpc-middleware`'s retry interceptor and is intended to be used on the
// client side.
//
// Behavior:
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
func UnaryClientInterceptor(cfg *Config) grpc.UnaryClientInterceptor {
	return retry.UnaryClientInterceptor(
		retry.WithCodes(codes.Unavailable, codes.DataLoss),
		retry.WithMax(uint(cfg.MaxAttempts())),
		retry.WithBackoff(retry.BackoffLinear(cfg.Backoff.Duration())),
		retry.WithPerRetryTimeout(cfg.Timeout.Duration()),
	)
}
