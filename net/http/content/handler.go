package content

import (
	"context"
	"net/http"

	"github.com/alexfalkowski/go-service/net/http/meta"
	"github.com/alexfalkowski/go-service/net/http/status"
	"github.com/alexfalkowski/go-service/types/ptr"
)

// RequestHandler is a handler with a generic request and response.
type RequestHandler[Req any, Res any] func(ctx context.Context, req *Req) (*Res, error)

// NewRequestHandler for content.
func NewRequestHandler[Req any, Res any](cont *Content, handler RequestHandler[Req, Res]) http.HandlerFunc {
	return newHandler(cont, func(ctx context.Context) (*Res, error) {
		req := ptr.Zero[Req]()

		encoder := meta.Encoder(ctx)
		request := meta.Request(ctx)

		if err := encoder.Decode(request.Body, req); err != nil {
			return nil, status.FromError(http.StatusBadRequest, err)
		}

		return handler(ctx, req)
	})
}

// Handler is a handler with a generic response.
type Handler[Res any] func(ctx context.Context) (*Res, error)

// NewHandler for content.
func NewHandler[Res any](cont *Content, handler Handler[Res]) http.HandlerFunc {
	return newHandler(cont, func(ctx context.Context) (*Res, error) {
		return handler(ctx)
	})
}

func newHandler[Res any](cont *Content, handler func(ctx context.Context) (*Res, error)) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		ctx = meta.WithRequest(ctx, req)
		ctx = meta.WithResponse(ctx, res)
		media := cont.NewFromRequest(req)

		ctx = meta.WithEncoder(ctx, media.Encoder)
		res.Header().Add(TypeKey, media.Type)

		data, err := handler(ctx)
		if err != nil {
			status.WriteError(ctx, res, err)

			return
		}

		if err := media.Encoder.Encode(res, data); err != nil {
			status.WriteError(ctx, res, err)

			return
		}
	}
}
