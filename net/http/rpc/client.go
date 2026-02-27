package rpc

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/client"
	"github.com/alexfalkowski/go-service/v2/time"
)

var (
	// ErrInvalidRequest is returned when the request payload is nil.
	ErrInvalidRequest = errors.New("rpc: invalid request")

	// ErrInvalidResponse is returned when the response target is nil.
	ErrInvalidResponse = errors.New("rpc: invalid response")
)

// ClientOption configures the RPC client helper constructed by NewClient.
//
// Options are applied in the order provided to NewClient. If multiple options configure the same
// field, the last one wins.
type ClientOption interface {
	apply(opts *clientOpts)
}

type clientOpts struct {
	roundTripper http.RoundTripper
	contentType  string
	timeout      time.Duration
}

type clientOptionFunc func(*clientOpts)

func (f clientOptionFunc) apply(o *clientOpts) {
	f(o)
}

// WithClientRoundTripper sets the underlying HTTP RoundTripper used by the RPC client.
//
// This is typically used to inject a transport that includes middleware such as retries, circuit
// breakers, authentication, or custom TLS.
func WithClientRoundTripper(rt http.RoundTripper) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.roundTripper = rt
	})
}

// WithClientContentType sets the Content-Type used for requests made by the RPC client.
//
// This value is passed through to the underlying content-aware HTTP client and is used to select the
// encoder/decoder for request/response bodies. Typical values include "application/json" or go-service
// protobuf media types.
func WithClientContentType(ct string) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.contentType = ct
	})
}

// WithClientTimeout sets the client timeout using a duration string (for example "1s" or "500ms").
//
// The duration string is parsed using time.MustParseDuration and will panic if it cannot be parsed.
// The resulting duration is applied to the underlying http.Client.Timeout.
func WithClientTimeout(timeout string) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.timeout = time.MustParseDuration(timeout)
	})
}

// NewClient constructs an RPC client backed by net/http/client.
//
// NewClient depends on package-level registration (see rpc.Register) for the content codecs (cont)
// and buffer pool (pool). Register must be called before NewClient; otherwise it will panic due to
// nil dependencies.
//
// The returned client issues RPC-style POST requests to the provided base url using the configured
// Content-Type and transport options. Redirect following is disabled by default (redirect responses
// are returned instead of being followed).
func NewClient(url string, opts ...ClientOption) *Client {
	os := options(opts...)
	client := client.NewClient(cont, pool,
		client.WithRoundTripper(os.roundTripper),
		client.WithTimeout(os.timeout),
		client.WithIgnoreRedirect(),
	)

	return &Client{client: client, url: url, contentType: os.contentType}
}

// Client is an RPC client that issues POST requests using the configured content codecs.
type Client struct {
	client      *client.Client
	contentType string
	url         string
}

// Post issues an RPC-style HTTP POST request to c.url+path.
//
// Request/response validation:
//   - returns ErrInvalidRequest when req is nil
//   - returns ErrInvalidResponse when res is nil
//
// Content-Type behavior:
// The request Content-Type is set to c.contentType and is used to select encoders/decoders via
// the underlying content-aware client.
//
// The res parameter is typically a pointer to the destination value (for example *MyResponse).
func (c *Client) Post(ctx context.Context, path string, req, res any) error {
	if req == nil {
		return ErrInvalidRequest
	}
	if res == nil {
		return ErrInvalidResponse
	}

	opts := &client.Options{ContentType: c.contentType, Request: req, Response: res}
	return c.client.Post(ctx, c.url+path, opts)
}

func options(opts ...ClientOption) *clientOpts {
	os := &clientOpts{}
	for _, o := range opts {
		o.apply(os)
	}
	return os
}
