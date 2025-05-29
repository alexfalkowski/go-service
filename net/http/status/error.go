package status

import (
	"context"
	"fmt"

	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net/http"
)

// WriteError will write the error to the response writer.
// This is based on http.Error.
func WriteError(ctx context.Context, res http.ResponseWriter, err error) {
	header := res.Header()
	header.Del("Content-Length")
	header.Set("Content-Type", "text/error; charset=utf-8")
	header.Set("X-Content-Type-Options", "nosniff")

	res.WriteHeader(Code(err))

	if _, err := fmt.Fprintln(res, err.Error()); err != nil {
		meta.WithAttribute(ctx, "writeError", meta.Error(err))
	}
}
