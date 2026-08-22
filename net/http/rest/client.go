package rest

import http "github.com/alexfalkowski/go-service/v2/net/http/client"

// Options is an alias for [client.Options].
type Options = http.Options

// NewClient constructs a REST client that uses httpClient.
//
// Construct httpClient with [client.NewClient] to configure codecs, response buffering, transport,
// timeout, and redirect policy. The client is safe to share with other HTTP helpers.
// To ignore redirects, configure httpClient with `client.WithRedirect(client.RedirectIgnore)`.
func NewClient(client *http.Client) *Client {
	return &Client{client}
}

// Client wraps [client.Client] for REST usage.
type Client struct {
	*http.Client
}
