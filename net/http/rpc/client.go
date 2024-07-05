package rpc

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/alexfalkowski/go-service/marshaller"
	"github.com/alexfalkowski/go-service/net/http/content"
	"github.com/alexfalkowski/go-service/runtime"
)

// NewClient for HTTP.
func NewClient[Req any, Res any](url, contentType string, client *http.Client, mar *marshaller.Map) *Client[Req, Res] {
	client.CheckRedirect = func(_ *http.Request, _ []*http.Request) error { return http.ErrUseLastResponse }
	ct := content.NewFromMedia(contentType)

	return &Client[Req, Res]{client: client, mar: ct.Marshaller(mar), url: url, ct: ct}
}

// Client for HTTP.
type Client[Req any, Res any] struct {
	mar    marshaller.Marshaller
	client *http.Client
	ct     *content.Type
	url    string
}

// Call for HTTP.
//
//nolint:nonamedreturns
func (c *Client[Req, Res]) Call(ctx context.Context, req *Req) (res *Res, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = runtime.ConvertRecover(r)
		}
	}()

	d, err := c.mar.Marshal(req)
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
		return nil, Error(response.StatusCode, strings.TrimSpace(string(body)))
	}

	var rp Res
	res = &rp

	err = c.mar.Unmarshal(body, res)
	runtime.Must(err)

	return //nolint:nakedret
}
