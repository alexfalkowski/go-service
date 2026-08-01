package stream

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/encoding/stream"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/budget"
	"github.com/alexfalkowski/go-service/v2/net/http/compress"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/net/http/meta"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
)

// NewHandler builds a handler for a send-only streaming response using sm to resolve streaming codecs.
//
// Content negotiation:
// The response encoder is resolved from the request Accept header, falling back to Content-Type,
// using [NewMediaFromAccept]. Unlike single-value handlers, an Accept that cannot be satisfied
// is rejected with 406 rather than falling back to JSON.
//
// Error contract:
// A handler error returned before the first successful Send is an ordinary HTTP error rendered
// through [status.WriteError]. A handler error returned after the first successful Send aborts the
// response via panic([http.ErrAbortHandler]) instead, since the response is already committed — see
// [Stream.Send] and [finalizeStream].
//
// Drain:
// When opts.Drain starts before the handler is invoked, this handler returns 503. Otherwise it cancels
// the handler context with [ErrDraining]. A handler that returns after observing ctx.Done ends a committed
// response cleanly; it must therefore select on ctx.Done while waiting on an upstream source.
//
// opts.WriteTimeout, when positive, is pushed forward as the response write deadline after every
// successful Send (see [http.ResponseController.SetWriteDeadline]), turning a whole-stream write
// timeout into a per-message inactivity budget. Zero disables this.
//
// Load control:
// If the request context carries a limiter (see [meta.WithLimiter], populated by
// [github.com/alexfalkowski/go-service/v2/transport/http/limiter.Handler] for a request that was already
// allowed to open the stream), every Send charges one additional token against that same limiter — see
// [Stream.Send]. A route reached without that middleware, or with no limiter configured, sees no
// per-message charging.
func NewHandler[Res any](cont *content.Content, sm *stream.Map, opts Options, handler Handler[Res]) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		if isDraining(opts.Drain) {
			_ = status.WriteError(ctx, res, status.SafeError(http.StatusServiceUnavailable, ErrDraining))

			return
		}

		limiter := meta.Limiter(ctx)

		resMedia, err := NewMediaFromAccept(req, sm)
		if err == nil && resMedia.NewEncoder == nil {
			err = ErrUnsupportedMedia
		}

		if err != nil {
			_ = status.WriteError(ctx, res, status.SafeError(http.StatusNotAcceptable, err))
			return
		}

		ctx = meta.WithContent(ctx, req, res, nil)
		res.Header().Set(content.TypeKey, resMedia.WithUTF8())
		res.Header().Set(compress.HeaderNoCompression, "1")

		buffer := cont.BorrowBuffer()
		defer cont.ReturnBuffer(buffer)

		writer := &commitWriter{res: res, buffer: buffer}
		controller := http.NewResponseController(res)
		ctx, cancel := withDrain(ctx, opts.Drain, nil)
		if cancel != nil {
			defer cancel(nil)
		}

		st := &Stream[Res]{
			ctx:          ctx,
			writer:       writer,
			encoder:      resMedia.NewEncoder(writer),
			controller:   controller,
			writeTimeout: opts.WriteTimeout,
			limiter:      limiter,
		}
		finalizeStream(ctx, res, st, handler(ctx, st))
	}
}

// Handler handles a send-only stream: the response streams, the request does not.
//
// The handler must propagate a [Stream.Send] error by returning it (directly, or wrapped) rather than
// ignoring it and continuing. Send is sticky — once it fails, every later call returns the same error
// immediately — but only the handler returning that error triggers [finalizeStream]'s abort/HTTP-error
// handling; a handler that swallows a Send error and returns nil degrades to a no-op loop over a dead
// connection instead of ending the request. When configured [Options.Drain] starts, ctx is
// canceled with [ErrDraining]; handlers that wait on an upstream source must select on ctx.Done and
// return its error so the framework can finish the response.
type Handler[Res any] func(ctx context.Context, stream *Stream[Res]) error

// NewRequestHandler builds a handler for a bidirectional stream using sm to resolve streaming codecs.
//
// HTTP/2 requirement:
// Bidirectional streaming requires HTTP/2 (including h2c): an HTTP/1.x request body is buffered ahead
// of the handler by intermediaries and the Go transport, so interleaving Recv/Send hangs rather than
// failing. Requests with req.ProtoMajor < 2 are rejected with 505 before the handler runs.
//
// Content negotiation:
// The request decoder is resolved from Content-Type via [NewMediaFromContentType], rejecting
// an unregistered or unparseable media type with 415. The response encoder is resolved from Accept
// (falling back to Content-Type) via [NewMediaFromAccept], rejecting an Accept that cannot be
// satisfied with 406 instead: Accept negotiates the response representation, not the request payload,
// so a failure there is answered as RFC 9110 §15.5.7 rather than §15.5.16.
//
// Inbound size limiting:
// opts.MaxReceiveSize bounds each value decoded by [RequestStream.Recv], not the request stream as a
// whole (see [github.com/alexfalkowski/go-service/v2/net/http/budget.Reader]): a request with many small
// values is never rejected for its cumulative size, only
// for a single value that exceeds opts.MaxReceiveSize. A value at or under the limit is a normal terminal
// [http.MaxBytesError] surfaced from Recv; the deviation is that no total byte ceiling exists for a
// streaming route the way [github.com/alexfalkowski/go-service/v2/net/http/body.NewHandler]'s buffered
// path enforces one. opts.MaxReceiveSize <= 0 disables the per-value cap entirely.
//
// Timeout:
// When opts.ReadTimeout is positive, every successful [RequestStream.Recv] extends the request read
// deadline; when opts.WriteTimeout is positive, every successful [Stream.Send] extends the response
// write deadline. This turns the server's whole-request read and write timeouts into per-message
// inactivity budgets. Zero disables the corresponding extension.
//
// See [NewHandler] for the error contract and load control: Recv charges one token per successfully
// decoded, in-cap value, the same way Send does, against the same request-scoped limiter. Its drain
// contract also applies: drain cancels the handler context and closes the request body to unblock an
// active Recv. A client can need to reconnect if an HTTP/2 peer observes that body close as a stream reset.
func NewRequestHandler[Req any, Res any](cont *content.Content, sm *stream.Map, opts Options, handler RequestHandler[Req, Res]) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		if isDraining(opts.Drain) {
			_ = status.WriteError(ctx, res, status.SafeError(http.StatusServiceUnavailable, ErrDraining))

			return
		}

		limiter := meta.Limiter(ctx)

		if req.ProtoMajor < 2 {
			_ = status.WriteError(ctx, res, status.SafeError(http.StatusHTTPVersionNotSupported, ErrBidiRequiresHTTP2))
			return
		}

		reqMedia, err := NewMediaFromContentType(req, sm)
		if err != nil {
			_ = status.WriteError(ctx, res, status.SafeError(http.StatusUnsupportedMediaType, err))
			return
		}

		resMedia, err := NewMediaFromAccept(req, sm)
		if err == nil && resMedia.NewEncoder == nil {
			err = ErrUnsupportedMedia
		}

		if err != nil {
			_ = status.WriteError(ctx, res, status.SafeError(http.StatusNotAcceptable, err))
			return
		}

		ctx = meta.WithContent(ctx, req, res, nil)
		res.Header().Set(content.TypeKey, resMedia.WithUTF8())
		res.Header().Set(compress.HeaderNoCompression, "1")

		buffer := cont.BorrowBuffer()
		defer cont.ReturnBuffer(buffer)

		writer := &commitWriter{res: res, buffer: buffer}
		controller := http.NewResponseController(res)
		ctx, cancel := withDrain(ctx, opts.Drain, func() {
			_ = req.Body.Close()
		})
		if cancel != nil {
			defer cancel(nil)
		}

		capped := budget.NewReader(req.Body, opts.MaxReceiveSize.Bytes())
		st := &RequestStream[Req, Res]{
			Stream: Stream[Res]{
				ctx:          ctx,
				writer:       writer,
				encoder:      resMedia.NewEncoder(writer),
				controller:   controller,
				readTimeout:  opts.ReadTimeout,
				writeTimeout: opts.WriteTimeout,
				limiter:      limiter,
			},
			decoder:        reqMedia.NewDecoder(capped),
			capped:         capped,
			maxReceiveSize: opts.MaxReceiveSize.Bytes(),
		}
		finalizeStream(ctx, res, st, handler(ctx, st))
	}
}

// RequestHandler handles a bidirectional stream: both the request and the response stream.
//
// The handler must propagate a [Stream.Send] error the same way [Handler] requires; the same
// obligation applies to a [RequestStream.Recv] error other than the terminal io.EOF, since Recv errors
// (an oversized value, a limiter denial, a decode failure) are likewise only acted on when the handler
// returns them. When configured [Options.Drain] starts, ctx is canceled with [ErrDraining] and
// an active Recv returns [ErrDraining].
type RequestHandler[Req any, Res any] func(ctx context.Context, stream *RequestStream[Req, Res]) error

// closer is satisfied by both [*Stream] and [*RequestStream] (which overrides close to also close its
// decoder), letting [finalizeStream] handle both handler shapes with one implementation.
type closer interface {
	committed() bool
	draining() bool
	close() error
}

// finalizeStream applies the streaming error contract after a stream handler returns.
//
// If the response was never committed, a handler error is rendered as an ordinary HTTP error via
// [status.WriteError] (a nil error leaves an empty, implicitly-200 response). If the response was
// committed, any handler error — or a failure finalizing the encoder on the success path — is
// recorded for operator diagnostics and the response is aborted via panic([http.ErrAbortHandler]),
// which [transport/http.recoveryHandler] and the access logger already special-case for a committed
// response.
func finalizeStream(ctx context.Context, res http.ResponseWriter, s closer, err error) {
	if s.draining() {
		if !s.committed() {
			_ = status.WriteError(ctx, res, status.SafeError(http.StatusServiceUnavailable, ErrDraining))

			return
		}

		err = s.close()
		if err == nil {
			return
		}
	}

	if !s.committed() {
		if err != nil {
			_ = status.WriteError(ctx, res, err)
		}

		return
	}

	if err == nil {
		err = s.close()
	}

	if err != nil {
		status.RecordError(ctx, err)
		span := tracer.SpanFromContext(ctx)
		span.RecordError(err)
		span.SetStatus(tracer.StatusCodeError, err.Error())
		panic(http.ErrAbortHandler)
	}
}
