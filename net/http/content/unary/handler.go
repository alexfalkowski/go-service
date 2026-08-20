package unary

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/encoding/codec"
	encodingerrors "github.com/alexfalkowski/go-service/v2/encoding/errors"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/net/http/meta"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/ptr"
)

// RequestHandler handles a decoded request and returns a response model.
type RequestHandler[Req any, Res any] func(ctx context.Context, req *Req) (*Res, error)

// NewRequestHandler builds a handler that decodes the request body and encodes the response.
//
// Context population:
// The handler attaches request-scoped values to the context via net/http/meta:
//   - the original *[http.Request],
//   - the [http.ResponseWriter].
//
// Content negotiation:
// Request-body decoding uses the request Content-Type, falling back to JSON when Content-Type is absent.
// An unparseable, unregistered, or intentionally undecodable Content-Type (see the decoder-bounds rule in
// the package documentation) is rejected with [ErrUnsupportedRequestMedia] rather than falling back to
// JSON; see [Content.NewFromRequestBody]. Unknown request members are discarded so API additions remain
// forward-compatible. Response encoding uses the request Accept header, falling back
// to Content-Type when Accept is absent. The response Content-Type header is set to the negotiated
// response media type. The selected response encoder remains internal to the handler.
//
// Errors:
// If request decoding fails, NewRequestHandler converts the decode error into a 400 Bad Request using
// net/http/status, allowing the response to be rendered consistently by [status.WriteError].
// If response encoding fails because the response value does not satisfy the negotiated encoder's
// required type (see [github.com/alexfalkowski/go-service/v2/encoding/errors.ErrInvalidType]), the encode
// error becomes a 406 Not Acceptable; any other encode failure keeps its default status.
//
// Successful responses are encoded into a pooled in-memory buffer before being written to the live
// response writer, so encode failures do not leak partial success bodies.
//
// Size limits:
// NewRequestHandler does not enforce request body size caps itself. In supported wiring, inbound request bodies
// are capped by the transport HTTP server before content handlers run. Direct users should wrap the handler with
// an equivalent request-size limit.
func NewRequestHandler[Req any, Res any](content *Content, handler RequestHandler[Req, Res]) http.HandlerFunc {
	return newHandler(content, func(ctx context.Context) (*Res, error) {
		req := ptr.Zero[Req]()

		request := meta.Request(ctx)
		mediaType, err := content.NewFromRequestBody(request)
		if err != nil {
			return nil, status.SafeError(http.StatusUnsupportedMediaType, err)
		}

		if err := mediaType.Encoder.Decode(request.Body, req, codec.WithDiscardUnknown()); err != nil {
			return nil, status.BadRequestError(err)
		}

		return handler(ctx, req)
	})
}

// Handler handles a request without a request body and returns a response model.
type Handler[Res any] func(ctx context.Context) (*Res, error)

// NewHandler builds a handler that encodes the response and writes errors using status helpers.
//
// Context population, response content negotiation, and encode-error handling are the same as
// NewRequestHandler (see its documentation).
//
// Successful responses are encoded into a pooled in-memory buffer before being written to the live
// response writer, so encode failures do not leak partial success bodies.
func NewHandler[Res any](content *Content, handler Handler[Res]) http.HandlerFunc {
	return newHandler(content, func(ctx context.Context) (*Res, error) {
		return handler(ctx)
	})
}

func newHandler[Res any](content *Content, handler func(ctx context.Context) (*Res, error)) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		http.AddVary(res.Header(), http.AcceptKey, http.ContentTypeKey)

		mediaType := content.NewFromAccept(req)
		ctx = meta.WithRequestResponse(ctx, req, res)
		res.Header().Set(http.ContentTypeKey, media.MustParse(mediaType.String()).WithUTF8())

		if mediaType.Encoder == nil {
			_ = status.WriteError(ctx, res, ErrUnsupportedMedia)

			return
		}

		data, err := handler(ctx)
		if err != nil {
			_ = status.WriteError(ctx, res, err)

			return
		}

		buffer := content.pool.Get()
		defer content.pool.Put(buffer)

		if err := mediaType.Encoder.Encode(buffer, data); err != nil {
			if errors.Is(err, encodingerrors.ErrInvalidType) {
				err = status.SafeError(http.StatusNotAcceptable, err)
			}

			_ = status.WriteError(ctx, res, err)

			return
		}

		_, _ = buffer.WriteTo(res)
	}
}
