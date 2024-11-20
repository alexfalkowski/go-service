package rest

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/maps"
	nh "github.com/alexfalkowski/go-service/net/http"
	"github.com/alexfalkowski/go-service/net/http/content"
	"github.com/alexfalkowski/go-service/net/http/status"
	"github.com/alexfalkowski/go-service/runtime"
	st "github.com/alexfalkowski/go-service/time"
)

// NoRequest to use for the client.
//
//nolint:revive
var NoRequest maps.StringAny = nil

// ClientOption for rest.
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

// WithSenderRoundTripper for rest.
func WithClientRoundTripper(rt http.RoundTripper) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.roundTripper = rt
	})
}

// WithClientContentType for rest.
func WithClientContentType(ct string) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.contentType = ct
	})
}

// WithClientTimeout for rest.
func WithClientTimeout(timeout string) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.timeout = st.MustParseDuration(timeout)
	})
}

// NewClient for rest.
func NewClient(opts ...ClientOption) *Client {
	os := options(opts...)
	client := &http.Client{
		Transport: os.roundTripper,
		Timeout:   os.timeout,
	}

	return &Client{client: client, ct: os.contentType}
}

// Client for rest.
type Client struct {
	client *http.Client
	ct     string
}

// Invoke for rest.
//
//nolint:nonamedreturns
func (c *Client) Invoke(ctx context.Context, method, url string, req maps.StringAny) (res maps.StringAny, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Prefix("rest", runtime.ConvertRecover(r))
		}
	}()

	b := pool.Get()
	defer pool.Put(b)

	ct := cont.NewFromMedia(c.ct)

	if req != nil {
		err = ct.Encoder.Encode(b, req)
		runtime.Must(err)
	}

	request, err := http.NewRequestWithContext(ctx, method, url, b)
	runtime.Must(err)

	request.Header.Set(content.TypeKey, ct.Media)

	response, err := c.client.Do(request)
	runtime.Must(err)

	defer response.Body.Close()

	b.Reset()

	_, err = io.Copy(b, response.Body)
	runtime.Must(err)

	// The server handlers return text on errors.
	ct = cont.NewFromMedia(response.Header.Get(content.TypeKey))
	if ct.IsText() {
		return nil, status.Error(response.StatusCode, strings.TrimSpace(b.String()))
	}

	//nolint:nilnil
	if len(b.Bytes()) == 0 {
		return nil, nil
	}

	res = make(maps.StringAny)

	err = ct.Encoder.Decode(b, &res)
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