package http_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/alexfalkowski/go-service/hooks"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/test"
	eh "github.com/alexfalkowski/go-service/transport/events/http"
	sh "github.com/alexfalkowski/go-service/transport/http"
	ht "github.com/alexfalkowski/go-service/transport/http/telemetry/tracer"
	events "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/protocol"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	webhooks "github.com/standard-webhooks/standard-webhooks/libraries/go"
	"go.uber.org/fx/fxtest"
)

func init() {
	tracer.Register()
}

type delRoundTripper struct {
	rt http.RoundTripper
}

func (r *delRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Del(webhooks.HeaderWebhookID)

	return r.rt.RoundTrip(req)
}

func TestSendReceive(t *testing.T) {
	Convey("Given I have a http event receiver", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		cfg := test.NewInsecureTransportConfig()

		m := test.NewMeter(lc)

		tcfg := test.NewOTLPTracerConfig()
		t, err := ht.NewTracer(ht.Params{Lifecycle: lc, Config: tcfg, Version: test.Version})
		So(err, ShouldBeNil)

		hs := test.NewHTTPServer(lc, logger, tcfg, cfg, m, nil)
		gs := test.NewGRPCServer(lc, logger, tcfg, cfg, false, m, nil, nil)

		test.RegisterTransport(lc, gs, hs)

		h, err := hooks.New(test.NewHook())
		So(err, ShouldBeNil)

		var event *events.Event

		err = eh.RegisterReceiver(context.Background(), test.Mux, "/events", func(_ context.Context, e events.Event) { event = &e }, eh.WithReceiverHook(h))
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I send an event", func() {
			rt, err := sh.NewRoundTripper(sh.WithClientLogger(logger), sh.WithClientTracer(t), sh.WithClientMetrics(m))
			So(err, ShouldBeNil)

			c, err := eh.NewSender(eh.WithSenderRoundTripper(rt), eh.WithSenderHook(h))
			So(err, ShouldBeNil)

			ctx := events.ContextWithTarget(context.Background(), "http://localhost:"+cfg.HTTP.Port+"/events")

			e := events.NewEvent()
			e.SetSource("example/uri")
			e.SetType("example.type")
			e.SetData(events.TextPlain, "test")

			result := c.Send(ctx, e)

			Convey("Then I should receive an event", func() {
				So(protocol.IsACK(result), ShouldBeTrue)
				So(event, ShouldNotBeNil)
				So(string(e.Data()), ShouldEqual, "test")
			})

			lc.RequireStop()
		})
	})
}

func TestSendNotReceive(t *testing.T) {
	Convey("Given I have a http event receiver", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		cfg := test.NewInsecureTransportConfig()
		m := test.NewMeter(lc)
		tcfg := test.NewOTLPTracerConfig()

		t, err := ht.NewTracer(ht.Params{Lifecycle: lc, Config: tcfg, Version: test.Version})
		So(err, ShouldBeNil)

		hs := test.NewHTTPServer(lc, logger, tcfg, cfg, m, nil)
		gs := test.NewGRPCServer(lc, logger, tcfg, cfg, false, m, nil, nil)

		test.RegisterTransport(lc, gs, hs)

		h, err := hooks.New(test.NewHook())
		So(err, ShouldBeNil)

		var event *events.Event

		err = eh.RegisterReceiver(context.Background(), test.Mux, "/events", func(_ context.Context, e events.Event) { event = &e }, eh.WithReceiverHook(h))
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I send an event", func() {
			rt, err := sh.NewRoundTripper(sh.WithClientLogger(logger), sh.WithClientTracer(t), sh.WithClientMetrics(m))
			So(err, ShouldBeNil)

			rt = &delRoundTripper{rt: rt}

			c, err := eh.NewSender(eh.WithSenderRoundTripper(rt), eh.WithSenderHook(h))
			So(err, ShouldBeNil)

			ctx := events.ContextWithTarget(context.Background(), "http://localhost:"+cfg.HTTP.Port+"/events")

			e := events.NewEvent()
			e.SetSource("example/uri")
			e.SetType("example.type")
			e.SetData(events.TextPlain, "test")

			result := c.Send(ctx, e)

			Convey("Then I should not receive an event", func() {
				So(protocol.IsNACK(result), ShouldBeTrue)
				So(event, ShouldBeNil)
			})

			lc.RequireStop()
		})
	})
}
