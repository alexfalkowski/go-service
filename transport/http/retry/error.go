package retry

import (
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/strings"
)

type responseError struct {
	resp *http.Response
	err  error
}

// Error intentionally hides the wrapped error text because callers only use this type for control flow.
func (e responseError) Error() string {
	return strings.Empty
}

// Unwrap returns the underlying retry reason.
func (e responseError) Unwrap() error {
	return e.err
}
