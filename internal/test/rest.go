package test

import (
	"cmp"
	"context"

	"github.com/alexfalkowski/go-service/meta"
	hm "github.com/alexfalkowski/go-service/net/http/meta"
)

// RestNoContent for test.
func RestNoContent(_ context.Context) (*Response, error) {
	return nil, nil
}

// RestRequestNoContent for test.
func RestRequestNoContent(_ context.Context, _ *Request) (*Response, error) {
	return nil, nil
}

// RestNoContent for test.
func RestContent(ctx context.Context) (*Response, error) {
	req := hm.Request(ctx)
	_ = hm.Response(ctx)
	name := cmp.Or(req.URL.Query().Get("name"), "Bob")
	s := "Hello " + name

	return &Response{Meta: meta.CamelStrings(ctx, ""), Greeting: s}, nil
}

// RestRequestContent for test.
func RestRequestContent(ctx context.Context, req *Request) (*Response, error) {
	name := cmp.Or(req.Name, "Bob")
	s := "Hello " + name

	return &Response{Meta: meta.CamelStrings(ctx, ""), Greeting: s}, nil
}

// RestError for test.
func RestError(_ context.Context) (*Response, error) {
	return nil, ErrInvalid
}

// RestRequestError for test.
func RestRequestError(_ context.Context, _ *Request) (*Response, error) {
	return nil, ErrInvalid
}
