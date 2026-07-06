package test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/hooks"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/events"
	transportevents "github.com/alexfalkowski/go-service/v2/transport/http/events"
	httphooks "github.com/alexfalkowski/go-service/v2/transport/http/hooks"
	"github.com/stretchr/testify/require"
)

// NewEvents builds a webhook-backed CloudEvents receiver and sender using the shared hook fixture.
func NewEvents(router *http.Router, rt http.RoundTripper, generator id.Generator) (*transportevents.Receiver, *transportevents.Sender, error) {
	h, err := hooks.NewHook(FS, NewHook())
	if err != nil {
		return nil, nil, err
	}

	receiver := transportevents.NewReceiver(router, httphooks.NewWebhook(h, generator))

	sender := transportevents.NewSender(httphooks.NewWebhook(h, generator), transportevents.WithSenderRoundTripper(rt))

	return receiver, sender, nil
}

// RequireACK requires result to be a CloudEvents ACK.
func RequireACK(tb testing.TB, result events.Result) {
	tb.Helper()

	require.True(tb, events.IsACK(result), "expected CloudEvents ACK: %v", result)
}

// RequireNACK requires result to be a CloudEvents NACK.
func RequireNACK(tb testing.TB, result events.Result) {
	tb.Helper()

	require.True(tb, events.IsNACK(result), "expected CloudEvents NACK: %v", result)
}

// RegisterEvents registers an `/events` receiver that stores the last delivered event on the world.
func (w *World) RegisterEvents(ctx context.Context) {
	w.Receiver.Register(ctx, "/events", func(_ context.Context, e events.Event) events.Result {
		w.Event = &e
		return nil
	})
}

// EventsContext returns a CloudEvents context targeting the world's HTTP `/events` endpoint.
func (w *World) EventsContext(ctx context.Context) context.Context {
	return events.ContextWithTarget(ctx, w.PathServerURL("http", "events"))
}
