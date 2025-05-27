package errors

import (
	"net/http"

	"github.com/alexfalkowski/go-service/v2/errors"
)

// ServerError returns nil if the err http.ErrServerClosed.
func ServerError(err error) error {
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}
