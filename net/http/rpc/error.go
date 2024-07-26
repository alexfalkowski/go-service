package rpc

import (
	"context"

	"github.com/alexfalkowski/go-service/net/http"
	hc "github.com/alexfalkowski/go-service/net/http/context"
	"github.com/alexfalkowski/go-service/net/http/status"
)

// WriteError for rpc.
func WriteError(ctx context.Context, err error) {
	http.WriteError(ctx, hc.Response(ctx), err, status.Code(err))
}
