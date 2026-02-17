package content

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/meta"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/types/ptr"
)

// RequestHandler handles a decoded request and returns a response model.
type RequestHandler[Req any, Res any] func(ctx context.Context, req *Req) (*Res, error)

// NewRequestHandler builds a handler that decodes the request body and encodes the response.
//
// The encoder is selected from the request Content-Type, and decode errors are turned into
// a 400 response via net/http/status.
func NewRequestHandler[Req any, Res any](cont *Content, handler RequestHandler[Req, Res]) http.HandlerFunc {
	return newHandler(cont, func(ctx context.Context) (*Res, error) {
		req := ptr.Zero[Req]()

		encoder := meta.Encoder(ctx)
		request := meta.Request(ctx)

		if err := encoder.Decode(request.Body, req); err != nil {
			return nil, status.BadRequestError(err)
		}

		return handler(ctx, req)
	})
}

// Handler handles a request without a request body and returns a response model.
type Handler[Res any] func(ctx context.Context) (*Res, error)

// NewHandler builds a handler that encodes the response and writes errors using status helpers.
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
