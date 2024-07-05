package http

import (
	"net/http"
)

// WriteError for HTTP.
func WriteError(ctx Context, err error) {
	http.Error(ctx.Response(), err.Error(), Code(err))
}
