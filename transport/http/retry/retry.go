package retry

import (
	"context"

	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/retry"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/hashicorp/go-retryablehttp"
)

// Config is an alias for retry.Config.
type Config = retry.Config

// ErrInvalidStatusCode for http retry.
var ErrInvalidStatusCode = errors.New("retry: invalid status code")

// NewRoundTripper for retry.
func NewRoundTripper(cfg *Config, hrt http.RoundTripper) *RoundTripper {
	return &RoundTripper{cfg: cfg, RoundTripper: hrt}
}

// RoundTripper for retry.
type RoundTripper struct {
	cfg *Config
	http.RoundTripper
}

func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	timeout := time.MustParseDuration(r.cfg.Timeout)

	var (
		res *http.Response
		err error
	)

	ctx := req.Context()
	operation := func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		res, err = r.RoundTripper.RoundTrip(req.WithContext(ctx))
		ok, perr := retryablehttp.ErrorPropagatedRetryPolicy(ctx, res, err)

		if ok {
			if perr != nil {
				return perr
			}

			if err != nil {
				return err
			}
		}

		return nil
	}

	_ = retry.Try(ctx, operation, r.cfg)

	return res, err
}
