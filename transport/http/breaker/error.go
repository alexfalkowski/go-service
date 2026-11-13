package breaker

import (
	"github.com/alexfalkowski/go-service/v3/net/http"
	"github.com/alexfalkowski/go-service/v3/strings"
)

type responseError struct {
	resp *http.Response
}

// Error satisfies the error interface without exposing response details.
func (e responseError) Error() string {
	return strings.Empty
}
