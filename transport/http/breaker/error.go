package breaker

import (
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/transport/strings"
)

type responseError struct {
	resp *http.Response
}

// Error is just used to satisfy the error interface. We don't want to expose the response details.
func (e responseError) Error() string {
	return strings.Empty
}
