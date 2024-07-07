package rpc

import (
	"context"
	"net/http"
)

// WriteError for rpc.
func WriteError(ctx context.Context, err error) {
	http.Error(Response(ctx), err.Error(), Code(err))
}
