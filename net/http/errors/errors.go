package errors

import (
	"errors"
	"net/http"
)

// ServerError returns nil if the err http.ErrServerClosed.
func ServerError(err error) error {
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}
