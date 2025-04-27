package content

import (
	"net/http"

	"github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/net/http/context"
	"github.com/alexfalkowski/go-service/net/http/status"
	"github.com/alexfalkowski/go-service/types/ptr"
)

// RequestHandler is a handler with a generic request and response.
type RequestHandler[Req any, Res any] func(ctx context.Context, req *Req) (*Res, error)

// NewRequestHandler for content.
func NewRequestHandler[Req any, Res any](cont *Content, prefix string, handler RequestHandler[Req, Res]) http.HandlerFunc {
	return newHandler(cont, prefix, func(ctx context.Context) (*Res, error) {
		req := ptr.Zero[Req]()

		encoder := context.Encoder(ctx)
		request := context.Request(ctx)

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
		ctx = context.WithRequest(ctx, req)
		ctx = context.WithResponse(ctx, res)
		media := cont.NewFromRequest(req)

		ctx = context.WithEncoder(ctx, media.Encoder)
		res.Header().Add(TypeKey, media.Type)

		data, err := handler(ctx)
		if err != nil {
			status.WriteError(res, errors.Prefix(prefix, err))

			return
		}

		if err := media.Encoder.Encode(res, data); err != nil {
			status.WriteError(res, errors.Prefix(prefix, err))

			return
		}
	}
}
