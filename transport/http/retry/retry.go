package retry

import (
	"context"
	"errors"
	"net/http"

	"github.com/alexfalkowski/go-service/retry"
	"github.com/alexfalkowski/go-service/time"
	rth "github.com/hashicorp/go-retryablehttp"
)

// ErrInvalidStatusCode for http retry.
var ErrInvalidStatusCode = errors.New("retry: invalid status code")

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
	d := time.MustParseDuration(r.cfg.Timeout)

	var (
		res *http.Response
		err error
	)

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
