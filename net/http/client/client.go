package client

import (
	"cmp"
	"context"
	"io"
	"net/http"

	"github.com/alexfalkowski/go-service/errors"
	nh "github.com/alexfalkowski/go-service/net/http"
	"github.com/alexfalkowski/go-service/net/http/content"
	"github.com/alexfalkowski/go-service/net/http/status"
	"github.com/alexfalkowski/go-service/strings"
	"github.com/alexfalkowski/go-service/sync"
	"github.com/alexfalkowski/go-service/time"
)

// ClientOption for http.
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

// WithRoundTripper for http.
func WithRoundTripper(rt http.RoundTripper) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.roundTripper = rt
	})
}

// WithTimeout for http.
func WithTimeout(timeout time.Duration) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.timeout = timeout
	})
}

// NewClient for http.
func NewClient(content *content.Content, pool *sync.BufferPool, opts ...ClientOption) *Client {
	os := options(opts...)
	client := &http.Client{
		Transport: os.roundTripper,
		Timeout:   os.timeout,
	}

	return &Client{client: client, content: content, pool: pool}
}

// Options for client.
type Options struct {
	Request     any
	Response    any
	ContentType string
}

// NoOptions for http.
var NoOptions = &Options{}

// Client for http.
type Client struct {
	client  *http.Client
	content *content.Content
	pool    *sync.BufferPool
}

// Delete the url.
func (c *Client) Delete(ctx context.Context, url string, opts *Options) error {
	return c.Do(ctx, http.MethodDelete, url, opts)
}

// Get the url.
func (c *Client) Get(ctx context.Context, url string, opts *Options) error {
	return c.Do(ctx, http.MethodGet, url, opts)
}

// Post the url.
func (c *Client) Post(ctx context.Context, url string, opts *Options) error {
	return c.Do(ctx, http.MethodPost, url, opts)
}

// Put the url.
func (c *Client) Put(ctx context.Context, url string, opts *Options) error {
	return c.Do(ctx, http.MethodPut, url, opts)
}

// Patch the url.
func (c *Client) Patch(ctx context.Context, url string, opts *Options) error {
	return c.Do(ctx, http.MethodPatch, url, opts)
}

// Do with method and url.
//
//nolint:cyclop
func (c *Client) Do(ctx context.Context, method, url string, opts *Options) error {
	buffer := c.pool.Get()
	defer c.pool.Put(buffer)

	mediaType := c.content.NewFromMedia(opts.ContentType)

	if opts.Request != nil {
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

	_, err = io.Copy(buffer, response.Body)
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

	if opts.Response != nil {
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
		os.roundTripper = nh.Transport(nil)
	}

	return os
}
