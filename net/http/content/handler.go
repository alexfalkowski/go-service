package content

import (
	"context"
	"net/http"

	"github.com/alexfalkowski/go-service/errors"
	nh "github.com/alexfalkowski/go-service/net/http"
	hc "github.com/alexfalkowski/go-service/net/http/context"
	"github.com/alexfalkowski/go-service/net/http/status"
	"github.com/alexfalkowski/go-service/runtime"
	"github.com/alexfalkowski/go-service/types/ptr"
)

// RequestHandler is a handler with a generic request and response.
type RequestHandler[Req any, Res any] func(ctx context.Context, req *Req) (*Res, error)

// NewRequestHandler for content.
func NewRequestHandler[Req any, Res any](cont *Content, prefix string, handler RequestHandler[Req, Res]) http.HandlerFunc {
	return contentHanler(cont, prefix, func(ctx context.Context) (*Res, error) {
		request := ptr.Zero[Req]()

		e := hc.Encoder(ctx)
		req := hc.Request(ctx)

		if err := e.Decode(req.Body, request); err != nil {
			return nil, status.Error(http.StatusBadRequest, err.Error())
		}

		rs, err := handler(ctx, request)
		if err != nil {
			return nil, err
		}

		return rs, nil
	})
}

// Handler is a handler with a generic response.
type Handler[Res any] func(ctx context.Context) (*Res, error)

// NewHandler for content.
func NewHandler[Res any](cont *Content, prefix string, handler Handler[Res]) http.HandlerFunc {
	return contentHanler(cont, prefix, func(ctx context.Context) (*Res, error) {
		return handler(ctx)
	})
}

func contentHanler[Res any](cont *Content, prefix string, handler func(ctx context.Context) (*Res, error)) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		ctx = hc.WithRequest(ctx, req)
		ctx = hc.WithResponse(ctx, res)

		defer func() {
			if r := recover(); r != nil {
				err := errors.Prefix(prefix, runtime.ConvertRecover(r))
				nh.WriteError(ctx, res, err, status.Code(err))
			}
		}()

		media := cont.NewFromRequest(req)

		ctx = hc.WithEncoder(ctx, media.Encoder)
		res.Header().Add(TypeKey, media.Type)

		data, err := handler(ctx)
		runtime.Must(err)

		err = media.Encoder.Encode(res, data)
		runtime.Must(err)
	}
}
