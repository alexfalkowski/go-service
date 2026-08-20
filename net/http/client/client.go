package client

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/encoding/codec"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content/stream"
	"github.com/alexfalkowski/go-service/v2/net/http/content/unary"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-sync"
)

// ClientOption configures the HTTP client wrapper constructed by NewClient.
//
// Options are applied in the order provided to NewClient. If multiple options configure the same
// field, the last one wins.
type ClientOption interface {
	apply(opts *clientOptions)
}

// Redirect configures how Client handles HTTP redirects.
type Redirect int

const (
	// RedirectFollow follows redirects using the standard library default policy.
	RedirectFollow Redirect = iota

	// RedirectIgnore returns redirect responses without following them.
	RedirectIgnore

	// RedirectSameOrigin follows redirects only when scheme and host are unchanged.
	RedirectSameOrigin
)

type clientOptions struct {
	roundTripper    http.RoundTripper
	timeout         time.Duration
	maxResponseSize bytes.Size
	redirect        Redirect
}

type clientOptionFunc func(*clientOptions)

func (f clientOptionFunc) apply(o *clientOptions) {
	f(o)
}

// WithRoundTripper sets the underlying RoundTripper used to execute requests.
//
// This is typically used to inject a transport that includes additional middleware such as:
// retries, circuit breakers, auth/token injection, custom TLS, etc.
//
// If not provided, NewClient uses [http.Transport](nil) (go-service's tuned default transport).
func WithRoundTripper(rt http.RoundTripper) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.roundTripper = rt
	})
}

// WithTimeout sets the overall timeout for unary [Client.Do] calls.
//
// The timeout includes connection time, redirects, and reading the response body. It is applied to
// the request context rather than [http.Client.Timeout], so [Stream] and [RequestStream] remain
// long-lived and are bounded only by their caller-provided contexts.
//
// If not provided, NewClient defaults to [time.DefaultTimeout].
func WithTimeout(timeout time.Duration) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.timeout = timeout
	})
}

// WithMaxResponseSize sets the maximum response body size buffered by [Client.Do].
//
// For streaming calls made through [Stream] or [RequestStream], the same size instead bounds each
// decoded value individually rather than the whole response: a stream has no natural total size, so a
// cumulative cap would either be meaningless for a long-lived stream or reject a legitimately large
// number of small values. This matches the inbound per-value cap applied to streaming request bodies
// on the server (see the transport/http/body streaming middleware), so operators use one size
// vocabulary for both directions.
//
// If not provided, NewClient defaults to [bytes.DefaultSize].
func WithMaxResponseSize(size bytes.Size) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.maxResponseSize = size
	})
}

// WithRedirect sets the redirect policy used by the underlying [http.Client].
//
// If not provided, NewClient uses RedirectSameOrigin so credential/signature
// middleware supplied through WithRoundTripper cannot be replayed to a different
// origin by an upstream redirect. Use RedirectFollow to opt into standard
// library cross-origin redirect behavior.
func WithRedirect(redirect Redirect) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.redirect = redirect
	})
}

// NewClient constructs a Client that encodes requests and decodes responses using unary and stream content.
//
// It reuses buffers from pool and applies the configured transport, unary timeout, and redirect policy.
//
// The underlying *[http.Client] is constructed via [http.NewClient] with no timeout. [WithTimeout]
// applies the configured deadline only to [Client.Do], preserving long-lived streaming calls.
//
// [Stream] and [RequestStream] resolve streaming media types through streamContent.
//
// Callers should treat the returned Client as safe for concurrent use.
func NewClient(uc *unary.Content, sc *stream.Content, pool *sync.BufferPool, opts ...ClientOption) *Client {
	clientOptions := options(opts...)

	client := http.NewClient(clientOptions.roundTripper, 0)

	switch clientOptions.redirect {
	case RedirectIgnore:
		client.CheckRedirect = http.IgnoreRedirect
	case RedirectSameOrigin:
		client.CheckRedirect = http.SameOriginRedirect
	}

	return &Client{
		client:          client,
		unaryContent:    uc,
		streamContent:   sc,
		pool:            pool,
		timeout:         clientOptions.timeout,
		maxResponseSize: clientOptions.maxResponseSize.Bytes(),
	}
}

// Options describes the request/response payloads and media types for a single call.
//
// ContentType is used to select the request encoder via net/http/content/unary. Accept is used to request
// a distinct response media type.
// Typical values are media types like "application/json" or go-service specific protobuf media types.
//
// Request and Response are optional:
//   - If Request is non-nil, it is encoded into the request body.
//   - If Response is non-nil, it is decoded from the response body after a successful (non-error)
//     response is received.
type Options struct {
	// Request is the request payload model to encode into the request body (optional).
	Request any

	// Response is the response payload model to decode into (optional).
	// It is typically a pointer to a struct or message.
	Response any

	// ContentType is the request Content-Type used for encoding and the fallback decoder selection
	// when the response carries no Content-Type header. Accept never selects the response decoder
	// itself; it is only sent on the request, so it acts on decoding indirectly, through the
	// Content-Type the server chooses to answer with.
	ContentType string

	// Accept is the response Accept media type sent on the request.
	Accept string
}

// HasRequest reports whether a request payload is set.
func (o Options) HasRequest() bool {
	return o.Request != nil
}

// HasResponse reports whether a response payload is set.
func (o Options) HasResponse() bool {
	return o.Response != nil
}

// Client wraps *[http.Client] with content-aware encoding and decoding helpers.
//
// It is intended for service-to-service calls where payload formats are selected by Content-Type.
// The Client uses a shared buffer pool to reduce allocations when encoding/decoding bodies.
type Client struct {
	client          *http.Client
	unaryContent    *unary.Content
	streamContent   *stream.Content
	pool            *sync.BufferPool
	timeout         time.Duration
	maxResponseSize int64
}

// Delete issues an HTTP DELETE request to url using opts.
//
// It is a convenience wrapper around Do.
func (c *Client) Delete(ctx context.Context, url string, opts Options) error {
	return c.Do(ctx, http.MethodDelete, url, opts)
}

// Get issues an HTTP GET request to url using opts.
//
// It is a convenience wrapper around Do.
func (c *Client) Get(ctx context.Context, url string, opts Options) error {
	return c.Do(ctx, http.MethodGet, url, opts)
}

// Post issues an HTTP POST request to url using opts.
//
// It is a convenience wrapper around Do.
func (c *Client) Post(ctx context.Context, url string, opts Options) error {
	return c.Do(ctx, http.MethodPost, url, opts)
}

// Put issues an HTTP PUT request to url using opts.
//
// It is a convenience wrapper around Do.
func (c *Client) Put(ctx context.Context, url string, opts Options) error {
	return c.Do(ctx, http.MethodPut, url, opts)
}

// Patch issues an HTTP PATCH request to url using opts.
//
// It is a convenience wrapper around Do.
func (c *Client) Patch(ctx context.Context, url string, opts Options) error {
	return c.Do(ctx, http.MethodPatch, url, opts)
}

// Do issues a request with method and url, encoding and decoding bodies via unary.
//
// Encoding:
//   - If opts.Request is non-nil, it is encoded into the request body using the encoder selected by
//     opts.ContentType.
//
// Request headers:
//   - The request Content-Type header is set to the negotiated media type.
//   - The request Accept header is set when opts.Accept is non-empty.
//
// Response handling:
//   - The response body is read into an internal buffer up to the configured response size limit.
//   - If the response Content-Type indicates an error payload (text/error), the body is treated as an
//     error message and returned as a net/http/status error.
//   - Otherwise, if the status code is in the 4xx/5xx range, a generic status error is returned.
//   - Otherwise, if opts.Response is non-nil, the response body is decoded into it using the encoder
//     selected by the response Content-Type (falling back to opts.ContentType).
//     Unknown response members are discarded so newer servers remain compatible with older clients.
//   - An empty successful response body is not decoded, leaving opts.Response unchanged.
//
// Notes:
//   - Callers may pass the zero Options value when no request/response bodies are needed.
//   - This method buffers response bodies in memory up to the configured limit.
func (c *Client) Do(ctx context.Context, method, url string, opts Options) error {
	request, err := c.newRequest(ctx, method, url, opts)
	if err != nil {
		return err
	}

	requestContext, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	request = request.WithContext(requestContext)

	response, err := c.client.Do(request)
	if err != nil {
		return errors.Prefix("http: do", err)
	}

	defer response.Body.Close()

	responseBody := c.pool.Get()
	defer c.pool.Put(responseBody)

	if err := c.readResponse(responseBody, response.Body); err != nil {
		return err
	}

	// The server handlers return text/error to indicate an error.
	responseMedia := c.unaryContent.NewFromMedia(responseContentType(response.Header, opts))
	if err := mediaStatusError(response, responseMedia, responseBody.String()); err != nil {
		return err
	}

	if opts.HasResponse() {
		// With the default encoder registry, responseMedia.Encoder is never nil here: the only media
		// type NewFromMedia resolves without an encoder is text/error, and mediaStatusError above always
		// turns that into an error before this point. That is not a structural guarantee, though —
		// nil-registering "json" makes NewFromMedia return a non-error Media with a nil Encoder. Response
		// decoding is intentionally not subject to the request-decode policy that guards
		// Content.NewFromRequestBody (a configured response endpoint is a different trust context from
		// an inbound request body); see the decoder-bounds rule in net/http/content/unary's documentation.
		if responseMedia.Encoder == nil {
			return errors.Prefix("http: decode", unary.ErrUnsupportedMedia)
		}

		if responseBody.Len() == 0 {
			return nil
		}

		if err := responseMedia.Encoder.Decode(responseBody, opts.Response, codec.WithDiscardUnknown()); err != nil {
			return errors.Prefix("http: decode", err)
		}
	}

	return nil
}

// newRequest builds an outgoing request for method/url, encoding opts.Request (if set) into the
// request body using the encoder selected by opts.ContentType, and setting the Content-Type/Accept
// headers.
//
// This is shared by Do and Stream; RequestStream builds its own pipe-backed request instead, since a
// bidirectional streaming request has no single-value body to encode upfront.
func (c *Client) newRequest(ctx context.Context, method, url string, opts Options) (*http.Request, error) {
	requestMedia := c.unaryContent.NewFromMedia(opts.ContentType)

	body := io.Reader(http.NoBody)
	if opts.HasRequest() {
		if requestMedia.Encoder == nil {
			return nil, errors.Prefix("http: encode", unary.ErrUnsupportedMedia)
		}

		var requestBody bytes.Buffer
		if err := requestMedia.Encoder.Encode(&requestBody, opts.Request); err != nil {
			return nil, errors.Prefix("http: encode", err)
		}

		// Do not use a pooled buffer for the request body: net/http may keep
		// reading it from a transport goroutine after Do returns.
		body = bytes.NewReader(requestBody.Bytes())
	}

	request, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, errors.Prefix("http: new request", err)
	}

	request.Header.Set(http.ContentTypeKey, requestMedia.String())
	if !strings.IsEmpty(opts.Accept) {
		request.Header.Set(http.AcceptKey, opts.Accept)
	}

	return request, nil
}

func (c *Client) readResponse(buffer *bytes.Buffer, body io.Reader) error {
	_, err := buffer.ReadFrom(io.LimitReader(body, c.maxResponseSize+1))
	if err != nil {
		return errors.Prefix("http: copy", err)
	}

	if int64(buffer.Len()) > c.maxResponseSize {
		return status.SafeError(http.StatusRequestEntityTooLarge, nil)
	}

	return nil
}

// checkResponseStatus applies [mediaStatusError] to response, reading a bounded prefix of
// response.Body into a pooled buffer only when the response uses the text/error media convention.
//
// Unlike Do, which always buffers the whole response body upfront, this only reads the body when the
// error signal actually requires a message, so a non-error streaming response is left untouched for
// the caller's decoder.
func (c *Client) checkResponseStatus(response *http.Response, opts Options) error {
	responseMedia := c.unaryContent.NewFromMedia(responseContentType(response.Header, opts))

	message := strings.Empty
	if responseMedia.IsError() {
		buffer := c.pool.Get()
		defer c.pool.Put(buffer)

		if err := c.readResponse(buffer, response.Body); err != nil {
			return err
		}

		message = buffer.String()
	}

	return mediaStatusError(response, responseMedia, message)
}

// mediaStatusError returns a [status.Error] when response indicates an error, either through the
// text/error media convention (responseMedia.IsError(), using message as the trimmed error body text)
// or a 4xx/5xx status code with no such convention. It returns nil when response is not an error by
// either signal.
//
// This is the shared decision behind Do's response handling and the equivalent check applied to the
// initial response of a [Stream] or [RequestStream] call, before either streams values from/to the
// caller-supplied handler. Factoring it out keeps both paths honoring the same text/error and status
// code contract instead of the streaming path silently dropping it (see [Client.checkResponseStatus]).
func mediaStatusError(response *http.Response, responseMedia unary.Media, message string) error {
	if responseMedia.IsError() {
		code := response.StatusCode
		if !isErrorStatus(code) {
			code = http.StatusInternalServerError
		}

		return status.Error(code, strings.TrimSpace(message))
	}

	if isErrorStatus(response.StatusCode) {
		return status.Error(response.StatusCode, defaultErrorMessage(response.StatusCode))
	}

	return nil
}

func isErrorStatus(code int) bool {
	return code >= 400 && code <= 599
}

func defaultErrorMessage(code int) string {
	if message := http.StatusText(code); !strings.IsEmpty(message) {
		return strings.ToLower(message)
	}

	return status.DefaultMessage(code)
}

func responseContentType(header http.Header, opts Options) string {
	if contentType := header.Get(http.ContentTypeKey); !strings.IsEmpty(contentType) {
		return contentType
	}

	return opts.ContentType
}

func options(opts ...ClientOption) *clientOptions {
	clientOptions := &clientOptions{redirect: RedirectSameOrigin}
	for _, o := range opts {
		o.apply(clientOptions)
	}

	if clientOptions.maxResponseSize <= 0 {
		clientOptions.maxResponseSize = bytes.DefaultSize
	}

	if clientOptions.timeout <= 0 {
		clientOptions.timeout = time.DefaultTimeout
	}

	if clientOptions.roundTripper == nil {
		clientOptions.roundTripper = http.Transport(nil)
	}

	return clientOptions
}
