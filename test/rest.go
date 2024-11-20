package test

import (
	"cmp"
	"context"

	"github.com/alexfalkowski/go-service/meta"
	nc "github.com/alexfalkowski/go-service/net/http/context"
)

// RestNoContent for test.
func RestNoContent(_ context.Context) any {
	return nil
}

// RestNoContent for test.
func RestContent(ctx context.Context) any {
	req := nc.Request(ctx)
	name := cmp.Or(req.URL.Query().Get("name"), "Bob")
	s := "Hello " + name

	return &Response{Meta: meta.CamelStrings(ctx, ""), Greeting: &s}
}
