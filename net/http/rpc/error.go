package rpc

import (
	"context"
	"net/http"

	"github.com/alexfalkowski/go-service/net/http/status"
)

// WriteError for rpc.
func WriteError(ctx context.Context, err error) {
	http.Error(Response(ctx), err.Error(), status.Code(err))
}
