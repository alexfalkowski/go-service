package rest

import (
	"net/http"
	"strings"
	"time"

	nh "github.com/alexfalkowski/go-service/net/http"
	"github.com/alexfalkowski/go-service/net/http/content"
	"github.com/alexfalkowski/go-service/net/http/status"
	st "github.com/alexfalkowski/go-service/time"
	"github.com/go-resty/resty/v2"
)

// Error will return an error if the response from the server is text (as server handlers return text on errors).
func Error(response *resty.Response) error {
	ct := cont.NewFromMedia(response.Header().Get(content.TypeKey.Value()))
	if ct.IsText() {
		return status.Error(response.StatusCode(), strings.TrimSpace(string(response.Body())))
	}

	return nil
}

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
		o.timeout = st.MustParseDuration(timeout)
	})
}

// NewClient for rest.
func NewClient(opts ...ClientOption) *resty.Client {
	os := options(opts...)
	client := &http.Client{Transport: os.roundTripper, Timeout: os.timeout}

	return resty.NewWithClient(client)
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
