package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/alexfalkowski/go-service/meta"
)

// ServerError returns nil if the err http.ErrServerClosed.
func ServerError(err error) error {
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

// WriteError for http.
func WriteError(ctx context.Context, res http.ResponseWriter, err error, status int) {
	meta.WithAttribute(ctx, "httpError", meta.Error(err))

	http.Error(res, err.Error(), status)
}
