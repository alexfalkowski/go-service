package retry

import (
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	config "github.com/alexfalkowski/go-service/v2/retry"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
)

// Config is an alias for `github.com/alexfalkowski/go-service/v2/retry.Config`.
//
// It describes the retry policy for gRPC unary client calls, including:
//   - `Attempts`: maximum number of attempts (initial call + retries).
//   - `Timeout`: per-attempt timeout duration string.
//   - `Backoff`: backoff duration string used by the chosen backoff strategy.
type Config = config.Config

// UnaryClientInterceptor returns a gRPC unary client interceptor that retries failed calls.
//
// The interceptor is built using `go-grpc-middleware`'s retry interceptor and is intended to be used on the
// client side.
//
// Behavior:
//   - It parses per-attempt timeout and backoff durations from cfg (`cfg.Timeout` and `cfg.Backoff`).
//   - It retries up to `cfg.Attempts` total attempts (including the initial attempt).
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
	timeout := time.MustParseDuration(cfg.Timeout)
	backoff := time.MustParseDuration(cfg.Backoff)

	return retry.UnaryClientInterceptor(
		retry.WithCodes(codes.Unavailable, codes.DataLoss),
		retry.WithMax(uint(cfg.Attempts)),
		retry.WithBackoff(retry.BackoffLinear(backoff)),
		retry.WithPerRetryTimeout(timeout),
	)
}
