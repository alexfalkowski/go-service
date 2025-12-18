package events_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	events "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/protocol"
	hooks "github.com/standard-webhooks/standard-webhooks/libraries/go"
	"github.com/stretchr/testify/require"
)

func TestSendReceiveWithRoundTripper(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldRoundTripper(http.DefaultTransport), test.WithWorldHTTP())
	world.Register()
	world.RequireStart()

	world.RegisterEvents(t.Context())

	ctx := world.EventsContext(t.Context())

	e := events.NewEvent()
	e.SetSource("example/uri")
	e.SetType("example.type")
	require.NoError(t, e.SetData(events.TextPlain, "test"))

	result := world.Sender.Send(ctx, e)
	require.True(t, protocol.IsACK(result))
	require.NotNil(t, world.Event)
	require.Equal(t, "test", bytes.String(e.Data()))

	world.RequireStop()
}

func TestSendReceiveWithoutRoundTripper(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())
	world.Register()
	world.RequireStart()

	world.RegisterEvents(t.Context())

	ctx := world.EventsContext(t.Context())

	e := events.NewEvent()
	e.SetSource("example/uri")
	e.SetType("example.type")
	require.NoError(t, e.SetData(events.TextPlain, "test"))

	result := world.Sender.Send(ctx, e)
	require.True(t, protocol.IsACK(result))
	require.NotNil(t, world.Event)
	require.Equal(t, "test", bytes.String(e.Data()))

	world.RequireStop()
}

func TestSendNotReceive(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldRoundTripper(&delRoundTripper{rt: http.DefaultTransport}), test.WithWorldHTTP())
	world.Register()
	world.RequireStart()

	world.RegisterEvents(t.Context())

	ctx := world.EventsContext(t.Context())

	e := events.NewEvent()
	e.SetSource("example/uri")
	e.SetType("example.type")
	require.NoError(t, e.SetData(events.TextPlain, "test"))

	result := world.Sender.Send(ctx, e)
	require.True(t, protocol.IsNACK(result))
	require.Nil(t, world.Event)

	world.RequireStop()
}

type delRoundTripper struct {
	rt http.RoundTripper
}

func (r *delRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Del(hooks.HeaderWebhookID)

	return r.rt.RoundTrip(req)
}
