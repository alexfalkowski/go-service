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

// ClientOption configures the RPC client helper.
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
func WithClientRoundTripper(rt http.RoundTripper) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.roundTripper = rt
	})
}

// WithClientContentType sets the Content-Type used for requests made by the RPC client.
func WithClientContentType(ct string) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.contentType = ct
	})
}

// WithClientTimeout sets the client timeout using a duration string (for example "1s" or "500ms").
//
// It uses time.MustParseDuration and will panic if the duration string cannot be parsed.
func WithClientTimeout(timeout string) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.timeout = time.MustParseDuration(timeout)
	})
}

// NewClient constructs an RPC client backed by net/http/client.
//
// NewClient depends on package-level registration (see rpc.Register) for the content codecs (cont)
// and buffer pool (pool). Register must be called before NewClient; otherwise it will panic due to nil dependencies.
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
// It returns ErrInvalidRequest when req is nil and ErrInvalidResponse when res is nil.
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
