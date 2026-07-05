package retry

import (
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/time"
	config "github.com/alexfalkowski/go-service/v2/transport/retry"
)

func retryInfoDelayExceedsBackoff(err error, backoff time.Duration) bool {
	return retryInfoDelay(err) > minimumJitteredBackoff(backoff)
}

func retryInfoDelay(err error) time.Duration {
	grpcStatus, ok := status.FromError(err)
	if !ok {
		return 0
	}

	for _, detail := range grpcStatus.Details() {
		if retryInfo, ok := detail.(*status.RetryInfo); ok {
			delay := retryInfo.GetRetryDelay()
			if delay == nil {
				return 0
			}

			return time.Duration(delay.AsDuration())
		}
	}

	return 0
}

func minimumJitteredBackoff(backoff time.Duration) time.Duration {
	return backoff - (backoff * time.Duration(config.DefaultJitterPercent) / 100)
}
