package rpc

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/alexfalkowski/go-service/errors"
	nh "github.com/alexfalkowski/go-service/net/http"
	"github.com/alexfalkowski/go-service/net/http/content"
	"github.com/alexfalkowski/go-service/net/http/status"
	"github.com/alexfalkowski/go-service/runtime"
	st "github.com/alexfalkowski/go-service/time"
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
		o.timeout = st.MustParseDuration(timeout)
	})
}

// NewClient for rpc.
func NewClient[Req any, Res any](url string, opts ...ClientOption) *Client[Req, Res] {
	os := options(opts...)
	client := &http.Client{
		Transport: os.roundTripper,
		Timeout:   os.timeout,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	return &Client[Req, Res]{client: client, url: url, ct: content.NewFromMedia(os.contentType)}
}

// Client for HTTP.
type Client[Req any, Res any] struct {
	client *http.Client
	ct     *content.Type
	url    string
}

// Invoke for rpc.
//
//nolint:nonamedreturns
func (c *Client[Req, Res]) Invoke(ctx context.Context, req *Req) (res *Res, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Prefix("rpc", runtime.ConvertRecover(r))
		}
	}()

	m, err := c.ct.Marshaller(enc)
	runtime.Must(err)

	d, err := m.Marshal(req)
	runtime.Must(err)

	request, err := http.NewRequestWithContext(ctx, "POST", c.url, bytes.NewBuffer(d))
	runtime.Must(err)

	request.Header.Set(content.TypeKey, c.ct.Media)

	response, err := c.client.Do(request)
	runtime.Must(err)

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	runtime.Must(err)

	// The server handlers return text on errors.
	ct := content.NewFromMedia(response.Header.Get(content.TypeKey))
	if ct.IsText() {
		return nil, status.Error(response.StatusCode, strings.TrimSpace(string(body)))
	}

	var rp Res
	res = &rp

	err = m.Unmarshal(body, res)
	runtime.Must(err)

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
