package rpc

// ClientOption configures the RPC client helper constructed by NewClient.
//
// Options are applied in the order provided to NewClient. If multiple options configure the same
// field, the last one wins.
type ClientOption interface {
	apply(opts *clientOptions)
}

type clientOptions struct {
	contentType string
	accept      string
}

type clientOptionFunc func(*clientOptions)

func (f clientOptionFunc) apply(o *clientOptions) {
	f(o)
}

// WithClientContentType sets the Content-Type used for requests made by the RPC client.
//
// This value is passed through to the underlying content-aware HTTP client and is used to select the
// request encoder. Typical values include "application/json", "application/hjson", or go-service
// protobuf media types.
func WithClientContentType(ct string) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.contentType = ct
	})
}

// WithClientAccept sets the Accept media type used for responses from the RPC client.
//
// This value is passed through to the underlying content-aware HTTP client and sent as the request
// Accept header.
func WithClientAccept(accept string) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.accept = accept
	})
}

func options(opts ...ClientOption) *clientOptions {
	os := &clientOptions{}
	for _, o := range opts {
		o.apply(os)
	}
	return os
}
