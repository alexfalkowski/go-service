package http

import (
	"context"
	"net/http"

	"github.com/alexfalkowski/go-service/meta"
)

// WriteError for http.
func WriteError(ctx context.Context, res http.ResponseWriter, err error, status int) {
	meta.WithAttribute(ctx, "httpError", meta.Error(err))

	http.Error(res, err.Error(), status)
}
