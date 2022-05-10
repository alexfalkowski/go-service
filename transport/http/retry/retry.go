package retry

import (
	"context"
	"errors"
	"net/http"

	retry "github.com/avast/retry-go/v3"
	"github.com/hashicorp/go-retryablehttp"
)

// ErrInvalidStatusCode for http retry.
var ErrInvalidStatusCode = errors.New("invalid status code")

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
	var (
		res *http.Response
		err error
	)

	ctx := req.Context()

	operation := func() error {
		tctx, cancel := context.WithTimeout(ctx, r.cfg.Timeout)
		defer cancel()

		res, err = r.RoundTripper.RoundTrip(req.WithContext(tctx))
		ok, perr := retryablehttp.ErrorPropagatedRetryPolicy(tctx, res, err)

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

	_ = retry.Do(operation, retry.Attempts(r.cfg.Attempts))

	return res, err
}
