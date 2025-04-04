package rpc

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/alexfalkowski/go-service/errors"
	nh "github.com/alexfalkowski/go-service/net/http"
	"github.com/alexfalkowski/go-service/net/http/content"
	"github.com/alexfalkowski/go-service/net/http/status"
	"github.com/alexfalkowski/go-service/time"
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

// WithSenderRoundTripper for rpc.
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
func NewClient[Req any, Res any](url string, opts ...ClientOption) *Client[Req, Res] {
	os := options(opts...)
	client := &http.Client{
		Transport: os.roundTripper,
		Timeout:   os.timeout,
	}

	return &Client[Req, Res]{client: client, url: url, mediaType: cont.NewFromMedia(os.contentType)}
}

// Client for rpc.
type Client[Req any, Res any] struct {
	client    *http.Client
	mediaType *content.Media
	url       string
}

// Invoke for rpc.
func (c *Client[Req, Res]) Invoke(ctx context.Context, req *Req) (*Res, error) {
	buffer := pool.Get()
	defer pool.Put(buffer)

	if err := c.mediaType.Encoder.Encode(buffer, req); err != nil {
		return nil, errors.Prefix("rpc", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, buffer)
	if err != nil {
		return nil, errors.Prefix("rpc", err)
	}

	request.Header.Set(content.TypeKey, c.mediaType.Type)

	response, err := c.client.Do(request)
	if err != nil {
		return nil, errors.Prefix("rpc", err)
	}

	defer response.Body.Close()

	buffer.Reset()

	_, err = io.Copy(buffer, response.Body)
	if err != nil {
		return nil, errors.Prefix("rpc", err)
	}

	// The server handlers return text on errors.
	media := cont.NewFromMedia(response.Header.Get(content.TypeKey))
	if media.IsText() {
		return nil, status.Error(response.StatusCode, strings.TrimSpace(buffer.String()))
	}

	var rp Res
	res := &rp

	if err := media.Encoder.Decode(buffer, res); err != nil {
		return nil, errors.Prefix("rpc", err)
	}

	return res, nil
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
