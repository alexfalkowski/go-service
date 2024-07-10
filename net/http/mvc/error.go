package mvc

import (
	"context"
	"net/http"

	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/net/http/status"
)

// WriteError for mvc.
func WriteError(ctx context.Context, view *View, res http.ResponseWriter, err error) {
	meta.WithAttribute(ctx, "mvcError", meta.Error(err))

	res.WriteHeader(status.Code(err))
	view.ExecuteFailure(res, err)
}
