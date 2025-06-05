package retry

import (
	config "github.com/alexfalkowski/go-service/v2/retry"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// Config is an alias retry.Config.
type Config = config.Config

// UnaryClientInterceptor for retry.
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
