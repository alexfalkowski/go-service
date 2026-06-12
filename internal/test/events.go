package test

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/hooks"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/events"
	transportevents "github.com/alexfalkowski/go-service/v2/transport/http/events"
	httphooks "github.com/alexfalkowski/go-service/v2/transport/http/hooks"
)

// NewEvents builds a webhook-backed CloudEvents receiver and sender using the shared hook fixture.
func NewEvents(mux *http.ServeMux, rt http.RoundTripper, generator id.Generator) (*transportevents.Receiver, *transportevents.Sender, error) {
	h, err := hooks.NewHook(FS, NewHook())
	if err != nil {
		return nil, nil, err
	}

	receiver := transportevents.NewReceiver(mux, httphooks.NewWebhook(h, generator))

	sender := transportevents.NewSender(httphooks.NewWebhook(h, generator), transportevents.WithSenderRoundTripper(rt))

	return receiver, sender, nil
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
