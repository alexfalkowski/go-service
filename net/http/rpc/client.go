package rpc

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/alexfalkowski/go-service/encoding"
	"github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/net/http/content"
	"github.com/alexfalkowski/go-service/runtime"
)

// NewClient for HTTP.
func NewClient[Req any, Res any](url, contentType string, client *http.Client, enc *encoding.Map) *Client[Req, Res] {
	client.CheckRedirect = func(_ *http.Request, _ []*http.Request) error { return http.ErrUseLastResponse }
	ct := content.NewFromMedia(contentType)

	return &Client[Req, Res]{client: client, enc: enc, url: url, ct: ct}
}

// Client for HTTP.
type Client[Req any, Res any] struct {
	enc    *encoding.Map
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

	m, err := c.ct.Marshaller(enc)
	runtime.Must(errors.Prefix("rpc marshaller", err))

	d, err := m.Marshal(req)
	runtime.Must(errors.Prefix("rpc marshal", err))

	request, err := http.NewRequestWithContext(ctx, "POST", c.url, bytes.NewBuffer(d))
	runtime.Must(err)

	request.Header.Set(content.TypeKey, c.ct.Media)

	response, err := c.client.Do(request)
	runtime.Must(errors.Prefix("rpc send", err))

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	runtime.Must(errors.Prefix("rpc read", err))

	// The server handlers return text on errors.
	ct := content.NewFromMedia(response.Header.Get(content.TypeKey))
	if ct.IsText() {
		return nil, errors.Prefix("rpc error", Error(response.StatusCode, strings.TrimSpace(string(body))))
	}

	var rp Res
	res = &rp

	err = m.Unmarshal(body, res)
	runtime.Must(errors.Prefix("rpc unmarshal", err))

	return //nolint:nakedret
}
