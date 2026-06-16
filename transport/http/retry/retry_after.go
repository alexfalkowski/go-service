package retry

import (
	"math"
	"strconv"

	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/time"
	config "github.com/alexfalkowski/go-service/v2/transport/retry"
)

const (
	retryAfterHeader     = "Retry-After"
	maxRetryAfterDelay   = time.Duration(math.MaxInt64)
	maxRetryAfterSeconds = uint64(maxRetryAfterDelay / time.Second)
)

func retryAfterDelayExceedsBackoff(res *http.Response, backoff time.Duration) bool {
	return retryAfterDelay(res) > minimumJitteredBackoff(backoff)
}

func retryAfterDelay(res *http.Response) time.Duration {
	value := res.Header.Get(retryAfterHeader)
	if value == "" {
		return 0
	}

	if delay, ok := parseRetryAfterSeconds(value); ok {
		return delay
	}

	return parseRetryAfterDate(value)
}

func minimumJitteredBackoff(backoff time.Duration) time.Duration {
	return backoff - (backoff * time.Duration(config.DefaultJitterPercent) / 100)
}

func parseRetryAfterSeconds(value string) (time.Duration, bool) {
	if seconds, err := strconv.ParseUint(value, 10, 64); err == nil {
		if seconds == 0 {
			return 0, true
		}
		if seconds > maxRetryAfterSeconds {
			return maxRetryAfterDelay, true
		}

		return time.Duration(seconds) * time.Second, true
	} else if isRetryAfterRangeError(err) {
		return maxRetryAfterDelay, true
	}

	return 0, false
}

func parseRetryAfterDate(value string) time.Duration {
	when, err := http.ParseTime(value)
	if err != nil {
		return 0
	}

	delay := time.Until(when)
	if delay <= 0 {
		return 0
	}

	return delay
}

func isRetryAfterRangeError(err error) bool {
	numErr, ok := errors.AsType[*strconv.NumError](err)
	return ok && errors.Is(numErr.Err, strconv.ErrRange)
}
