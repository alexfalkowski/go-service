package retry

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/alexfalkowski/go-service/retry"
	rth "github.com/hashicorp/go-retryablehttp"
)

// ErrInvalidStatusCode for http retry.
var ErrInvalidStatusCode = errors.New("invalid status code")

// NewRoundTripper for retry.
func NewRoundTripper(cfg *retry.Config, hrt http.RoundTripper) *RoundTripper {
	return &RoundTripper{cfg: cfg, RoundTripper: hrt}
}

// RoundTripper for retry.
type RoundTripper struct {
	cfg *retry.Config
	http.RoundTripper
}

func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	d, err := time.ParseDuration(r.cfg.Timeout)
	if err != nil {
		return nil, err
	}

	var res *http.Response

	ctx := req.Context()
	operation := func() error {
		tctx, cancel := context.WithTimeout(ctx, d)
		defer cancel()

		res, err = r.RoundTripper.RoundTrip(req.WithContext(tctx))
		ok, perr := rth.ErrorPropagatedRetryPolicy(tctx, res, err)

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

	_ = retry.Try(operation, r.cfg)

	return res, err
}
