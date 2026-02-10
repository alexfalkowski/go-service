package rest

import (
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/client"
	"github.com/alexfalkowski/go-service/v2/time"
)

// NoOptions is an alias for client.NoOptions.
var NoOptions = client.NoOptions

// Options is an alias for client.Options.
type Options = client.Options

// ClientOption configures the REST client helper.
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

// WithClientRoundTripper sets the underlying HTTP RoundTripper used by the REST client.
func WithClientRoundTripper(rt http.RoundTripper) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.roundTripper = rt
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

// NewClient constructs a REST client backed by net/http/client.
//
// NewClient depends on package-level registration (see rest.Register) for the content codecs (cont)
// and buffer pool (pool). Register must be called before NewClient; otherwise it will panic due to nil dependencies.
func NewClient(opts ...ClientOption) *Client {
	os := options(opts...)
	client := client.NewClient(cont, pool,
		client.WithRoundTripper(os.roundTripper),
		client.WithTimeout(os.timeout),
		client.WithIgnoreRedirect(),
	)

	return &Client{client}
}

// Client wraps client.Client for REST usage.
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
