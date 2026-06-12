package events_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/hooks"
	"github.com/alexfalkowski/go-service/v2/id/uuid"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/events"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
	transportevents "github.com/alexfalkowski/go-service/v2/transport/http/events"
	httphooks "github.com/alexfalkowski/go-service/v2/transport/http/hooks"
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
	test.RequireACK(t, result)
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
	test.RequireACK(t, result)
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
	test.RequireNACK(t, result)
	require.Nil(t, world.Event)
}

func TestReceiveCanReportProcessingFailure(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())
	world.Receiver.Register(t.Context(), "/events", func(_ context.Context, _ events.Event) events.Result {
		return events.NewHTTPResult(http.StatusServiceUnavailable, "processor unavailable")
	})

	ctx := world.EventsContext(t.Context())

	e := events.NewEvent()
	e.SetSource("example/uri")
	e.SetType("example.type")
	require.NoError(t, e.SetData(events.TextPlain, "test"))

	result := world.Sender.Send(ctx, e)
	test.RequireNACK(t, result)
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

	sender := transportevents.NewSender(httphooks.NewWebhook(hook, uuid.NewGenerator()))

	e := events.NewEvent()
	e.SetSource("example/uri")
	e.SetType("example.type")
	require.NoError(t, e.SetData(events.TextPlain, "test"))

	result := sender.Send(events.ContextWithTarget(t.Context(), trusted.URL+"/events"), e)
	test.RequireNACK(t, result)
	require.NotEmpty(t, trustedSignature)
	require.Empty(t, attackerSignature)
}

func TestSenderUsesStructuredEncoding(t *testing.T) {
	contentTypes := make(chan string, 1)
	specVersions := make(chan string, 1)
	bodies := make(chan string, 1)
	readErrors := make(chan error, 1)

	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		data, _, err := io.ReadAll(req.Body)
		contentTypes <- req.Header.Get("Content-Type")
		specVersions <- req.Header.Get("Ce-Specversion")
		bodies <- bytes.String(data)
		readErrors <- err

		res.WriteHeader(http.StatusNoContent)
	}))
	t.Cleanup(server.Close)

	sender := transportevents.NewSender(nil)

	e := events.NewEvent()
	e.SetID("event-1")
	e.SetSource("example/uri")
	e.SetType("example.type")
	require.NoError(t, e.SetData(events.TextPlain, "test"))

	result := sender.Send(events.ContextWithTarget(t.Context(), server.URL+"/events"), e)
	require.NoError(t, <-readErrors)
	test.RequireACK(t, result)
	require.Equal(t, "application/cloudevents+json", <-contentTypes)
	require.Empty(t, <-specVersions)
	body := <-bodies
	require.Contains(t, body, `"id":"event-1"`)
	require.Contains(t, body, `"type":"example.type"`)
}

func TestSenderUsesDefaultTimeout(t *testing.T) {
	start, deadline := sendWithDeadline(t)
	remaining := time.Duration(deadline.Sub(start))

	require.LessOrEqual(t, remaining, time.DefaultTimeout+time.Second)
	require.Greater(t, remaining, 29*time.Second)
}

func TestSenderUsesConfiguredTimeout(t *testing.T) {
	start, deadline := sendWithDeadline(t, transportevents.WithSenderTimeout(time.Second))
	remaining := time.Duration(deadline.Sub(start))

	require.LessOrEqual(t, remaining, time.Second+100*time.Millisecond)
	require.Greater(t, remaining, 900*time.Millisecond)
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

func sendWithDeadline(t *testing.T, opts ...transportevents.SenderOption) (time.Time, time.Time) {
	t.Helper()

	deadlines := make(chan time.Time, 1)
	rt := test.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		deadline, ok := req.Context().Deadline()
		require.True(t, ok)
		deadlines <- deadline

		return &http.Response{StatusCode: http.StatusNoContent, Body: http.NoBody, Header: http.Header{}}, nil
	})
	opts = append(opts, transportevents.WithSenderRoundTripper(rt))
	sender := transportevents.NewSender(nil, opts...)

	e := events.NewEvent()
	e.SetSource("example/uri")
	e.SetType("example.type")
	require.NoError(t, e.SetData(events.TextPlain, "test"))

	start := time.Now()
	result := sender.Send(events.ContextWithTarget(t.Context(), "http://example.com/events"), e)
	test.RequireACK(t, result)

	return start, <-deadlines
}
