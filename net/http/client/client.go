package client

import (
	"cmp"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
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
	apply(opts *clientOpts)
}

// Redirect configures how Client handles HTTP redirects.
type Redirect int

// RedirectFollow follows redirects using the standard library default policy.
const RedirectFollow Redirect = iota

// RedirectIgnore returns redirect responses without following them.
const RedirectIgnore Redirect = 1

// RedirectSameOrigin follows redirects only when scheme and host are unchanged.
const RedirectSameOrigin Redirect = 2

type clientOpts struct {
	roundTripper    http.RoundTripper
	timeout         time.Duration
	maxResponseSize bytes.Size
	redirect        Redirect
}

type clientOptionFunc func(*clientOpts)

func (f clientOptionFunc) apply(o *clientOpts) {
	f(o)
}

// WithRoundTripper sets the underlying RoundTripper used to execute requests.
//
// This is typically used to inject a transport that includes additional middleware such as:
// retries, circuit breakers, auth/token injection, custom TLS, etc.
//
// If not provided, NewClient uses http.Transport(nil) (go-service's tuned default transport).
func WithRoundTripper(rt http.RoundTripper) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.roundTripper = rt
	})
}

// WithTimeout sets the overall request timeout on the underlying http.Client.
//
// The timeout value is assigned to http.Client.Timeout (total time limit for a request, including
// connection time, redirects, and reading the response body).
//
// If not provided, NewClient defaults to time.DefaultTimeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.timeout = timeout
	})
}

// WithMaxResponseSize sets the maximum response body size buffered by Client.Do.
//
// If not provided, NewClient defaults to bytes.DefaultSize.
func WithMaxResponseSize(size bytes.Size) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.maxResponseSize = size
	})
}

// WithRedirect sets the redirect policy used by the underlying http.Client.
//
// If not provided, NewClient uses RedirectFollow, which preserves standard library redirect behavior.
func WithRedirect(redirect Redirect) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.redirect = redirect
	})
}

// NewClient constructs a Client that encodes requests and decodes responses using content.
//
// It reuses buffers from pool and applies the configured transport, timeout, and redirect policy.
//
// The underlying *http.Client is constructed via http.NewClient, which sets http.Client.Timeout and
// instruments requests with OpenTelemetry when tracing or metrics are enabled.
//
// Callers should treat the returned Client as safe for concurrent use.
func NewClient(content *content.Content, pool *sync.BufferPool, opts ...ClientOption) *Client {
	os := options(opts...)
	client := http.NewClient(os.roundTripper, os.timeout)

	switch os.redirect {
	case RedirectIgnore:
		client.CheckRedirect = http.IgnoreRedirect
	case RedirectSameOrigin:
		client.CheckRedirect = http.SameOriginRedirect
	}

	return &Client{client: client, content: content, pool: pool, maxResponseSize: os.maxResponseSize.Bytes()}
}

// Options describes the request/response payloads and content type for a single call.
//
// ContentType is used to select an encoder/decoder via net/http/content. Typical values are media
// types like "application/json" or go-service specific protobuf media types.
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
	// when the response does not provide a Content-Type header.
	ContentType string
}

// HasRequest reports whether a request payload is set.
func (o Options) HasRequest() bool {
	return o.Request != nil
}

// HasResponse reports whether a response payload is set.
func (o Options) HasResponse() bool {
	return o.Response != nil
}

// Client wraps *http.Client with content-aware encoding and decoding helpers.
//
// It is intended for service-to-service calls where payload formats are selected by Content-Type.
// The Client uses a shared buffer pool to reduce allocations when encoding/decoding bodies.
type Client struct {
	client          *http.Client
	content         *content.Content
	pool            *sync.BufferPool
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

// Do issues a request with method and url, encoding and decoding bodies via content.
//
// Encoding:
//   - If opts.Request is non-nil, it is encoded into the request body using the encoder selected by
//     opts.ContentType.
//
// Request headers:
//   - The request Content-Type header is set to the negotiated media type.
//
// Response handling:
//   - The response body is read into an internal buffer up to the configured response size limit.
//   - If the response Content-Type indicates an error payload (text/error), the body is treated as an
//     error message and returned as a net/http/status error.
//   - Otherwise, if the status code is in the 4xx/5xx range, a generic status error is returned.
//   - Otherwise, if opts.Response is non-nil, the response body is decoded into it using the encoder
//     selected by the response Content-Type (falling back to opts.ContentType if absent).
//
// Notes:
//   - callers may pass the zero Options value when no request/response bodies are needed.
//   - This method buffers response bodies in memory up to the configured limit.
//
//nolint:cyclop
func (c *Client) Do(ctx context.Context, method, url string, opts Options) error {
	mediaType := c.content.NewFromMedia(opts.ContentType)

	body := io.Reader(http.NoBody)
	if opts.HasRequest() {
		var requestBody bytes.Buffer
		if err := mediaType.Encoder.Encode(&requestBody, opts.Request); err != nil {
			return errors.Prefix("http: encode", err)
		}

		// Do not use a pooled buffer for the request body: net/http may keep
		// reading it from a transport goroutine after Do returns.
		body = bytes.NewReader(requestBody.Bytes())
	}

	request, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return errors.Prefix("http: new request", err)
	}

	request.Header.Set(content.TypeKey, mediaType.Type)

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

	// If for some reason the server does not return it, default to opts.
	contentType := cmp.Or(response.Header.Get(content.TypeKey), opts.ContentType)

	// The server handlers return text/error to indicate an error.
	media := c.content.NewFromMedia(contentType)
	if media.IsError() {
		code := response.StatusCode
		if !isErrorStatus(code) {
			code = http.StatusInternalServerError
		}

		return status.Error(code, strings.TrimSpace(responseBody.String()))
	}

	if isErrorStatus(response.StatusCode) {
		return status.Error(response.StatusCode, strings.ToLower(http.StatusText(response.StatusCode)))
	}

	if opts.HasResponse() {
		if err := media.Encoder.Decode(responseBody, opts.Response); err != nil {
			return errors.Prefix("http: decode", err)
		}
	}

	return nil
}

func isErrorStatus(code int) bool {
	return code >= 400 && code <= 599
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

func options(opts ...ClientOption) *clientOpts {
	os := &clientOpts{}
	for _, o := range opts {
		o.apply(os)
	}

	if os.timeout <= 0 {
		os.timeout = time.DefaultTimeout
	}

	if os.maxResponseSize <= 0 {
		os.maxResponseSize = bytes.DefaultSize
	}

	if os.roundTripper == nil {
		os.roundTripper = http.Transport(nil)
	}

	return os
}
