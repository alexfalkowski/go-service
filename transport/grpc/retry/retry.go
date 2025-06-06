package retry

import (
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	config "github.com/alexfalkowski/go-service/v2/retry"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
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
