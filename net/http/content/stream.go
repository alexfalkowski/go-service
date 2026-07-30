package content

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/encoding/stream"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/meta"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/ptr"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/klauspost/compress/gzhttp"
)

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
// using [Content.NewStreamFromAccept]. Unlike single-value handlers, an unregistered or unparseable
// media type is rejected with 415 rather than falling back to JSON.
//
// Error contract:
// A handler error returned before the first successful Send is an ordinary HTTP error rendered
// through [status.WriteError]. A handler error returned after the first successful Send aborts the
// response via panic([http.ErrAbortHandler]) instead, since the response is already committed — see
// [Stream.Send] and [finalizeStream].
//
// timeout, when positive, is pushed forward as the response write deadline after every successful
// Send (see [http.ResponseController.SetWriteDeadline]), turning a whole-stream write timeout into a
// per-message inactivity budget. Pass zero to disable this.
//
// Load control:
// If the request context carries a limiter (see [meta.WithLimiter], populated by
// [github.com/alexfalkowski/go-service/v2/transport/http/limiter.Handler] for a request that was already
// allowed to open the stream), every Send charges one additional token against that same limiter — see
// [Stream.Send]. A route reached without that middleware, or with no limiter configured, sees no
// per-message charging.
func NewStreamHandler[Res any](cont *Content, timeout time.Duration, handler StreamHandler[Res]) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		limiter := meta.Limiter(ctx)

		resMedia, err := cont.NewStreamFromAccept(req)
		if err == nil && resMedia.NewEncoder == nil {
			err = ErrUnsupportedStreamMedia
		}

		if err != nil {
			_ = status.WriteError(ctx, res, status.SafeError(http.StatusUnsupportedMediaType, err))
			return
		}

		ctx = meta.WithContent(ctx, req, res, nil)
		res.Header().Set(TypeKey, resMedia.WithUTF8())
		res.Header().Set(gzhttp.HeaderNoCompression, "1")

		buffer := cont.pool.Get()
		defer cont.pool.Put(buffer)

		writer := &commitWriter{res: res, buffer: buffer}
		st := &Stream[Res]{
			ctx:        ctx,
			writer:     writer,
			encoder:    resMedia.NewEncoder(writer),
			controller: http.NewResponseController(res),
			timeout:    timeout,
			limiter:    limiter,
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
// The request decoder is resolved from Content-Type via [Content.NewStreamFromContentType]; the
// response encoder is resolved from Accept (falling back to Content-Type) via
// [Content.NewStreamFromAccept]. Both reject an unregistered or unparseable media type with 415.
//
// Inbound size limiting:
// maxReceiveSize bounds each value decoded by [RequestStream.Recv], not the request stream as a whole
// (see [capReader]): a request with many small values is never rejected for its cumulative size, only
// for a single value that exceeds maxReceiveSize. A value at or under the limit is a normal terminal
// [http.MaxBytesError] surfaced from Recv; the deviation is that no total byte ceiling exists for a
// streaming route the way [github.com/alexfalkowski/go-service/v2/net/http/body.NewHandler]'s buffered
// path enforces one. maxReceiveSize <= 0 disables the per-value cap entirely.
//
// Timeout:
// When timeout is positive, every successful [RequestStream.Recv] extends the request read deadline
// and every successful [Stream.Send] extends the response write deadline. This turns the server's
// whole-request read and write timeouts into per-message inactivity budgets. Pass zero to disable this.
//
// See [NewStreamHandler] for the error contract and load control: Recv charges one token per successfully
// decoded, in-cap value, the same way Send does, against the same request-scoped limiter.
func NewRequestStreamHandler[Req any, Res any](cont *Content, timeout time.Duration, maxReceiveSize bytes.Size, handler RequestStreamHandler[Req, Res]) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		limiter := meta.Limiter(ctx)

		if req.ProtoMajor < 2 {
			_ = status.WriteError(ctx, res, status.SafeError(http.StatusHTTPVersionNotSupported, ErrBidiRequiresHTTP2))
			return
		}

		reqMedia, err := cont.NewStreamFromContentType(req)
		if err == nil && reqMedia.NewDecoder == nil {
			err = ErrUnsupportedStreamMedia
		}

		if err != nil {
			_ = status.WriteError(ctx, res, status.SafeError(http.StatusUnsupportedMediaType, err))
			return
		}

		resMedia, err := cont.NewStreamFromAccept(req)
		if err == nil && resMedia.NewEncoder == nil {
			err = ErrUnsupportedStreamMedia
		}

		if err != nil {
			_ = status.WriteError(ctx, res, status.SafeError(http.StatusUnsupportedMediaType, err))
			return
		}

		ctx = meta.WithContent(ctx, req, res, nil)
		res.Header().Set(TypeKey, resMedia.WithUTF8())
		res.Header().Set(gzhttp.HeaderNoCompression, "1")

		buffer := cont.pool.Get()
		defer cont.pool.Put(buffer)

		writer := &commitWriter{res: res, buffer: buffer}
		capped := &capReader{r: req.Body, max: maxReceiveSize.Bytes()}
		st := &RequestStream[Req, Res]{
			Stream: Stream[Res]{
				ctx:        ctx,
				writer:     writer,
				encoder:    resMedia.NewEncoder(writer),
				controller: http.NewResponseController(res),
				timeout:    timeout,
				limiter:    limiter,
			},
			decoder: reqMedia.NewDecoder(capped),
			capped:  capped,
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
	timeout    time.Duration
}

// Send encodes res and flushes it to the client.
//
// The first successful Send commits the response: before that, the encoded value is held in a scratch
// buffer so a failed first encode never writes a partial success body. Send is sticky: once it returns
// a non-nil error, every later call returns that same error immediately, so a handler that ignores the
// error degrades to a no-op rather than an infinite producer.
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

	if err := s.extendDeadline(); err != nil {
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

func (s *Stream[Res]) extendDeadline() error {
	if s.timeout <= 0 {
		return nil
	}

	return s.controller.SetWriteDeadline(time.Now().Add(s.timeout.Duration()))
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
	capped  *capReader
	Stream[Res]
}

// Recv decodes the next request value.
//
// Recv returns io.EOF once the client half-closes the request stream, matching the gRPC idiom and the
// terminal behavior of the underlying stream decoders.
//
// Recv is bounded independently by the per-value cap configured on [NewRequestStreamHandler]: the cap
// resets before every Recv rather than accumulating across the whole request stream (see [capReader]).
// A value over the cap surfaces as a [http.MaxBytesError], the same error type
// [github.com/alexfalkowski/go-service/v2/net/http/body.NewHandler]'s buffered path already produces,
// so it maps to the same 413 through [github.com/alexfalkowski/go-service/v2/net/http/status.Code].
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
	s.capped.resetValue(bufferedLen(s.decoder))

	req := ptr.Zero[Req]()
	if err := s.decoder.Decode(req); err != nil {
		if s.capped.err != nil {
			return nil, s.capped.err
		}

		return nil, err
	}

	if s.capped.exceededBy(bufferedLen(s.decoder)) {
		return nil, s.capped.exceededError()
	}

	if err := s.takeLimiterToken(); err != nil {
		return nil, err
	}

	if err := s.extendReadDeadline(); err != nil {
		return nil, err
	}

	return req, nil
}

// IsFinished reports whether err is the terminal [io.EOF] Recv returns once the client half-closes the
// request stream, letting a handler's Recv loop end cleanly without importing errors/io itself.
func (s *RequestStream[Req, Res]) IsFinished(err error) bool {
	return errors.Is(err, io.EOF)
}

func (s *RequestStream[Req, Res]) extendReadDeadline() error {
	if s.timeout <= 0 {
		return nil
	}

	return s.controller.SetReadDeadline(time.Now().Add(s.timeout.Duration()))
}

func (s *RequestStream[Req, Res]) close() error {
	err := s.Stream.close()

	if decErr := s.decoder.Close(); err == nil {
		err = decErr
	}

	return err
}

// capReader wraps an [io.Reader] with a byte counter that [RequestStream.Recv] resets before every
// decoded value, so a stream decoder bound to it for the whole request body (see [stream.Decoder]) gets
// a per-value cap rather than a cumulative one. This is the inbound mirror of
// [github.com/alexfalkowski/go-service/v2/net/http/client.capReader] — same reset-before/check-after
// sequencing around Decode — placed in this package rather than shared with the client because the
// reset call has to happen at the exact point a decoded value boundary is known, and RequestStream is
// that point on the server side just as [github.com/alexfalkowski/go-service/v2/net/http/client.ResponseStream]
// is on the client side.
//
// capReader never truncates or discards a Read call's data, for the same reason documented on the
// client-side type: every byte pulled from the underlying request body is always returned to the
// decoder, so the decoder's own internal buffering can never be desynchronized by this guard. Read still
// refuses to pull more data once the budget from a previous Read within the same value is already
// spent, bounding a single pathological value across repeated Read calls; Recv's post-Decode exceeded
// check additionally catches a value whose entire content arrived in one Read call that already
// exceeded the budget.
//
// Read-ahead correction: a decoder such as [encoding/json.Decoder] fills its own internal buffer in
// chunks, so one Decode call can pull bytes belonging to a later value out of the underlying reader
// (confirmed by a direct repro: decoding the first of two small NDJSON lines read the whole remaining
// buffer, attributing both lines' bytes to the first value and rejecting it as oversized even though
// neither line alone was near the cap). [Recv] corrects for this via [bufferedLen], which subtracts
// whatever the decoder still holds unconsumed in its own buffer (see [stream.Decoder]'s wrapped
// [encoding/json.Decoder.Buffered]) from the raw bytes capReader counted for this call, so only bytes
// actually consumed to produce the value just decoded count against its cap. This works for every kind
// reachable from content negotiation today (only "json"/NDJSON is registered in [streamKinds]); a kind
// whose decoder does not expose a compatible Buffered method gets no correction, and falls back to the
// same "bound on reads attributed to one value" precision this type always documented.
//
// One deliberate difference from the client-side type: max <= 0 disables the cap outright (Read never
// refuses to read and exceeded always reports false), rather than assuming a caller-enforced positive
// default. The client's capReader may assume [github.com/alexfalkowski/go-service/v2/net/http/client.Client]
// always resolves a positive maxResponseSize before construction; the inbound side is wired through DI
// from [config/server.Config.GetMaxReceiveSize], which can be projected to zero when the HTTP transport
// itself is disabled (see [github.com/alexfalkowski/go-service/v2/transport/http.maxReceiveSize]), and a
// disabled transport has no request to cap.
type capReader struct {
	r    io.Reader
	err  error
	max  int64
	read int64
}

// resetValue rearms the byte counter for the next decoded value, clearing any size-limit error latched
// by the previous value's Read calls. bufferedAhead has already been read from the underlying stream by
// the decoder while handling the previous value, so it starts this value's accounting rather than being
// discarded when the counter resets.
func (r *capReader) resetValue(bufferedAhead int64) {
	r.read = bufferedAhead
	r.err = nil
}

// exceededBy reports whether the bytes read since the last resetValue, minus bufferedAhead, exceed the
// configured cap. bufferedAhead is the decoder's own read-ahead for a later value (see [bufferedLen]);
// subtracting it corrects for a decoder that pulled more than one value's worth of bytes out of the
// reader in a single underlying Read. The corrected count is clamped to zero, and this always reports
// false when the cap is disabled (max <= 0). This is the post-Decode check; Read below applies its own
// uncorrected, live check to bound a single pathological value across repeated Read calls.
func (r *capReader) exceededBy(bufferedAhead int64) bool {
	if r.max <= 0 {
		return false
	}

	attributed := max(r.read-bufferedAhead, 0)

	return attributed > r.max
}

// Read implements [io.Reader]. Once the bytes read since the last resetValue already reach the
// configured cap, Read refuses to read more and returns exceededError, and every subsequent call
// returns that same error until resetValue is called. Read never refuses to read when the cap is
// disabled (max <= 0).
func (r *capReader) Read(p []byte) (int, error) {
	if r.err != nil {
		return 0, r.err
	}

	if r.max > 0 && r.read >= r.max {
		r.err = r.exceededError()
		return 0, r.err
	}

	n, err := r.r.Read(p)
	r.read += int64(n)

	return n, err
}

// exceededError returns the error surfaced for a value that exceeds the configured cap: a raw
// [http.MaxBytesError], matching what [github.com/alexfalkowski/go-service/v2/net/http/body.NewHandler]'s
// buffered path already produces for an oversized request body, so both paths resolve to the same 413
// through [github.com/alexfalkowski/go-service/v2/net/http/status.Code].
func (r *capReader) exceededError() error {
	return &http.MaxBytesError{Limit: r.max}
}

// bufferedLen returns the number of bytes decoder has already pulled from its underlying reader for a
// value it has not decoded yet, so a caller can subtract them from the raw bytes [capReader] counted
// for the value it just decoded (see [capReader.exceededBy]).
//
// This only recognizes decoders whose Buffered method matches [encoding/json.Decoder.Buffered]'s shape
// and returns a concrete [bytes.Reader] — true today for every [stream.Decoder] reachable from content
// negotiation, since only the "json" kind (NDJSON) is registered in [streamKinds]. It returns 0 for any
// decoder that does not match, which is always safe: the correction is skipped, not a wrong non-zero
// answer, so the check falls back to the same "bound on reads attributed to one value" behavior
// documented on [capReader] rather than any incorrect result.
func bufferedLen(decoder stream.Decoder) int64 {
	buffered, ok := decoder.(interface{ Buffered() io.Reader })
	if !ok {
		return 0
	}

	reader, ok := buffered.Buffered().(*bytes.Reader)
	if !ok {
		return 0
	}

	return int64(reader.Len())
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
