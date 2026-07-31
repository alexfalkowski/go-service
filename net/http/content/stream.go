package content

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/encoding/stream"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/budget"
	"github.com/alexfalkowski/go-service/v2/net/http/compress"
	"github.com/alexfalkowski/go-service/v2/net/http/meta"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/ptr"
	"github.com/alexfalkowski/go-service/v2/time"
)

// StreamOptions configures the streaming route helpers built by [NewStreamHandler] and
// [NewRequestStreamHandler].
//
// The zero value disables every knob: no deadline extension and no per-value receive cap, matching
// today's behavior when no options are supplied.
type StreamOptions struct {
	// ReadTimeout is the per-message inactivity budget applied to [RequestStream.Recv]'s request read
	// deadline. Zero disables read deadline extension. Unused by [NewStreamHandler]'s send-only stream.
	ReadTimeout time.Duration

	// WriteTimeout is the per-message inactivity budget applied to [Stream.Send]'s response write
	// deadline. Zero disables write deadline extension.
	WriteTimeout time.Duration

	// MaxReceiveSize bounds each value decoded by [RequestStream.Recv], not the request stream as a
	// whole (see [github.com/alexfalkowski/go-service/v2/net/http/budget.Reader]). Zero disables the
	// per-value cap. Unused by [NewStreamHandler]'s send-only stream.
	MaxReceiveSize bytes.Size
}

// StreamHandler handles a send-only stream: the response streams, the request does not.
//
// The handler must propagate a [Stream.Send] error by returning it (directly, or wrapped) rather than
// ignoring it and continuing. Send is sticky — once it fails, every later call returns the same error
// immediately — but only the handler returning that error triggers [finalizeStream]'s abort/HTTP-error
// handling; a handler that swallows a Send error and returns nil degrades to a no-op loop over a dead
// connection instead of ending the request.
type StreamHandler[Res any] func(ctx context.Context, stream *Stream[Res]) error

// RequestStreamHandler handles a bidirectional stream: both the request and the response stream.
//
// The handler must propagate a [Stream.Send] error the same way [StreamHandler] requires; the same
// obligation applies to a [RequestStream.Recv] error other than the terminal io.EOF, since Recv errors
// (an oversized value, a limiter denial, a decode failure) are likewise only acted on when the handler
// returns them.
type RequestStreamHandler[Req any, Res any] func(ctx context.Context, stream *RequestStream[Req, Res]) error

// NewStreamHandler builds a handler for a send-only streaming response.
//
// Content negotiation:
// The response encoder is resolved from the request Accept header, falling back to Content-Type,
// using [Content.NewStreamFromAccept]. Unlike single-value handlers, an Accept that cannot be satisfied
// is rejected with 406 rather than falling back to JSON.
//
// Error contract:
// A handler error returned before the first successful Send is an ordinary HTTP error rendered
// through [status.WriteError]. A handler error returned after the first successful Send aborts the
// response via panic([http.ErrAbortHandler]) instead, since the response is already committed — see
// [Stream.Send] and [finalizeStream].
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
func NewStreamHandler[Res any](cont *Content, opts StreamOptions, handler StreamHandler[Res]) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		limiter := meta.Limiter(ctx)

		resMedia, err := cont.NewStreamFromAccept(req)
		if err == nil && resMedia.NewEncoder == nil {
			err = ErrUnsupportedStreamMedia
		}

		if err != nil {
			_ = status.WriteError(ctx, res, status.SafeError(http.StatusNotAcceptable, err))
			return
		}

		ctx = meta.WithContent(ctx, req, res, nil)
		res.Header().Set(TypeKey, resMedia.WithUTF8())
		res.Header().Set(compress.HeaderNoCompression, "1")

		buffer := cont.pool.Get()
		defer cont.pool.Put(buffer)

		writer := &commitWriter{res: res, buffer: buffer}
		st := &Stream[Res]{
			ctx:          ctx,
			writer:       writer,
			encoder:      resMedia.NewEncoder(writer),
			controller:   http.NewResponseController(res),
			writeTimeout: opts.WriteTimeout,
			limiter:      limiter,
		}
		finalizeStream(ctx, res, st, handler(ctx, st))
	}
}

// NewRequestStreamHandler builds a handler for a bidirectional stream.
//
// HTTP/2 requirement:
// Bidirectional streaming requires HTTP/2 (including h2c): an HTTP/1.x request body is buffered ahead
// of the handler by intermediaries and the Go transport, so interleaving Recv/Send hangs rather than
// failing. Requests with req.ProtoMajor < 2 are rejected with 505 before the handler runs.
//
// Content negotiation:
// The request decoder is resolved from Content-Type via [Content.NewStreamFromContentType], rejecting
// an unregistered or unparseable media type with 415. The response encoder is resolved from Accept
// (falling back to Content-Type) via [Content.NewStreamFromAccept], rejecting an Accept that cannot be
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
// See [NewStreamHandler] for the error contract and load control: Recv charges one token per successfully
// decoded, in-cap value, the same way Send does, against the same request-scoped limiter.
func NewRequestStreamHandler[Req any, Res any](cont *Content, opts StreamOptions, handler RequestStreamHandler[Req, Res]) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		limiter := meta.Limiter(ctx)

		if req.ProtoMajor < 2 {
			_ = status.WriteError(ctx, res, status.SafeError(http.StatusHTTPVersionNotSupported, ErrBidiRequiresHTTP2))
			return
		}

		reqMedia, err := cont.NewStreamFromContentType(req)
		if err != nil {
			_ = status.WriteError(ctx, res, status.SafeError(http.StatusUnsupportedMediaType, err))
			return
		}

		resMedia, err := cont.NewStreamFromAccept(req)
		if err == nil && resMedia.NewEncoder == nil {
			err = ErrUnsupportedStreamMedia
		}

		if err != nil {
			_ = status.WriteError(ctx, res, status.SafeError(http.StatusNotAcceptable, err))
			return
		}

		ctx = meta.WithContent(ctx, req, res, nil)
		res.Header().Set(TypeKey, resMedia.WithUTF8())
		res.Header().Set(compress.HeaderNoCompression, "1")

		buffer := cont.pool.Get()
		defer cont.pool.Put(buffer)

		writer := &commitWriter{res: res, buffer: buffer}
		capped := budget.NewReader(req.Body, opts.MaxReceiveSize.Bytes())
		st := &RequestStream[Req, Res]{
			Stream: Stream[Res]{
				ctx:          ctx,
				writer:       writer,
				encoder:      resMedia.NewEncoder(writer),
				controller:   http.NewResponseController(res),
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

// Stream is a send-only streaming response handle.
//
// Stream is not safe for concurrent use: Send is expected to be called from one goroutine at a time,
// matching the gRPC server-streaming convention.
type Stream[Res any] struct {
	ctx        context.Context
	encoder    stream.Encoder
	err        error
	writer     *commitWriter
	controller *http.ResponseController
	limiter    meta.RateLimiter
	// readTimeout is zero for a send-only [Stream] built by [NewStreamHandler]; [NewRequestStreamHandler]
	// sets it, letting Send extend both deadlines for a bidirectional stream (see [Stream.Send]).
	readTimeout  time.Duration
	writeTimeout time.Duration
}

// Send encodes res and flushes it to the client.
//
// The first successful Send commits the response: before that, the encoded value is held in a scratch
// buffer so a failed first encode never writes a partial success body. Send is sticky: once it returns
// a non-nil error, every later call returns that same error immediately, so a handler that ignores the
// error degrades to a no-op rather than an infinite producer.
//
// Deadlines:
// Send always extends the write deadline (see [NewStreamHandler]/[NewRequestStreamHandler]'s timeout
// semantics). On a bidirectional stream built by [NewRequestStreamHandler], Send also extends the
// request read deadline, the same way [RequestStream.Recv] extends the write deadline: activity in
// either direction proves the peer is alive, so both deadlines move forward together and only genuine
// inactivity in both directions severs the stream. This is a no-op for a send-only [Stream].
//
// Load control:
// If a limiter was captured at construction (see [NewStreamHandler]), Send charges one limiter token
// before doing any work, mirroring gRPC's per-SendMsg charge (see
// [github.com/alexfalkowski/go-service/v2/transport/grpc/limiter.serverStream.SendMsg]). Before the first
// successful Send this surfaces as an ordinary 429 [status.Error] through the pre-commit error path;
// after commit it aborts the response the same way any other post-commit Send error does,
// since a mid-stream 429 cannot be written once headers are sent. A denial or a limiter error is sticky
// like any other Send error.
func (s *Stream[Res]) Send(res *Res) error {
	if s.err != nil {
		return s.err
	}

	if err := s.takeLimiterToken(); err != nil {
		s.err = err
		return err
	}

	if err := s.extendWriteDeadline(); err != nil {
		s.err = err
		return err
	}

	if err := s.extendReadDeadline(); err != nil {
		s.err = err
		return err
	}

	if err := s.encoder.Encode(res); err != nil {
		s.err = err
		return err
	}

	if err := s.writer.commit(); err != nil {
		s.err = err
		return err
	}

	if err := s.controller.Flush(); err != nil {
		s.err = err
		return err
	}

	return nil
}

// takeLimiterToken charges one limiter token for one Send or Recv message, on top of the token already
// charged when the stream opened (see [github.com/alexfalkowski/go-service/v2/transport/http/limiter.Handler.ServeHTTP]).
//
// It is a no-op, with no overhead beyond the nil check, when no limiter was captured at construction — the
// common case, since most routes have no limiter configured, or reach a stream handler outside that
// middleware. A denial maps to the same 429 [status.Error] the limiter middleware itself returns for a
// rejected request; a Take failure maps to the same [status.InternalServerError] the middleware returns
// for a limiter error.
//
// The RateLimit and RateLimit-Policy response headers are not re-emitted here: they describe only the
// stream-open decision, since HTTP response headers cannot change once a streaming response is committed.
func (s *Stream[Res]) takeLimiterToken() error {
	if s.limiter == nil {
		return nil
	}

	allowed, _, err := s.limiter.Take(s.ctx)
	if err != nil {
		return status.InternalServerError(err)
	}

	if !allowed {
		return status.Error(http.StatusTooManyRequests, http.StatusText(http.StatusTooManyRequests))
	}

	return nil
}

func (s *Stream[Res]) extendWriteDeadline() error {
	if s.writeTimeout <= 0 {
		return nil
	}

	return s.controller.SetWriteDeadline(time.Now().Add(s.writeTimeout.Duration()))
}

// extendReadDeadline extends the request read deadline. It is a no-op for a send-only [Stream] built
// by [NewStreamHandler] (readTimeout stays zero there); [RequestStream.Recv] uses it directly, and
// [Stream.Send] uses it too so a bidirectional stream's write-side activity also keeps the read side
// alive (see [RequestStream]).
func (s *Stream[Res]) extendReadDeadline() error {
	if s.readTimeout <= 0 {
		return nil
	}

	return s.controller.SetReadDeadline(time.Now().Add(s.readTimeout.Duration()))
}

func (s *Stream[Res]) committed() bool {
	return s.writer.committed
}

func (s *Stream[Res]) close() error {
	return s.encoder.Close()
}

// RequestStream is a bidirectional streaming handle: both the request and the response stream.
//
// RequestStream is not safe for arbitrary concurrent use. The supported pattern, matching gRPC's bidi
// streaming convention, is one goroutine calling Recv and one goroutine calling Send.
type RequestStream[Req any, Res any] struct {
	decoder stream.Decoder
	capped  *budget.Reader
	Stream[Res]
	maxReceiveSize int64
}

// Recv decodes the next request value.
//
// Recv returns io.EOF once the client half-closes the request stream, matching the gRPC idiom and the
// terminal behavior of the underlying stream decoders.
//
// Recv is bounded independently by the per-value cap configured on [NewRequestStreamHandler]: the cap
// resets before every Recv rather than accumulating across the whole request stream (see
// [github.com/alexfalkowski/go-service/v2/net/http/budget.Reader]). A value over the cap surfaces as a
// [http.MaxBytesError], the same error type
// [github.com/alexfalkowski/go-service/v2/net/http/body.NewHandler]'s buffered path already produces,
// so it maps to the same 413 through [github.com/alexfalkowski/go-service/v2/net/http/status.Code].
//
// Deadlines:
// Recv extends the read deadline before attempting Decode, so the wait for the next value — including
// the first — is bounded by the configured timeout rather than only the calls after it. A successful
// Recv extends the read deadline again (refreshing the wait for the value after this one) and also
// extends the write deadline, the same way [Stream.Send] extends the read deadline: activity in either
// direction proves the peer is alive, so both deadlines move forward together and only genuine
// inactivity in both directions severs the stream.
//
// Load control:
// If a limiter was captured at construction, Recv charges one limiter token after a value decodes
// successfully and within the per-value cap, mirroring gRPC's per-RecvMsg charge (see
// [github.com/alexfalkowski/go-service/v2/transport/grpc/limiter.serverStream.RecvMsg]) — a stream-end
// io.EOF or an oversized value is not a received message and is never charged. See [Stream.Send] for how
// a denial or limiter error is mapped and surfaced.
//
// A decoder is not guaranteed to return capped's Read-time error unwrapped: it may fold it into its own
// error value instead (observed with the standard library JSON decoder under GOEXPERIMENT=jsonv2, which
// wraps a mid-scan Read error into a *json.SyntaxError with no Unwrap method). Recv therefore checks
// capped's own latched error directly once Decode fails, rather than trusting Decode's returned error
// to still be classifiable as the size-limit error.
func (s *RequestStream[Req, Res]) Recv() (*Req, error) {
	if err := s.extendReadDeadline(); err != nil {
		return nil, err
	}

	s.capped.Reset(budget.BufferedLen(s.decoder))

	req := ptr.Zero[Req]()
	if err := s.decoder.Decode(req); err != nil {
		if s.capped.Err() != nil {
			return nil, &http.MaxBytesError{Limit: s.maxReceiveSize}
		}

		return nil, err
	}

	if s.capped.Exceeds(budget.BufferedLen(s.decoder)) {
		return nil, &http.MaxBytesError{Limit: s.maxReceiveSize}
	}

	if err := s.takeLimiterToken(); err != nil {
		return nil, err
	}

	if err := s.extendReadDeadline(); err != nil {
		return nil, err
	}

	if err := s.extendWriteDeadline(); err != nil {
		return nil, err
	}

	return req, nil
}

// IsFinished reports whether err is the terminal [io.EOF] Recv returns once the client half-closes the
// request stream, letting a handler's Recv loop end cleanly without importing errors/io itself.
func (s *RequestStream[Req, Res]) IsFinished(err error) bool {
	return errors.Is(err, io.EOF)
}

func (s *RequestStream[Req, Res]) close() error {
	err := s.Stream.close()

	if decErr := s.decoder.Close(); err == nil {
		err = decErr
	}

	return err
}

// commitWriter buffers writes until commit is called, after which writes go straight to the live
// response writer.
//
// Buffering only the first value, rather than swapping the encoder to a new writer after commit, keeps
// one encoder instance bound to one writer for the life of the stream — required for codecs (yaml) that
// carry document-separator state tied to their writer.
type commitWriter struct {
	res       http.ResponseWriter
	buffer    *bytes.Buffer
	committed bool
}

// Write implements [io.Writer]. Before commit, writes accumulate in buffer; after commit, they go
// straight to res.
func (w *commitWriter) Write(p []byte) (int, error) {
	if w.committed {
		return w.res.Write(p)
	}

	return w.buffer.Write(p)
}

// commit flushes any buffered bytes to the live response writer. It is a no-op once already committed.
func (w *commitWriter) commit() error {
	if w.committed {
		return nil
	}

	w.committed = true
	_, err := w.buffer.WriteTo(w.res)

	return err
}

// closer is satisfied by both [*Stream] and [*RequestStream] (which overrides close to also close its
// decoder), letting [finalizeStream] handle both handler shapes with one implementation.
type closer interface {
	committed() bool
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
		panic(http.ErrAbortHandler)
	}
}
