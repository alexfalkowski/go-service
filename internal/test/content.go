package test

import "github.com/alexfalkowski/go-service/v2/net/http/client"

// NewContentClient constructs an HTTP client with the shared test codecs and buffer pool.
func NewContentClient(opts ...client.ClientOption) *client.Client {
	return client.NewClient(UnaryContent, StreamContent, Pool, opts...)
}
