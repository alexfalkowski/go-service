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

// WriteResponse handles the error and adds it to the context with meta.WithAttribute.
func WriteResponse(ctx context.Context, res http.ResponseWriter, b []byte) {
	if _, err := res.Write(b); err != nil {
		meta.WithAttribute(ctx, "httpError", meta.Error(err))
		res.WriteHeader(http.StatusInternalServerError)
	}
}
