package rpc

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/alexfalkowski/go-service/marshaller"
	"github.com/alexfalkowski/go-service/runtime"
	ct "github.com/elnormous/contenttype"
)

// NewClient for HTTP.
func NewClient[Req any, Res any](url, contentType string, client *http.Client, mar *marshaller.Map) *Client[Req, Res] {
	client.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
		return http.ErrUseLastResponse
	}

	ct, kind := kindFromContentType(contentType)
	ma := mar.Get(kind)

	return &Client[Req, Res]{client: client, mar: ma, url: url, contentType: ct, kind: kind}
}

// Client for HTTP.
type Client[Req any, Res any] struct {
	client      *http.Client
	mar         marshaller.Marshaller
	url         string
	contentType string
	kind        string
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

	request.Header.Set(contentTypeKey, c.contentType)

	response, err := c.client.Do(request)
	runtime.Must(err)

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	runtime.Must(err)

	t, err := ct.ParseMediaType(response.Header.Get(contentTypeKey))
	runtime.Must(err)

	// The server handlers return text on errors.
	if isText(t) {
		return nil, Error(response.StatusCode, strings.TrimSpace(string(body)))
	}

	var rp Res
	res = &rp

	err = c.mar.Unmarshal(body, res)
	runtime.Must(err)

	return //nolint:nakedret
}

func kindFromContentType(contentType string) (string, string) {
	t, err := ct.ParseMediaType(contentType)
	if err != nil {
		return "application/json", "json"
	}

	return t.String(), t.Subtype
}

func isText(t ct.MediaType) bool {
	return t.Type == "text" && t.Subtype == "plain"
}
