package rpc

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/net/http/content"
	"github.com/alexfalkowski/go-service/net/http/status"
	"github.com/alexfalkowski/go-service/runtime"
)

// NewClient for HTTP.
func NewClient[Req any, Res any](url, ct string, client *http.Client) *Client[Req, Res] {
	client.CheckRedirect = func(_ *http.Request, _ []*http.Request) error { return http.ErrUseLastResponse }

	return &Client[Req, Res]{client: client, url: url, ct: content.NewFromMedia(ct)}
}

// Client for HTTP.
type Client[Req any, Res any] struct {
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

	return //nolint:nakedret
}
