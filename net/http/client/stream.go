package client

import (
	"github.com/alexfalkowski/go-service/v2/context"
	encodingstream "github.com/alexfalkowski/go-service/v2/encoding/stream"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/budget"
	contentstream "github.com/alexfalkowski/go-service/v2/net/http/content/stream"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// StreamHandler handles a client-side receive-only stream: the response streams, the request does
// not. It is the client-side counterpart of a call to a send-only streaming route (see
// [github.com/alexfalkowski/go-service/v2/net/http/content/stream.Handler]).
//
// Like [Options.Response], a streamed value has no static Go type known to this package: handle it by
// calling [ResponseStream.Recv] with a pointer to whatever type the caller expects, the same way
// [Client.Do] decodes into opts.Response.
type StreamHandler func(ctx context.Context, stream *ResponseStream) error

// RequestStreamHandler handles a client-side bidirectional stream: both the request and the response
// stream. It is the client-side counterpart of a call to a bidirectional streaming route (see
// [github.com/alexfalkowski/go-service/v2/net/http/content/stream.RequestHandler]).
type RequestStreamHandler func(ctx context.Context, stream *RequestResponseStream) error

// StreamGet issues a send-only streaming HTTP GET request to url using opts.
//
// It is a convenience wrapper around Stream.
func (c *Client) StreamGet(ctx context.Context, url string, opts Options, handler StreamHandler) error {
	return c.Stream(ctx, http.MethodGet, url, opts, handler)
}

// StreamPost issues a bidirectional streaming HTTP POST request to url using opts.
//
// It is a convenience wrapper around RequestStream.
func (c *Client) StreamPost(ctx context.Context, url string, opts Options, handler RequestStreamHandler) error {
	return c.RequestStream(ctx, http.MethodPost, url, opts, handler)
}

// StreamPut issues a bidirectional streaming HTTP PUT request to url using opts.
//
// It is a convenience wrapper around RequestStream.
func (c *Client) StreamPut(ctx context.Context, url string, opts Options, handler RequestStreamHandler) error {
	return c.RequestStream(ctx, http.MethodPut, url, opts, handler)
}

// StreamPatch issues a bidirectional streaming HTTP PATCH request to url using opts.
//
// It is a convenience wrapper around RequestStream.
func (c *Client) StreamPatch(ctx context.Context, url string, opts Options, handler RequestStreamHandler) error {
	return c.RequestStream(ctx, http.MethodPatch, url, opts, handler)
}

// Stream issues a request with method and url and drives handler over the streamed response.
//
// Unlike [Client.Do], the response is not buffered: values are decoded and delivered to handler as
// they arrive on the wire. opts.Request/opts.ContentType are encoded into the request body exactly as
// [Client.Do] would (Stream has no streamed request body; use RequestStream for that). opts.Response is
// ignored.
//
// Content negotiation:
// The response streaming decoder is resolved from the response Content-Type header, falling back to
// opts.ContentType, via the streaming registry passed to [NewClient] (see
// [contentstream.Content.NewFromMedia]). An unregistered or unparseable streaming media type, or a
// Client whose streaming registry has no matching codec, is returned as an error and handler is never called —
// but only once the response exists: unlike a media type Stream could reject before dialing at all,
// there is no way to know the response's actual Content-Type without making the request first, so a
// misconfigured Client still pays for one round trip before failing.
//
// Error contract:
// Before handler is called, the initial response is checked exactly as [Client.Do] checks it: a
// text/error response or a 4xx/5xx status code is returned as a [github.com/alexfalkowski/go-service/v2/net/http/status.Error]
// and handler is never called (see [Client.checkResponseStatus]). Once handler is called, any error it
// returns is returned as-is by Stream: there is no separate abort signal on the client side, since the
// client is not committing an HTTP response of its own. The response body is always closed before
// Stream returns, whether handler succeeds or fails.
//
// Retries:
// Stream has no retry logic of its own, but its request is an ordinary [Client.Do]-shaped request with
// no streamed body, so a retry-capable [http.RoundTripper] configured via [WithRoundTripper] (for
// example [github.com/alexfalkowski/go-service/v2/transport/http/retry]) still governs whether the
// initial request is retried, exactly as it would for Client.Do. There is no retry once handler starts
// receiving values: a failure at that point is returned as-is, and it is the caller's responsibility to
// retry the whole call if appropriate. Contrast [Client.RequestStream], whose pipe-backed request body
// has no [http.Request.GetBody] and is therefore never retryable regardless of the configured
// RoundTripper.
//
// Size limiting:
// See [WithMaxResponseSize] for the per-value (not cumulative) response size cap applied to each
// decoded value.
//
// Timeouts:
// The underlying [http.Client] never has a [http.Client.Timeout] (see [NewClient]); bound the call
// with ctx instead.
func (c *Client) Stream(ctx context.Context, method, url string, opts Options, handler StreamHandler) error {
	request, err := c.newRequest(ctx, method, url, opts)
	if err != nil {
		return err
	}

	response, err := c.client.Do(request)
	if err != nil {
		return errors.Prefix("http: do", err)
	}

	ready := make(chan responseResult, 1)
	ready <- responseResult{response: response}

	stream := &ResponseStream{c: c, opts: opts, ready: ready}
	if err := stream.ensure(); err != nil {
		_ = stream.close()
		return err
	}

	panicValue, handlerErr := callHandler(func() error { return handler(ctx, stream) })
	closeErr := stream.close()

	if panicValue != nil {
		panic(panicValue)
	}

	if handlerErr != nil {
		return handlerErr
	}

	return closeErr
}

// RequestStream issues a request with method and url over a pipe-backed, chunked request body and
// drives handler over both directions: handler sends request values via
// [RequestResponseStream.Send] and receives response values via [RequestResponseStream.Recv].
//
// Pipe lifecycle:
// RequestStream owns the [io.Pipe] for the life of the call; handler never sees or controls it
// directly. The underlying HTTP round trip is started on a separate goroutine as soon as the pipe and
// request are built, so it can proceed concurrently with handler's Send calls — this allows genuine h2
// bidi interleaving, where the server can respond to earlier values while the client is still writing
// later ones. Once handler returns:
//   - if handler returned an error, the pipe writer is closed with that error
//     ([io.PipeWriter.CloseWithError]), so a transport still reading the request body observes a
//     failure rather than a clean end of stream;
//   - otherwise, the request stream encoder is finalized (Close), and the pipe writer is closed
//     with that finalize error if it failed, or closed cleanly (EOF) if it succeeded.
//
// In every case RequestStream then waits for the underlying response to resolve (if handler never
// called Recv, this is the first time the response is awaited) and closes the response body exactly
// once before returning, so the connection is never leaked regardless of how handler exits.
//
// req.ContentLength is set to -1 so the transport uses chunked encoding for the pipe-backed body.
//
// HTTP/2 requirement:
// This mirrors the server's bidirectional streaming requirement (see
// [github.com/alexfalkowski/go-service/v2/net/http/content/stream.NewRequestHandler]): calling a bidi
// route over HTTP/1.1 hangs rather than failing, because the h1 transport buffers the request body
// ahead of the handler. Callers must dial an HTTP/2 (including h2c) endpoint.
//
// Content negotiation:
// The request streaming encoder is resolved from opts.ContentType; the response streaming decoder is
// resolved the same way [Client.Stream] resolves it. Both require the streaming registry passed to
// [NewClient]; an unregistered or unparseable streaming media type is
// returned as an error before any request is sent.
//
// Error contract:
// Unlike [Client.Stream], RequestStream cannot check the initial response before calling handler: the
// response only exists once the server has read enough of the request stream to answer, so handler is
// always called. A text/error response or a 4xx/5xx status code (see [Client.checkResponseStatus]) is
// instead surfaced as the error returned by the first [RequestResponseStream.Recv] call, exactly like
// any other Recv error; a handler that never calls Recv never observes it as a Recv error, but
// RequestStream still resolves and closes the response before returning.
//
// Retries:
// RequestStream's request body has no [http.Request.GetBody] (it is pipe-backed), so
// [github.com/alexfalkowski/go-service/v2/transport/http/retry]'s canRetry check already excludes it —
// streaming requests are never retried. This is a confirmed consequence, not a defect: retrying would
// require replaying already-sent request values.
//
// Size limiting:
// See [WithMaxResponseSize] for the per-value (not cumulative) response size cap applied to each
// value received via Recv.
//
// Timeouts:
// See [Client.Stream]: the underlying [http.Client] never has a [http.Client.Timeout].
func (c *Client) RequestStream(ctx context.Context, method, url string, opts Options, handler RequestStreamHandler) error {
	reqMedia, err := c.streamContent.NewFromMedia(opts.ContentType)
	if err != nil {
		return errors.Prefix("http: stream media", err)
	}
	if reqMedia.NewEncoder == nil {
		return errors.Prefix("http: stream media", contentstream.ErrUnsupportedMedia)
	}

	reader, writer := io.Pipe()

	request, err := http.NewRequestWithContext(ctx, method, url, reader)
	if err != nil {
		_ = writer.Close()
		return errors.Prefix("http: new request", err)
	}

	// ContentLength = -1 forces chunked transfer encoding: the pipe has no known total length, and a
	// declared Content-Length would make the transport try to read exactly that many bytes upfront.
	request.ContentLength = -1
	request.Header.Set(http.ContentTypeKey, reqMedia.String())

	if !strings.IsEmpty(opts.Accept) {
		request.Header.Set(http.AcceptKey, opts.Accept)
	}

	ready := make(chan responseResult, 1)

	go func() {
		response, doErr := c.client.Do(request)
		ready <- responseResult{response: response, err: doErr}
	}()

	stream := &RequestResponseStream{
		ResponseStream: ResponseStream{c: c, opts: opts, ready: ready},
		encoder:        reqMedia.NewEncoder(writer),
		writer:         writer,
	}

	panicValue, handlerErr := callHandler(func() error { return handler(ctx, stream) })
	if panicValue != nil {
		// finish branches on a nil handlerErr by finalizing the encoder and closing the pipe cleanly,
		// as if the handler had succeeded. A panic is not success: force the error branch
		// (CloseWithError) so a transport still reading the request body sees a failure, not a clean
		// end of stream, matching how a returned handler error is already treated. ConvertRecover
		// preserves the panic value for the reader while marking the error as recovered.
		handlerErr = runtime.ConvertRecover(panicValue)
	}

	finishErr := stream.finish(handlerErr)

	if panicValue != nil {
		panic(panicValue)
	}

	return finishErr
}

// callHandler invokes fn, recovering a panic instead of letting it unwind immediately, so the caller can
// still run its pipe/response cleanup (which only happens after the handler returns) before re-raising
// the original panic. Without this, a panicking [StreamHandler] or [RequestStreamHandler] would skip
// [ResponseStream.close]/[RequestResponseStream.finish] entirely, leaking the response body and, for
// [Client.RequestStream], leaving its background [Client.client.Do] goroutine blocked on an unclosed
// pipe.
func callHandler(fn func() error) (panicValue any, err error) {
	defer func() {
		panicValue = recover()
	}()

	err = fn()

	return nil, err
}

// responseResult carries the outcome of resolving the HTTP response for a streaming call: either a
// response, or the error [Client.client.Do] returned instead of one.
type responseResult struct {
	response *http.Response
	err      error
}

// ResponseStream is a client-side receive-only streaming handle: it decodes values from a streamed
// response without writing a request body itself.
//
// ResponseStream is not safe for concurrent use: Recv is expected to be called from one goroutine at a
// time, matching [github.com/alexfalkowski/go-service/v2/net/http/content/stream.Stream]'s server-side
// concurrency contract.
//
// The response is resolved lazily, on the first call to Recv (or, if Recv is never called, when the
// driving [Client.Stream] or [Client.RequestStream] call finishes): this lets [Client.RequestStream]
// hand a stream to handler immediately, before the response exists, while [Client.Stream] forces
// resolution itself before calling handler so a known-bad response never reaches handler at all. See
// [ResponseStream.ensure].
type ResponseStream struct {
	err      error
	c        *Client
	ready    <-chan responseResult
	response *http.Response
	dec      *responseDecoder
	opts     Options
}

// ensure resolves the response exactly once: waiting for ready, applying the same status/media error
// check [Client.Do] applies (see [Client.checkResponseStatus]), and selecting the streaming decoder.
// Recv and close call it serially: ResponseStream is not safe for concurrent use.
func (s *ResponseStream) ensure() error {
	if s.ready == nil {
		return s.err
	}

	result := <-s.ready
	s.ready = nil
	if result.err != nil {
		s.err = errors.Prefix("http: do", result.err)
		return s.err
	}

	s.response = result.response

	if err := s.c.checkResponseStatus(s.response, s.opts); err != nil {
		s.err = err
		return err
	}

	dec, err := newResponseDecoder(s.c, s.response, s.opts)
	if err != nil {
		s.err = err
		return err
	}

	s.dec = dec

	return nil
}

// Recv decodes the next response value into v, a pointer to whatever type the caller expects — the
// same convention [Client.Do] uses for opts.Response.
//
// Recv returns io.EOF once the response stream ends, matching
// [github.com/alexfalkowski/go-service/v2/net/http/content/stream.RequestStream.Recv]'s server-side terminal
// behavior and the gRPC idiom. The first call to Recv resolves the underlying response if it has not
// been resolved yet (see [ResponseStream.ensure]), including the text/error and status code checks
// [Client.Do] applies to a non-streaming response; a response identified as an error surfaces here as
// a [github.com/alexfalkowski/go-service/v2/net/http/status.Error], not as a decode failure.
//
// Each call to Recv is bounded independently by the size configured with [WithMaxResponseSize]: the
// cap resets for every value rather than accumulating across the whole stream (see
// [WithMaxResponseSize]).
func (s *ResponseStream) Recv(v any) error {
	if err := s.ensure(); err != nil {
		return err
	}

	return s.dec.recv(v)
}

// IsFinished reports whether err is the terminal [io.EOF] Recv returns once the response stream ends,
// letting a handler's Recv loop end cleanly without importing errors/io itself. IsFinished is promoted
// to [RequestResponseStream] through the embedded ResponseStream.
func (s *ResponseStream) IsFinished(err error) bool {
	return errors.Is(err, io.EOF)
}

// close finalizes the decoder (if the response ever resolved to one) and closes the underlying
// response body. It resolves the response first if that has not happened yet, so a call that never
// invoked Recv still closes its connection instead of leaking it.
func (s *ResponseStream) close() error {
	err := s.ensure()

	if s.dec != nil {
		if closeErr := s.dec.close(); err == nil {
			err = closeErr
		}
	}

	if s.response != nil {
		_ = s.response.Body.Close()
	}

	return err
}

// RequestResponseStream is a client-side bidirectional streaming handle: it both encodes values into
// a pipe-backed request body via Send and decodes values from the response via the embedded
// [ResponseStream.Recv].
//
// RequestResponseStream is not safe for arbitrary concurrent use. The supported pattern, matching
// gRPC's bidi streaming convention and
// [github.com/alexfalkowski/go-service/v2/net/http/content/stream.RequestStream]'s server-side contract, is
// one goroutine calling Send and one goroutine calling Recv.
type RequestResponseStream struct {
	encoder encodingstream.Encoder
	sendErr error
	writer  *io.PipeWriter
	ResponseStream
}

// Send encodes v, a pointer to whatever type the caller is sending, and writes it to the pipe-backed
// request body.
//
// Send is sticky: once it returns a non-nil error, every later call returns that same error
// immediately, matching
// [github.com/alexfalkowski/go-service/v2/net/http/content/stream.Stream.Send]'s server-side contract — a
// handler that ignores the error degrades to a no-op rather than writing into a broken pipe forever.
func (s *RequestResponseStream) Send(v any) error {
	if s.sendErr != nil {
		return s.sendErr
	}

	if err := s.encoder.Encode(v); err != nil {
		s.sendErr = err
		return err
	}

	return nil
}

// finish closes the request-side pipe according to handlerErr, then always resolves and closes the
// response, returning the first of handlerErr, an encoder finalize error, or a response-resolution
// error, in that priority order.
func (s *RequestResponseStream) finish(handlerErr error) error {
	if handlerErr != nil {
		_ = s.writer.CloseWithError(handlerErr)
	} else if closeErr := s.encoder.Close(); closeErr != nil {
		_ = s.writer.CloseWithError(closeErr)
		handlerErr = closeErr
	} else {
		_ = s.writer.Close()
	}

	closeErr := s.close()
	if handlerErr != nil {
		return handlerErr
	}

	return closeErr
}

// responseDecoder holds the resolved decode-side machinery shared by [Client.Stream] and
// [Client.RequestStream] once a response is available: a streaming decoder bound to a per-value
// size-capped reader over the response body.
type responseDecoder struct {
	decoder encodingstream.Decoder
	capped  *budget.Reader
}

func newResponseDecoder(c *Client, response *http.Response, opts Options) (*responseDecoder, error) {
	resMedia, err := c.streamContent.NewFromMedia(responseContentType(response.Header, opts))
	if err == nil && resMedia.NewDecoder == nil {
		err = contentstream.ErrUnsupportedMedia
	}

	if err != nil {
		return nil, errors.Prefix("http: stream media", err)
	}

	capped := budget.NewReader(response.Body, c.maxResponseSize)

	return &responseDecoder{decoder: resMedia.NewDecoder(capped), capped: capped}, nil
}

// recv decodes the next value into v and, on success, checks it against the per-value cap. The check
// runs after Decode returns (mirroring [Client.readResponse]'s read-then-check pattern) rather than
// truncating capped's reads, because truncating mid-read would silently drop bytes the decoder never
// asked to give back — those bytes belong to the underlying response body, not to this one value, and
// dropping them would desynchronize every later value in the stream. See
// [github.com/alexfalkowski/go-service/v2/net/http/budget.Reader] for the read-time half of the guard,
// which still bounds a single pathological value across repeated Read calls.
//
// A decoder is not guaranteed to return capped's Read-time error unwrapped: it may fold it into its own
// error value instead (observed with the standard library JSON decoder under GOEXPERIMENT=jsonv2, which
// wraps a mid-scan Read error into a *json.SyntaxError with no Unwrap method). recv therefore checks
// capped's own latched error directly once Decode fails, rather than trusting Decode's returned error
// to still be classifiable as the size-limit error.
func (d *responseDecoder) recv(v any) error {
	d.capped.Reset(budget.BufferedLen(d.decoder))

	if err := d.decoder.Decode(v); err != nil {
		if d.capped.Err() != nil {
			return status.SafeError(http.StatusRequestEntityTooLarge, nil)
		}

		return err
	}

	if d.capped.Exceeds(budget.BufferedLen(d.decoder)) {
		return status.SafeError(http.StatusRequestEntityTooLarge, nil)
	}

	return nil
}

func (d *responseDecoder) close() error {
	return d.decoder.Close()
}
