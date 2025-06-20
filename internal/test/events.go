package test

import (
	"context"
	"net/url"

	"github.com/alexfalkowski/go-service/v2/hooks"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/transport/http/events"
	hh "github.com/alexfalkowski/go-service/v2/transport/http/hooks"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client"
)

// NewEvents for test.
func NewEvents(mux *http.ServeMux, rt http.RoundTripper, id *id.UUID) (*events.Receiver, client.Client) {
	h, err := hooks.New(FS, NewHook())
	runtime.Must(err)

	receiver := events.NewReceiver(mux, hh.NewWebhook(h, id))

	sender, err := events.NewSender(hh.NewWebhook(h, id), events.WithSenderRoundTripper(rt))
	runtime.Must(err)

	return receiver, sender
}

// RegisterEvents for world.
func (w *World) RegisterEvents(ctx context.Context) {
	w.Receiver.Register(ctx, "/events", func(_ context.Context, e cloudevents.Event) { w.Event = &e })
}

// EventsContext for world.
func (w *World) EventsContext(ctx context.Context) context.Context {
	url, err := url.JoinPath(w.ServerURL("http"), "events")
	runtime.Must(err)

	return cloudevents.ContextWithTarget(ctx, url)
}
