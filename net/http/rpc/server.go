package rpc

import (
	"context"

	"github.com/alexfalkowski/go-service/net/http/content"
)

// Handler for rpc.
type Handler[Req any, Res any] func(ctx context.Context, req *Req) (*Res, error)

// Route for rpc.
func Route[Req any, Res any](path string, handler Handler[Req, Res]) {
	h := cont.NewHandler("rpc", func(ctx context.Context) (any, error) {
		var rq Req
		ptr := &rq

		err := content.Decode(ctx, ptr)
		if err != nil {
			return nil, err
		}

		rs, err := handler(ctx, ptr)
		if err != nil {
			return nil, err
		}

		return rs, nil
	})

	mux.HandleFunc("POST "+path, h)
}
