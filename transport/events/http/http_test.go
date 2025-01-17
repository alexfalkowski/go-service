//nolint:varnamelen
package http_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/alexfalkowski/go-service/hooks"
	"github.com/alexfalkowski/go-service/test"
	eh "github.com/alexfalkowski/go-service/transport/events/http"
	sh "github.com/alexfalkowski/go-service/transport/http"
	events "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/protocol"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	h "github.com/standard-webhooks/standard-webhooks/libraries/go"
	"go.uber.org/fx/fxtest"
)

func TestSendReceiveWithRoundTripper(t *testing.T) {
	Convey("Given I have a http event receiver", t, func() {
		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		cfg := test.NewInsecureTransportConfig()
		m := test.NewOTLPMeter(lc)
		tc := test.NewOTLPTracerConfig()

		s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Mux: mux}
		s.Register()

		h, err := hooks.New(test.NewHook())
		So(err, ShouldBeNil)

		r := eh.NewReceiver(mux, h)

		var event *events.Event

		r.Register(context.Background(), "/events", func(_ context.Context, e events.Event) { event = &e })
		lc.RequireStart()

		Convey("When I send an event", func() {
			tracer := test.NewTracer(lc, tc, logger)
			rt := sh.NewRoundTripper(sh.WithClientLogger(logger), sh.WithClientTracer(tracer), sh.WithClientMetrics(m))

			c, err := eh.NewSender(h, eh.WithSenderRoundTripper(rt))
			So(err, ShouldBeNil)

			ctx := events.ContextWithTarget(context.Background(), fmt.Sprintf("http://%s/events", cfg.HTTP.Address))

			e := events.NewEvent()
			e.SetSource("example/uri")
			e.SetType("example.type")

			err = e.SetData(events.TextPlain, "test")
			So(err, ShouldBeNil)

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

func TestSendReceiveWithoutRoundTripper(t *testing.T) {
	Convey("Given I have a http event receiver", t, func() {
		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		cfg := test.NewInsecureTransportConfig()
		m := test.NewOTLPMeter(lc)
		tc := test.NewOTLPTracerConfig()

		s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Mux: mux}
		s.Register()

		h, err := hooks.New(test.NewHook())
		So(err, ShouldBeNil)

		r := eh.NewReceiver(mux, h)

		var event *events.Event

		r.Register(context.Background(), "/events", func(_ context.Context, e events.Event) { event = &e })
		lc.RequireStart()

		Convey("When I send an event", func() {
			c, err := eh.NewSender(h)
			So(err, ShouldBeNil)

			ctx := events.ContextWithTarget(context.Background(), fmt.Sprintf("http://%s/events", cfg.HTTP.Address))

			e := events.NewEvent()
			e.SetSource("example/uri")
			e.SetType("example.type")

			err = e.SetData(events.TextPlain, "test")
			So(err, ShouldBeNil)

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
		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		cfg := test.NewInsecureTransportConfig()
		m := test.NewOTLPMeter(lc)
		tc := test.NewOTLPTracerConfig()

		s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Mux: mux}
		s.Register()

		h, err := hooks.New(test.NewHook())
		So(err, ShouldBeNil)

		r := eh.NewReceiver(mux, h)

		var event *events.Event

		r.Register(context.Background(), "/events", func(_ context.Context, e events.Event) { event = &e })
		lc.RequireStart()

		Convey("When I send an event", func() {
			tracer := test.NewTracer(lc, tc, logger)
			rt := sh.NewRoundTripper(sh.WithClientLogger(logger), sh.WithClientTracer(tracer), sh.WithClientMetrics(m))
			rt = &delRoundTripper{rt: rt}

			c, err := eh.NewSender(h, eh.WithSenderRoundTripper(rt))
			So(err, ShouldBeNil)

			ctx := events.ContextWithTarget(context.Background(), fmt.Sprintf("http://%s/events", cfg.HTTP.Address))

			e := events.NewEvent()
			e.SetSource("example/uri")
			e.SetType("example.type")

			err = e.SetData(events.TextPlain, "test")
			So(err, ShouldBeNil)

			result := c.Send(ctx, e)

			Convey("Then I should not receive an event", func() {
				So(protocol.IsNACK(result), ShouldBeTrue)
				So(event, ShouldBeNil)
			})

			lc.RequireStop()
		})
	})
}

type delRoundTripper struct {
	rt http.RoundTripper
}

func (r *delRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Del(h.HeaderWebhookID)

	return r.rt.RoundTrip(req)
}
