package client

import (
	"cmp"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/sync"
	"github.com/alexfalkowski/go-service/v2/time"
)

// ClientOption configures the HTTP client wrapper constructed by NewClient.
//
// Options are applied in the order provided to NewClient. If multiple options configure the same
// field, the last one wins.
type ClientOption interface {
	apply(opts *clientOpts)
}

type clientOpts struct {
	roundTripper   http.RoundTripper
	timeout        time.Duration
	ignoreRedirect bool
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
// If not provided, NewClient defaults to 30 seconds.
func WithTimeout(timeout time.Duration) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.timeout = timeout
	})
}

// WithIgnoreRedirect disables automatic redirect following.
//
// When set, the underlying http.Client.CheckRedirect is configured to return http.ErrUseLastResponse,
// causing the client to return the redirect response (3xx) instead of following it.
//
// This can be useful when you want to inspect Location headers or preserve redirect responses as-is.
func WithIgnoreRedirect() ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.ignoreRedirect = true
	})
}

// NewClient constructs a Client that encodes requests and decodes responses using content.
//
// It reuses buffers from pool and applies the configured transport, timeout, and redirect policy.
//
// The underlying *http.Client is constructed via http.NewClient, which instruments requests with
// OpenTelemetry and sets http.Client.Timeout.
//
// Callers should treat the returned Client as safe for concurrent use.
func NewClient(content *content.Content, pool *sync.BufferPool, opts ...ClientOption) *Client {
	os := options(opts...)
	client := http.NewClient(os.roundTripper, os.timeout)

	if os.ignoreRedirect {
		client.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	return &Client{client: client, content: content, pool: pool}
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

	// ContentType is the request Content-Type used for encoding, and is also used as a fallback
	// decoder selection when the response does not provide a Content-Type header.
	ContentType string
}

// HasRequest reports whether a request payload is set.
func (o *Options) HasRequest() bool {
	return o.Request != nil
}

// HasResponse reports whether a response payload is set.
func (o *Options) HasResponse() bool {
	return o.Response != nil
}

// NoOptions is a reusable empty Options value.
//
// It can be passed when no request/response bodies are needed and defaults are acceptable.
var NoOptions = &Options{}

// Client wraps *http.Client with content-aware encoding and decoding helpers.
//
// It is intended for service-to-service calls where payload formats are selected by Content-Type.
// The Client uses a shared buffer pool to reduce allocations when encoding/decoding bodies.
type Client struct {
	client  *http.Client
	content *content.Content
	pool    *sync.BufferPool
}

// Delete issues an HTTP DELETE request to url using opts.
//
// It is a convenience wrapper around Do.
func (c *Client) Delete(ctx context.Context, url string, opts *Options) error {
	return c.Do(ctx, http.MethodDelete, url, opts)
}

// Get issues an HTTP GET request to url using opts.
//
// It is a convenience wrapper around Do.
func (c *Client) Get(ctx context.Context, url string, opts *Options) error {
	return c.Do(ctx, http.MethodGet, url, opts)
}

// Post issues an HTTP POST request to url using opts.
//
// It is a convenience wrapper around Do.
func (c *Client) Post(ctx context.Context, url string, opts *Options) error {
	return c.Do(ctx, http.MethodPost, url, opts)
}

// Put issues an HTTP PUT request to url using opts.
//
// It is a convenience wrapper around Do.
func (c *Client) Put(ctx context.Context, url string, opts *Options) error {
	return c.Do(ctx, http.MethodPut, url, opts)
}

// Patch issues an HTTP PATCH request to url using opts.
//
// It is a convenience wrapper around Do.
func (c *Client) Patch(ctx context.Context, url string, opts *Options) error {
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
//   - The full response body is read into an internal buffer.
//   - If the response Content-Type indicates an error payload (text/error), the body is treated as an
//     error message and returned as a net/http/status error.
//   - Otherwise, if the status code is in the 4xx/5xx range, a generic status error is returned.
//   - Otherwise, if opts.Response is non-nil, the response body is decoded into it using the encoder
//     selected by the response Content-Type (falling back to opts.ContentType if absent).
//
// Notes:
//   - opts must be non-nil; callers may pass NoOptions.
//   - This method buffers the entire response body in memory.
//
//nolint:cyclop
func (c *Client) Do(ctx context.Context, method, url string, opts *Options) error {
	buffer := c.pool.Get()
	defer c.pool.Put(buffer)

	mediaType := c.content.NewFromMedia(opts.ContentType)

	if opts.HasRequest() {
		if err := mediaType.Encoder.Encode(buffer, opts.Request); err != nil {
			return errors.Prefix("http: encode", err)
		}
	}

	request, err := http.NewRequestWithContext(ctx, method, url, buffer)
	if err != nil {
		return errors.Prefix("http: new request", err)
	}

	request.Header.Set(content.TypeKey, mediaType.Type)

	response, err := c.client.Do(request)
	if err != nil {
		return errors.Prefix("http: do", err)
	}

	defer response.Body.Close()

	buffer.Reset()

	_, err = buffer.ReadFrom(response.Body)
	if err != nil {
		return errors.Prefix("http: copy", err)
	}

	// If for some reason the server does not return it, default to opts.
	contentType := cmp.Or(response.Header.Get(content.TypeKey), opts.ContentType)

	// The server handlers return text/error to indicate an error.
	media := c.content.NewFromMedia(contentType)
	if media.IsError() {
		return status.Error(response.StatusCode, strings.TrimSpace(buffer.String()))
	}

	if response.StatusCode >= 400 && response.StatusCode <= 599 {
		return status.Error(response.StatusCode, strings.ToLower(http.StatusText(response.StatusCode)))
	}

	if opts.HasResponse() {
		if err := media.Encoder.Decode(buffer, opts.Response); err != nil {
			return errors.Prefix("http: decode", err)
		}
	}

	return nil
}

func options(opts ...ClientOption) *clientOpts {
	os := &clientOpts{}
	for _, o := range opts {
		o.apply(os)
	}

	if os.timeout == 0 {
		os.timeout = 30 * time.Second
	}

	if os.roundTripper == nil {
		os.roundTripper = http.Transport(nil)
	}

	return os
}
