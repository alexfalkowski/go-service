package rpc

import (
	"context"
	"net/http"

	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/net/http/status"
)

// WriteError for rpc.
func WriteError(ctx context.Context, err error) {
	ctx = meta.WithAttribute(ctx, "rpcError", meta.Error(err))

	http.Error(Response(ctx), err.Error(), status.Code(err))
}
