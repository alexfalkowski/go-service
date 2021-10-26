package retry

import (
	"context"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"time"

	retry "github.com/avast/retry-go/v3"
)

var (
	// A regular expression to match the error returned by net/http when the
	// configured number of redirects is exhausted. This error isn't typed
	// specifically so we resort to matching on the error string.
	redirectsErrorRe = regexp.MustCompile(`stopped after \d+ redirects\z`)

	// A regular expression to match the error returned by net/http when the
	// scheme specified in the URL is invalid. This error isn't typed
	// specifically so we resort to matching on the error string.
	schemeErrorRe = regexp.MustCompile(`unsupported protocol scheme`)
)

// NewRoundTripper for retry.
func NewRoundTripper(hrt http.RoundTripper) *RoundTripper {
	return &RoundTripper{RoundTripper: hrt}
}

// RoundTripper for retry.
type RoundTripper struct {
	http.RoundTripper
}

func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var (
		res *http.Response
		err error
	)

	ctx := req.Context()

	operation := func() error {
		tctx, cancel := context.WithTimeout(ctx, 2*time.Second) // nolint:gomnd
		defer cancel()

		res, err = r.RoundTripper.RoundTrip(req.WithContext(tctx)) // nolint:bodyclose
		if err != nil {
			if v, ok := err.(*url.Error); ok {
				// Don't retry if the error was due to too many redirects.
				if redirectsErrorRe.MatchString(v.Error()) {
					return nil
				}

				// Don't retry if the error was due to an invalid protocol scheme.
				if schemeErrorRe.MatchString(v.Error()) {
					return nil
				}

				// Don't retry if the error was due to TLS cert verification failure.
				if _, ok := v.Err.(x509.UnknownAuthorityError); ok {
					return nil
				}
			}

			// The error is likely recoverable so retry.
			return err
		}

		// Check the response code. We retry on 500-range responses to allow
		// the server time to recover, as 500's are typically not permanent
		// errors and may relate to outages on the server side. This will catch
		// invalid response codes as well, like 0 and 999.
		if res.StatusCode == 0 || (res.StatusCode >= 500 && res.StatusCode != 501) {
			return fmt.Errorf("invalid status code %d", res.StatusCode)
		}

		return nil
	}

	// We don't need to check the error as it's only used to retry. We save the last error in err.
	retry.Do(operation, retry.Attempts(3)) // nolint:errcheck,gomnd

	return res, err
}
