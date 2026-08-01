package stream

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/encoding/stream"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/budget"
	"github.com/alexfalkowski/go-service/v2/net/http/meta"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/ptr"
	"github.com/alexfalkowski/go-service/v2/time"
)

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
	// readTimeout is zero for a send-only [Stream] built by [NewHandler]; [NewRequestHandler]
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
// Send always extends the write deadline (see [NewHandler]/[NewRequestHandler]'s timeout
// semantics). On a bidirectional stream built by [NewRequestHandler], Send also extends the
// request read deadline, the same way [RequestStream.Recv] extends the write deadline: activity in
// either direction proves the peer is alive, so both deadlines move forward together and only genuine
// inactivity in both directions severs the stream. This is a no-op for a send-only [Stream].
//
// Load control:
// If a limiter was captured at construction (see [NewHandler]), Send charges one limiter token
// before doing any work, mirroring gRPC's per-SendMsg charge (see
// [github.com/alexfalkowski/go-service/v2/transport/grpc/limiter.serverStream.SendMsg]). Before the first
// successful Send this surfaces as an ordinary 429 [status.Error] through the pre-commit error path;
// after commit it aborts the response the same way any other post-commit Send error does,
// since a mid-stream 429 cannot be written once headers are sent. A denial or a limiter error is sticky
// like any other Send error.
//
// If the configured drain signal has started, Send returns [ErrDraining] without sending another value.
func (s *Stream[Res]) Send(res *Res) error {
	if s.err != nil {
		return s.err
	}

	if s.draining() {
		s.err = ErrDraining

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
// by [NewHandler] (readTimeout stays zero there); [RequestStream.Recv] uses it directly, and
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

func (s *Stream[Res]) draining() bool {
	return errors.Is(context.Cause(s.ctx), ErrDraining)
}

func (s *Stream[Res]) close() error {
	return s.encoder.Close()
}

// requestEncoder keeps bidirectional stream finalization on the response stream while closing both codecs.
type requestEncoder struct {
	stream.Encoder
	decoder stream.Decoder
}

func (e *requestEncoder) Close() error {
	if err := e.Encoder.Close(); err != nil {
		_ = e.decoder.Close()

		return err
	}

	return e.decoder.Close()
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
// Recv is bounded independently by the per-value cap configured on [NewRequestHandler]: the cap
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
//
// If the configured drain signal starts, Recv returns [ErrDraining]. An active Recv is unblocked by
// closing the request body; this can be observed as a stream reset by an HTTP/2 client, which should
// reconnect to a non-draining server.
func (s *RequestStream[Req, Res]) Recv() (*Req, error) {
	req, err := s.decode()
	if err != nil {
		return nil, err
	}

	if s.draining() {
		return nil, ErrDraining
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

func (s *RequestStream[Req, Res]) decode() (*Req, error) {
	if s.draining() {
		return nil, ErrDraining
	}

	if err := s.extendReadDeadline(); err != nil {
		return nil, err
	}

	s.capped.Reset(budget.BufferedLen(s.decoder))

	req := ptr.Zero[Req]()
	if err := s.decoder.Decode(req); err != nil {
		if s.draining() {
			return nil, ErrDraining
		}

		if s.capped.Err() != nil {
			return nil, &http.MaxBytesError{Limit: s.maxReceiveSize}
		}

		return nil, err
	}

	if s.capped.Exceeds(budget.BufferedLen(s.decoder)) {
		return nil, &http.MaxBytesError{Limit: s.maxReceiveSize}
	}

	return req, nil
}

// IsFinished reports whether err is the terminal [io.EOF] Recv returns once the client half-closes the
// request stream, letting a handler's Recv loop end cleanly without importing errors/io itself.
func (s *RequestStream[Req, Res]) IsFinished(err error) bool {
	return errors.Is(err, io.EOF)
}
