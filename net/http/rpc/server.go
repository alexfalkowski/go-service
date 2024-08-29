package rpc

import (
	"context"

	"github.com/alexfalkowski/go-service/net/http/content"
	hc "github.com/alexfalkowski/go-service/net/http/context"
	"github.com/alexfalkowski/go-service/runtime"
)

// Handler for rpc.
type Handler[Req any, Res any] func(ctx context.Context, req *Req) (*Res, error)

// Route for rpc.
func Route[Req any, Res any](path string, handler Handler[Req, Res]) {
	h := content.NewHandler("rpc", enc, func(ctx context.Context) any {
		e := hc.Encoder(ctx)
		req := hc.Request(ctx)

		var rq Req
		ptr := &rq

		err := e.Decode(req.Body, ptr)
		runtime.Must(err)

		rs, err := handler(ctx, ptr)
		runtime.Must(err)

		return rs
	})

	mux.HandleFunc("POST "+path, h)
}
