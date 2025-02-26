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
	return newHandler(cont, prefix, func(ctx context.Context) (*Res, error) {
		req := ptr.Zero[Req]()

		encoder := hc.Encoder(ctx)
		request := hc.Request(ctx)

		if err := encoder.Decode(request.Body, req); err != nil {
			return nil, status.Error(http.StatusBadRequest, err.Error())
		}

		return handler(ctx, req)
	})
}

// Handler is a handler with a generic response.
type Handler[Res any] func(ctx context.Context) (*Res, error)

// NewHandler for content.
func NewHandler[Res any](cont *Content, prefix string, handler Handler[Res]) http.HandlerFunc {
	return newHandler(cont, prefix, func(ctx context.Context) (*Res, error) {
		return handler(ctx)
	})
}

func newHandler[Res any](cont *Content, prefix string, handler func(ctx context.Context) (*Res, error)) http.HandlerFunc {
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
