package test

import (
	"cmp"
	"context"

	"github.com/alexfalkowski/go-service/meta"
	nc "github.com/alexfalkowski/go-service/net/http/context"
)

// RestNoContent for test.
//
//nolint:nilnil
func RestNoContent(_ context.Context) (any, error) {
	return nil, nil
}

// RestNoContent for test.
func RestContent(ctx context.Context) (any, error) {
	req := nc.Request(ctx)
	name := cmp.Or(req.URL.Query().Get("name"), "Bob")
	s := "Hello " + name

	return &Response{Meta: meta.CamelStrings(ctx, ""), Greeting: &s}, nil
}
