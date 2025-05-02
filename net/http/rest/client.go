package rest

import (
	"net/http"

	"github.com/alexfalkowski/go-service/net/http/client"
	"github.com/alexfalkowski/go-service/time"
)

// NoOptions is just an alias for client.NoOptions.
var NoOptions = client.NoOptions

// Options is just an alias for client.Options.
type Options = client.Options

// ClientOption for rest.
type ClientOption interface {
	apply(opts *clientOpts)
}

type clientOpts struct {
	roundTripper http.RoundTripper
	timeout      time.Duration
}

type clientOptionFunc func(*clientOpts)

func (f clientOptionFunc) apply(o *clientOpts) {
	f(o)
}

// WithSenderRoundTripper for rest.
func WithClientRoundTripper(rt http.RoundTripper) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.roundTripper = rt
	})
}

// WithClientTimeout for rest.
func WithClientTimeout(timeout string) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.timeout = time.MustParseDuration(timeout)
	})
}

// NewClient for rest.
func NewClient(opts ...ClientOption) *Client {
	os := options(opts...)
	client := client.NewClient(cont, pool,
		client.WithRoundTripper(os.roundTripper),
		client.WithTimeout(os.timeout),
	)

	return &Client{client}
}

// Client for rest.
type Client struct {
	*client.Client
}

func options(opts ...ClientOption) *clientOpts {
	os := &clientOpts{}
	for _, o := range opts {
		o.apply(os)
	}

	return os
}
