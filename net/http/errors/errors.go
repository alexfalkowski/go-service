package errors

import (
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net/http"
)

// ServerError normalizes expected HTTP server shutdown errors.
//
// The standard library returns http.ErrServerClosed from (*http.Server).Serve / ServeTLS when the server
// is shut down normally via (*http.Server).Shutdown or Close. In many go-service components this is not
// considered a real error and should not be surfaced as a startup/serve failure.
//
// ServerError returns nil when err is http.ErrServerClosed, otherwise it returns err unchanged.
func ServerError(err error) error {
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}
