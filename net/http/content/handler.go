package content

import (
	"context"
	"net/http"

	"github.com/alexfalkowski/go-service/errors"
	nh "github.com/alexfalkowski/go-service/net/http"
	hc "github.com/alexfalkowski/go-service/net/http/context"
	"github.com/alexfalkowski/go-service/net/http/status"
	"github.com/alexfalkowski/go-service/runtime"
)

type (
	handler func(ctx context.Context) (any, error)

	// RequestHandler is a handler with a generic request and response.
	RequestHandler[Req any, Res any] func(ctx context.Context, req *Req) (*Res, error)

	// Handler is a handler with a generic response.
	Handler[Res any] func(ctx context.Context) (*Res, error)
)

// NewRequestHandler for content.
func NewRequestHandler[Req any, Res any](cont *Content, prefix string, handler RequestHandler[Req, Res]) http.HandlerFunc {
	return cont.handler(prefix, func(ctx context.Context) (any, error) {
		var rq Req
		ptr := &rq

		e := hc.Encoder(ctx)
		req := hc.Request(ctx)

		if err := e.Decode(req.Body, ptr); err != nil {
			return nil, status.Error(http.StatusBadRequest, err.Error())
		}

		rs, err := handler(ctx, ptr)
		if err != nil {
			return nil, err
		}

		return rs, nil
	})
}

// NewHandler for content.
func NewHandler[Res any](cont *Content, prefix string, handler Handler[Res]) http.HandlerFunc {
	return cont.handler(prefix, func(ctx context.Context) (any, error) {
		return handler(ctx)
	})
}

func (c *Content) handler(prefix string, handler handler) http.HandlerFunc {
	h := func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		ctx = hc.WithRequest(ctx, req)
		ctx = hc.WithResponse(ctx, res)

		defer func() {
			if r := recover(); r != nil {
				err := errors.Prefix(prefix, runtime.ConvertRecover(r))
				nh.WriteError(ctx, res, err, status.Code(err))
			}
		}()

		ct := c.NewFromRequest(req)

		ctx = hc.WithEncoder(ctx, ct.Encoder)
		res.Header().Add(TypeKey, ct.Type)

		data, err := handler(ctx)
		runtime.Must(err)

		err = ct.Encoder.Encode(res, data)
		runtime.Must(err)
	}

	return h
}
