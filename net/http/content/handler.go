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
	// Handler for content.
	Handler func(ctx context.Context) (any, error)

	// RequestResponseHandler is a handler with a generic request and response.
	RequestResponseHandler[Req any, Res any] func(ctx context.Context, req *Req) (*Res, error)

	// ResponseHandler is a handler with a generic response.
	ResponseHandler[Res any] func(ctx context.Context) (*Res, error)
)

// NewRequestResponseHandler for content.
func NewRequestResponseHandler[Req any, Res any](cont *Content, prefix string, handler RequestResponseHandler[Req, Res]) http.HandlerFunc {
	return cont.NewHandler(prefix, func(ctx context.Context) (any, error) {
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

// NewResponseHandler for content.
func NewResponseHandler[Res any](cont *Content, prefix string, handler ResponseHandler[Res]) http.HandlerFunc {
	return cont.NewHandler(prefix, func(ctx context.Context) (any, error) {
		return handler(ctx)
	})
}

// NewHandler for content.
func (c *Content) NewHandler(prefix string, handler Handler) http.HandlerFunc {
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
