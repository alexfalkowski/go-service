package rpc

import (
	"context"
	"net/http"

	"github.com/alexfalkowski/go-service/errors"
	nh "github.com/alexfalkowski/go-service/net/http"
	"github.com/alexfalkowski/go-service/net/http/content"
	hc "github.com/alexfalkowski/go-service/net/http/context"
	"github.com/alexfalkowski/go-service/net/http/status"
	"github.com/alexfalkowski/go-service/runtime"
)

// Handler for rpc.
type Handler[Req any, Res any] func(ctx context.Context, req *Req) (*Res, error)

// Route for rpc.
func Route[Req any, Res any](path string, handler Handler[Req, Res]) {
	h := func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		ctx = hc.WithRequest(ctx, req)
		ctx = hc.WithResponse(ctx, res)

		defer func() {
			if r := recover(); r != nil {
				err := errors.Prefix("rpc", runtime.ConvertRecover(r))
				nh.WriteError(ctx, res, err, status.Code(err))
			}
		}()

		ct := content.NewFromRequest(req, enc)
		res.Header().Add(content.TypeKey, ct.Media)

		var rq Req
		ptr := &rq

		err := ct.Encoder.Decode(req.Body, ptr)
		runtime.Must(err)

		rs, err := handler(ctx, ptr)
		runtime.Must(err)

		err = ct.Encoder.Encode(res, rs)
		runtime.Must(err)
	}

	mux.HandleFunc("POST "+path, h)
}
