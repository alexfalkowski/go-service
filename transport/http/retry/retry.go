package retry

import (
	"context"
	"errors"
	"net/http"

	retry "github.com/avast/retry-go/v3"
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
		if err != nil {
			return err
		}

		if res.StatusCode == 0 || (res.StatusCode >= 500 && res.StatusCode != 501) {
			return ErrInvalidStatusCode
		}

		return nil
	}

	retry.Do(operation, retry.Attempts(r.cfg.Attempts)) // nolint:errcheck

	return res, err
}
