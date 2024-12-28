package test

import (
	"cmp"
	"context"

	"github.com/alexfalkowski/go-service/meta"
)

// RestNoContent for test.
//
//nolint:nilnil
func RestNoContent(_ context.Context, _ *Request) (*Response, error) {
	return nil, nil
}

// RestContent for test.
func RestContent(ctx context.Context, req *Request) (*Response, error) {
	name := cmp.Or(req.Name, "Bob")
	s := "Hello " + name

	return &Response{Meta: meta.CamelStrings(ctx, ""), Greeting: s}, nil
}

// RestError for test.
func RestError(_ context.Context, _ *Request) (*Response, error) {
	return nil, ErrInvalid
}
