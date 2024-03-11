package http_test

import (
	"context"
	"testing"

	"github.com/alexfalkowski/go-service/hooks"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/test"
	eh "github.com/alexfalkowski/go-service/transport/events/http"
	sh "github.com/alexfalkowski/go-service/transport/http"
	ht "github.com/alexfalkowski/go-service/transport/http/telemetry/tracer"
	events "github.com/cloudevents/sdk-go/v2"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
)

func init() {
	tracer.Register()
}

func TestSendReceive(t *testing.T) {
	Convey("Given I have a http event receiver", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		cfg := test.NewInsecureTransportConfig()

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		tcfg := test.NewDefaultTracerConfig()
		t, err := ht.NewTracer(ht.Params{Lifecycle: lc, Config: tcfg, Version: test.Version})
		So(err, ShouldBeNil)

		hs := test.NewHTTPServer(lc, logger, tcfg, cfg, m, nil)
		gs := test.NewGRPCServer(lc, logger, tcfg, cfg, false, m, nil, nil)

		test.RegisterTransport(lc, gs, hs)

		h, err := hooks.New(test.NewHook())
		So(err, ShouldBeNil)

		var event *events.Event

		err = eh.RegisterReceiver(context.Background(), hs.Mux, "/events", func(_ context.Context, e events.Event) { event = &e }, eh.WithReceiverHook(h))
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

			c.Send(ctx, e)

			Convey("Then I should receive an event", func() {
				So(event, ShouldNotBeNil)
				So(string(e.Data()), ShouldEqual, "test")
			})

			lc.RequireStop()
		})
	})
}
