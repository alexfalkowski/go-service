package rpc

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/errors"
	http "github.com/alexfalkowski/go-service/v2/net/http/client"
	"github.com/alexfalkowski/go-service/v2/reflect"
)

// ErrInvalidRequest is returned when the request payload is nil.
var ErrInvalidRequest = errors.New("rpc: invalid request")

// ErrInvalidResponse is returned when the response target is nil.
var ErrInvalidResponse = errors.New("rpc: invalid response")

// NewClient constructs an RPC client that uses httpClient.
//
// The returned client issues RPC-style POST requests to url using the configured Content-Type and
// Accept options. Construct httpClient with [client.NewClient] to configure codecs, response buffering,
// transport, timeout, and redirect policy. The client is safe to share with other HTTP helpers.
// To ignore redirects, configure httpClient with `client.WithRedirect(client.RedirectIgnore)`.
func NewClient(url string, httpClient *http.Client, opts ...ClientOption) *Client {
	os := options(opts...)

	return &Client{client: httpClient, url: url, contentType: os.contentType, accept: os.accept}
}

// Client is an RPC client that issues POST requests using the configured content codecs.
type Client struct {
	client      *http.Client
	contentType string
	accept      string
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
// the underlying content-aware client. When configured, Accept is sent as the response media
// preference and fallback decoder selection.
//
// The res parameter is typically a pointer to the destination value (for example *MyResponse).
func (c *Client) Post(ctx context.Context, path string, req, res any) error {
	if reflect.IsNil(req) {
		return ErrInvalidRequest
	}
	if reflect.IsNil(res) {
		return ErrInvalidResponse
	}

	opts := http.Options{ContentType: c.contentType, Accept: c.accept, Request: req, Response: res}
	return c.client.Post(ctx, c.url+path, opts)
}
