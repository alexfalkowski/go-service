package http

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/alexfalkowski/go-service/marshaller"
	ct "github.com/elnormous/contenttype"
)

// NewClient for HTTP.
func NewClient[Req any, Res any](client *http.Client, mar *marshaller.Map) *Client[Req, Res] {
	return &Client[Req, Res]{client: client, mar: mar}
}

// Client for HTTP.
type Client[Req any, Res any] struct {
	client *http.Client
	mar    *marshaller.Map
}

// Call for HTTP.
func (c *Client[Req, Res]) Call(ctx context.Context, url, contentType string, req *Req) (*Res, error) {
	cType, kind := c.kind(contentType)
	ma := c.mar.Get(kind)

	d, err := ma.Marshal(req)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(d))
	if err != nil {
		return nil, err
	}

	request.Header.Set(contentTypeKey, cType)

	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	t, err := ct.ParseMediaType(response.Header.Get(contentTypeKey))
	if err != nil {
		return nil, err
	}

	// The server handlers return text on errors.
	if isText(t) {
		return nil, Error(response.StatusCode, strings.TrimSpace(string(body)))
	}

	var rp Res
	ptr := &rp

	if err := ma.Unmarshal(body, ptr); err != nil {
		return nil, err
	}

	return ptr, nil
}

func (c *Client[Req, Res]) kind(contentType string) (string, string) {
	t, err := ct.ParseMediaType(contentType)
	if err != nil {
		return "application/json", "json"
	}

	return t.String(), t.Subtype
}

func isText(t ct.MediaType) bool {
	return t.Type == "text" && t.Subtype == "plain"
}
