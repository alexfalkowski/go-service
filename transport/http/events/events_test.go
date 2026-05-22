package events_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/hooks"
	"github.com/alexfalkowski/go-service/v2/id/uuid"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/strings"
	transportevents "github.com/alexfalkowski/go-service/v2/transport/http/events"
	httphooks "github.com/alexfalkowski/go-service/v2/transport/http/hooks"
	events "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/protocol"
	webhooks "github.com/standard-webhooks/standard-webhooks/libraries/go"
	"github.com/stretchr/testify/require"
)

func TestSendReceiveWithRoundTripper(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldRoundTripper(http.DefaultTransport), test.WithWorldHTTP())

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
}

func TestSendReceiveWithoutRoundTripper(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())

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
}

func TestSendNotReceive(t *testing.T) {
	world := test.NewStartedWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldRoundTripper(&test.HeaderDeletingRoundTripper{
			RoundTripper: http.DefaultTransport,
			Header:       webhooks.HeaderWebhookID,
		}),
		test.WithWorldHTTP(),
	)

	world.RegisterEvents(t.Context())

	ctx := world.EventsContext(t.Context())

	e := events.NewEvent()
	e.SetSource("example/uri")
	e.SetType("example.type")
	require.NoError(t, e.SetData(events.TextPlain, "test"))

	result := world.Sender.Send(ctx, e)
	require.True(t, protocol.IsNACK(result))
	require.Nil(t, world.Event)
}

func TestSenderWithWebhookDoesNotFollowCrossOriginRedirect(t *testing.T) {
	var attackerSignature string
	attacker := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
		attackerSignature = req.Header.Get(webhooks.HeaderWebhookSignature)
	}))
	t.Cleanup(attacker.Close)

	var trustedSignature string
	trusted := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		trustedSignature = req.Header.Get(webhooks.HeaderWebhookSignature)
		res.Header().Set("Location", attacker.URL+"/events")
		res.WriteHeader(http.StatusTemporaryRedirect)
	}))
	t.Cleanup(trusted.Close)

	hook, err := hooks.NewHook(test.FS, test.NewHook())
	require.NoError(t, err)

	sender, err := transportevents.NewSender(httphooks.NewWebhook(hook, uuid.NewGenerator()))
	require.NoError(t, err)

	e := events.NewEvent()
	e.SetSource("example/uri")
	e.SetType("example.type")
	require.NoError(t, e.SetData(events.TextPlain, "test"))

	result := sender.Send(events.ContextWithTarget(t.Context(), trusted.URL+"/events"), e)
	require.True(t, protocol.IsNACK(result))
	require.NotEmpty(t, trustedSignature)
	require.Empty(t, attackerSignature)
}

func TestReceiveUsesServerMaxReceiveSizeBeforeWebhookVerification(t *testing.T) {
	cfg := test.NewInsecureTransportConfig()
	cfg.HTTP.MaxReceiveSize = 64

	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldTransportConfig(cfg), test.WithWorldHTTP())
	world.RegisterEvents(t.Context())
	world.Start()

	body := strings.Repeat("a", 256)
	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, world.PathServerURL("http", "events"), strings.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/cloudevents+json")

	hook, err := hooks.NewHook(test.FS, test.NewHook())
	require.NoError(t, err)
	require.NoError(t, httphooks.NewWebhook(hook, uuid.NewGenerator()).Sign(req))

	res, err := world.Do(req)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, res.Body.Close()) })

	require.Equal(t, http.StatusRequestEntityTooLarge, res.StatusCode)
	require.Nil(t, world.Event)
}
