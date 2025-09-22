package rpc

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/client"
	"github.com/alexfalkowski/go-service/v2/time"
)

var (
	// ErrInvalidRequest when we pass nil.
	ErrInvalidRequest = errors.New("rpc: invalid request")

	// ErrInvalidResponse when we pass nil.
	ErrInvalidResponse = errors.New("rpc: invalid response")
)

// ClientOption for rpc.
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

// WithClientRoundTripper for rpc.
func WithClientRoundTripper(rt http.RoundTripper) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.roundTripper = rt
	})
}

// WithClientContentType for rpc.
func WithClientContentType(ct string) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.contentType = ct
	})
}

// WithClientTimeout for rpc.
func WithClientTimeout(timeout string) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.timeout = time.MustParseDuration(timeout)
	})
}

// NewClient for rpc.
func NewClient(url string, opts ...ClientOption) *Client {
	os := options(opts...)
	client := client.NewClient(cont, pool,
		client.WithRoundTripper(os.roundTripper),
		client.WithTimeout(os.timeout),
		client.WithIgnoreRedirect(),
	)

	return &Client{client: client, url: url, contentType: os.contentType}
}

// Client for rpc.
type Client struct {
	client      *client.Client
	contentType string
	url         string
}

// Post for rpc.
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
